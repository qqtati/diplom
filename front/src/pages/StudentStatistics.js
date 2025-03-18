// src/pages/StudentStatistics.js
import React, { useState } from 'react';
import { DatePicker, Table, Tag } from 'antd';
import moment from 'moment';

const { RangePicker } = DatePicker;

const StudentStatistics = () => {
    const [range, setRange] = useState([moment().startOf('month'), moment().endOf('month')]);

    const studentsData = [
        {
            key: '1',
            name: 'Иван Иванов',
            rating: 4.8,
            absences: 2,
            hours: 40,
            revenue: 15000,
        },
        {
            key: '2',
            name: 'Анна Смирнова',
            rating: 4.5,
            absences: 3,
            hours: 35,
            revenue: 13000,
        },
        {
            key: '3',
            name: 'Петр Петров',
            rating: 2.8,
            absences: 5,
            hours: 20,
            revenue: 8000,
        },
        // Добавьте больше данных по мере необходимости
    ];

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
                <Tag color={rating >= 4.5 ? 'green' : rating < 3.0 ? 'red' : 'orange'}>
                    {rating}
                </Tag>
            ),
        },
        {
            title: 'Пропуски',
            dataIndex: 'absences',
            key: 'absences',
        },
        {
            title: 'Часы',
            dataIndex: 'hours',
            key: 'hours',
        },
        {
            title: 'Доход (₽)',
            dataIndex: 'revenue',
            key: 'revenue',
        },
    ];

    const onDateChange = (dates, dateStrings) => {
        setRange(dates);
        // Добавьте здесь логику для фильтрации данных на основе выбранного временного интервала
    };

    return (
        <div style={{ padding: '20px' }}>
            <h2>Статистика учеников</h2>
            <RangePicker
                value={range}
                onChange={onDateChange}
                style={{ marginBottom: '20px' }}
            />
            <Table
                columns={columns}
                dataSource={studentsData}
                bordered
                pagination={{ pageSize: 5 }}
            />
        </div>
    );
};

export default StudentStatistics;