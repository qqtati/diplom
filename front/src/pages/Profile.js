// src/pages/Profile.js
import React, { useEffect, useState } from 'react';
import { Card, List, Avatar, Button, message } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { fetchProfile } from '../api';

const Profile = () => {
    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const accessToken = localStorage.getItem('accessToken'); // Получаем токен из локального хранилища
        if (!accessToken) {
            message.error('Токен авторизации не найден');
            return;
        }

        const loadProfile = async () => {
            try {
                const {result} = await fetchProfile(accessToken);
                setProfile(result);
            } catch (error) {
                message.error('Не удалось загрузить профиль');
            } finally {
                setLoading(false);
            }
        };

        loadProfile();
    }, []);

    const handleLogout = () => {
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        message.success('Вы успешно вышли');
    };

    if (loading) {
        return <div>Загрузка...</div>;
    }

    if (!profile) {
        return <div>Данные профиля не найдены</div>;
    }

    return (
        <div style={{ padding: '20px' }}>
            <Card
                style={{ maxWidth: 400, margin: '0 auto', marginBottom: '20px' }}
                actions={[
                    <Button type="primary" danger onClick={handleLogout} block>
                        Выйти
                    </Button>,
                ]}
            >
                <Card.Meta
                    avatar={<Avatar icon={<UserOutlined />} />}
                    title={profile.name}
                    description={
                        <>
                            <p><strong>ID:</strong> {profile.id}</p>
                            <p><strong>Имя пользователя:</strong> {profile.username}</p>
                            <p><strong>Роль:</strong> {profile.role === 0 ? 'Репетитор' : 'Ученик'}</p>
                            <p><strong>Код приглашения:</strong> {profile.invite_code}</p>
                        </>
                    }
                />

            </Card>
            {/* Здесь можно добавить другие элементы интерфейса, такие как список уроков */}
        </div>
    );
};

export default Profile;