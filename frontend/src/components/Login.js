import React, { useState } from 'react';

const LoginPage = () => {
    const [formData, setFormData] = useState({
        login: '',
        password: '',
    });

    const [message, setMessage] = useState('');

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            const response = await fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });

            if (!response.ok) {
                throw new Error('Ошибка подключения к сети');
            }

            const data = await response.json();

            if (response.status === 200) {
                localStorage.setItem('token', data.token)
                setMessage('Успешная авторизация');
                setTimeout(() => {
                    window.location.reload();
                }, 1000);
            } else if (response.status === 404) {
                setMessage('Пользователь не найден. Проверьте данные');
            } else if (response.status === 401) {
                setMessage('Неправильный логин или пароль');
            } else if (response.status === 500) {
                setMessage('Ошибка сервера. Повторите попытку позже');
            }
        } catch (error) {
            console.error('Error during login:', error);
            setMessage('Ошибка авторизации');
        }
    };

    return (
        <div className="text-center">
            <div className="container justify-content-center mt-5">
                {message && <p className="alert alert-success">{message}</p>}
                <form onSubmit={handleSubmit}>
                    <h2 className="h2 mb-3 font-weight-normal">Авторизация в системе</h2>
                    <p className="text-muted font-italic mb-0"><small>1. Введите логин / пароль</small></p>
                    <p className="text-muted font-italic mb-5 mt-0"><small>2. Нажмите "Вход"</small></p>
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
                    <button className="btn btn-lg btn-primary btn-block mt-2" type="submit">Вход</button>
                </form>
            </div>
        </div>
    );
};

export default LoginPage;


// import React, { useState } from 'react';
//
// const LoginPage = () => {
//     const [formData, setFormData] = useState({
//         login: '',
//         password: '',
//         // email: '',
//     });
//
//     const [message, setMessage] = useState('');
//
//     const handleChange = (e) => {
//         const { name, value } = e.target;
//         setFormData({ ...formData, [name]: value });
//     };
//
//     const handleSubmit = (e) => {
//         e.preventDefault();
//
//         // Отправляем данные на сервер в формате JSON
//         try {
//             const response = fetch('http://localhost:8080/login', {
//                 method: 'POST',
//                 headers: {
//                     'Content-Type': 'application/json',
//                 },
//                 body: JSON.stringify(formData),
//             });
//
//             if (!response.ok) {
//                 throw new Error('Ошибка подключения к сети');
//             }
//
//             if (response.status === 200) {
//                 setMessage('Успешная авторизация');
//             } else if (response.status === 404) {
//                 setMessage('Пользователь не найден. Проверьте данные');
//             } else if (response.status === 401) {
//                 setMessage('Пользователь не авторизирован. Проверьте данные');
//             } else if (response.status === 500) {
//                 setMessage('Ошибка токена. Проверьте данные');
//             }
//         } catch (error) {
//             console.error('Error during registration:', error);
//             setMessage('Ошибка авторизации')
//         }
//     };
//
//     return (
//         <div className="text-center">
//             <div className="container justify-content-center mt-5">
//                 {message && <p className="alert alert-success">{message}</p>}
//                 <form onSubmit={handleSubmit}>
//                     <h2 className="h2 mb-3 font-weight-normal">Авторизация в системе</h2>
//                     <p className="text-muted font-italic mb-0"><small>1. Введите логин / пароль</small></p>
//                     <p className="text-muted font-italic mb-5 mt-0"><small>2. Нажмите "Вход"</small></p>
//                     <div className="row">
//                         <div className="col">
//                             <input
//                                 className="form-control"
//                                 type="text"
//                                 id="login"
//                                 name="login"
//                                 placeholder="Логин"
//                                 value={formData.login}
//                                 onChange={handleChange}
//                                 required
//                             />
//                         </div>
//                         <div className="col">
//                             <input
//                                 className="form-control"
//                                 type="password"
//                                 id="password"
//                                 name="password"
//                                 placeholder="Пароль"
//                                 value={formData.password}
//                                 onChange={handleChange}
//                                 required
//                             />
//                         </div>
//                     </div>
//                     {/*<label className="sr-only" htmlFor="login">Логин:</label>*/}
//                     {/*<label className="sr-only" htmlFor="password">Пароль:</label>*/}
//                     {/*<label className="sr-only" htmlFor="email">Электронная почта:</label>*/}
//                     {/*<input*/}
//                     {/*    className="form-control"*/}
//                     {/*    type="email"*/}
//                     {/*    id="email"*/}
//                     {/*    name="email"*/}
//                     {/*    placeholder="Электронная почта"*/}
//                     {/*    value={formData.email}*/}
//                     {/*    onChange={handleChange}*/}
//                     {/*/>*/}
//                     <button className="btn btn-lg btn-primary btn-block mt-2" type="submit">Вход</button>
//                 </form>
//             </div>
//         </div>
//     );
// };
//
// export default LoginPage;
