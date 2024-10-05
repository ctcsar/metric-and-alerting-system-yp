package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	testCases := []struct {
		method       string
		expectedCode int
		Url          string
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, Url: "gauge/AllLoc/234"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "localhost:8080/update/"+tc.Url, nil)
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			Webhook(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым"+r.PathValue("name"))
			// проверим корректность полученного тела ответа, если мы его ожидаем
		})
	}
}
