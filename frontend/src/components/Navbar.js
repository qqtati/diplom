import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import { Menu, Layout } from 'antd';
import AuthContext from '../context/AuthContext';
import { 
    HomeOutlined, 
    CalendarOutlined, 
    VideoCameraOutlined, 
    UserOutlined,
    BookOutlined
} from '@ant-design/icons';

const { Header } = Layout;

const Navbar = () => {
    const { user, logout } = useContext(AuthContext);
    console.log(user);

    const getMenuItems = () => {
        if (!user) {
            return [
                {
                    key: '5',
                    label: <Link to="/tutor/signup">Регистрация</Link>
                },
                {
                    key: '6',
                    label: <Link to="/tutor/login">Вход</Link>
                }
            ];
        }

        const items = [
            {
                key: 'calendar',
                icon: <CalendarOutlined />,
                label: <Link to="/calendar">Календарь</Link>
            }
        ];

        if (user?.result?.role !== 1) {
            items.push({
                key: 'stats',
                label: <Link to="/tutor/stats">Статистика учеников</Link>
            });
        }

        items.push(
            {
                key: 'whiteboard',
                icon: <BookOutlined />,
                label: <Link to="/whiteboard">Доска</Link>
            },
            {
                key: 'homework',
                icon: <BookOutlined />,
                label: <Link to="/homework">Домашние задания</Link>
            },
            {
                key: 'video',
                icon: <VideoCameraOutlined />,
                label: <Link to="/video-chat/new">Видеочат</Link>
            },
            {
                key: 'profile',
                label: <Link to="/profile">Профиль</Link>
            },
            {
                key: 'logout',
                style: { marginLeft: 'auto' },
                label: <a onClick={logout}>Выйти</a>
            }
        );

        return items;
    };

    return (
        <Header>
            <Menu 
                theme="dark" 
                mode="horizontal" 
                defaultSelectedKeys={['1']}
                items={getMenuItems()}
            />
        </Header>
    );
};

export default Navbar;