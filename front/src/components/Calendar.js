// src/components/MyCalendar.js
import React, { useEffect, useState } from 'react';
import {Calendar, Calendar as BigCalendar, momentLocalizer} from 'react-big-calendar';
import moment from 'moment';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import {fetchEvents, createEvent, updateEvent, deleteEvent, fetchProfile, fetchStudents} from '../api';
import {Drawer, Button, Form, Input, DatePicker, message, Select, InputNumber} from 'antd';
import withDragAndDrop from "react-big-calendar/lib/addons/dragAndDrop";

const localizer = momentLocalizer(moment);
const DnDCalendar = withDragAndDrop(Calendar);

const MyCalendar = () => {
    const [events, setEvents] = useState([]);
    const [drawerVisible, setDrawerVisible] = useState(false);
    const [currentEvent, setCurrentEvent] = useState(null);
    const [profile, setProfile] = useState(null);
    const [students, setStudents] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadEvents()
        loadProfile()
        loadStudents()
        // setInterval(loadEvents, 1000)
    }, []);

    const loadProfile = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken'); // Получаем токен из локального хранилища
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            const {result} = await fetchProfile(accessToken);
            setProfile(result);
        } catch (error) {
            message.error('Не удалось загрузить профиль');
        } finally {
            setLoading(false);
        }
    };

    const loadStudents = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken'); // Получаем токен из локального хранилища
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            const result = await fetchStudents(accessToken);
            setStudents(result);
            console.log(result);
        } catch (error) {
            message.error('Не удалось загрузить профиль');
        } finally {
            setLoading(false);
        }
    };

    const loadEvents = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken'); // Получаем токен из локального хранилища
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            const eventsData = await fetchEvents(accessToken);

            setEvents(eventsData.map(event => ({
                ...event,
                start: new Date(event.start_time),
                end: new Date(event.end_time),
                start_time: moment(event.start_time, "YYYY-MM-DDTHH:mm:ssZ"),
                end_time: moment(event.end_time, "YYYY-MM-DDTHH:mm:ssZ"),
                title: event.description,
            })));
        } catch (error) {
            message.error(error.message);
        }
    };

    const handleSelectEvent = (event) => {
        setCurrentEvent(event);
        setDrawerVisible(true);
    };

    const handleSelectSlot = ({ start, end }) => {
        setCurrentEvent({ start, end, start_time: moment(start), end_time: moment(end)});
        setDrawerVisible(true);
    };

    const handleEventDrop = async ({ event, start, end }) => {
        const updatedEvent = { ...event, start, end, start_time: start, end_time: end };
        try {
            const accessToken = localStorage.getItem('accessToken'); // Получаем токен из локального хранилища
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            await updateEvent(updatedEvent, accessToken);
            loadEvents();
            message.success('Событие обновлено');
        } catch (error) {
            message.error('Ошибка при обновлении события');
        }
    };

    const handleCreateUpdateEvent = async (values) => {
        const eventData = {
            ...currentEvent,
            start_time: currentEvent.start.toISOString(),
            end_time: currentEvent.end.toISOString(),
            description: values.description,
            teacher_id: (profile.role === 0 ? profile.id : null),
            student_id: values.student_id,
            skipped: values.skipped ? values.skipped : false,
            price: values.price,
        };
        try {

            const accessToken = localStorage.getItem('accessToken'); // Получаем токен из локального хранилища
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            if (eventData.id) {
                await updateEvent(eventData, accessToken);
                message.success('Событие обновлено');
            } else {
                await createEvent(eventData, accessToken);
                message.success('Событие создано');
            }
            setDrawerVisible(false);
            loadEvents();
        } catch (error) {
            message.error(error.message);
        }
    };

    const handleDeleteEvent = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken'); // Получаем токен из локального хранилища
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            await deleteEvent(currentEvent.id, accessToken);
            message.success('Событие удалено');
            setDrawerVisible(false);
            loadEvents();
        } catch (error) {
            message.error(error.message);
        }
    };

    const eventRender = (date) => {
        return (
            <div>{events.filter((event) => moment(event.start_time).isSame(date, 'day')).map((event) => (
                <div key={event.id}>
                    <Button type="link" onClick={() => {
                        setCurrentEvent(event);
                        setDrawerVisible(true);
                    }}>
                        {event.description}/{event.price}₽
                    </Button>
                </div>
            ))}</div>
        );
    };

    return (
        <div>
            <DnDCalendar
                localizer={localizer}
                events={events}
                startAccessor="start"
                endAccessor="end"
                defaultView="week"
                views={['day', 'week', 'month']}
                style={{ height: 500 }}
                selectable
                onSelectEvent={handleSelectEvent}
                onSelectSlot={handleSelectSlot}
                onEventDrop={handleEventDrop}
                dateCellRender={eventRender}
                resizable
                draggableAccessor={() => true}
            />
            <Drawer
                title={currentEvent?.id ? 'Редактировать событие' : 'Новое событие'}
                width={320}
                onClose={() => setDrawerVisible(false)}
                visible={drawerVisible}
            >
                <Form
                    initialValues={currentEvent || {}}
                    onFinish={handleCreateUpdateEvent}
                >
                    <Form.Item
                        name="description"
                        label="Описание"
                    >
                        <Input.TextArea />
                    </Form.Item>
                    <Form.Item
                        name="student_id"
                        label="Выберите ученика"
                    >
                        <Select options={students && students.map(e => ({ value: e.id, label: e.name }))}/>
                    </Form.Item>
                    <Form.Item
                        name="price"
                        label="Цена"
                    >
                        <InputNumber />
                    </Form.Item>
                    <Button type="primary" htmlType="submit">
                        {currentEvent?.id ? 'Обновить' : 'Создать'}
                    </Button>
                    {currentEvent?.id && (
                        <Button type="danger" onClick={handleDeleteEvent}>
                            Удалить
                        </Button>
                    )}
                </Form>
            </Drawer>
        </div>
    );
};

export default MyCalendar;