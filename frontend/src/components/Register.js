import React, { useState } from 'react';
import { Link, Navigate } from 'react-router-dom';

const RegistrationPage = () => {
    const [formData, setFormData] = useState({
        login: '',
        password: '',
    });

    const [message, setMessage] = useState('');
    const [isRegistered, setIsRegistered] = useState(false);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        // Отправляем данные на сервер в формате JSON
        fetch('http://localhost:8080/registration', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибки подключения к сети');
                }
                if (response.status === 201) {
                    setMessage('Успешная регистрация');

                    // Если не ставить задержку, он моментально оказывается там, поэтому я установил 1.5 секунды
                    setTimeout(() => {
                        setIsRegistered(true);
                    }, 1500);
                }
                if (response.status === 400) {
                    setMessage('Проверьте логин и пароль. Возможны пустые поля');
                }
            })
            .then(data => {
                console.log('Response data:', data);
                // Я обрабатываю только заголовки. Это неправильно. Лучше делать разбор данных присланных с сервера
                // В этом блоке кода как раз можно прописать подобную логику. Но для экономии времени я ее не рализовал
            })
            .catch(error => {
                console.error('Error during registration:', error);
                setMessage('Ошибка регистрации')
            });
    };

    // Если пользователь успешно зарегистрирован, перенаправляем его на страницу авторизации
    if (isRegistered) {
        return <Navigate to="/login" />;
    }

    return (
        <div className="text-center">
            <div className="container justify-content-center mt-5">
                    {message && <p className="alert alert-success">{message}</p>}
                    <form onSubmit={handleSubmit}>
                        <h2 className="h2 mb-3 font-weight-normal">Регистрация</h2>
                        <p className="text-muted font-italic mb-0"><small>1. Придумайте логин / пароль</small></p>
                        <p className="text-muted font-italic mb-5 mt-0"><small>2. Нажмите создать</small></p>
                        <div className="row">
                            <div className="col">
                                <input
                                    className="form-control"
                                    type="text"
                                    id="login"
                                    name="login"
                                    placeholder="Логин"
                                    value={formData.login}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                            <div className="col">
                                <input
                                    className="form-control"
                                    type="password"
                                    id="password"
                                    name="password"
                                    placeholder="Пароль"
                                    value={formData.password}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>
                        <button className="btn btn-lg btn-primary btn-block mt-2" type="submit">Создать</button>
                    </form>
            </div>
        </div>
    );
};

export default RegistrationPage;
