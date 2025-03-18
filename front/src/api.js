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

    if (!response.ok) {
        throw new Error('Ошибка при получении данных профиля');
    }

    return response.json();
};

const BASE_URL = 'http://localhost:9001';

export const fetchEvents = async (accessToken) => {
    const response = await fetch(`${BASE_URL}/event/`, {
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
    const response = await fetch(`${BASE_URL}/event/`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken,
        },
        body: JSON.stringify(eventData),
    });
    if (!response.ok) {
        throw new Error('Ошибка при создании события');
    }
    return response.json();
};

export const updateEvent = async (updatedEvent, accessToken) => {
    const response = await fetch(`${BASE_URL}/event/`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken,
        },
        body: JSON.stringify(updatedEvent),
    });
    if (!response.ok) {
        throw new Error('Ошибка при обновлении события');
    }
    return response.json();
};

export const deleteEvent = async (eventId, accessToken) => {
    const response = await fetch(`${BASE_URL}/event/`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'AccessToken': accessToken,
        },
        body: JSON.stringify({ id: eventId }),
    });
    if (!response.ok) {
        throw new Error('Ошибка при удалении события');
    }
    return response.json();
};