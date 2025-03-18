import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import { Menu, Layout } from 'antd';
import AuthContext from '../context/AuthContext';

const { Header } = Layout;

const Navbar = () => {
    const { user, logout } = useContext(AuthContext);

    return (
        <Header>
            <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['1']}>
                {user && (
                    <>
                        <Menu.Item key="1">
                            <Link to="/calendar">Календарь</Link>
                        </Menu.Item>
                        {user.role === 1 && ( // Отображать только для роли "Репетитор"
                            <Menu.Item key="2">
                                <Link to="/tutor/stats">Статистика учеников</Link>
                            </Menu.Item>
                        )}
                        <Menu.Item key="3">
                            <Link to="/whiteboard">Интерактивная доска</Link>
                        </Menu.Item>
                        <Menu.Item key="4">
                            <Link to="/video">Видеочат</Link>
                        </Menu.Item>
                        <Menu.Item key="9">
                            <Link to="/profile">Профиль</Link>
                        </Menu.Item>
                        <Menu.Item key="logout">
                            <button onClick={logout} style={{ background: 'none', border: 'none', color: 'white' }}>
                                Выйти
                            </button>
                        </Menu.Item>
                    </>
                )}
                {!user && (
                    <>
                        <Menu.Item key="5">
                            <Link to="/tutor/signup">Регистрация</Link>
                        </Menu.Item>
                        <Menu.Item key="6">
                            <Link to="/tutor/login">Вход</Link>
                        </Menu.Item>
                        {/*<Menu.Item key="7">*/}
                        {/*    <Link to="/student/signup">Регистрация ученика</Link>*/}
                        {/*</Menu.Item>*/}
                        {/*<Menu.Item key="8">*/}
                        {/*    <Link to="/student/login">Вход ученика</Link>*/}
                        {/*</Menu.Item>*/}
                    </>
                )}
            </Menu>
        </Header>
    );
};

export default Navbar;