package http

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	ws "backend/internal/models/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
}

type Room struct {
	ID             string
	Clients        map[string]*Client
	ChangeHistory  []ws.Message
	MaxHistorySize int
	mu             sync.Mutex
}

type WebSocketHandler struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		rooms: make(map[string]*Room),
	}
}

func (h *WebSocketHandler) HandleWebSocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		roomID := c.Params("roomId")
		if roomID == "" {
			return
		}

		// Получаем имя пользователя из контекста
		username := c.Locals("username").(string)
		clientID := c.Query("clientId")
		if clientID == "" {
			clientID = username
		}

		// Получаем или создаем комнату
		room := h.getOrCreateRoom(roomID)

		// Создаем нового клиента
		client := &Client{
			ID:       clientID,
			Username: username,
			Conn:     c,
		}

		// Добавляем клиента в комнату
		room.mu.Lock()
		room.Clients[client.ID] = client
		room.mu.Unlock()

		// Отправляем историю изменений новому клиенту
		if err := h.sendChangeHistory(client, room); err != nil {
			log.Printf("error sending change history: %v", err)
		}

		// Отправляем сообщение о присоединении
		joinMsg := ws.Message{
			Type:      ws.MessageTypeJoin,
			RoomID:    roomID,
			ClientID:  client.ID,
			Username:  client.Username,
			Timestamp: time.Now().UnixNano(),
		}
		h.broadcastToRoom(room, joinMsg)

		defer func() {
			// Удаляем клиента при отключении
			room.mu.Lock()
			delete(room.Clients, client.ID)
			room.mu.Unlock()

			// Отправляем сообщение об отключении
			leaveMsg := ws.Message{
				Type:      ws.MessageTypeLeave,
				RoomID:    roomID,
				ClientID:  client.ID,
				Username:  client.Username,
				Timestamp: time.Now().UnixNano(),
			}
			h.broadcastToRoom(room, leaveMsg)

			// Если комната пуста, удаляем ее
			if len(room.Clients) == 0 {
				h.mu.Lock()
				delete(h.rooms, roomID)
				h.mu.Unlock()
			}
		}()

		// Основной цикл обработки сообщений
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("error reading message: %v", err)
				break
			}

			var message ws.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("error unmarshaling message: %v", err)
				continue
			}

			message.ClientID = client.ID
			message.Username = client.Username
			message.RoomID = roomID
			message.Timestamp = time.Now().UnixNano()

			// Если это сообщение с изменениями, сохраняем его в истории
			if message.Type == ws.MessageTypeChanges {
				room.mu.Lock()
				room.ChangeHistory = append(room.ChangeHistory, message)
				// Ограничиваем размер истории
				if len(room.ChangeHistory) > room.MaxHistorySize {
					room.ChangeHistory = room.ChangeHistory[1:]
				}
				room.mu.Unlock()
			}

			// Рассылаем сообщение всем клиентам в комнате
			h.broadcastToRoom(room, message)
		}
	})
}

func (h *WebSocketHandler) sendChangeHistory(client *Client, room *Room) error {
	room.mu.Lock()
	history := room.ChangeHistory
	room.mu.Unlock()

	if len(history) == 0 {
		return nil
	}

	// Создаем специальное сообщение с историей изменений
	historyMsg := ws.Message{
		Type:      ws.MessageTypeChanges,
		RoomID:    room.ID,
		ClientID:  "system",
		Username:  "system",
		Timestamp: time.Now().UnixNano(),
		Payload: map[string]interface{}{
			"history":       history,
			"isInitialSync": true,
		},
	}

	messageJSON, err := json.Marshal(historyMsg)
	if err != nil {
		return err
	}

	return client.Conn.WriteMessage(websocket.TextMessage, messageJSON)
}

func (h *WebSocketHandler) getOrCreateRoom(roomID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, exists := h.rooms[roomID]; exists {
		return room
	}

	room := &Room{
		ID:             roomID,
		Clients:        make(map[string]*Client),
		ChangeHistory:  make([]ws.Message, 0),
		MaxHistorySize: 1000, // Максимальное количество сохраняемых изменений
	}
	h.rooms[roomID] = room
	return room
}

func (h *WebSocketHandler) broadcastToRoom(room *Room, message ws.Message) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	for _, client := range room.Clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
			log.Printf("error sending message to client: %v", err)
		}
	}
}

func MapWebSocketRoutes(r fiber.Router, h *WebSocketHandler, mw *UserHandler) {
	r.Get("/:roomId", mw.CheckAccessToken(), h.HandleWebSocket())
}
