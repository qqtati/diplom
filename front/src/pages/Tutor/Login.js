// src/pages/Tutor/Login.js
import React, {useContext, useState} from 'react';
import { signIn } from '../../api';
import { Form, Input, Button, message } from 'antd';
import authContext from "../../context/AuthContext";
import AuthContext from "../../context/AuthContext";

const TutorLogin = () => {
    const [loading, setLoading] = useState(false);
    const authContext = useContext(AuthContext);

    const onFinish = async (values) => {
        setLoading(true);
        try {
            const { result } = await signIn(values.username, values.password);
            const { access_token, refresh_token } = result;
            // Сохраните токены в локальное хранилище или глобальное состояние
            authContext.login(null, result);
            message.success('Успешный вход в систему')
        } catch (error) {
            message.error('Не удалось выполнить вход');
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
            <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                    Войти
                </Button>
            </Form.Item>
        </Form>
    );
};

export default TutorLogin;