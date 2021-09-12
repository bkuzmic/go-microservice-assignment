package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"go-microservice-assignment/app/models"
	"net/http"
	"strings"
	"testing"
	"time"
)

const personId = "123"

// mock for http.ResponseWriter
type rwMock struct {
	mock.Mock
}

func (rw *rwMock) Header() http.Header {
	rw.Called()
	return http.Header{}
}

func (rw *rwMock) Write(b []byte) (int, error) {
	args := rw.Called(b)
	return args.Int(0), args.Error(1)
}

func (rw *rwMock) WriteHeader(statusCode int) {
	rw.Called(statusCode)
}

// mock for storage.RedisDB
type redisMock struct {
	mock.Mock
}

func (redis *redisMock) GetPerson(ctx context.Context, id string) (*models.Person, error) {
	args := redis.Called(ctx, id)
	return args.Get(0).(*models.Person), args.Error(1)
}

func (redis *redisMock) CreatePerson(ctx context.Context, p *models.Person) error {
	args := redis.Called(ctx, p)
	return args.Error(0)
}

func (redis *redisMock) UpdatePersonOptimistic(ctx context.Context, p *models.Person) (*models.Person, error) {
	args := redis.Called(ctx, p)
	return args.Get(0).(*models.Person), args.Error(1)
}

func (redis *redisMock) UpdatePersonPessimistic(ctx context.Context, p *models.Person) (*models.Person, error) {
	args := redis.Called(ctx, p)
	return args.Get(0).(*models.Person), args.Error(1)
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

func TestGetPersonHandler_OkResponse(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("Header").Return(http.Header{})
	mockResponseWriter.On("WriteHeader", http.StatusOK)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	mockRedis := redisMock{}
	dummyPerson := models.Person{
		Name: "Test123",
		Address: "Berlin 123",
		DateOfBirth: models.JSONDate(time.Date(1981, time.November, 29, 0, 0, 0, 0, time.UTC)),
	}
	mockRedis.On("GetPerson", mock.Anything, personId).Return(&dummyPerson, nil)

	app := New(&mockRedis)
	handler := app.GetPersonHandler()
	handler.ServeHTTP(&mockResponseWriter, createTestGetRequest(false))

	mockRedis.AssertNumberOfCalls(t, "GetPerson", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Header", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestGetPersonHandler_MissingId(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusBadRequest)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	app := New(nil)
	handler := app.GetPersonHandler()
	handler.ServeHTTP(&mockResponseWriter, createTestGetRequest(true))

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)
	mockResponseWriter.AssertExpectations(t)
}

func TestGetPersonHandler_PersonNotFound(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusNotFound)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	mockRedis := redisMock{}
	dummyPerson := models.Person{}
	mockRedis.On("GetPerson", mock.Anything, personId).Return(&dummyPerson, errors.New("redis: nil"))

	app := New(&mockRedis)
	handler := app.GetPersonHandler()
	handler.ServeHTTP(&mockResponseWriter, createTestGetRequest(false))

	mockRedis.AssertNumberOfCalls(t, "GetPerson", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestGetPersonHandler_ServerError(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusInternalServerError)

	mockRedis := redisMock{}
	dummyPerson := models.Person{}
	mockRedis.On("GetPerson", mock.Anything, personId).Return(&dummyPerson, errors.New("server error"))

	app := New(&mockRedis)
	handler := app.GetPersonHandler()
	handler.ServeHTTP(&mockResponseWriter, createTestGetRequest(false))

	mockRedis.AssertNumberOfCalls(t, "GetPerson", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestCreatePersonHandler_OkResponse(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("Header").Return(http.Header{})
	mockResponseWriter.On("WriteHeader", http.StatusCreated)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	dummyPerson := models.Person{
		Name: "Test123",
		Address: "Berlin 123",
		DateOfBirth: models.JSONDate(time.Date(1981, time.November, 29, 0, 0, 0, 0, time.UTC)),
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("POST", "/api/v1/person/123", body)

	mockRedis := redisMock{}
	mockRedis.On("CreatePerson", context.TODO(), mock.AnythingOfType("*models.Person")).Return(nil)

	app := New(&mockRedis)
	handler := app.CreatePersonHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockRedis.AssertNumberOfCalls(t, "CreatePerson", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Header", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestCreatePersonHandler_BodyInvalidJson(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusBadRequest)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	body := strings.NewReader("{\"wrong\":\"body\"")
	testRequest,_ := http.NewRequest("POST", "/api/v1/person/123", body)

	app := New(nil)
	handler := app.CreatePersonHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockResponseWriter.AssertExpectations(t)
}

func TestCreatePersonHandler_CreatePerson_ServerError(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusInternalServerError)

	dummyPerson := models.Person{
		Name: "Test123",
		Address: "Berlin 123",
		DateOfBirth: models.JSONDate(time.Date(1981, time.November, 29, 0, 0, 0, 0, time.UTC)),
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("POST", "/api/v1/person/123", body)

	mockRedis := redisMock{}
	mockRedis.On("CreatePerson", context.TODO(), mock.AnythingOfType("*models.Person")).Return(errors.New("server error"))

	app := New(&mockRedis)
	handler := app.CreatePersonHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockRedis.AssertNumberOfCalls(t, "CreatePerson", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonOptimisticHandler_OkResponse(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("Header").Return(http.Header{})
	mockResponseWriter.On("WriteHeader", http.StatusOK)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	dummyPerson := models.Person{
		Id: "testId",
		Name: "Test123",
		Address: "Berlin 123",
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	mockRedis := redisMock{}
	mockRedis.On("UpdatePersonOptimistic", context.TODO(), mock.AnythingOfType("*models.Person")).Return(&dummyPerson, nil)

	app := New(&mockRedis)
	handler := app.UpdatePersonOptimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockRedis.AssertNumberOfCalls(t, "UpdatePersonOptimistic", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Header", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonOptimisticHandler_BodyInvalidJson(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusBadRequest)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	body := strings.NewReader("{\"wrong\":\"body\"")
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	app := New(nil)
	handler := app.UpdatePersonOptimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonOptimisticHandler_MissingPersonId(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusBadRequest)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	dummyPerson := models.Person{
		Name: "Test123",
		Address: "Berlin 123",
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	app := New(nil)
	handler := app.UpdatePersonOptimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)
	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonOptimisticHandler_UpdatePersonOptimistic_ServerError(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusInternalServerError)

	dummyPerson := models.Person{
		Id: "testId",
		Name: "Test123",
		Address: "Berlin 123",
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	mockRedis := redisMock{}
	mockRedis.On("UpdatePersonOptimistic", context.TODO(), mock.AnythingOfType("*models.Person")).Return(&dummyPerson, errors.New("server error"))

	app := New(&mockRedis)
	handler := app.UpdatePersonOptimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockRedis.AssertNumberOfCalls(t, "UpdatePersonOptimistic", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonPessimisticHandler_OkResponse(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("Header").Return(http.Header{})
	mockResponseWriter.On("WriteHeader", http.StatusOK)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	dummyPerson := models.Person{
		Id: "testId",
		Name: "Test123",
		Address: "Berlin 123",
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	mockRedis := redisMock{}
	mockRedis.On("UpdatePersonPessimistic", context.TODO(), mock.AnythingOfType("*models.Person")).Return(&dummyPerson, nil)

	app := New(&mockRedis)
	handler := app.UpdatePersonPessimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockRedis.AssertNumberOfCalls(t, "UpdatePersonPessimistic", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Header", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonPessimisticHandler_BodyInvalidJson(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusBadRequest)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	body := strings.NewReader("{\"wrong\":\"body\"")
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	app := New(nil)
	handler := app.UpdatePersonPessimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)

	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonPessimisticHandler_MissingPersonId(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusBadRequest)
	mockResponseWriter.On("Write", mock.Anything).Return(1, nil)

	dummyPerson := models.Person{
		Name: "Test123",
		Address: "Berlin 123",
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	app := New(nil)
	handler := app.UpdatePersonPessimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "Write", 1)
	mockResponseWriter.AssertExpectations(t)
}

func TestUpdatePersonPessimisticHandler_UpdatePersonPessimistic_ServerError(t *testing.T) {
	mockResponseWriter := rwMock{}
	mockResponseWriter.On("WriteHeader", http.StatusInternalServerError)

	dummyPerson := models.Person{
		Id: "testId",
		Name: "Test123",
		Address: "Berlin 123",
	}
	request, _ := json.Marshal(&dummyPerson)
	body := strings.NewReader(string(request))
	testRequest,_ := http.NewRequest("PATCH", "/api/v1/person/123", body)

	mockRedis := redisMock{}
	mockRedis.On("UpdatePersonPessimistic", context.TODO(), mock.AnythingOfType("*models.Person")).Return(&dummyPerson, errors.New("server error"))

	app := New(&mockRedis)
	handler := app.UpdatePersonPessimisticHandler()
	handler.ServeHTTP(&mockResponseWriter, testRequest)

	mockRedis.AssertNumberOfCalls(t, "UpdatePersonPessimistic", 1)
	mockResponseWriter.AssertNumberOfCalls(t, "WriteHeader", 1)

	mockRedis.AssertExpectations(t)
	mockResponseWriter.AssertExpectations(t)
}

func createTestGetRequest(useEmptyVars bool) *http.Request {
	var vars map[string]string
	if useEmptyVars {
		vars = map[string]string{}
	} else {
		vars = map[string]string {
			"id": personId,
		}
	}
	testRequest,_ := http.NewRequest("GET", "/api/v1/person/123", nil)
	testRequest = mux.SetURLVars(testRequest, vars)
	return testRequest
}
