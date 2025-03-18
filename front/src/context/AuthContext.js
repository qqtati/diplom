// src/context/AuthContext.js
import React, { createContext, useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import {fetchProfile} from "../api";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const history = useHistory();

    const navigate = (path) => {
        history.push(path)
    }

    useEffect(() => {
        (async () => {
            // Проверка существования токена при первой загрузке
            const accessToken = localStorage.getItem('accessToken');
            if (accessToken) {
                // Предполагаем, что данные пользователя можно извлечь локально или через запрашиваемое API
                const storedUser = await fetchProfile(accessToken);
                setUser(storedUser);
            }
        })()
    }, []);

    const login = async (userData, tokens) => {
        localStorage.setItem('accessToken', tokens.access_token);
        localStorage.setItem('refreshToken', tokens.refresh_token);
        userData = await fetchProfile(tokens.access_token);
        localStorage.setItem('user', JSON.stringify(userData));
        setUser(userData);
        navigate('/calendar');
    };

    const logout = () => {
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('user');
        setUser(null);
        navigate('/tutor/login');
    };

    return (
        <AuthContext.Provider value={{ user, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContext;