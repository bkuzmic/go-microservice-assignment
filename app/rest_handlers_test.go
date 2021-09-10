package app

import (
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type rwMock struct {
	mock.Mock
}

func (rw *rwMock) Header() http.Header {
	rw.Called()
	return nil
}

func (rw *rwMock) Write(b []byte) (int, error) {
	args := rw.Called(b)
	return args.Int(0), args.Error(1)
}

func (rw *rwMock) WriteHeader(statusCode int) {
	rw.Called(statusCode)
}

func TestIndexHandler(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	app := New(nil)
	handler := app.IndexHandler()
	handler.ServeHTTP(&mockResponseWriter, nil)

	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)
	mockResponseWriter.AssertExpectations(t)
}

func TestHealthHandler(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusOK)

	app := New(nil)
	handler := app.HealthHandler()
	handler.ServeHTTP(&mockResponseWriter, nil)

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertExpectations(t)
}

func TestReadinessHandler(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusOK)

	app := New(nil)
	handler := app.ReadinessHandler()
	handler.ServeHTTP(&mockResponseWriter, nil)

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertExpectations(t)
}
