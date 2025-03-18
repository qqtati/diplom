import React, { useState, useEffect } from 'react';
import { Input, Button, Space, Card, Typography, message } from 'antd';
import { useHistory, useLocation } from 'react-router-dom';
import Whiteboard from '../components/Whiteboard';

const { Title } = Typography;

const WhiteboardPage = () => {
    const [roomId, setRoomId] = useState('');
    const [isJoined, setIsJoined] = useState(false);
    const history = useHistory();
    const location = useLocation();

    // Проверяем URL при загрузке страницы
    useEffect(() => {
        const params = new URLSearchParams(location.search);
        const room = params.get('room');
        if (room) {
            setRoomId(room);
            setIsJoined(true);
        }
    }, [location]);

    const handleJoinRoom = () => {
        if (!roomId.trim()) {
            message.error('Введите ID комнаты');
            return;
        }
        
        // Добавляем ID комнаты в URL
        history.push(`/whiteboard?room=${roomId}`);
        setIsJoined(true);
    };

    const handleLeaveRoom = () => {
        setIsJoined(false);
        setRoomId('');
        history.push('/whiteboard');
    };

    if (isJoined) {
        return (
            <div style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
                <div style={{ padding: '16px', background: '#fff', borderBottom: '1px solid #f0f0f0' }}>
                    <Space>
                        <Title level={5} style={{ margin: 0 }}>Комната: {roomId}</Title>
                        <Button onClick={handleLeaveRoom} type="primary" danger>
                            Покинуть комнату
                        </Button>
                    </Space>
                </div>
                <div style={{ flex: 1, position: 'relative' }}>
                    <Whiteboard roomId={roomId} />
                </div>
            </div>
        );
    }

    return (
        <div style={{ 
            height: '100vh', 
            display: 'flex', 
            justifyContent: 'center', 
            alignItems: 'center', 
            background: '#f0f2f5' 
        }}>
            <Card style={{ width: 400 }}>
                <Title level={3} style={{ textAlign: 'center', marginBottom: 24 }}>
                    Интерактивная доска
                </Title>
                <Space direction="vertical" style={{ width: '100%' }}>
                    <Input
                        placeholder="Введите ID комнаты"
                        value={roomId}
                        onChange={(e) => setRoomId(e.target.value)}
                        onPressEnter={handleJoinRoom}
                        size="large"
                    />
                    <Button 
                        type="primary" 
                        onClick={handleJoinRoom}
                        size="large"
                        block
                    >
                        Присоединиться
                    </Button>
                </Space>
            </Card>
        </div>
    );
};

export default WhiteboardPage;