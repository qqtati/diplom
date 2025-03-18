package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn     *websocket.Conn
	RoomID   string
	UserID   string
	IsTutor  bool
	PeerConn *PeerConnection
}

type PeerConnection struct {
	mu sync.Mutex
	// Здесь можно добавить дополнительные поля для хранения состояния соединения
}

type Room struct {
	mu      sync.Mutex
	Clients map[string]*Client
}

var (
	rooms = make(map[string]*Room)
	mu    sync.RWMutex
)

func getOrCreateRoom(roomID string) *Room {
	mu.Lock()
	defer mu.Unlock()

	if room, exists := rooms[roomID]; exists {
		return room
	}

	room := &Room{
		Clients: make(map[string]*Client),
	}
	rooms[roomID] = room
	return room
}

func HandleWebRTC(c *websocket.Conn) {
	roomID := c.Query("room_id")
	userID := c.Query("user_id")
	isTutor := c.Query("is_tutor") == "true"

	log.Printf("Новое WebRTC подключение: room_id=%s, user_id=%s, is_tutor=%v", roomID, userID, isTutor)

	client := &Client{
		Conn:     c,
		RoomID:   roomID,
		UserID:   userID,
		IsTutor:  isTutor,
		PeerConn: &PeerConnection{},
	}

	room := getOrCreateRoom(roomID)
	room.mu.Lock()
	room.Clients[userID] = client
	log.Printf("Клиент добавлен в комнату %s. Всего клиентов: %d", roomID, len(room.Clients))
	room.mu.Unlock()

	defer func() {
		room.mu.Lock()
		delete(room.Clients, userID)
		log.Printf("Клиент удален из комнаты %s. Осталось клиентов: %d", roomID, len(room.Clients))
		room.mu.Unlock()
		c.Close()
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Ошибка чтения сообщения от клиента %s: %v", userID, err)
			break
		}

		// Логируем входящее сообщение
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			log.Printf("Получено сообщение от клиента %s: тип=%v, содержимое=%+v", userID, msg["type"], msg)
		} else {
			log.Printf("Ошибка разбора сообщения от клиента %s: %v", userID, err)
		}

		// Пересылаем сообщение всем клиентам в комнате, кроме отправителя
		room.mu.Lock()
		for _, otherClient := range room.Clients {
			if otherClient.UserID != userID {
				log.Printf("Пересылка сообщения от %s к %s: %s", userID, otherClient.UserID, string(message))
				if err := otherClient.Conn.WriteMessage(messageType, message); err != nil {
					log.Printf("Ошибка отправки сообщения клиенту %s: %v", otherClient.UserID, err)
				}
			}
		}
		room.mu.Unlock()
	}
}

func SetupWebRTCRoutes(app *fiber.App) {
	fmt.Println("Настройка WebRTC маршрутов")

	app.Use("/ws/webrtc", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			fmt.Printf("WebSocket upgrade запрос от %s", c.IP())
			c.Locals("allowed", true)
			return c.Next()
		}
		fmt.Printf("Отклонен не-WebSocket запрос от %s", c.IP())
		return fiber.ErrUpgradeRequired
	})

	app.Get("/webrtc", websocket.New(HandleWebRTC))
	fmt.Println("WebRTC маршруты настроены")
}
