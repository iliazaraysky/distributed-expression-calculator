package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/maja42/goval"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Message struct {
	Text string `json:"text"`
}

const worker_name = "worker2"

type WorkerStatus struct {
	WorkerName     string        `json:"worker_name"`
	UniqueId       string        `json:"last_unique_id"`
	QueryText      string        `json:"last_query_text"`
	Result         int           `json:"last_result"`
	CreationTime   time.Time     `json:"creation_time"`
	LastUpdated    time.Time     `json:"last_updated"`
	ExecutionTime  time.Duration `json:"execution_time"'`
	CurrentTimeOut int           `json:"current_timeout"`
	Status         string        `json:"status"`
	RequestStatus  string        `json:"request_status"`
}

type LastWorkerStatus struct {
	WorkerName       string    `json:"worker_name"`
	UniqueId         string    `json:"last_task"`
	Status           string    `json:"status"`
	LastTimeoutSetup time.Time `json:"last_timeout_setup"`
	CurrentTimeout   int       `json:"current_timeout"`
}

type MessageFromQueue struct {
	UniqueId     string    `json:"unique_id"`
	QueryText    string    `json:"query_text"`
	CreationTime time.Time `json:"creation_time"`
}

var message Message
var workerStatus WorkerStatus

// Функция для подключения к базе данных
func connectToDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Обновление статуса последней задачи поступившей в воркер
func updateStatusLastTask(db *sql.DB, workerStatus WorkerStatus) error {
	workerStatus.Status = "ready"
	workerStatus.LastUpdated = time.Now()

	_, err := db.Exec("UPDATE workers SET timer_setup_date = $1, status = $2, last_task= $3 WHERE name = $4",
		workerStatus.LastUpdated, workerStatus.Status, workerStatus.UniqueId, worker_name)
	return err
}

// Функция для обновления данных в таблице workers для worker2
// Необходима для обозначения статуса для страницы состояния воркеров
func updateBeforeWork(db *sql.DB, workerStatus WorkerStatus) error {
	workerStatus.Status = "work"
	_, err := db.Exec("UPDATE workers SET last_task = $1, status = $2, timeout = $3  WHERE name = $4",
		workerStatus.UniqueId, workerStatus.Status, workerStatus.CurrentTimeOut, worker_name)
	if err != nil {
		log.Println("Worker status update error, updateBeforeWork: ", err)
		return err
	}

	_, err = db.Exec("UPDATE requests SET status = $1, server_name = $2 WHERE unique_id = $3",
		"In progress", worker_name, workerStatus.UniqueId)
	if err != nil {
		log.Println("Request status update error, updateBeforeWork: ", err)
		return err
	}

	return nil
}

func updateRequestsTable(db *sql.DB, workerStatus WorkerStatus) error {
	completion_time := time.Now()
	execution_time := completion_time.Sub(workerStatus.CreationTime).Seconds()

	_, err := db.Exec("UPDATE requests SET completion_time = $1, execution_time = $2, result = $3, status = $4 WHERE unique_id = $5",
		completion_time, execution_time, workerStatus.Result, workerStatus.RequestStatus, workerStatus.UniqueId)
	if err != nil {
		return err
	}
	return nil
}

// Получаем сообщение из брокера, если нет заданий ждем
func consumeMessage(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Println("Channel closed, exiting consume loop")
				return
			}

			// Получили сообщение из очереди
			log.Printf("Received a message: %s", msg.Body)

			var receivedMessage MessageFromQueue

			err := json.Unmarshal(msg.Body, &receivedMessage)

			if err != nil {
				log.Println("Failed to decode JSON:", err)
				continue
			}

			// Точка входа для расчета поступившего запроса
			err = ProcessAndStoreData(receivedMessage)

			if err != nil {
				log.Println("Error writing to PostgreSQL: ", err)
				return
			}

			// Подключение к базе данных
			db, err := connectToDB()
			if err != nil {
				log.Println("Failed to connect to the database: ", err)
				return
			}

			defer db.Close()

			// Обновляем статус worker1 в таблице workers
			err = updateStatusLastTask(db, workerStatus)
			if err != nil {
				log.Println("Failed to update data in the database, updateStatusLastTask: ", err)
				return
			}

			// Обновляем таблицу requests
			err = updateRequestsTable(db, workerStatus)
			if err != nil {
				log.Println("Failed to update data in the database, updateRequestsTable: ", err)
				return
			}

		// Ходим в брокер сообщений каждый 5 секунд
		case <-time.After(5 * time.Second):
			workerStatus.WorkerName = worker_name
			workerStatus.Status = "ready"
			workerStatus.RequestStatus = "Done"
			workerStatus.LastUpdated = time.Now()
			log.Println("No messages in the queue, waiting")
		}
	}
}

// Забираем текущий таймаут
func getCurrentTimeout(db *sql.DB) (int, error) {
	var timeoutData int
	rows, err := db.Query("SELECT timeout FROM workers WHERE name=$1", worker_name)
	if err != nil {
		log.Println("Timeout request error: ", err)
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&timeoutData)
		if err != nil {
			log.Println("Rows scan error: ", err)
			return 0, err
		}
	}
	return timeoutData, nil
}

// Проверка на ошибки. Для брокера
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// goval используется для подсчета арифметических выражений
func govalCalculate(query_text string) error {
	log.Println("Expression evaluation has begun: ", workerStatus.UniqueId)
	log.Println("Current timeout: ", workerStatus.CurrentTimeOut)
	time.Sleep(time.Duration(workerStatus.CurrentTimeOut) * time.Second)

	// Поспали, посчитаем
	eval := goval.NewEvaluator()
	result, err := eval.Evaluate(query_text, nil, nil)

	if err != nil {
		result = "Calculation error"
	}

	if intValue, ok := result.(int); ok {
		workerStatus.Result = intValue
		workerStatus.RequestStatus = "Done"
	} else {
		log.Println("Failed to cast to int")
		workerStatus.RequestStatus = "Error"
	}
	return nil
}

// Считаем и записываем в базу данных, таблица requests
func ProcessAndStoreData(data MessageFromQueue) error {
	// Подключение к БД
	db, err := connectToDB()
	if err != nil {
		log.Println("Database connection error: ", err)
		return err
	}
	defer db.Close()

	// Забираем текущую установку таймаута
	timeoutData, err := getCurrentTimeout(db)
	if err != nil {
		log.Println("Timeout request error: ", err)
		return err
	}

	workerStatus.Status = "work"
	workerStatus.UniqueId = data.UniqueId
	workerStatus.CreationTime = data.CreationTime
	workerStatus.LastUpdated = time.Now()
	workerStatus.CurrentTimeOut = timeoutData
	workerStatus.Result = 0
	workerStatus.QueryText = data.QueryText
	workerStatus.ExecutionTime = time.Duration(0)

	err = updateBeforeWork(db, workerStatus)
	if err != nil {
		log.Println("updateBeforeWork error: ", err)
		return err
	}

	// Получаем результат арифметического выражения
	err = govalCalculate(data.QueryText)
	if err != nil {
		log.Println("govalCalculate error: ", err)
		return err
	}
	return nil
}

// Страница приветствия. Главная страница
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	} else {
		message = Message{Text: "Hello, worker2"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	}
}

// Страница для получения статуса по HTTP
func workerStatusHandler(w http.ResponseWriter, r *http.Request) {
	var lastWorkerStatus LastWorkerStatus

	// Подключение к БД
	db, err := connectToDB()
	if err != nil {
		log.Println("Database connection error: ", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, timer_setup_date, status, last_task, timeout FROM workers WHERE name=$1", worker_name)
	if err != nil {
		log.Println("Timeout request error: ", err)
		return
	}

	for rows.Next() {
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
	}
	defer rows.Close()

	// Возвращаем JSON-ответ с информацией о статусе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&lastWorkerStatus)
}

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

func main() {
	// Инициализация соединения с брокером сообщений
	// username и password прописаны в Docker-compose
	var conn *amqp.Connection
	var err error

	// Повторяем попытку подключения, если RabbitMQ недоступен или не успел прогрузится
	for {
		conn, err = amqp.Dial("amqp://user:password@rabbitmq:5672")
		if err == nil {
			// Выходим из цикла, если подключение успешно
			break
		}
		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks", true, false, false, false, nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Запуск горутины для обработки сообщений из очереди
	go consumeMessage(ch, q.Name)

	// Настройка HTTP-сервера и запуск горутины для обработки HTTP-запросов
	http.Handle("/", corsHandler(http.HandlerFunc(helloHandler)))
	http.HandleFunc("/workerStatus", workerStatusHandler)
	go func() {
		fmt.Println("Worker2 HTTP server started on :8082")
		http.ListenAndServe(":8082", nil)
	}()

	// main продолжает выполнение, ожидая сигнал завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
