// src/api.js
export const signIn = async (username, password) => {
    const response = await fetch('http://localhost:9001/user/sign_in', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
        throw new Error('Ошибка при авторизации');
    }

    return response.json();
};

export const signUp = async (data) => {
    const response = await fetch('http://localhost:9001/user/sign_up', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        throw new Error('Ошибка при регистрации');
    }

    return response.json();
};

export const fetchProfile = async (accessToken) => {
    const response = await fetch('http://localhost:9001/user/me', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken,
        },
    });

    return response.json();
};

export const updateProfile = async (profileData, accessToken) => {
    const response = await fetch(`${BASE_URL}/user/profile`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken,
        },
        body: JSON.stringify(profileData),
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.description || 'Ошибка при обновлении профиля');
    }

    return response.json();
};

const BASE_URL = 'http://localhost:9001';

export const fetchEvents = async (accessToken) => {
    const response = await fetch(`${BASE_URL}/event`, {
        headers: {
            'AccessToken': accessToken
        }
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.description || 'Ошибка при получении событий');
    }

    const data = await response.json();
    return data.result;
};

export const fetchStudents = async (accessToken) => {
    const response = await fetch(`${BASE_URL}/user/students`, {
        headers: {
            'AccessToken': accessToken,
        }
    });
    if (!response.ok) {
        throw new Error('Ошибка при получении событий');
    }
    const {result} = await response.json();
    return result;
};

export const createEvent = async (eventData, accessToken) => {
    const response = await fetch(`${BASE_URL}/event`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken
        },
        body: JSON.stringify(eventData)
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.description || 'Ошибка при создании события');
    }

    return response.json();
};

export const updateEvent = async (eventData, accessToken) => {
    const response = await fetch(`${BASE_URL}/event`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken
        },
        body: JSON.stringify(eventData)
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.description || 'Ошибка при обновлении события');
    }

    return response.json();
};

export const deleteEvent = async (eventId, accessToken) => {
    const response = await fetch(`${BASE_URL}/event`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken
        },
        body: JSON.stringify({ id: Number(eventId) })
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.description || 'Ошибка при удалении события');
    }

    return response.json();
};

export const fetchStudentStats = async (days, accessToken) => {
    const response = await fetch(`${BASE_URL}/user/students/stats?days=${days}`, {
        headers: {
            'AccessToken': accessToken,
        }
    });
    if (!response.ok) {
        throw new Error('Ошибка при получении статистики студентов');
    }
    const { result } = await response.json();
    return result;
};

export const fetchTeacherEvents = async (accessToken) => {
    const response = await fetch(`${BASE_URL}/event/teacher`, {
        headers: {
            'AccessToken': accessToken
        }
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.description || 'Ошибка при получении событий учителя');
    }

    const data = await response.json();
    return data.result;
};

export const fetchHomeworks = async (accessToken) => {
    const response = await fetch(`${BASE_URL}/homework/`, {
        headers: {
            'AccessToken': accessToken,
        }
    });
    if (!response.ok) {
        throw new Error('Ошибка при получении домашних заданий');
    }
    const {result} = await response.json();
    if (!result) {
        return [];
    }
    return result;
};

export const createHomework = async (homeworkData, accessToken) => {
    const response = await fetch(`${BASE_URL}/homework/`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken,
        },
        body: JSON.stringify(homeworkData),
    });
    if (!response.ok) {
        throw new Error('Ошибка при создании домашнего задания');
    }
    return response.json();
};

export const updateHomeworkRating = async (homeworkId, rating, accessToken) => {
    const response = await fetch(`${BASE_URL}/homework/${homeworkId}/rating`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken,
        },
        body: JSON.stringify({
            homework_id: homeworkId,
            rating: rating
        }),
    });
    if (!response.ok) {
        throw new Error('Ошибка при обновлении оценки');
    }
    return response.json();
};

export const uploadHomeworkFile = async (homeworkId, file, accessToken) => {
    const formData = new FormData();
    formData.append('file', file);
    // formData.append('homework_id', homeworkId);

    const response = await fetch(`${BASE_URL}/homework/${homeworkId}/file`, {
        method: 'POST',
        headers: {
            'AccessToken': accessToken,
        },
        body: formData,
    });
    if (!response.ok) {
        throw new Error('Ошибка при загрузке файла');
    }
    return response.json();
};

export const getHomeworkFiles = async (homeworkId, accessToken) => {
    const response = await fetch(`${BASE_URL}/homework/${homeworkId}/file`, {
        headers: {
            'AccessToken': accessToken,
        }
    });
    if (!response.ok) {
        throw new Error('Ошибка при получении файлов');
    }
    const {result} = await response.json();
    return result;
};

export const downloadHomeworkFile = async (fileId, fName, accessToken) => {
    try {
        const response = await fetch(`${BASE_URL}/homework/file/${fileId}`, {
            method: 'GET',
            headers: {
                'AccessToken': accessToken
            }
        });

        if (!response.ok) {
            throw new Error('Не удалось скачать файл');
        }

        const blob = await response.blob();
        const contentDisposition = response.headers.get('content-disposition');
        const fileName = contentDisposition
            ? contentDisposition.split('filename=')[1].replace(/"/g, '')
            : fName;

        // Создаем ссылку для скачивания
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    } catch (error) {
        throw new Error(`Ошибка при скачивании файла: ${error.message}`);
    }
};

export const getProfile = async () => {
    try {
        const response = await fetch(`${BASE_URL}/profile`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('accessToken')}`
            }
        });
        return await response.json();
    } catch (error) {
        console.error('Error fetching profile:', error);
        throw error;
    }
};
