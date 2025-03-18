import React, { useState, useEffect, useContext } from 'react';
import { 
    Select, 
    List, 
    Card, 
    Upload, 
    Button, 
    Rate, 
    message, 
    Modal,
    Form,
    Input,
    DatePicker
} from 'antd';
import { UploadOutlined, FileOutlined, PlusOutlined, DownloadOutlined } from '@ant-design/icons';
import AuthContext from '../context/AuthContext';
import { 
    fetchHomeworks, 
    createHomework, 
    updateHomeworkRating, 
    uploadHomeworkFile, 
    getHomeworkFiles,
    fetchStudents,
    downloadHomeworkFile
} from '../api';

const { Option } = Select;
const { TextArea } = Input;

const Homework = () => {
    const [students, setStudents] = useState([]);
    const [selectedStudent, setSelectedStudent] = useState(null);
    const [homeworks, setHomeworks] = useState([]);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [form] = Form.useForm();
    const { user } = useContext(AuthContext);
    const isTeacher = user?.result?.role === 0;
    const accessToken = localStorage.getItem('accessToken');

    const loadData = async () => {
        try {
            if (isTeacher) {
                const studentsData = await fetchStudents(accessToken);
                console.log(studentsData)
                setStudents(studentsData);
            } else {
                const homeworksData = await fetchHomeworks(accessToken);
                console.log(homeworksData)
                setHomeworks(homeworksData);
            }
        } catch (error) {
            message.error(error.message);
        }
    };

    useEffect(() => {
        loadData();
    }, [isTeacher, user]);


    const loadStudentHomeworks = async () => {
        if (selectedStudent) {
            try {
                const homeworksData = await fetchHomeworks(accessToken);
                console.log(homeworksData.filter(hw => hw.studentId === selectedStudent), selectedStudent);
                setHomeworks(homeworksData.filter(hw => hw.studentId === selectedStudent));
            } catch (error) {
                message.error(error.message);
            }
        }
    };

    useEffect(() => {
        loadStudentHomeworks();
    }, [selectedStudent, user]);

    const handleFileUpload = async (file, homeworkId) => {
        try {
            await uploadHomeworkFile(homeworkId, file, accessToken);
            await loadData();
            message.success('Файл успешно загружен');
        } catch (error) {
            message.error(error.message);
        }
    };

    const handleRateChange = async (value, homeworkId) => {
        try {
            await updateHomeworkRating(homeworkId, value, accessToken);
            setHomeworks(homeworks.map(hw => 
                hw.id === homeworkId ? { ...hw, rating: value } : hw
            ));
            message.success('Оценка обновлена');
        } catch (error) {
            message.error(error.message);
        }
    };

    const handleCreateHomework = async (values) => {
        try {
            const homeworkData = {
                ...values,
                student_id: selectedStudent,
                due_date: values.due_date.format('YYYY-MM-DD')
            };
            await createHomework(homeworkData, accessToken);
            const homeworksData = await fetchHomeworks(accessToken);
            console.log(homeworkData)
            setHomeworks(homeworksData.filter(hw => hw.student_id === selectedStudent));
            setIsModalVisible(false);
            form.resetFields();
            message.success('Домашнее задание создано');
        } catch (error) {
            message.error(error.message);
        }
    };

    const handleDownloadFile = async (fileId, fileName) => {
        try {
            await downloadHomeworkFile(fileId, fileName, accessToken);
        } catch (error) {
            message.error(error.message);
        }
    };

    const renderHomeworkCard = (homework) => (
        <Card 
            key={homework.id}
            title={`Домашнее задание #${homework.id}`}
            extra={
                isTeacher && (
                    <Rate 
                        value={homework.rating} 
                        onChange={(value) => handleRateChange(value, homework.id)}
                    />
                )
            }
        >
            <p>{homework.description}</p>
            <p>Срок сдачи: {homework.dueDate.substr(0, 10)}</p>
            {!isTeacher && (
                <Upload
                    customRequest={({ file }) => handleFileUpload(file, homework.id)}
                    showUploadList={false}
                >
                    <Button icon={<UploadOutlined />}>Загрузить решение</Button>
                </Upload>
            )}
            {homework.files && homework.files.length > 0 && (
                <div style={{ marginTop: 16 }}>
                    <h4>Прикрепленные файлы:</h4>
                    <List
                        dataSource={homework.files}
                        renderItem={file => (
                            <List.Item
                                actions={[
                                    <Button
                                        type="link"
                                        icon={<DownloadOutlined />}
                                        onClick={() => handleDownloadFile(file.id, file.fileName)}
                                    >
                                        Скачать
                                    </Button>
                                ]}
                            >
                                <List.Item.Meta
                                    avatar={<FileOutlined />}
                                    title={file.fileName}
                                    description={`Загружен: ${new Date(file.uploadedAt).toLocaleString()}`}
                                />
                            </List.Item>
                        )}
                    />
                </div>
            )}
        </Card>
    );

    return (
        <div style={{ padding: 24 }}>
            {isTeacher ? (
                <>
                    <Select
                        style={{ width: 200, marginBottom: 16 }}
                        placeholder="Выберите ученика"
                        onChange={setSelectedStudent}
                        value={selectedStudent}
                    >
                        {students.map(student => (
                            <Option key={student.id} value={student.id}>
                                {student.name}
                            </Option>
                        ))}
                    </Select>
                    {selectedStudent && (
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => setIsModalVisible(true)}
                            style={{ marginLeft: 16 }}
                        >
                            Создать задание
                        </Button>
                    )}
                </>
            ) : null}
            
            <List
                grid={{ gutter: 16, column: 1 }}
                dataSource={homeworks}
                renderItem={renderHomeworkCard}
            />

            <Modal
                title="Создание домашнего задания"
                open={isModalVisible}
                onCancel={() => setIsModalVisible(false)}
                footer={null}
            >
                <Form
                    form={form}
                    onFinish={handleCreateHomework}
                    layout="vertical"
                >
                    <Form.Item
                        name="description"
                        label="Описание"
                        rules={[{ required: true, message: 'Введите описание задания' }]}
                    >
                        <TextArea rows={4} />
                    </Form.Item>
                    <Form.Item
                        name="due_date"
                        label="Срок сдачи"
                        rules={[{ required: true, message: 'Выберите срок сдачи' }]}
                    >
                        <DatePicker style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit">
                            Создать
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default Homework; 