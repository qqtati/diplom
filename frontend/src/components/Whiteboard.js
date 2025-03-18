import React, { useEffect, useState, useCallback, useContext } from 'react';
import { Tldraw, useEditor, createTLStore } from 'tldraw';
import 'tldraw/tldraw.css';
import { message } from 'antd';
import AuthContext from '../context/AuthContext';

const WEBSOCKET_URL = 'ws://localhost:9001/ws';

const Whiteboard = ({ roomId }) => {
    const [store] = useState(() => createTLStore());
    const [ws, setWs] = useState(null);
    const { user } = useContext(AuthContext);

    // Инициализация WebSocket соединения
    useEffect(() => {
        if (!user) {
            message.error('Необходимо авторизоваться');
            return;
        }

        const accessToken = localStorage.getItem('accessToken');
        if (!accessToken) {
            message.error('Токен авторизации не найден');
            return;
        }

        const socket = new WebSocket(`${WEBSOCKET_URL}/${roomId}`);
        
        socket.onopen = () => {
            console.log('WebSocket соединение установлено');
        };

        socket.onclose = (event) => {
            if (event.code === 1000) {
                console.log('WebSocket соединение закрыто');
            } else {
                message.error('Соединение прервано. Проверьте авторизацию.');
            }
        };

        socket.onerror = (error) => {
            console.error('WebSocket ошибка:', error);
            message.error('Ошибка соединения с сервером');
        };

        setWs(socket);

        return () => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.close();
            }
        };
    }, [roomId, user]);

    // Обработка входящих сообщений
    useEffect(() => {
        if (!ws) return;

        ws.onmessage = (event) => {
            try {
                console.log('Raw message received:', event.data);
                const data = JSON.parse(event.data);
                console.log('Parsed message:', data);
                
                if (data.type === 'BOARD_UPDATE' && data.payload) {
                    console.log('Processing BOARD_UPDATE:', data.payload);
                    store.mergeRemoteChanges(() => {
                        if (typeof data.payload === 'object' && data.payload !== null) {
                            // Сначала собираем все изменения
                            const toAdd = [];
                            const toRemove = [];

                            Object.entries(data.payload).forEach(([id, change]) => {
                                console.log('Processing change for id:', id, 'change:', change);
                                if (change && (change.type === 'create' || change.type === 'update')) {
                                    console.log('Adding/Updating record:', change.record);
                                    toAdd.push(change.record);
                                } else if (change && change.type === 'delete') {
                                    console.log('Removing id:', id);
                                    toRemove.push(id);
                                }
                            });

                            console.log('Collected changes - toAdd:', toAdd, 'toRemove:', toRemove);

                            // Затем применяем их пакетно
                            if (toAdd.length > 0) {
                                console.log('Putting records:', toAdd);
                                store.put(toAdd);
                            }
                            if (toRemove.length > 0) {
                                console.log('Removing records:', toRemove);
                                store.remove(toRemove);
                            }

                            // Проверяем состояние после изменений
                            console.log('Store state after changes:', store.allRecords());
                        } else {
                            console.error('Invalid payload format:', data.payload);
                        }
                    });
                } else if (data.type === 'ERROR') {
                    message.error(data.message || 'Произошла ошибка');
                }
            } catch (error) {
                console.error('Ошибка при обработке сообщения:', error);
                console.error('Raw message:', event.data);
            }
        };
    }, [ws, store]);

    // Отправка изменений на сервер
    const handleStoreChange = useCallback((update) => {
        if (!ws || ws.readyState !== WebSocket.OPEN) return;

        console.log('Store change update:', update);
        const changes = {};
        
        // Обработка добавленных элементов
        if (update.changes.added && Object.keys(update.changes.added).length > 0) {
            console.log('Processing added changes:', update.changes.added);
            Object.entries(update.changes.added).forEach(([id, record]) => {
                changes[id] = {
                    type: 'create',
                    record: record
                };
            });
        }

        // Обработка обновленных элементов
        if (update.changes.updated && Object.keys(update.changes.updated).length > 0) {
            console.log('Processing updated changes:', update.changes.updated);
            Object.entries(update.changes.updated).forEach(([id, record]) => {
                if (id == 'pointer:pointer' || id == 'camera:page:page') {
                    return ;
                }
                changes[id] = {
                    type: 'update',
                    record: record[0]
                };
            });
        }

        // Обработка удаленных элементов
        if (update.changes.removed && Object.keys(update.changes.removed).length > 0) {
            console.log('Processing removed changes:', update.changes.removed);
            Object.entries(update.changes.removed).forEach(([id]) => {
                changes[id] = {
                    type: 'delete'
                };
            });
        }

        // Отправляем изменения только если они есть
        if (Object.keys(changes).length > 0) {
            const message = {
                type: 'BOARD_UPDATE',
                room: roomId,
                payload: changes
            };
            
            try {
                console.log('Отправка сообщения на сервер:', message);
                ws.send(JSON.stringify(message));
                console.log('Сообщение успешно отправлено');
            } catch (error) {
                console.error('Ошибка при отправке изменений:', error);
                message.error('Не удалось отправить изменения');
            }
        } else {
            console.log('Нет изменений для отправки');
        }
    }, [ws, roomId]);

    // Подписка на изменения store
    useEffect(() => {
        if (!store) return;

        const unsubscribe = store.listen(handleStoreChange, { source: 'user' });
        return () => {
            unsubscribe();
        };
    }, [store, handleStoreChange]);

    if (!user) {
        return null;
    }

    return (
        <div style={{ position: 'fixed', inset: 0 }}>
            <Tldraw
                store={store}
                autoFocus
                showMenu={false}
                showPages={false}
                showTools={true}
                showZoom={true}
                showStyles={true}
                showUI={true}
            />
        </div>
    );
};

export default Whiteboard; 
