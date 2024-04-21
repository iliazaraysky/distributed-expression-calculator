// src/App.js
import React, {createContext, useEffect, useState} from 'react';
import {BrowserRouter as Router, Link, Route, Routes, Navigate} from 'react-router-dom';
import Home from './components/Home';
import ExpressionList from './components/ExpressionList';
import OperationList from './components/OperationList';
import WorkersList from './components/WorkersList';
import RequestDetails from "./components/RequestDetails";
import RegistrationPage from "./components/Register";
import LoginPage from "./components/Login";
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import Container from 'react-bootstrap/Container';
import {jwtDecode} from 'jwt-decode';
import UserRequestDetails from "./components/UserRequestDetails";

export const UserContext = createContext();

function App() {
    const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));
    const [username, setUsername] = useState('');


    useEffect(() => {
        document.title = "Калькулятор by LMS";

        // Декодируем Token из хранилища, чтобы забрать Login
        const token = localStorage.getItem('token');
        if (token) {
            const decodedToken = jwtDecode(token);
            const login = decodedToken.login;
            setUsername(login);

            // Проверяем не истек ли токен
            const currentTimeInSecond = Math.floor(Date.now() / 1000);
            const exp = decodedToken.exp;
            if (currentTimeInSecond > exp) {
                localStorage.removeItem('token');
                window.location.reload()
            }
        }
    }, []);

    // Функция для выхода из учетной записи
    const handleLogout = () => {
        localStorage.removeItem('token');
        setIsLoggedIn(false);
        window.location.reload();
    };

    return (
        <Router>
            <div>
                <Navbar bg="dark" variant="dark">
                    <Container>
                        <Navbar.Brand href="/home">Калькулятор by <b className="text-danger">LMS</b></Navbar.Brand>
                        <Nav className="mr-auto">
                            <Nav.Link as={Link} to="/home">Главная</Nav.Link>
                            <Nav.Link as={Link} to="/expression-list">Список выражений</Nav.Link>
                            <Nav.Link as={Link} to="/operation-list">Список операций</Nav.Link>
                            <Nav.Link as={Link} to="/workers-list">Вычислительные мощности</Nav.Link>
                        </Nav>
                        <Nav>
                            {isLoggedIn ? (
                                <Nav.Link as={Link} to={`/get-operation-by-user-id/${username}`}>{username}</Nav.Link>
                            ) : (
                                <Nav.Link as={Link} to="/registration">Регистрация</Nav.Link>
                            )}
                            {isLoggedIn ? (
                                <Nav.Link onClick={handleLogout}>Выход</Nav.Link>
                            ) : (
                                <Nav.Link as={Link} to="/login">Вход в кабинет</Nav.Link>
                            )}
                        </Nav>
                    </Container>
                </Navbar>

                <UserContext.Provider value={{username}}>
                    <Routes>
                        <Route path="/" element={<Home/>}/>
                        <Route path="/home" element={<Home/>}/>
                        <Route path="/expression-list" element={<ExpressionList/>}/>
                        <Route path="/operation-list" element={<OperationList/>}/>
                        <Route path="/workers-list" element={<WorkersList/>}/>
                        <Route path="/get-request-by-id/:uuid" element={<RequestDetails/>}/>
                        <Route path="/get-operation-by-user-id/:username" element={<UserRequestDetails/>}/>
                        <Route path="/registration" element={<RegistrationPage/>}/>
                        <Route path="/login" element={<LoginPage/>}/>
                    </Routes>
                </UserContext.Provider>
            </div>
        </Router>
    );
}

export default App;
