// RequestDetails.js
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

function RequestDetails() {
    const { uuid } = useParams();
    const [requestData, setRequestData] = useState(null);

    useEffect(() => {
        const fetchRequestData = async () => {
            try {
                const response = await fetch(`http://localhost:8080/get-request-by-id/${uuid}`);
                const data = await response.json();
                setRequestData(data);
            } catch (error) {
                console.error('Error fetching request data:', error);
            }
        };

        fetchRequestData();
    }, [uuid]);


    return (
        <div className="container justify-content-center mt-5">
            <h2 className="text-center mb-4">Детали запроса</h2>
            {requestData ? (
                <div className="card mb-4" key={requestData.unique_id}>
                    <ul className="list-group list-group-flush">
                        <li className="list-group-item"><strong>UUID запроса:</strong> {requestData.unique_id}</li>
                        <li className="list-group-item"><strong>Запрос:</strong> {requestData.query_text}</li>
                        <li className="list-group-item"><strong>Сервер:</strong> {requestData.server_name}</li>
                        <li className="list-group-item"><strong>Результат:</strong> {requestData.result}</li>
                        <li className="list-group-item"><strong>Статус:</strong> {requestData.status}</li>
                    </ul>
                </div>
            ) : (
                <p>Loading data...</p>
            )}
        </div>
    );
}

export default RequestDetails;
