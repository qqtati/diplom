// src/pages/StudentStatistics.js
import React, { useState, useEffect } from 'react';
import { DatePicker, Table, Tag, Card, Row, Col, Statistic, Typography } from 'antd';
import moment from 'moment';
import { fetchStudentStats } from '../api';

const StudentStatistics = () => {
    const [selectedDate, setSelectedDate] = useState(moment().startOf('month'));
    const [students, setStudents] = useState(null);
    const [error, setError] = useState(null);
    const [days, setDays] = useState(30);

    useEffect(() => {
        getStudents(days);
    }, [days]);

    const getStudents = async (cnt) => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            const result = await fetchStudentStats(cnt, accessToken);
            setStudents(result);
        } catch (e) {
            setError(e);
        }
    };

    const columns = [
        {
            title: 'Имя',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: 'Рейтинг',
            dataIndex: 'rating',
            key: 'rating',
            render: (rating) => (
                <Tag color={rating >= 4.0 ? 'green' : rating < 3.0 ? 'red' : 'orange'}>
                    {Math.round(rating)}
                </Tag>
            ),
        },
        {
            title: 'Пропуски',
            dataIndex: 'skip_count',
            key: 'skip_count',
        },
        {
            title: 'Занятия',
            dataIndex: 'event_count',
            key: 'event_count',
        },
        {
            title: 'Среднее время (мин)',
            dataIndex: 'avg_duration',
            key: 'avg_duration',
            render: (duration) => Math.floor(duration),
        },
        {
            title: 'Средняя цена (₽)',
            key: 'avg_price',
            render: (_, record) => record.event_count === 0 ? '-' : Math.floor(record.total_income / record.event_count),
        },
        {
            title: 'Доход (₽)',
            dataIndex: 'total_income',
            key: 'total_income',
        },
    ];

    const onDateChange = (date) => {
        setSelectedDate(date);
        if (date) {
            const daysDiff = moment().diff(date, 'days');
            setDays(daysDiff);
            getStudents(daysDiff);
        }
    };

    const totalLessons = students ? students.reduce((val, el) => val + el.event_count, 0) : 0;
    const totalIncome = students ? students.reduce((val, el) => val + el.total_income, 0) : 0;
    const totalStudents = students ? students.length : 0;

    return (
        <div style={{ padding: '20px' }}>
            <h2>Статистика учеников</h2>
            <Row style={{ marginBottom: '20px' }}>
                <Col>
                    <Typography.Text>Выберите начальную дату: </Typography.Text>
                    <DatePicker
                        value={selectedDate}
                        onChange={onDateChange}
                        label="Выберите начальную дату"
                    />
                </Col>
            </Row>
            <Row gutter={16} style={{ marginBottom: '20px' }}>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="Всего проведено занятий"
                            value={totalLessons}
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="Заработок"
                            value={totalIncome}
                            suffix="₽"
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="Всего учеников"
                            value={totalStudents}
                        />
                    </Card>
                </Col>
            </Row>
            <Table
                columns={columns}
                dataSource={students}
                bordered
                pagination={{ pageSize: 10 }}
                rowKey="id"
            />
        </div>
    );
};

export default StudentStatistics;