import React from 'react';
import { Form, Input, Button } from 'antd';

const StudentLogin = () => {
    return (
        <Form layout="vertical" style={{ maxWidth: '600px', margin: 'auto', padding: '2rem' }}>
            <Form.Item label="Email / Телефон">
                <Input placeholder="Введите ваш Email или телефон" />
            </Form.Item>
            <Form.Item label="Пароль">
                <Input.Password placeholder="Введите ваш пароль" />
            </Form.Item>
            <Button type="primary">Войти</Button>
        </Form>
    );
};

export default StudentLogin;