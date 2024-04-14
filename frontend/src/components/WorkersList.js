// src/components/WorkersList.js
import React, { useEffect, useState } from 'react';

function WorkersList() {
    const initialWorkers = {
        'http://localhost:8081': {
            name: 'Worker 1',
            status: 'Online',
            lastChecked: null,
        },
        'http://localhost:8082': {
            name: 'Worker 2',
            status: 'Online',
            lastChecked: null,
        },
        'http://localhost:8083': {
            name: 'Worker 3',
            status: 'Online',
            lastChecked: null,
        },
    }
    const [workersData, setWorkersData] = useState(null);
    const [selectedWorker, setSelectedWorker] = useState('worker1');
    const [timeoutValue, setTimeoutValue] = useState('');
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

    useEffect(() => {
        const urls = ['http://localhost:8081', 'http://localhost:8082', 'http://localhost:8083'];

        const fetchWorkersData = async () => {
            try {
                const response = await fetch('http://localhost:8080/setup-workers');
                const data = await response.json();
                setWorkersData(data);
            } catch (error) {
                console.error('Error fetching workers data:', error);
            }
        };

        urls.forEach((url) => checkServerStatus(url));
        fetchWorkersData();

        const intervalWorkersData = setInterval(() => {
            fetchWorkersData();
            urls.forEach((url) => checkServerStatus(url));
        }, 6000);
            return () => clearInterval(intervalWorkersData)
        }, []);



    const handleSendRequest = (e) => {
        e.preventDefault();

        const requestData = {
            worker_name: selectedWorker,
            timeout_data: parseInt(timeoutValue, 10)
        };

        fetch('http://localhost:8080/setup-workers', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
        })
            .then(response => {
                console.log(response);
                return response.json();
            })
            .then(data => {
                console.log('Success:', data);
                setShowAlert(true);

                setTimeout(() => {
                    setShowAlert(false);
                }, 2000);
            })
            .catch(error => {
                console.error('Error', error);
            });
    };

    return (
        <div className="container">
            <h2 className="text-center mt-2 mb-4">Список Workers</h2>
            <div className="row">
                <div className="col-md-6">
                    <div>
                        <h4 className="text-center mb-3">Выберите Worker на котором хотите поменять таймаут</h4>
                        <select className="form-control" value={selectedWorker} onChange={(e) => setSelectedWorker(e.target.value)}>
                            <option value="worker1">Worker 1</option>
                            <option value="worker2">Worker 2</option>
                            <option value="worker3">Worker 3</option>
                        </select>
                        <input type="number" className="form-control mt-2" placeholder="Введите число. Таймаут в секундах" value={timeoutValue} onChange={(e) => setTimeoutValue(e.target.value)} />
                        <button className="btn btn-primary mt-2 mb-5" onClick={handleSendRequest}>Отправить</button>
                    </div>
                    {showAlert && (
                        <div className="alert alert-success mt-3" role="alert">
                            Таймаут {selectedWorker} обновлен
                        </div>
                    )}
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

                <div className="col-md-6">
                    <h4 className="text-center mb-3">Текущие настройки Workers</h4>
                    {workersData ? (
                        <div>
                            {workersData.map((worker) => (
                                <div className="card mb-5" key={worker.worker_name}>
                                    <ul className="list-group list-group-flush">
                                        <li className="list-group-item"><strong>Сервер:</strong> {worker.worker_name}</li>
                                        <li className="list-group-item"><strong>Последняя задача:</strong> {worker.last_task}</li>
                                        <li className="list-group-item"><strong>Статус:</strong> {worker.status}</li>
                                        <li className="list-group-item"><strong>Время установки таймаута:</strong> {worker.last_timeout_setup}</li>
                                        <li className="list-group-item"><strong>Текущий таймаут:</strong> {worker.current_timeout}</li>
                                    </ul>
                                </div>
                            ))}
                        </div>
                    ) : (

                        <p>Loading workers data...</p>
                    )}
                </div>
            </div>
        </div>
    );
}

export default WorkersList;
