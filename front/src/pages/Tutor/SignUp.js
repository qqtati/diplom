// src/pages/Tutor/SignUp.js
import React, { useState } from 'react';
import { signUp } from '../../api';
import { Form, Input, Button, message, Checkbox } from 'antd';

const TutorSignUp = () => {
    const [loading, setLoading] = useState(false);

    const onFinish = async (values) => {
        setLoading(true);
        try {
            const { is_teacher, invite_code, ...rest } = values;
            const data = {
                ...rest,
                is_teacher: is_teacher ? 1 : 0,
                invite_code: invite_code || undefined, // Убираем invite_code, если оно пустое
            };
            const response = await signUp(data);
            if (response.success) {
                message.success('Регистрация успешна');
                // Вы можете перенаправить пользователя на другую страницу
                // history.push('/login'); // Используйте нужную вам страницу
            }
        } catch (error) {
            message.error('Не удалось зарегистрироваться');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Form onFinish={onFinish} style={{ maxWidth: 400, margin: '15px auto' }}>
            <Form.Item name="username" rules={[{ required: true, message: 'Введите имя пользователя' }]}>
                <Input placeholder="Имя пользователя" />
            </Form.Item>
            <Form.Item name="password" rules={[{ required: true, message: 'Введите пароль' }]}>
                <Input.Password placeholder="Пароль" />
            </Form.Item>
            <Form.Item name="name" rules={[{ required: true, message: 'Введите имя' }]}>
                <Input placeholder="Имя" />
            </Form.Item>
            <Form.Item name="is_teacher" valuePropName="checked">
                <Checkbox>Я репетитор</Checkbox>
            </Form.Item>
            <Form.Item name="invite_code">
                <Input placeholder="Код приглашения (если есть)" />
            </Form.Item>
            <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                    Зарегистрироваться
                </Button>
            </Form.Item>
        </Form>
    );
};

export default TutorSignUp;