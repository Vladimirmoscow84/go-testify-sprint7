package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenRequareIsCorrected(t *testing.T) {
	// создается запрос к сервису
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)

	// создается ResponseRecorder для записи ответа
	responseRecorder := httptest.NewRecorder()

	// вызывается обработчик mainHandle() в который передается  ResponseRecorder и req
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// проверка корректного запроса
	status := responseRecorder.Code
	//  если код неверный, то закрываем тест
	require.Equal(t, status, http.StatusOK, "Status code should be OK")
	//проверка непустого тела
	body := responseRecorder.Body
	assert.NotEmpty(t, body, "Body should be not empty")

}

func TestMainHandlerWhenNameCityNotSupported(t *testing.T) {
	// создается запрос к сервису
	req := httptest.NewRequest("GET", "/cafe?count=4&city=tula", nil)

	// создается ResponseRecorder для записи ответа
	responseRecorder := httptest.NewRecorder()

	// вызывается обработчик mainHandle() в который передается  ResponseRecorder и req
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// проверяется правильность введенного города
	status := responseRecorder.Code
	require.Equal(t, status, http.StatusBadRequest)
	body := responseRecorder.Body.String()
	assert.Equal(t, body, "wrong city value")

}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {

	// создается запрос к сервису
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()

	// вызывается обработчик mainHandle() в который передается  ResponseRecorder и req
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// проверка количества кафе
	totalCount := 4

	status := responseRecorder.Code
	//  проверяем код ответа
	require.Equal(t, status, http.StatusOK, "Status code should be OK")

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")
	assert.Lenf(t, totalCount, len(list), "Expected cafe count %d, got %d", totalCount, len(list))

}
