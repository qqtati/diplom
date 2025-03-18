import React from 'react';
import { Form, Input, Button } from 'antd';

const StudentSignUp = () => {
    return (
        <Form layout="vertical" style={{ maxWidth: '600px', margin: 'auto', padding: '2rem' }}>
            <Form.Item label="Имя">
                <Input placeholder="Введите ваше имя" />
            </Form.Item>
            <Form.Item label="Email">
                <Input type="email" placeholder="Введите ваш Email" />
            </Form.Item>
            <Form.Item label="Телефон">
                <Input placeholder="Введите ваш телефон" />
            </Form.Item>
            <Form.Item label="Предпочитаемые предметы">
                <Input placeholder="Выберите предметы" />
            </Form.Item>
            <Button type="primary">Зарегистрироваться</Button>
        </Form>
    );
};

export default StudentSignUp;