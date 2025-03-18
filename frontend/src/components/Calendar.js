// src/components/MyCalendar.js
import React, { useEffect, useState, useContext } from 'react';
import {Calendar, Calendar as BigCalendar, momentLocalizer} from 'react-big-calendar';
import moment, { duration } from 'moment';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import {fetchEvents, createEvent, updateEvent, deleteEvent, fetchProfile, fetchStudents, fetchTeacherEvents} from '../api';
import {Drawer, Button, Form, Input, DatePicker, message, Select, InputNumber, Switch, Flex, Typography, Space, Modal, TimePicker} from 'antd';
import AuthContext from '../context/AuthContext';

import withDragAndDrop from "react-big-calendar/lib/addons/dragAndDrop";
import './CalendarStyles.css'
const { Text, Link } = Typography;

const localizer = momentLocalizer(moment);
const DnDCalendar = withDragAndDrop(Calendar);

const MyCalendar = () => {
    const [events, setEvents] = useState([]);
    const [teacherEvents, setTeacherEvents] = useState([]);
    const [drawerVisible, setDrawerVisible] = useState(false);
    const [currentEvent, setCurrentEvent] = useState(null);
    const [profile, setProfile] = useState(null);
    const [students, setStudents] = useState(null);
    const [loading, setLoading] = useState(true);
    const [showTeacherEvents, setShowTeacherEvents] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [selectedDate, setSelectedDate] = useState(null);
    const [form] = Form.useForm();
    const { user } = useContext(AuthContext);
    const isTeacher = user?.result?.role === 0;
    const [isRecurringModalVisible, setIsRecurringModalVisible] = useState(false);
    const [isRecurringDeleteModalVisible, setIsRecurringDeleteModalVisible] = useState(false);
    const [editingEvent, setEditingEvent] = useState(null);
    const [recurringOption, setRecurringOption] = useState('single');
    const [workingHours, setWorkingHours] = useState(null);
    console.log(user, isTeacher)
    useEffect(() => {
        loadEvents()
        loadProfile()
        loadStudents()
        if (profile?.role === 1) {
            loadTeacherEvents()
        }
        if (user?.result?.role === 0) {
            loadWorkingHours();
        }
    }, []);

    const loadProfile = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
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
            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            const result = await fetchStudents(accessToken);
            setStudents(result);
        } catch (error) {
            message.error('Не удалось загрузить профиль');
        } finally {
            setLoading(false);
        }
    };

    const loadEvents = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
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

    const loadTeacherEvents = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }
            const eventsData = await fetchTeacherEvents(accessToken);
            setTeacherEvents(eventsData.map(event => ({
                ...event,
                start: new Date(event.start_time),
                end: new Date(event.end_time),
                title: event.description,
                backgroundColor: '#ffd591',
                borderColor: '#ffd591',
            })));
        } catch (error) {
            message.error('Ошибка при загрузке занятости учителя');
        }
    };

    const loadWorkingHours = async () => {
        try {
            const { result } = await fetchProfile(localStorage.getItem('accessToken'));
            if (result && result.working_hours) {
                setWorkingHours(JSON.parse(result.working_hours));
            }
        } catch (error) {
            console.error('Error loading working hours:', error);
        }
    };

    const getWorkingHoursForDate = (date) => {
        if (!workingHours) return null;
        const dayOfWeek = moment(date).format('dddd').toLowerCase();
        return workingHours[dayOfWeek];
    };

    useEffect(() => {
        if (!currentEvent) return;
        form.setFieldsValue({...currentEvent,
            date: currentEvent.start ? moment(currentEvent.start) : moment(), 
            time: currentEvent.start ? moment(currentEvent.start) : moment(),
            duration: currentEvent?.start && currentEvent?.end ? moment(currentEvent.end).diff(moment(currentEvent.start), 'minutes') : 30})
        console.log("updating curr event", currentEvent)
    }, [form, currentEvent])

    const handleSelectEvent = (event) => {
        setCurrentEvent({
            ...event,
            recurrence_end_date: event.recurrence_end_date ? moment(event.recurrence_end_date) : null
        });
        setDrawerVisible(true);
    };

    const handleSelectSlot = ({ start, end }) => {
        setCurrentEvent({ start, end: moment(end), start_time: moment(start), end_time: moment(end) });
        setDrawerVisible(true);
    };

    const isWithinWorkingHours = (startTime, endTime) => {
        if (!profile?.working_hours) return true;
        
        const workingHours = JSON.parse(profile.working_hours);
        const dayOfWeek = moment(startTime).format('dddd').toLowerCase();
        const dayHours = workingHours[dayOfWeek];
        
        if (!dayHours?.start || !dayHours?.end) return false;
        
        const start = moment(startTime);
        const end = moment(endTime);
        const workStart = moment(dayHours.start, 'HH:mm');
        const workEnd = moment(dayHours.end, 'HH:mm');
        
        return start.isSameOrAfter(workStart) && end.isSameOrBefore(workEnd);
    };

    const handleCreateUpdateEvent = async (values) => {
        try {
            console.log(values)
            const eventData = {
                id: currentEvent?.id || null,
                teacher_id: currentEvent?.teacher_id || user?.result?.id,
                student_id: values.student_id,
                start_time: currentEvent.start_time.toISOString(),
                end_time: currentEvent.end_time.toISOString(),
                description: values.description,
                price: values.price,
                is_recurring: values.is_recurring,
                recurrence_pattern: values.recurrence_pattern,
                recurrence_end_date: values.recurrence_end_date ? moment(values.recurrence_end_date).endOf('day').toISOString() : null,
                parent_event_id: currentEvent?.parent_event_id || null,
                approved_by_teacher: isTeacher,
            };

            // Проверяем рабочее время
            if (isTeacher === 0) { // Если текущий пользователь - учитель
                const { result: teacherProfile } = await fetchProfile(localStorage.getItem('accessToken'));
                if (teacherProfile && teacherProfile.working_hours) {
                    const workingHours = JSON.parse(teacherProfile.working_hours);
                    const dayOfWeek = values.start_time.format('dddd').toLowerCase();
                    const eventStartTime = values.start_time.format('HH:mm');
                    const eventEndTime = moment(values.end_time).format('HH:mm');

                    if (!workingHours[dayOfWeek] || 
                        !workingHours[dayOfWeek].enabled || 
                        eventStartTime < workingHours[dayOfWeek].start || 
                        eventEndTime > workingHours[dayOfWeek].end) {
                        message.error('Время урока не соответствует рабочему времени');
                        return;
                    }
                }
            }

            if (eventData.id) {
                await updateEvent(eventData, localStorage.getItem('accessToken'));
            } else {
                const createdEvent = await createEvent(eventData, localStorage.getItem('accessToken'));
                eventData.id = createdEvent.id;

                // Создаем повторяющиеся события
                if (eventData.is_recurring && eventData.recurrence_end_date) {
                    const events = [];
                    const startDate = moment(eventData.start_time);
                    const endDate = moment(eventData.recurrence_end_date).endOf('day');
                    const duration = moment(eventData.end_time).diff(moment(eventData.start_time));

                    let currentDate = startDate.clone().add(1, 'day');
                    while (currentDate.isBefore(endDate)) {
                        const dayOfWeek = currentDate.format('dddd').toLowerCase();
                        const { result: teacherProfile } = await fetchProfile(localStorage.getItem('accessToken'));
                        const workingHours = JSON.parse(teacherProfile.working_hours);
                        
                        if (workingHours[dayOfWeek] && workingHours[dayOfWeek].enabled) {
                            const eventStartTime = moment(currentDate).format('HH:mm');
                            const eventEndTime = moment(currentDate).add(duration, 'ms').format('HH:mm');
                            
                            if (eventStartTime >= workingHours[dayOfWeek].start && 
                                eventEndTime <= workingHours[dayOfWeek].end) {
                                events.push({
                                    ...eventData,
                                    id: null,
                                    start_time: currentDate.toISOString(),
                                    end_time: moment(currentDate).add(duration, 'ms').toISOString(),
                                    parent_event_id: eventData.id
                                });
                            }
                        }

                        switch (eventData.recurrence_pattern) {
                            case 'weekly':
                                currentDate.add(1, 'week');
                                break;
                            case 'biweekly':
                                currentDate.add(2, 'weeks');
                                break;
                            case 'monthly':
                                currentDate.add(1, 'month');
                                break;
                        }
                    }

                    if (events.length > 0) {
                        await Promise.all(events.map(event => createEvent(event, localStorage.getItem('accessToken'))));
                        message.success('Создано повторяющихся событий: ' + events.length);
                    }
                }
            }

            await loadEvents();
            setDrawerVisible(false);
            setCurrentEvent(null);
        } catch (error) {
            if ((error + '').includes(`Unexpected token 'e', "error chec"... is not valid JSON`)) {
                message.error('Это нерабочее время')
            } else {
                message.error('Ошибка при создании/обновлении события');
            }
        }
    };

    const handleRecurringUpdate = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }

            if (recurringOption === 'single') {
                const updatedEvent = {
                    ...editingEvent,
                    start_time: moment(editingEvent.start).toISOString(),
                    end_time: moment(editingEvent.end).toISOString(),
                    is_recurring: false,
                    recurrence_pattern: null,
                    recurrence_end_date: null
                };
                console.log('Updating single event:', updatedEvent);
                await updateEvent(updatedEvent, accessToken);
                message.success('Событие обновлено');
            } else {
                // Обновляем все  события в серии
                const events = await fetchEvents(accessToken);
                console.log('All events:', events);
                console.log('Editing event:', editingEvent);
                
                const futureEvents = events.filter(event => 
                    (event.parent_event_id === editingEvent.parent_event_id || event.id === editingEvent.parent_event_id || event.parent_event_id === editingEvent.id) && event.is_recurring
                );
                
                console.log('Future events to update:', futureEvents);

                // Вычисляем смещение времени для каждого события
                const originalStart = moment(editingEvent.start_time);
                const newStart = moment(editingEvent.start);
                const timeOffset = newStart.diff(originalStart);
                console.log('Original start:', originalStart.format());
                console.log('New start:', newStart.format());
                console.log('Time offset:', timeOffset);
                
                // Сначала обновляем текущее событие
                const currentEventUpdate = {
                    ...editingEvent,
                    start_time: moment(editingEvent.start).toISOString(),
                    end_time: moment(editingEvent.end).toISOString()
                };
                await updateEvent(currentEventUpdate, accessToken);
                
                // Затем обновляем все  события
                await Promise.all(futureEvents.map(event => {
                    const newStartTime = moment(event.start_time).add(timeOffset, 'milliseconds');
                    const newEndTime = moment(event.end_time).add(timeOffset, 'milliseconds');
                    
                    const updatedEvent = {
                        ...event,
                        start_time: newStartTime.toISOString(),
                        end_time: newEndTime.toISOString(),
                        description: editingEvent.description,
                        price: editingEvent.price,
                        student_id: editingEvent.student_id,
                        teacher_id: editingEvent.teacher_id,
                        skipped: editingEvent.skipped,
                        is_recurring: event.is_recurring,
                        recurrence_pattern: event.recurrence_pattern,
                        recurrence_end_date: event.recurrence_end_date,
                        parent_event_id: event.parent_event_id
                    };
                    
                    console.log('Updating event:', updatedEvent);
                    return updateEvent(updatedEvent, accessToken);
                }));
                message.success('Все  события обновлены');
            }

            setIsRecurringModalVisible(false);
            setEditingEvent(null);
            setDrawerVisible(false);
            loadEvents();
        } catch (error) {
            console.error('Error updating events:', error);
            console.log(error)
            if ((error + '').includes(`Unexpected token 'e', "error chec"... is not valid JSON`)) {
                message.error('Это нерабочее время');
            } else {
                message.error('Ошибка при обновлении событий');
            }
        }
    };

    const handleDeleteEvent = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }

            if (currentEvent.is_recurring) {
                setEditingEvent(currentEvent);
                setIsRecurringDeleteModalVisible(true);
                return;
            }

            await deleteEvent(currentEvent.id, accessToken);
            message.success('Событие удалено');
            setDrawerVisible(false);
            loadEvents();
        } catch (error) {
            message.error('Ошибка при удалении события');
        }
    };

    const handleRecurringDelete = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            if (!accessToken) {
                message.error('Токен авторизации не найден');
                return;
            }

            if (recurringOption === 'single') {
                await deleteEvent(editingEvent.id, accessToken);
                message.success('Событие удалено');
            } else {
                // Удаляем все  события в серии, включая текущее
                const events = await fetchEvents(accessToken);
                const futureEvents = events.filter(event => 
                    (event.parent_event_id === editingEvent.id || event.id === editingEvent.id || event.parent_event_id === editingEvent.parent_event_id) && event.is_recurring
                );

                await Promise.all(futureEvents.map(event => 
                    deleteEvent(event.id, accessToken)
                ));
                message.success('Все  события удалены');
            }

            setIsRecurringModalVisible(false);
            setEditingEvent(null);
            setDrawerVisible(false);
            loadEvents();
        } catch (error) {
            message.error('Ошибка при удалении событий');
        }
    };

    const handleShowTeacherEvents = async () => {
        setShowTeacherEvents(true);
        await loadTeacherEvents();
    };

    const handleHideTeacherEvents = () => {
        setShowTeacherEvents(false);
        setTeacherEvents([]);
    };

    const eventRender = (props) => {
        const event = props.event;
        const teacher = profile && profile.role === 1 ? profile : null;
        const student = students?.find(s => s.id === event.student_id);
        const isTeacher = user.role === 1;
        const isStudent = user.role === 1;

        return (
            <div className="event-content" style={{backgroundColor: (event.approved_by_teacher ? (event.skipped ? 'red' : '#1890ff') : '#9999aa')}}>
                <div className="event-title">{event.is_recurring && <span> </span>}{event.title}</div>
                <div className="event-price">{event.price}₽</div>
            </div>
        );
    };

    const handleEventDrop = async ({ event, start, end }) => {
        if (user.role === 1 && !isWithinWorkingHours(start, end)) {
            message.error('Время занятия должно быть в пределах рабочих часов');
            return;
        }

        const dayHours = getWorkingHoursForDate(start);

        if (user.role === 1 && dayHours && !dayHours.enabled) {
            message.warning('В этот день нет рабочих часов');
            return;
        }

        if (user.role === 1 && dayHours) {
            const eventStartTime = moment(start).format('HH:mm');
            const eventEndTime = moment(end).format('HH:mm');
            
            if (eventStartTime < dayHours.start || eventEndTime > dayHours.end) {
                message.warning('Время события не соответствует рабочему времени');
                return;
            }
        }

        const updatedEvent = { 
            ...event, 
            start: moment(start), 
            end: moment(end), 
        };

        if (updatedEvent.is_recurring) {
            setEditingEvent(updatedEvent);
            setIsRecurringModalVisible(true);
        } else {
            try {
                updatedEvent.start_time = moment(start).toISOString();
                updatedEvent.end_time = moment(end).toISOString();
                const accessToken = localStorage.getItem('accessToken');
                if (!accessToken) {
                    message.error('Токен авторизации не найден');
                    return;
                }
                await updateEvent(updatedEvent, accessToken);
                loadEvents();
                message.success('Событие обновлено');
            } catch (error) {
                console.log(error)
                if ((error + '').includes(`Unexpected token 'e', "error chec"... is not valid JSON`)) {
                    message.error('Это нерабочие часы')
                } else {
                    message.error('Ошибка при обновлении события');
                }
            }
        }
    };

    const handleEventResize = async (info) => {
        const event = info.event;
        const newStart = event.start;
        const newEnd = event.end;
        const dayHours = getWorkingHoursForDate(newStart);

        if (user.role === 1 && dayHours && !dayHours.enabled) {
            message.warning('В этот день нет рабочих часов');
            info.revert();
            return;
        }

        if (user.role === 1 && dayHours) {
            const eventStartTime = moment(newStart).format('HH:mm');
            const eventEndTime = moment(newEnd).format('HH:mm');
            
            if (eventStartTime < dayHours.start || eventEndTime > dayHours.end) {
                message.warning('Время события не соответствует рабочему времени');
                info.revert();
                return;
            }
        }

        try {
            const updatedEvent = {
                ...event.extendedProps,
                start_time: newStart.toISOString(),
                end_time: newEnd.toISOString()
            };

            if (event.extendedProps.is_recurring) {
                setEditingEvent(updatedEvent);
                setIsRecurringModalVisible(true);
            } else {
                await updateEvent(updatedEvent, localStorage.getItem('accessToken'));
                await loadEvents();
            }
        } catch (error) {
            console.error('Error updating event:', error);
            message.error('Ошибка при обновлении события');
            info.revert();
        }
    };

    const getBusinessHours = () => {
        if (!workingHours || user?.result?.role !== 0) return null;

        return Object.entries(workingHours)
            .filter(([_, hours]) => hours.enabled)
            .map(([day, hours]) => ({
                daysOfWeek: [getDayNumber(day)],
                startTime: hours.start,
                endTime: hours.end
            }));
    };

    const getMinMaxTime = () => {
        if (!workingHours || user?.result?.role !== 0) {
            return {
                min: new Date(0, 0, 0, 8, 0, 0),
                max: new Date(0, 0, 0, 22, 0, 0)
            };
        }

        const enabledHours = Object.values(workingHours).filter(h => h.enabled);
        if (enabledHours.length === 0) {
            return {
                min: new Date(0, 0, 0, 8, 0, 0),
                max: new Date(0, 0, 0, 22, 0, 0)
            };
        }

        const minTime = enabledHours.reduce((min, h) => {
            const [hours, minutes] = h.start.split(':').map(Number);
            return new Date(0, 0, 0, hours, minutes, 0) < min ? 
                new Date(0, 0, 0, hours, minutes, 0) : min;
        }, new Date(0, 0, 0, 23, 59, 0));

        const maxTime = enabledHours.reduce((max, h) => {
            const [hours, minutes] = h.end.split(':').map(Number);
            return new Date(0, 0, 0, hours, minutes, 0) > max ? 
                new Date(0, 0, 0, hours, minutes, 0) : max;
        }, new Date(0, 0, 0, 0, 0, 0));

        return { min: minTime, max: maxTime };
    };

    const { min, max } = getMinMaxTime();

    const TimeSlotWrapper = ({ children, value, resource }) => {
        const isWorkingHours = () => {
            if (!workingHours || user?.result?.role !== 0) return true;
            
            const dayOfWeek = moment(value).format('dddd').toLowerCase();
            const dayHours = workingHours[dayOfWeek];
            
            if (!dayHours || !dayHours.enabled) return false;
            
            const time = moment(value).format('HH:mm');
            return time >= dayHours.start && time < dayHours.end;
        };

        const className = isWorkingHours() ? 'working-hours' : 'non-working-hours';
        
        return (
            <div className={`rbc-time-slot ${className}`}>
                {children}
            </div>
        );
    };

    return (
        <div className="calendar-container">
            {!isTeacher && <div className="calendar-header">
                <Button onClick={handleShowTeacherEvents} disabled={showTeacherEvents}>
                    Показать занятость учителя
                </Button>
                <Button onClick={handleHideTeacherEvents} disabled={!showTeacherEvents}>
                    Скрыть занятость учителя
                </Button>
            </div>}
            <DnDCalendar
                localizer={localizer}
                events={[...events, ...(showTeacherEvents ? teacherEvents : [])]}
                startAccessor="start"
                endAccessor="end"
                defaultView="week"
                min={min}
                max={max}
                businessHours={getBusinessHours()}
                // style={{ height: 500 }}
                onSelectEvent={handleSelectEvent}
                onSelectSlot={handleSelectSlot}
                selectable
                onEventDrop={handleEventDrop}
                onEventResize={handleEventResize}
                components={{
                    event: eventRender,
                    timeSlotWrapper: TimeSlotWrapper
                }}
                views={['month', 'week', 'day']}
                step={30}
                timeslots={2}
            />
            <Drawer
                title={currentEvent?.id ? "Редактировать событие" : "Создать событие"}
                placement="right"
                onClose={() => setDrawerVisible(false)}
                open={drawerVisible}
                width={400}
            >
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleCreateUpdateEvent}
                    initialValues={currentEvent}
                >
                    <Form.Item
                        name="description"
                        label="Описание"
                        rules={[{ required: true, message: 'Пожалуйста, введите описание' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        name="date"
                        label="Дата"
                        rules={[{ required: true, message: 'Пожалуйста, выберите дату' }]}
                        initialValue={currentEvent?.start_time ? moment(currentEvent.start_time) : null}
                    >
                        <DatePicker style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                        name="time"
                        label="Время начала"
                        rules={[{ required: true, message: 'Пожалуйста, выберите время' }]}
                        initialValue={currentEvent?.start_time ? moment(currentEvent.start_time) : null}
                    >
                        <TimePicker format="HH:mm" style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                        name="duration"
                        label="Длительность (минуты)"
                        rules={[{ required: true, message: 'Пожалуйста, введите длительность' }]}
                        initialValue={currentEvent?.start && currentEvent?.end ? moment(currentEvent.end).diff(moment(currentEvent.start), 'minutes') : 30}
                    >
                        <InputNumber min={15} step={15} style={{ width: '100%' }} />
                    </Form.Item>
                    {isTeacher &&
                    <Form.Item
                        name="student_id"
                        label="Студент"
                        rules={[{ required: true, message: 'Пожалуйста, выберите студента' }]}
                    >
                        <Select>
                            {students?.map(student => (
                                <Select.Option key={student.id} value={student.id}>
                                    {student.name}
                                </Select.Option>
                            ))}
                        </Select>
                    </Form.Item>}
                    <Form.Item
                        name="price"
                        label="Цена (руб/час)"
                        rules={[{ required: true, message: 'Пожалуйста, введите цену' }]}
                    >
                        <InputNumber min={0} />
                    </Form.Item>
                    <Form.Item
                        name="is_recurring"
                        label="Повторяющееся событие"
                        valuePropName="checked"
                    >
                        <Switch />
                    </Form.Item>
                    <Form.Item
                        noStyle
                        shouldUpdate={(prevValues, currentValues) => prevValues.is_recurring !== currentValues.is_recurring}
                    >
                        {({ getFieldValue }) => 
                            getFieldValue('is_recurring') && (
                                <>
                                    <Form.Item
                                        name="recurrence_pattern"
                                        label="Периодичность"
                                        rules={[{ required: true, message: 'Пожалуйста, выберите периодичность' }]}
                                    >
                                        <Select>
                                            <Select.Option value="weekly">Еженедельно</Select.Option>
                                            {/* <Select.Option value="biweekly">Раз в две недели</Select.Option>
                                            <Select.Option value="monthly">Ежемесячно</Select.Option> */}
                                        </Select>
                                    </Form.Item>
                                    <Form.Item
                                        name="recurrence_end_date"
                                        label="Дата окончания"
                                        rules={[{ required: true, message: 'Пожалуйста, выберите дату окончания' }]}
                                    >
                                        <DatePicker />
                                    </Form.Item>
                                </>
                            )
                        }
                    </Form.Item>
                    <Form.Item
                        name="skipped"
                        label="Пропущено"
                        valuePropName="checked"
                    >
                        <Switch />
                    </Form.Item>
                    <Form.Item>
                        <Space>
                            <Button type="primary" htmlType="submit">
                                {currentEvent?.id ? "Обновить" : "Создать"}
                            </Button>
                            {currentEvent?.id && (
                                <Button type="primary" danger onClick={handleDeleteEvent}>
                                    Удалить
                                </Button>
                            )}
                        </Space>
                    </Form.Item>
                </Form>
            </Drawer>
            <Modal
                title="Обновление повторяющегося события"
                open={isRecurringModalVisible}
                onOk={handleRecurringUpdate}
                onCancel={() => {
                    setIsRecurringModalVisible(false);
                    setEditingEvent(null);
                }}
            >
                <Form layout="vertical">
                    <Form.Item label="Выберите действие">
                        <Select
                            value={recurringOption}
                            onChange={setRecurringOption}
                        >
                            <Select.Option value="single">Только это событие</Select.Option>
                            <Select.Option value="future">Это и все  события</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>
            <Modal
                title="Удаление повторяющегося события"
                open={isRecurringDeleteModalVisible}
                onOk={handleRecurringDelete}
                onCancel={() => {
                    setIsRecurringDeleteModalVisible(false);
                    setEditingEvent(null);
                }}
            >
                <Form layout="vertical">
                    <Form.Item label="Выберите действие">
                        <Select
                            value={recurringOption}
                            onChange={setRecurringOption}
                        >
                            <Select.Option value="single">Только это событие</Select.Option>
                            <Select.Option value="future">Это и все  события</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

// Вспомогательная функция для преобразования дня недели в число
const getDayNumber = (day) => {
    const days = {
        monday: 1,
        tuesday: 2,
        wednesday: 3,
        thursday: 4,
        friday: 5,
        saturday: 6,
        sunday: 0
    };
    return days[day.toLowerCase()];
};

export default MyCalendar;