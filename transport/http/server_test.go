package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/way0218/way/logger"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

type testData struct {
	Path string `json:"path"`
}

func TestServer(t *testing.T) {
	logger := logger.NewLogger(logger.WithConsoleEncoder())
	srv := NewServer(Logger(logger))
	router := srv.Router()

	router.HandleFunc("/test", HealthCheckHandler)
	router.HandleFunc("/test/index?a=1&b=2", HealthCheckHandler)
	router.HandleFunc("/test/home", HealthCheckHandler)
	router.HandleFunc("/test/products/:id", HealthCheckHandler)

	// group := router.Group("/test")
	// {
	// 	group.GET("/", fn)
	// 	group.HEAD("/index?a=1&b=2", fn)
	// 	group.OPTIONS("/home", fn)
	// 	group.PUT("/products/:id", fn)
	// 	group.POST("/products/:id", fn)
	// 	group.PATCH("/products/:id", fn)
	// 	group.DELETE("/products/:id", fn)
	// }

	time.AfterFunc(time.Second, func() {
		defer srv.Stop()
		testClient(t, srv)
	})
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
}

func testClient(t *testing.T, srv *Server) {
	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/test"},
		{"PUT", "/test/products/1?a=1&b=2"},
		{"POST", "/test/products/2"},
		{"PATCH", "/test/products/3"},
		{"DELETE", "/test/products/4"},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.path, nil)
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, req)
		result := w.Result()
		defer result.Body.Close()
		body, _ := ioutil.ReadAll(result.Body)
		t.Logf("%s", body)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	//创建一个请求
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 我们创建一个 ResponseRecorder (which satisfies http.ResponseWriter)来记录响应
	rr := httptest.NewRecorder()

	//直接使用HealthCheckHandler，传入参数rr,req
	HealthCheckHandler(rr, req)

	// 检测返回的状态码
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 检测返回的数据
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
