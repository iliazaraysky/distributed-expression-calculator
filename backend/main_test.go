package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Подключение к базе данных. DATABASE_URL прописан при развертывании БД
func testConnectToDB() (*sql.DB, error) {
	connStr := "postgres://calculator_db_user:calculator_db_password@postgres:5432/calculator_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// TestHelloHandler - тест для проверки корректности GET-запроса к стартовой странице
func TestHelloHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"text":"Привет! Cтартовая страница backend","user":""}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestGetLoginHandler - тест для проверки корректности GET-запроса к странице авторизации
func TestGetLoginHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"text":"Добро пожаловать на страницу авторизации","user":""}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestGetRegistrationHandler - тест для проверки корректности GET-запроса к странице регистрации
func TestGetRegistrationHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/registration", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"text":"Добро пожаловать на страницу регистрации","user":""}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestGetExpressionsHandler - тест для проверки данных со страницы список выражений
func TestGetExpressionsHandler(t *testing.T) {
	expressionListURL := "http://localhost:8080/get-expressions"

	resp, err := http.Get(expressionListURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response struct {
		Data []struct {
			UniqueID  string `json:"unique_id"`
			QueryText string `json:"query_text"`
			Result    struct {
				String string `json:"String"`
				Valid  bool   `json:"Valid"`
			} `json:"result"`
		} `json:"data"`
		TotalItems   int `json:"total_items"`
		TotalPages   int `json:"total_pages"`
		CurrentPage  int `json:"current_page"`
		ItemsPerPage int `json:"items_per_page"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal("JSON data error", err)
	}

	if len(response.Data) == 0 {
		t.Fatal("Get an empty list")
	}

	expectedTotalItems := 3
	if response.TotalItems != expectedTotalItems {
		t.Errorf("Incorrect number of elements: got %v want %v",
			response.TotalItems, expectedTotalItems)
	}
}

// TestRegistrationHandler - тест для проверки создания пользователя в базе данных
func TestRegisterHandler(t *testing.T) {
	// Тестовые данные которые будем использовать
	registerURL := "http://localhost:8080/registration"

	testUser := map[string]string{
		"login":    "testLMSUser",
		"password": "test",
	}

	// Кодируем JSON
	requestBody, err := json.Marshal(testUser)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", registerURL, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Errorf("Wrong status code: got %v, want: %v", response.StatusCode, http.StatusCreated)
	}

	// Дополнительно проверяем, что пользователь появился в базе данных
	if !userExists("testLMSUser") {
		t.Errorf("User not created in database")
	}

	// Удаляем тестового пользователя после создания
	deleteUserFromDatabase("testLMSUser")
}

// TestLoginHandler - тест для проверки авторизации пользователя
func TestLoginHandler(t *testing.T) {
	loginURL := "http://localhost:8080/login"

	testUser := map[string]string{
		"login":    "test_user",
		"password": "longPassswordForTest",
	}

	// Кодируем JSON
	requestBody, err := json.Marshal(testUser)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", loginURL, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want: %v", response.StatusCode, http.StatusOK)
	}

	expectedContentType := "application/json"
	if contentType := response.Header.Get("Content-Type"); expectedContentType != contentType {
		t.Errorf("Wrong Content-Type: got %v want %v", contentType, expectedContentType)
	}

	//Проверяем, что пришло в теле ответа
	var responseBody map[string]string
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	if _, ok := responseBody["token"]; !ok {
		t.Errorf("Token not found in response: %v", err)
	}
}

func userExists(username string) bool {
	// Подключение к БД
	connStr := "postgres://calculator_db_user:calculator_db_password@localhost:5432/calculator_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var userCheck string
	err = db.QueryRow("SELECT login FROM users WHERE login=$1", username).Scan(&userCheck)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if strings.TrimSpace(userCheck) != strings.TrimSpace(username) {
		log.Fatalf("userCheck: userCheck %s, %s\n", username, userCheck)
		return false
	}

	return true
}

func deleteUserFromDatabase(username string) {
	// Подключение к БД
	connStr := "postgres://calculator_db_user:calculator_db_password@localhost:5432/calculator_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE login=$1", username)
	if err != nil {
		log.Fatal(err)
	}
}
