// src/Home.js

import React, {useEffect, useState, useContext} from 'react';
import { UserContext } from "../App";

function Home() {
    const { username, tokenExpiration } = useContext(UserContext);


   const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));
   const initialWorkers = {
        'http://localhost:8081': {
            name: 'Worker 1',
            status: 'Offline',
            lastChecked: null,
        },
        'http://localhost:8082': {
            name: 'Worker 2',
            status: 'Offline',
            lastChecked: null,
        },
        'http://localhost:8083': {
            name: 'Worker 3',
            status: 'Offline',
            lastChecked: null,
        },
  }
  const [expression, setExpression] = useState('');
  const [showAlert, setShowAlert] = useState(false);
  const [statusList, setStatusList] = useState(initialWorkers);

  const checkServerStatus = async (url) => {
        try {
            const response = await fetch(url);
            const currentTime = new Date();

            if (response.ok) {
                setStatusList((prevWorkers) => ({
                    ...prevWorkers,
                    [url]: {
                        ...prevWorkers[url],
                        status: 'Online',
                        lastChecked: currentTime.toLocaleTimeString(),
                    },
                }));
            } else {
                setStatusList((prevWorkers) => ({
                    ...prevWorkers,
                    [url]: {
                        ...prevWorkers[url],
                        status: 'Offline',
                        lastChecked: currentTime.toLocaleTimeString(),
                    },
                }));
            }
        } catch (error) {
            setStatusList((prevWorkers) => ({
                ...prevWorkers,
                [url]: {
                    ...prevWorkers[url],
                    status: 'Offline',
                    lastChecked: new Date().toLocaleTimeString(),
                },
            }));
        }
  };

  const handleExpressionChange = (e) => {
    setExpression(e.target.value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    // Отправка данных на эндпоинт
    fetch('http://localhost:8080/add-expression', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ text: expression, user: username }),
    })
      .then(response => response.json())
      .then(data => {
        console.log('Success:', data);
        setShowAlert(true);

        // Спрятать оповещение через 2 секунды
        setTimeout(() => {
          setShowAlert(false);
        }, 2000);
      })
      .catch(error => {
        console.error('Error:', error);
      });
  };

  useEffect(() => {
      const urls = ['http://localhost:8081', 'http://localhost:8082', 'http://localhost:8083'];
      urls.forEach((url) => checkServerStatus(url));

      const intervalWorkersData = setInterval(() => {
          urls.forEach((url) => checkServerStatus(url));
      }, 6000);
      return () => clearInterval(intervalWorkersData)
  }, []);

  return (
      <div className="container">
        <h2 className="text-center mb-5 mt-5">Распределенный вычислитель арифметических выражений</h2>
        <div className="row mb-5">
            <div>
                <h4>Строка отправки выражения</h4>
                {isLoggedIn ?(
                    <form onSubmit={handleSubmit} className="form-floating">
                        <div className="mb-3">
                            <input
                                type="text"
                                className="form-control"
                                id="expressionInput"
                                value={expression}
                                onChange={handleExpressionChange}
                            />
                        </div>
                        <button type="submit" className="btn btn-primary">Отправить</button>
                    </form>
                ) : (
                    <h4 className="text-danger">Отправка выражений доступна только после авторизации</h4>
                )}

                {showAlert && (
                    <div className="alert alert-success mt-3" role="alert">
                        Ваше выражение успешно отправлено!
                    </div>
                )}

            </div>
        </div>

        <div className="row">
            <div>
                <h4>Статус Workers</h4>
                <div className="card mb-5">
                    <ul className="list-group list-group-flush">
                        {Object.entries(statusList).map(([url, worker]) => (
                            <li className="list-group-item" key={url}>
                                {worker.name} :: <a href={url} target="_blank" rel="noreferrer">{url}</a> :: <strong style={{ color: worker.status === 'Online' ? 'green' : 'red' }}>{worker.status}</strong> :: Последняя проверка {worker.lastChecked}
                            </li>
                        ))}
                    </ul>
                </div>
            </div>
        </div>
      </div>
  );
}

export default Home;