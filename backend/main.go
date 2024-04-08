package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type User struct {
	Login    string
	Password string
}

type Message struct {
	Text string `json:"text"`
}

type MessageForQueue struct {
	UniqueId     string    `json:"unique_id"`
	QueryText    string    `json:"query_text"`
	CreationTime time.Time `json:"creation_time"`
}

type requestById struct {
	UniqueID   string `json:"unique_id"`
	QueryText  string `json:"query_text"`
	ServerName string `json:"server_name"`
	Result     string `json:"result"`
	Status     string `json:"status"`
}

type WorkerControl struct {
	WorkerName   string    `json:"worker_name"`
	TimeoutData  int       `json:"timeout_data"`
	CreationTime time.Time `json:"creation_time"`
}

type GetAllResults struct {
	UniqueID       string         `json:"unique_id"`
	QueryText      string         `json:"query_text"`
	CreationTime   sql.NullTime   `json:"creation_time"`
	CompletionTime sql.NullTime   `json:"completion_time"`
	ExecutionTime  string         `json:"execution_time"`
	ServerName     sql.NullString `json:"server_name"`
	Result         sql.NullString `json:"result"`
	Status         string         `json:"status"`
}

type PageData struct {
	Data         []GetAllResults `json:"data"`
	TotalItems   int             `json:"total_items"`
	TotalPages   int             `json:"total_pages"`
	CurrentPage  int             `json:"current_page"`
	ItemsPerPage int             `json:"items_per_page"`
}

type LastWorkerStatus struct {
	WorkerName       string    `json:"worker_name"`
	UniqueId         string    `json:"last_task"`
	Status           string    `json:"status"`
	LastTimeoutSetup time.Time `json:"last_timeout_setup"`
	CurrentTimeout   int       `json:"current_timeout"`
}

type GetAllExpression struct {
	UniqueID  string         `json:"unique_id"`
	QueryText string         `json:"query_text"`
	Result    sql.NullString `json:"result"`
}

type PageDataExpression struct {
	Data         []GetAllExpression `json:"data"`
	TotalItems   int                `json:"total_items"`
	TotalPages   int                `json:"total_pages"`
	CurrentPage  int                `json:"current_page"`
	ItemsPerPage int                `json:"items_per_page"`
}

// Подключение к базе данных. DATABASE_URL прописан при развертывании БД
func connectToDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Функция обновления настроек брокера
func updateWorkerControl(db *sql.DB, workerControl WorkerControl) error {
	workerControl.CreationTime = time.Now()
	_, err := db.Exec("UPDATE workers SET timeout = $1, timer_setup_date = $2 WHERE name = $3",
		workerControl.TimeoutData, workerControl.CreationTime, workerControl.WorkerName)
	return err
}

// Добавляем данные в таблицу requests
func insertRequestData(queue MessageForQueue) error {
	db, err := connectToDB()
	if err != nil {
		log.Println("Database connection error, when recording status: ", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO requests (unique_id, query_text, creation_time, status) VALUES ($1, $2, $3, $4)",
		queue.UniqueId, queue.QueryText, queue.CreationTime, "In queue")
	return err
}

// Выборка всех данных из таблицы requests
func selectAllFromRequests(db *sql.DB) ([]GetAllResults, error) {
	rows, err := db.Query("SELECT unique_id, query_text, creation_time, completion_time, server_name, result, COALESCE (execution_time, '00:00:00'::interval) AS execution_time, status FROM requests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []GetAllResults

	for rows.Next() {
		var request GetAllResults
		err := rows.Scan(
			&request.UniqueID,
			&request.QueryText,
			&request.CreationTime,
			&request.CompletionTime,
			&request.ServerName,
			&request.Result,
			&request.ExecutionTime,
			&request.Status)

		if err != nil {
			log.Println("Request error, selectAllFromRequests: ", err)
			return nil, err
		}
		// Валидацию добавил, так как если поле еще не заполнено SQL выдает ошибку NULL
		if !request.CreationTime.Valid {
			request.CreationTime.Time = time.Time{}
		}
		if !request.CompletionTime.Valid {
			request.CompletionTime.Time = time.Time{}
		}
		if !request.ServerName.Valid {
			request.ServerName.String = "N/A"
		}
		if !request.Result.Valid {
			request.Result.String = "N/A"
		}
		results = append(results, request)
	}
	return results, nil
}

// Стартовая страница
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	} else {
		message := Message{Text: "Привет! Cтартовая страница backend"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	}
}

// Отправка сообщения в брокер
func SendMessageToQueue(message string) error {
	conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks", true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uniqueID := uuid.New().String()

	// Создание экземпляра MessageForQueue
	messageForQueue := MessageForQueue{
		UniqueId:     uniqueID,
		QueryText:    message,
		CreationTime: time.Now(),
	}

	// Преобразование в JSON
	jsonBody, err := json.Marshal(messageForQueue)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.PublishWithContext(ctx,
		"", q.Name, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonBody,
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		return fmt.Errorf("Failed to publish a message: %v", err)
	}

	err = insertRequestData(messageForQueue)
	if err != nil {
		log.Println("Insert into database query status, insertRequestData: ", err)
		return err
	}
	return nil
}

// Отправляем задачу в брокер сообщений
func addExpressionHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody Message

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Отправляем сообщение в брокер
	err = SendMessageToQueue(requestBody.Text)
	if err != nil {
		log.Printf("Error message: %v", err)
		http.Error(w, "Failed to send message to the queue: %v", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON-ответ с информацией о запросе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requestBody)

}

// Установка таймаутов для воркеров, если пришел POST - запрос
// Получение настроек воркеров, если пришел GET - запрос
func setupWorkers(w http.ResponseWriter, r *http.Request) {

	// Проверяем, что нам поступил POST - запрос
	if r.Method == http.MethodPost {
		var workerControl WorkerControl
		workerControl.CreationTime = time.Now()

		err := json.NewDecoder(r.Body).Decode(&workerControl)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Подключение к БД
		db, err := connectToDB()
		if err != nil {
			log.Println("Database connection error: ", err)
			return
		}
		defer db.Close()

		err = updateWorkerControl(db, workerControl)

		if err != nil {
			log.Println("Record update error: ", err)
		}

		message := Message{Text: "Данные обновлены"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	}

	// Проверяем, что нам поступил GET - запрос
	if r.Method == http.MethodGet {
		var results []LastWorkerStatus
		// Подключение к БД
		db, err := connectToDB()
		if err != nil {
			log.Println("Database connection error: ", err)
			return
		}
		defer db.Close()

		// Формируем запрос в БД
		rows, err := db.Query("SELECT name, timer_setup_date, status, last_task, timeout FROM workers")
		if err != nil {
			log.Println("Timeout request error: ", err)
			return
		}

		for rows.Next() {
			var lastWorkerStatus LastWorkerStatus
			err = rows.Scan(
				&lastWorkerStatus.WorkerName,
				&lastWorkerStatus.LastTimeoutSetup,
				&lastWorkerStatus.Status,
				&lastWorkerStatus.UniqueId,
				&lastWorkerStatus.CurrentTimeout,
			)
			if err != nil {
				log.Println("Get data from DB error, workerStatusHandler: ", err)
				return
			}
			results = append(results, lastWorkerStatus)
		}
		defer rows.Close()

		// Возвращаем JSON-ответ с информацией о статусе
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	}
}

// Детальная информация по всем запросам. UniqueID, QueryText, CreationTime, CompletionTime, ServerName, Result, ExecutionTime, Status
func getOperationsHandler(w http.ResponseWriter, r *http.Request) {
	// Подключение к БД
	db, err := connectToDB()
	if err != nil {
		log.Println("Database connection error: ", err)
		return
	}
	defer db.Close()

	// Формируем запрос в БД
	results, err := selectAllFromRequests(db)
	if err != nil {
		log.Println("Get data error from requests, selectAllFromRequests: ", err)
		return
	}

	// Получение номера страницы из параметра запроса (если не указано, используется 1)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Разбивка данных на страницы
	itemsPerPage := 5
	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > len(results) {
		endIndex = len(results)
	}

	// Формирование JSON для текущей страницы
	pageData := PageData{
		Data:         results[startIndex:endIndex],
		TotalItems:   len(results),
		TotalPages:   (len(results) + itemsPerPage - 1) / itemsPerPage,
		CurrentPage:  page,
		ItemsPerPage: itemsPerPage,
	}

	jsonData, err := json.Marshal(pageData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Отправка JSON в ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Все выражения без детальной информации. Только UniqueID, QueryText, Result
func getExpressionHandler(w http.ResponseWriter, r *http.Request) {
	var results []GetAllExpression

	// Подключение к БД
	db, err := connectToDB()
	if err != nil {
		log.Println("Database connection error, getExpressionHandler: ", err)
		return
	}
	defer db.Close()

	// Формируем запрос в БД
	rows, err := db.Query("SELECT unique_id, query_text, result FROM requests")
	if err != nil {
		log.Println("Database connection error, getExpressionHandler: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var request GetAllExpression
		err := rows.Scan(
			&request.UniqueID,
			&request.QueryText,
			&request.Result)

		if err != nil {
			log.Println("Request error, getExpressionHandler: ", err)
			return
		}

		if !request.Result.Valid {
			request.Result.String = "N/A"
		}

		results = append(results, request)
	}

	// Получение номера страницы из параметра запроса (если не указано, используется 1)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Разбивка данных на страницы
	itemsPerPage := 5
	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > len(results) {
		endIndex = len(results)
	}

	// Формирование JSON для текущей страницы
	pageData := PageDataExpression{
		Data:         results[startIndex:endIndex],
		TotalItems:   len(results),
		TotalPages:   (len(results) + itemsPerPage - 1) / itemsPerPage,
		CurrentPage:  page,
		ItemsPerPage: itemsPerPage,
	}

	jsonData, err := json.Marshal(pageData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Отправка JSON в ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Получаем сведения о конкретной операции по UniqueID
func getResultByID(w http.ResponseWriter, r *http.Request) {
	// Получение идентификатора из URL
	uniqueId := strings.TrimPrefix(r.URL.Path, "/get-request-by-id/")

	// Подключение к БД
	db, err := connectToDB()
	if err != nil {
		log.Println("Database connection error: ", err)
		return
	}
	defer db.Close()

	var requestData requestById

	// Формируем запрос в БД
	row := db.QueryRow("SELECT unique_id, query_text, server_name, result, status FROM requests WHERE unique_id = $1", uniqueId)
	err = row.Scan(&requestData.UniqueID, &requestData.QueryText, &requestData.ServerName, &requestData.Result, &requestData.Status)

	if err != nil {
		log.Println("SQL error: ", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	// Возвращаем JSON-ответ с информацией о статусе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&requestData)
}

// Обработчик CORS
func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

//func userRegister(w http.ResponseWriter, r *http.Request) {
//	var user User
//
//	err := json.NewDecoder(r.Body).Decode(&user)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Проверка, что логин и пароль не пустые
//	if user.Login == "" || user.Password == "" {
//		http.Error(w, "Username and are requiered", http.StatusBadRequest)
//		return
//	}
//
//	// Подключение к БД
//	db, err := connectToDB()
//	if err != nil {
//		log.Println("Database connection error: ", err)
//		return
//	}
//	defer db.Close()
//
//	// Обращаемся к базе данных, проверяем наличие пользователя
//	var count int
//	err = db.Query("SELECT COUNT(*) FROM users WHERE login=$1", user.Login).Scan(&count)
//	if err != nil {
//		log.Println("Timeout request error: ", err)
//		return
//	}
//
//	if count > 0 {
//	}
//	log.Println("Username already exists", http.StatusConflict)
//	return
//
//	// Добавляем пользователя в БД
//	_, err = db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)",
//		user.Login, user.Password)
//	return
//
//	w.WriteHeader(http.StatusCreated)
//	log.Println(w, "User register successfully")
//}

func main() {
	http.HandleFunc("/", helloHandler)
	http.Handle("/setup-workers", corsHandler(http.HandlerFunc(setupWorkers)))
	http.Handle("/get-request-by-id/", corsHandler(http.HandlerFunc(getResultByID)))
	http.Handle("/get-operations", corsHandler(http.HandlerFunc(getOperationsHandler)))
	http.Handle("/add-expression", corsHandler(http.HandlerFunc(addExpressionHandler)))
	http.Handle("/get-expressions", corsHandler(http.HandlerFunc(getExpressionHandler)))

	fmt.Println("Server is listening on :8080")
	http.ListenAndServe(":8080", nil)
}
