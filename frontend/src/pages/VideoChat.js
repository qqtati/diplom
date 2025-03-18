// src/pages/VideoChat.js
import React, { useState, useRef, useEffect } from 'react';
import { Button, Space, message, Input, Card, Typography, Avatar } from 'antd';
import { useParams, useHistory } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { UserOutlined, VideoCameraOutlined, VideoCameraFilled } from '@ant-design/icons';

const { Title } = Typography;

const VideoChat = () => {
    // Получаем параметры из URL и контекст авторизации
    const { roomId: urlRoomId } = useParams();
    const { user } = useAuth();
    const history = useHistory();

    // Состояния
    const [roomId, setRoomId] = useState(urlRoomId || '');
    const [isJoined, setIsJoined] = useState(!!urlRoomId);
    const [isConnected, setIsConnected] = useState(false);
    const [isMicOn, setIsMicOn] = useState(true);
    const [isAudioOn, setIsAudioOn] = useState(true);
    const [isScreenSharing, setIsScreenSharing] = useState(false);
    const [isVideoOn, setIsVideoOn] = useState(true);
    const [error, setError] = useState(null);

    // Рефы
    const localVideoRef = useRef(null);
    const remoteVideoRef = useRef(null);
    const peerConnectionRef = useRef(null);
    const wsRef = useRef(null);
    const localStreamRef = useRef(null);
    const screenStreamRef = useRef(null);

    // Конфигурация WebRTC
    const rtcConfiguration = {
        iceServers: [
            { urls: 'stun:stun.l.google.com:19302' },
            { urls: 'stun:stun1.l.google.com:19302' }
        ]
    };

    // Добавляем эффект для отслеживания изменений isJoined
    useEffect(() => {
        console.log('isJoined', isJoined);
        console.log('roomId', roomId);
        console.log('user', user);
        console.log('localStreamRef.current', localStreamRef.current);
        if (isJoined && roomId && user?.result && localStreamRef.current) {
            console.log('Состояние готово для подключения:', {
                isJoined,
                roomId,
                hasUser: !!user,
                hasUserResult: !!user?.result,
                hasLocalStream: !!localStreamRef.current
            });
            connectWebSocket();
        }
    }, [isJoined, roomId, user, localStreamRef.current]);

    // Добавляем проверку и полифилл для mediaDevices
    useEffect(() => {
        // Проверяем поддержку mediaDevices
        if (!navigator.mediaDevices) {
            console.log('mediaDevices не поддерживается, добавляем полифилл');
            // Полифилл для старых браузеров
            navigator.mediaDevices = {};
        }

        // Полифилл для getUserMedia
        if (!navigator.mediaDevices.getUserMedia) {
            navigator.mediaDevices.getUserMedia = function(constraints) {
                const getUserMedia = navigator.webkitGetUserMedia || navigator.mozGetUserMedia;

                if (!getUserMedia) {
                    return Promise.reject(new Error('getUserMedia не поддерживается в этом браузере'));
                }

                return new Promise((resolve, reject) => {
                    getUserMedia.call(navigator, constraints, resolve, reject);
                });
            };
        }

        // Полифилл для enumerateDevices
        if (!navigator.mediaDevices.enumerateDevices) {
            navigator.mediaDevices.enumerateDevices = function() {
                return new Promise((resolve) => {
                    // Для старых браузеров возвращаем пустой массив
                    resolve([]);
                });
            };
        }

        // Проверяем доступность камеры и микрофона
        const checkDevices = async () => {
            try {
                // Сначала запрашиваем разрешение на доступ к медиа устройствам
                const stream = await navigator.mediaDevices.getUserMedia({ 
                    video: true, 
                    audio: true 
                });
                
                // Останавливаем тестовый поток
                stream.getTracks().forEach(track => track.stop());

                // Теперь можем безопасно перечислить устройства
                const devices = await navigator.mediaDevices.enumerateDevices();
                const hasVideo = devices.some(device => device.kind === 'videoinput');
                const hasAudio = devices.some(device => device.kind === 'audioinput');
                
                console.log('Доступные устройства:', {
                    hasVideo,
                    hasAudio,
                    devices: devices.map(d => ({
                        kind: d.kind,
                        label: d.label || 'Неизвестное устройство',
                        deviceId: d.deviceId
                    }))
                });

                if (!hasVideo || !hasAudio) {
                    setError('Камера или микрофон недоступны. Пожалуйста, проверьте настройки браузера.');
                }
            } catch (error) {
                console.error('Ошибка при проверке устройств:', error);
                if (error.name === 'NotAllowedError') {
                    setError('Доступ к камере и микрофону запрещен. Пожалуйста, разрешите доступ в настройках браузера.');
                } else if (error.name === 'NotFoundError') {
                    setError('Камера или микрофон не найдены. Пожалуйста, проверьте подключение устройств.');
                } else {
                    setError('Не удалось получить доступ к медиа устройствам: ' + error.message);
                }
            }
        };

        checkDevices();
    }, []);

    // Обработчики событий
    const handleJoinRoom = async () => {
        if (!roomId.trim()) {
            message.error('Введите ID комнаты');
            return;
        }

        if (!user?.result) {
            message.error('Не удалось получить информацию о пользователе');
            return;
        }

        if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
            message.error('Ваш браузер не поддерживает доступ к камере и микрофону');
            return;
        }

        try {
            console.log('Начало подключения к комнате:', {
                roomId,
                userId: user.result.id,
                userRole: user.result.role,
                hasMediaDevices: !!navigator.mediaDevices,
                hasGetUserMedia: !!navigator.mediaDevices.getUserMedia
            });

            // Получаем доступ к медиа устройствам
            const stream = await navigator.mediaDevices.getUserMedia({
                video: {
                    width: { ideal: 640 },
                    height: { ideal: 480 },
                    facingMode: 'user',
                    frameRate: { ideal: 30 }
                },
                audio: {
                    echoCancellation: false,
                    noiseSuppression: false,
                    autoGainControl: false
                }
            }).catch(error => {
                console.error('Ошибка при получении доступа к медиа устройствам:', error);
                if (error.name === 'NotAllowedError') {
                    throw new Error('Доступ к камере и микрофону запрещен. Пожалуйста, разрешите доступ в настройках браузера.');
                } else if (error.name === 'NotFoundError') {
                    throw new Error('Камера или микрофон не найдены. Пожалуйста, проверьте подключение устройств.');
                } else {
                    throw error;
                }
            });

            console.log('Получен доступ к медиа устройствам:', {
                videoTracks: stream.getVideoTracks().map(track => ({
                    enabled: track.enabled,
                    muted: track.muted,
                    readyState: track.readyState,
                    label: track.label
                })),
                audioTracks: stream.getAudioTracks().map(track => ({
                    enabled: track.enabled,
                    muted: track.muted,
                    readyState: track.readyState,
                    label: track.label
                }))
            });

            // Сохраняем локальный поток
            localStreamRef.current = stream;

            // Устанавливаем поток в видео элемент
            if (localVideoRef.current) {
                localVideoRef.current.srcObject = stream;
                // Добавляем обработчики событий для видео элемента
                localVideoRef.current.onloadedmetadata = () => {
                    console.log('Локальное видео метаданные загружены');
                    localVideoRef.current.play().catch(error => {
                        console.error('Ошибка воспроизведения локального видео:', error);
                    });
                };
                localVideoRef.current.onerror = (error) => {
                    console.error('Ошибка локального видео:', error);
                };
            }

            // Обновляем URL
            history.push(`/video-chat/${roomId}`);
            
            // Устанавливаем состояние isJoined в конце
            setIsJoined(true);
        } catch (error) {
            console.error('Ошибка доступа к медиа устройствам:', error);
            message.error(error.message || 'Не удалось получить доступ к камере или микрофону');
        }
    };

    useEffect(() => {
        if (urlRoomId) {
            handleJoinRoom();
        }
    }, [user]);

    // Обработчик изменения ID комнаты
    const handleRoomIdChange = (e) => {
        const newRoomId = e.target.value;
        setRoomId(newRoomId);
        // Если мы уже в комнате, но меняем ID, сбрасываем состояние
        if (isJoined) {
            setIsJoined(false);
            cleanup();
        }
    };

    const handleLeaveRoom = () => {
        cleanup();
        setIsJoined(false);
        setRoomId('');
        history.push('/video-chat/new');
    };

    const toggleMicrophone = () => {
        if (localStreamRef.current) {
            const audioTrack = localStreamRef.current.getAudioTracks()[0];
            if (audioTrack) {
                audioTrack.enabled = !isMicOn;
                setIsMicOn(!isMicOn);
                console.log('Состояние микрофона изменено:', {
                    enabled: !isMicOn,
                    muted: audioTrack.muted,
                    readyState: audioTrack.readyState
                });
            }
        }
    };

    const toggleVideo = () => {
        if (localStreamRef.current) {
            const videoTrack = localStreamRef.current.getVideoTracks()[0];
            if (videoTrack) {
                videoTrack.enabled = !isVideoOn;
                setIsVideoOn(!isVideoOn);
                console.log('Состояние видео изменено:', {
                    enabled: !isVideoOn,
                    muted: videoTrack.muted,
                    readyState: videoTrack.readyState
                });
            }
        }
    };

    const toggleRemoteAudio = () => {
        if (remoteVideoRef.current) {
            remoteVideoRef.current.muted = !isAudioOn;
            setIsAudioOn(!isAudioOn);
            console.log('Состояние удаленного аудио изменено:', {
                muted: !isAudioOn,
                volume: remoteVideoRef.current.volume
            });
        }
    };

    const startScreenShare = async () => {
        try {
            const screenStream = await navigator.mediaDevices.getDisplayMedia({
                video: {
                    width: { ideal: 1920 },
                    height: { ideal: 1080 }
                }
            });

            screenStreamRef.current = screenStream;
            const videoTrack = screenStream.getVideoTracks()[0];

            // Заменяем только видео трек в peer connection, сохраняя аудио
            const sender = peerConnectionRef.current?.getSenders()
                .find(s => s.track?.kind === 'video');
            
            if (sender) {
                await sender.replaceTrack(videoTrack);
            }

            // Создаем новый MediaStream только с видео из screenStream
            const screenVideoStream = new MediaStream([videoTrack]);
            
            // Отображаем экран в локальном видео
            if (localVideoRef.current) {
                // Сохраняем аудио трек из локального потока
                const audioTrack = localStreamRef.current?.getAudioTracks()[0];
                if (audioTrack) {
                    // Создаем новый поток, объединяя видео с экрана и аудио с микрофона
                    const combinedStream = new MediaStream([videoTrack, audioTrack]);
                    localVideoRef.current.srcObject = combinedStream;
                } else {
                    localVideoRef.current.srcObject = screenVideoStream;
                }
            }

            setIsScreenSharing(true);

            // Обработчик завершения демонстрации экрана
            videoTrack.onended = () => {
                stopScreenShare();
            };
        } catch (error) {
            console.error('Ошибка при демонстрации экрана:', error);
            message.error('Не удалось начать демонстрацию экрана');
        }
    };

    const stopScreenShare = () => {
        if (screenStreamRef.current) {
            const videoTrack = screenStreamRef.current.getVideoTracks()[0];
            const sender = peerConnectionRef.current?.getSenders()
                .find(s => s.track?.kind === 'video');
            
            if (sender && localStreamRef.current) {
                const localVideoTrack = localStreamRef.current.getVideoTracks()[0];
                if (localVideoTrack) {
                    sender.replaceTrack(localVideoTrack);
                    localVideoTrack.enabled = isVideoOn;
                }
            }

            // Останавливаем только видео треки из screenStream
            screenStreamRef.current.getVideoTracks().forEach(track => track.stop());
            
            // Восстанавливаем локальный поток
            if (localVideoRef.current) {
                localVideoRef.current.srcObject = localStreamRef.current;
            }
            
            setIsScreenSharing(false);
        }
    };

    // WebSocket соединение
    const connectWebSocket = () => {
        console.log('Проверка параметров подключения:', {
            isJoined,
            roomId,
            hasUser: !!user,
            hasUserResult: !!user?.result,
            userRole: user?.result?.role,
            userId: user?.result?.id,
            hasLocalStream: !!localStreamRef.current
        });

        // Проверяем все необходимые условия
        if (!isJoined || !roomId || !user?.result || !localStreamRef.current) {
            console.log('Отсутствуют необходимые параметры для подключения:', {
                isJoined: isJoined ? 'да' : 'нет',
                roomId: roomId || 'не указан',
                hasUser: user ? 'да' : 'нет',
                hasUserResult: user?.result ? 'да' : 'нет',
                hasLocalStream: localStreamRef.current ? 'да' : 'нет',
                userRole: user?.result?.role,
                userId: user?.result?.id
            });
            return;
        }

        const wsUrl = `ws://localhost:9001/webrtc?room_id=${roomId}&user_id=${user.result.id}&is_tutor=${user.result.role === 0}`;
        console.log('Подключение к WebSocket:', wsUrl);

        try {
            const ws = new WebSocket(wsUrl);
            wsRef.current = ws;

            ws.onopen = () => {
                console.log('WebSocket соединение установлено');
                setError(null);
                initializePeerConnection();
            };

            ws.onmessage = async (event) => {
                try {
                    const message = JSON.parse(event.data);
                    await handleSignalingMessage(message);
                } catch (error) {
                    console.error('Ошибка при обработке сообщения:', error);
                }
            };

            ws.onerror = (error) => {
                console.error('WebSocket ошибка:', error);
                setError('Ошибка соединения с сервером');
            };

            ws.onclose = () => {
                console.log('WebSocket соединение закрыто');
                cleanup();
            };
        } catch (error) {
            console.error('Ошибка при создании WebSocket:', error);
            setError('Не удалось создать соединение');
        }
    };

    // WebRTC соединение
    const initializePeerConnection = () => {
        console.log('Инициализация PeerConnection');
        const pc = new RTCPeerConnection(rtcConfiguration);
        peerConnectionRef.current = pc;

        // Обработчики событий PeerConnection
        pc.onicecandidate = (event) => {
            if (event.candidate) {
                console.log('Новый ICE кандидат:', event.candidate);
                wsRef.current?.send(JSON.stringify({
                    type: 'candidate',
                    ice: event.candidate,
                    roomId,
                    userId: user.result.id
                }));
            }
        };

        pc.oniceconnectionstatechange = () => {
            console.log('ICE состояние:', pc.iceConnectionState);
            if (pc.iceConnectionState === 'connected') {
                setIsConnected(true);
            } else if (pc.iceConnectionState === 'disconnected' || 
                      pc.iceConnectionState === 'failed' || 
                      pc.iceConnectionState === 'closed') {
                setIsConnected(false);
            }
        };

        pc.ontrack = (event) => {
            console.log('Получен удаленный трек:', {
                kind: event.track.kind,
                enabled: event.track.enabled,
                muted: event.track.muted,
                readyState: event.track.readyState
            });

            if (event.streams && event.streams[0]) {
                if (remoteVideoRef.current) {
                    remoteVideoRef.current.srcObject = event.streams[0];
                    remoteVideoRef.current.muted = !isAudioOn;
                    
                    const audioTrack = event.streams[0].getAudioTracks()[0];
                    if (audioTrack) {
                        audioTrack.onmute = () => console.log('Удаленный аудио трек приглушен');
                        audioTrack.onunmute = () => console.log('Удаленный аудио трек включен');
                        audioTrack.onended = () => console.log('Удаленный аудио трек завершен');
                    }
                }
            }
        };

        // Добавляем локальные треки
        if (localStreamRef.current) {
            localStreamRef.current.getTracks().forEach(track => {
                console.log('Добавление локального трека:', track.kind);
                pc.addTrack(track, localStreamRef.current);
            });
        }

        // Если мы инициатор (репетитор), создаем предложение
        if (user?.result?.role === 0) {
            createOffer();
        }
    };

    const createOffer = async () => {
        try {
            const offer = await peerConnectionRef.current.createOffer({
                offerToReceiveAudio: true,
                offerToReceiveVideo: true
            });

            await peerConnectionRef.current.setLocalDescription(offer);

            wsRef.current?.send(JSON.stringify({
                type: 'offer',
                sdp: {
                    type: offer.type,
                    sdp: offer.sdp
                },
                roomId,
                userId: user.result.id
            }));
        } catch (error) {
            console.error('Ошибка при создании предложения:', error);
            message.error('Не удалось создать предложение подключения');
        }
    };

    const handleSignalingMessage = async (message) => {
        console.log('Получено сигнальное сообщение:', message.type);
        const pc = peerConnectionRef.current;

        if (!pc) {
            console.warn('PeerConnection не инициализирован');
            return;
        }

        try {
            switch (message.type) {
                case 'offer':
                    await pc.setRemoteDescription(new RTCSessionDescription(message.sdp));
                    const answer = await pc.createAnswer();
                    await pc.setLocalDescription(answer);
                    wsRef.current?.send(JSON.stringify({
                        type: 'answer',
                        sdp: {
                            type: answer.type,
                            sdp: answer.sdp
                        },
                        roomId,
                        userId: user.result.id
                    }));
                    break;

                case 'answer':
                    await pc.setRemoteDescription(new RTCSessionDescription(message.sdp));
                    break;

                case 'candidate':
                    if (message.ice) {
                        await pc.addIceCandidate(new RTCIceCandidate(message.ice));
                    }
                    break;

                default:
                    console.warn('Неизвестный тип сообщения:', message.type);
            }
        } catch (error) {
            console.error('Ошибка при обработке сигнального сообщения:', error);
        }
    };

    const cleanup = () => {
        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }

        if (peerConnectionRef.current) {
            peerConnectionRef.current.close();
            peerConnectionRef.current = null;
        }

        if (localStreamRef.current) {
            localStreamRef.current.getTracks().forEach(track => track.stop());
            localStreamRef.current = null;
        }

        if (screenStreamRef.current) {
            screenStreamRef.current.getTracks().forEach(track => track.stop());
            screenStreamRef.current = null;
        }

        if (localVideoRef.current) {
            localVideoRef.current.srcObject = null;
        }

        if (remoteVideoRef.current) {
            remoteVideoRef.current.srcObject = null;
        }

        setIsConnected(false);
        setError(null);
    };

    // Очистка при размонтировании компонента
    useEffect(() => {
        return () => {
            cleanup();
        };
    }, []);

    // Если не присоединились к комнате, показываем форму входа
    if (!isJoined) {
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
                        Видеочат
                    </Title>
                    <Space direction="vertical" style={{ width: '100%' }}>
                        <Input
                            placeholder="Введите ID комнаты"
                            value={roomId}
                            onChange={handleRoomIdChange}
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
    }

    // Основной интерфейс видеочата
    return (
        <div style={{ padding: '20px', textAlign: 'center' }}>
            <div style={{ marginBottom: '20px' }}>
                <Space direction="vertical" size="small">
                    <Space>
                        <Title level={4} style={{ margin: 0 }}>
                            Видеочат {isConnected ? '(Подключено)' : '(Ожидание подключения...)'}
                        </Title>
                        <Button onClick={handleLeaveRoom} type="primary" danger>
                            Покинуть комнату
                        </Button>
                    </Space>
                    {error && (
                        <Typography.Text type="danger">
                            {error}
                        </Typography.Text>
                    )}
                </Space>
            </div>

            <div style={{ display: 'flex', justifyContent: 'center', gap: '20px', marginBottom: '20px' }}>
                <div>
                    <h3>Ваше видео</h3>
                    <div style={{ 
                        width: 640, 
                        height: 480, 
                        backgroundColor: '#f0f2f5',
                        position: 'relative',
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                        borderRadius: '8px',
                        overflow: 'hidden'
                    }}>
                        {(isVideoOn || isScreenSharing) ? (
                            <video
                                ref={localVideoRef}
                                autoPlay
                                playsInline
                                muted
                                style={{ 
                                    width: '100%',
                                    height: '100%',
                                    objectFit: isScreenSharing ? 'contain' : 'cover'
                                }}
                            />
                        ) : (
                            <div style={{ 
                                display: 'flex', 
                                flexDirection: 'column', 
                                alignItems: 'center',
                                gap: '16px'
                            }}>
                                <Avatar 
                                    size={120} 
                                    icon={<UserOutlined />} 
                                    style={{ 
                                        backgroundColor: '#1890ff',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center'
                                    }}
                                />
                                <Typography.Text type="secondary">
                                    Видео выключено
                                </Typography.Text>
                            </div>
                        )}
                    </div>
                    {localStreamRef.current && (
                        <div style={{ marginTop: '10px' }}>
                            <Space direction="vertical" size="small">
                                <Typography.Text type="secondary">
                                    {isScreenSharing ? 'Демонстрация экрана' : `Камера: ${localStreamRef.current.getVideoTracks()[0]?.label || 'Неизвестно'}`}
                                </Typography.Text>
                                <Typography.Text type="secondary">
                                    Микрофон: {localStreamRef.current.getAudioTracks()[0]?.label || 'Неизвестно'}
                                </Typography.Text>
                            </Space>
                        </div>
                    )}
                </div>
                <div>
                    <h3>Удаленное видео</h3>
                    <div style={{ 
                        width: 640, 
                        height: 480, 
                        backgroundColor: '#f0f2f5',
                        position: 'relative',
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                        borderRadius: '8px',
                        overflow: 'hidden'
                    }}>
                        <video
                            ref={remoteVideoRef}
                            autoPlay
                            playsInline
                            muted={!isAudioOn}
                            style={{ 
                                width: '100%',
                                height: '100%',
                                objectFit: 'cover',
                                display: remoteVideoRef.current?.srcObject ? 'block' : 'none'
                            }}
                        />
                        {!remoteVideoRef.current?.srcObject && (
                            <div style={{ 
                                display: 'flex', 
                                flexDirection: 'column', 
                                alignItems: 'center',
                                gap: '16px'
                            }}>
                                <Avatar 
                                    size={120} 
                                    icon={<UserOutlined />} 
                                    style={{ 
                                        backgroundColor: '#1890ff',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center'
                                    }}
                                />
                                <Typography.Text type="secondary">
                                    Ожидание подключения
                                </Typography.Text>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            <Space>
                <Button 
                    onClick={toggleMicrophone} 
                    type={isMicOn ? "primary" : "default"}
                    danger={!isMicOn}
                    icon={<UserOutlined />}
                >
                    {isMicOn ? 'Выключить микрофон' : 'Включить микрофон'}
                </Button>
                <Button 
                    onClick={toggleVideo} 
                    type={isVideoOn ? "primary" : "default"}
                    danger={!isVideoOn}
                    icon={isVideoOn ? <VideoCameraFilled /> : <VideoCameraOutlined />}
                    // disabled={isScreenSharing}
                >
                    {isVideoOn ? 'Выключить видео' : 'Включить видео'}
                </Button>
                <Button 
                    onClick={toggleRemoteAudio} 
                    type={isAudioOn ? "primary" : "default"}
                    danger={!isAudioOn}
                >
                    {isAudioOn ? 'Выключить звук' : 'Включить звук'}
                </Button>
                <Button 
                    onClick={isScreenSharing ? stopScreenShare : startScreenShare} 
                    type={isScreenSharing ? "primary" : "default"}
                    danger={isScreenSharing}
                >
                    {isScreenSharing ? 'Остановить демонстрацию' : 'Поделиться экраном'}
                </Button>
            </Space>

            <div style={{ marginTop: '20px' }}>
                <Space direction="vertical" size="small">
                    <Typography.Text type={isMicOn ? "success" : "danger"}>
                        Микрофон: {isMicOn ? 'Включен' : 'Выключен'}
                    </Typography.Text>
                    <Typography.Text type={isAudioOn ? "success" : "danger"}>
                        Звук: {isAudioOn ? 'Включен' : 'Выключен'}
                    </Typography.Text>
                    {isConnected && (
                        <Typography.Text type="success">
                            Соединение установлено
                        </Typography.Text>
                    )}
                </Space>
            </div>
        </div>
    );
};

export default VideoChat;
