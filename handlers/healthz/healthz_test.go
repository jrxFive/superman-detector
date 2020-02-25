package healthz

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestHealthz_GetHealthzPass(t *testing.T) {
	e := echo.New()
	defer e.Close()

	mp := &MockPinger{}
	statsdClient := &statsd.NoOpClient{}
	h := NewHealthz(mp, statsdClient)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mp.On("Ping").Return(nil)

	if assert.NoError(t, h.GetHealthz(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"status": "ok"}`, rec.Body.String())
	}
}

func TestHealthz_GetHealthzFail(t *testing.T) {
	e := echo.New()
	defer e.Close()

	mp := &MockPinger{}
	statsdClient := &statsd.NoOpClient{}
	h := NewHealthz(mp, statsdClient)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	rec.Code = 0

	c := e.NewContext(req, rec)

	mp.On("Ping").Return(errors.New("no pong"))
	err := h.GetHealthz(c)

	if assert.Error(t, err) {
		he, _ := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusServiceUnavailable, he.Code)
	}
}

func TestHealthz_HeadHealthz(t *testing.T) {
	e := echo.New()
	defer e.Close()

	mp := &MockPinger{}
	statsdClient := &statsd.NoOpClient{}
	h := NewHealthz(mp, statsdClient)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mp.On("Ping").Return(nil)

	if assert.NoError(t, h.HeadHealthz(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestHealthz_HeadHealthzFail(t *testing.T) {
	e := echo.New()
	defer e.Close()

	mp := &MockPinger{}
	statsdClient := &statsd.NoOpClient{}
	h := NewHealthz(mp, statsdClient)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	rec.Code = 0

	c := e.NewContext(req, rec)

	mp.On("Ping").Return(errors.New("no pong"))
	err := h.HeadHealthz(c)

	if assert.Error(t, err) {
		he, _ := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusServiceUnavailable, he.Code)
	}
}
