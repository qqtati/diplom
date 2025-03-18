import React, { useState, useEffect, useContext } from 'react';
import { Form, Input, Button, message, TimePicker, Card, Typography, Space, Switch, Row, Col } from 'antd';
import { updateProfile, fetchProfile, getProfile } from '../api';
import AuthContext from '../context/AuthContext';
import moment from 'moment';

const { Title } = Typography;

const Profile = () => {
    const [form] = Form.useForm();
    const { user } = useContext(AuthContext);
    const [loading, setLoading] = useState(false);
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
        loadProfile();
    }, []);

    const loadProfile = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            const { result } = await fetchProfile(accessToken);
            form.setFieldsValue(result);
            if (result.working_hours) {
                setWorkingHours(JSON.parse(result.working_hours));
            }
        } catch (error) {
            message.error('Не удалось загрузить профиль');
        }
    };

    const handleSubmit = async (values) => {
        try {
            setLoading(true);
            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }

            const updatedProfile = {
                ...values,
                working_hours: JSON.stringify(workingHours)
            };

            await updateProfile(updatedProfile, accessToken);
            message.success('Профиль обновлен');
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

    return (
        <div style={{ maxWidth: 600, margin: '0 auto', padding: 24 }}>
            <Title level={2}>Профиль</Title>
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
                    <Input />
                </Form.Item>
                <Form.Item label="Рабочее время">
                    {renderWorkingHours()}
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        Сохранить
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
};

export default Profile; 