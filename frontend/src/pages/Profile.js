// src/pages/Profile.js
import React, { useEffect, useState } from 'react';
import { Card, List, Avatar, Button, message, Form, Input, Switch, Row, Col, TimePicker } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { fetchProfile, getProfile, updateProfile } from '../api';
import { useHistory } from 'react-router-dom';
import moment from 'moment';

const Profile = () => {
    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [form] = Form.useForm();
    const history = useHistory();
    const [workingHours, setWorkingHours] = useState({
        monday: { enabled: false, start: '09:00', end: '18:00' },
        tuesday: { enabled: false, start: '09:00', end: '18:00' },
        wednesday: { enabled: false, start: '09:00', end: '18:00' },
        thursday: { enabled: false, start: '09:00', end: '18:00' },
        friday: { enabled: false, start: '09:00', end: '18:00' },
        saturday: { enabled: false, start: '09:00', end: '18:00' },
        sunday: { enabled: false, start: '09:00', end: '18:00' }
    });

    useEffect(() => {
        const accessToken = localStorage.getItem('accessToken');
        if (!accessToken) {
            message.error('Токен авторизации не найден');
            return;
        }

        const loadProfile = async () => {
            try {
                const {result} = await fetchProfile(accessToken);
                setProfile(result);
                form.setFieldsValue({
                    name: result.name,
                    email: result.username
                });

                if (result.working_hours) {
                    try {
                        const parsedHours = JSON.parse(result.working_hours);
                        // Проверяем, что все дни недели присутствуют в данных
                        const defaultHours = {
                            monday: { enabled: false, start: '09:00', end: '18:00' },
                            tuesday: { enabled: false, start: '09:00', end: '18:00' },
                            wednesday: { enabled: false, start: '09:00', end: '18:00' },
                            thursday: { enabled: false, start: '09:00', end: '18:00' },
                            friday: { enabled: false, start: '09:00', end: '18:00' },
                            saturday: { enabled: false, start: '09:00', end: '18:00' },
                            sunday: { enabled: false, start: '09:00', end: '18:00' }
                        };

                        // Объединяем сохраненные часы с дефолтными значениями
                        const mergedHours = Object.keys(defaultHours).reduce((acc, day) => {
                            acc[day] = {
                                ...defaultHours[day],
                                ...(parsedHours[day] || {})
                            };
                            return acc;
                        }, {});

                        setWorkingHours(mergedHours);
                    } catch (e) {
                        console.error('Ошибка при парсинге рабочего времени:', e);
                        message.warning('Не удалось загрузить настройки рабочего времени');
                    }
                }
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

    const handleSubmit = async (values) => {
        try {
            setLoading(true);

            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            const profileData = {
                name: values.name,
                email: values.email,
                working_hours: JSON.stringify(workingHours)
            };

            await updateProfile(profileData, accessToken);
            message.success('Профиль успешно обновлен');
        } catch (error) {
            message.error('Ошибка при обновлении профиля');
        } finally {
            setLoading(false);
        }
    };

    const handleWorkingHoursChange = (day, field, value) => {
        setWorkingHours(prev => ({
            ...prev,
            [day]: {
                ...prev[day],
                [field]: value
            }
        }));
    };

    const renderWorkingHours = () => {
        const days = {
            monday: 'Понедельник',
            tuesday: 'Вторник',
            wednesday: 'Среда',
            thursday: 'Четверг',
            friday: 'Пятница',
            saturday: 'Суббота',
            sunday: 'Воскресенье'
        };

        return Object.entries(days).map(([day, label]) => (
            <Row key={day} gutter={16} style={{ marginBottom: 16 }}>
                <Col span={6}>
                    <Switch
                        checked={workingHours[day].enabled}
                        onChange={(checked) => handleWorkingHoursChange(day, 'enabled', checked)}
                    />
                    <span style={{ marginLeft: 8 }}>{label}</span>
                </Col>
                <Col span={9}>
                    <TimePicker
                        format="HH:mm"
                        value={moment(workingHours[day].start, 'HH:mm')}
                        onChange={(time) => handleWorkingHoursChange(day, 'start', time.format('HH:mm'))}
                        disabled={!workingHours[day].enabled}
                    />
                </Col>
                <Col span={9}>
                    <TimePicker
                        format="HH:mm"
                        value={moment(workingHours[day].end, 'HH:mm')}
                        onChange={(time) => handleWorkingHoursChange(day, 'end', time.format('HH:mm'))}
                        disabled={!workingHours[day].enabled}
                    />
                </Col>
            </Row>
        ));
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
            <div style={{ maxWidth: 600, margin: '0 auto', padding: '24px' }}>
                <h1>Профиль</h1>
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                >
                    <Form.Item
                        name="name"
                        label="Имя"
                        rules={[{ required: true, message: 'Пожалуйста, введите имя' }]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        name="email"
                        label="Email"
                        rules={[
                            { required: true, message: 'Пожалуйста, введите email' },
                            { type: 'email', message: 'Введите корректный email' }
                        ]}
                    >
                        <Input disabled />
                    </Form.Item>

                    <h2>Рабочее время</h2>
                    {renderWorkingHours()}

                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={loading}>
                            Сохранить
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        </div>
    );
};

export default Profile;