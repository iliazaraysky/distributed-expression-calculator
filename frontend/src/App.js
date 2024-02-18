// src/App.js
import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './components/Home';
import ExpressionList from './components/ExpressionList';
import OperationList from './components/OperationList';
import WorkersList from './components/WorkersList';
import RequestDetails from "./components/RequestDetails";
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';

function App() {
  return (
    <Router>
      <div>
        <Navbar bg="dark" variant="dark">
          <Nav className="mr-auto">
            <Nav.Link as={Link} to="/home">Главная</Nav.Link>
            <Nav.Link as={Link} to="/expression-list">Список выражений</Nav.Link>
            <Nav.Link as={Link} to="/operation-list">Список операций</Nav.Link>
            <Nav.Link as={Link} to="/workers-list">Вычислительные мощности</Nav.Link>
          </Nav>
        </Navbar>

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/home" element={<Home />} />
          <Route path="/expression-list" element={<ExpressionList />} />
          <Route path="/operation-list" element={<OperationList />} />
          <Route path="/workers-list" element={<WorkersList />} />
          <Route path="/get-request-by-id/:uuid" element={<RequestDetails />} /> {/* Добавьте новый Route */}
        </Routes>
      </div>
    </Router>
  );
}

export default App;
