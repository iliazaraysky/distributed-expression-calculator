// src/OperationList.js
import React, {useEffect, useState} from 'react';
import { Link } from "react-router-dom";

function OperationList() {
    const [pageData, setPageData] = useState({
        data: null,
        totalItems: 0,
        totalPages: 1,
        currentPage: 1,
        itemsPerPage: 5,
    });
    const handlePageChange = (newPage) => {
        setPageData({
            ...pageData,
            currentPage: newPage,
        });
    };
    useEffect(() => {
        // Функция для загрузки данных операций
        const fetchOperations = async () => {
            try {
                console.log(localStorage.getItem('token'))
                const token = localStorage.getItem('token')
                const response = await fetch(`http://localhost:8080/get-operations?page=${pageData.currentPage}`, {
                    method: 'GET',
                    headers: new Headers({
                        'Authorization': `Bearer ${token}`,
                    }),
                });
                const responseData = await response.json();
                console.log('Ответ сервера', responseData)

                setPageData({
                    data: responseData.data,
                    totalItems: responseData.total_items,
                    totalPages: responseData.total_pages,
                    currentPage: responseData.current_page,
                    itemsPerPage: responseData.items_per_page,
                });
                console.log(pageData)
            } catch (error) {
                console.error('Error fetching operations:', error);
            }
        };
        fetchOperations();
    }, [pageData.currentPage]);

    return (
        <div>
            <div className="container justify-content-center mt-5">
                <h2 className="text-center mb-4 ">Список операций</h2>
                {pageData.data && pageData.data.length > 0 ? (
                    <div>
                        {pageData.data.map((item) => (
                            <div className="card mb-4" key={item.unique_id}>
                                <ul className="list-group list-group-flush">
                                    <li className="list-group-item"><strong>UUID запроса: </strong><Link to={`/get-request-by-id/${item.unique_id}`} target="_blank" rel="noopener noreferrer">{item.unique_id}</Link></li>
                                    <li className="list-group-item"><strong>Запрос: </strong>{item.query_text}</li>
                                    <li className="list-group-item"><strong>Время
                                        создания: </strong>{item.creation_time.Time}</li>
                                    <li className="list-group-item"><strong>Время
                                        завершения: </strong>{item.completion_time?.Time}</li>
                                    <li className="list-group-item"><strong>Время выполнения: </strong>{item.execution_time}
                                    </li>
                                    <li className="list-group-item"><strong>Сервер: </strong>{item.server_name?.String}</li>
                                    <li className="list-group-item"><strong>Результат: </strong>{item.result?.String}</li>
                                    <li className="list-group-item"><strong>Статус: </strong>{item.status}</li>
                                </ul>
                            </div>
                        ))}
                    </div>
                ) : (
                    // <p>{pageData.data === null ? 'Loading data...' : 'No data available.'}</p>
                    <div className="container justify-content-center">
                        <h2 className="h2 mb-3 font-weight-normal text-center">Это раздел доступен только авторизированным пользователям</h2>
                        <p className="text-center"><i>Зарегистрируйтесь или войдите в систему</i></p>
                        <div className="row text-center">
                            <Link to={`/registration`}>
                                <button className="btn btn-primary mb-2">Регистрация</button>
                            </Link>
                            <Link to={`/login`}>
                                <button className="btn btn-secondary">Войти в кабинет</button>
                            </Link>
                        </div>
                    </div>
                )}
            </div>
            <div className="container justify-content-center mt-2 mb-5">
            {/* Пагинация */}
                {pageData.data && pageData.totalPages > 1 &&(
                    <nav>
                        <ul className="pagination justify-content-center">
                            {Array.from({ length: pageData.totalPages }, (_, index) => (
                                <li key={index + 1}
                                    className={`page-item ${pageData.currentPage === index + 1 ? 'active' : ''}`}>
                                    <button
                                        className="page-link"
                                        onClick={() => handlePageChange(index + 1)}
                                    >
                                        {index + 1}
                                    </button>
                                </li>
                            ))}
                        </ul>
                    </nav>
                )}
            </div>
        </div>
    );
}

export default OperationList;
