// src/pages/VideoChat.js
import React, { useState, useRef } from 'react';
import Webcam from 'react-webcam';
import { Button, Space } from 'antd';

const VideoChat = () => {
    const [isMicOn, setIsMicOn] = useState(true);
    const [isScreenSharing, setIsScreenSharing] = useState(false);
    const webcamRef = useRef(null);
    const userVideoRef = useRef(null);
    const screenStreamRef = useRef(null);

    const toggleMicrophone = () => {
        setIsMicOn(!isMicOn);
        const stream = webcamRef.current.stream;
        const audioTrack = stream.getAudioTracks()[0];
        audioTrack.enabled = isMicOn;
    };

    const startScreenShare = async () => {
        try {
            screenStreamRef.current = await navigator.mediaDevices.getDisplayMedia({ video: true });
            userVideoRef.current.srcObject = screenStreamRef.current;
            setIsScreenSharing(true);

            screenStreamRef.current.getVideoTracks()[0].onended = () => {
                stopScreenShare();
            };
        } catch (error) {
            console.error('Error sharing screen:', error);
        }
    };

    const stopScreenShare = () => {
        if (screenStreamRef.current) {
            screenStreamRef.current.getTracks().forEach(track => track.stop());
            webcamRef.current.video.srcObject = webcamRef.current.stream;
        }
        setIsScreenSharing(false);
    };

    return (
        <div style={{ padding: '20px', textAlign: 'center' }}>
            <h2>Видеочат</h2>
            <Webcam
                audio
                ref={webcamRef}
                screenshotFormat="image/jpeg"
                style={{ width: 640, height: 480, marginBottom: '20px' }}
            />
            <video
                ref={userVideoRef}
                autoPlay
                playsInline
                style={{ width: 640, height: 480, display: isScreenSharing ? 'block' : 'none', marginBottom: '20px' }}
            />
            <Space>
                <Button onClick={toggleMicrophone} type="primary">
                    {isMicOn ? 'Выключить микрофон' : 'Включить микрофон'}
                </Button>
                <Button onClick={isScreenSharing ? stopScreenShare : startScreenShare} type="primary">
                    {isScreenSharing ? 'Остановить демонстрацию' : 'Поделиться экраном'}
                </Button>
            </Space>
        </div>
    );
};

export default VideoChat;