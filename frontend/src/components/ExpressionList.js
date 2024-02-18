// src/ExpressionList.js
import React, {useEffect, useState} from 'react';

function ExpressionList() {
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
        const response = await fetch(`http://localhost:8080/get-expressions?page=${pageData.currentPage}`);
        const responseData = await response.json();

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
          <h2 className="text-center mb-4 ">Список выражений</h2>
          {pageData.data && pageData.data.length > 0 ? (
              <div>
                {pageData.data.map((item) => (
                    <div className="card mb-4" key={item.unique_id}>
                      <ul className="list-group list-group-flush">
                        <li className="list-group-item"><strong>UUID запроса: </strong> {item.unique_id}</li>
                        <li className="list-group-item"><strong>Запрос: </strong>{item.query_text} = {item.result.String} </li>
                      </ul>
                    </div>
                ))}
              </div>
          ) : (
              <p>{pageData.data === null ? 'Loading data...' : 'No data available.'}</p>
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

export default ExpressionList;
