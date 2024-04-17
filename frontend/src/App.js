// src/App.js
import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
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
import { jwtDecode } from 'jwt-decode';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));
  const [username, setUsername] = useState('');

  useEffect(() => {
    document.title = "Калькулятор by LMS";
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = jwtDecode(token);
      const login = decodedToken.login;
      setUsername(login);
    }
  }, []);

  // Функция для выхода из учетной записи
  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
    window.location.reload()
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
              {isLoggedIn ?(
                  <Nav.Link>{username}</Nav.Link>
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

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/home" element={<Home />} />
          <Route path="/expression-list" element={<ExpressionList />} />
          <Route path="/operation-list" element={<OperationList />} />
          <Route path="/workers-list" element={<WorkersList />} />
          <Route path="/get-request-by-id/:uuid" element={<RequestDetails />} /> {/* Добавьте новый Route */}
          <Route path="/registration" element={<RegistrationPage />} />
          <Route path="/login" element={<LoginPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
