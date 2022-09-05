package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
)

func TestGetHandler(t *testing.T) {
	//Создаем запрос с указанием нашего хандлера. Так как мы тетируем GET-ендпоинт
	//то нам не нужно передавать тело, поэтому третьим аргументом передаем nil
	req, err := http.NewRequest("GET", "/?fileExt=.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	//Мы создаем ResponseRecorder (реализует интерфейс  http.ResponseWriter)
	//и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := &Handler{}
	//Наш хандлер соответствует интерфейсу http.Handler, а значит
	//мы можем использовать ServeHTTP и напрямую указать Request и ResponseRecorder
	handler.ServeHTTP(rr, req)
	//Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	//Проверяем тело ответа
	//expected := `Parsed query-param with key "name": John`
	expected := "file2.go"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
