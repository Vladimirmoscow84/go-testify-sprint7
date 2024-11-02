package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

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
	body := responseRecorder.Body.String()
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
	assert.Equal(t, status, http.StatusBadRequest)
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

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")
	assert.Lenf(t, totalCount, len(list), "Expected cafe count %d, got %d", totalCount, len(list))

}
