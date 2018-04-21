package customer

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) FindByName(name string) (*Customer, error) {
	args := r.Called(name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Customer), args.Error(1)
}

func (r *repositoryMock) Create(customer *Customer) error {
	args := r.Called(customer)
	return args.Error(0)
}

type validatorMock struct {
	mock.Mock
}

func (v *validatorMock) validate(customer Customer) (ValidationErrors, error) {
	args := v.Called(customer)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(ValidationErrors), args.Error(1)
}

func setUpContextGetUser() (echo.Context, *httptest.ResponseRecorder) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/customer/:name")
	ctx.SetParamNames("name")
	ctx.SetParamValues("test-1")

	return ctx, rec
}

func setUpContextCreateUser(body string) (echo.Context, *httptest.ResponseRecorder) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/customer")

	return ctx, rec
}

func TestGetUser_Error(t *testing.T) {
	ctx, _ := setUpContextGetUser()

	r := new(repositoryMock)
	r.On("FindByName", "test-1").Return(nil, errors.New("db error"))

	c := &controller{
		repository: r,
	}

	assert.Error(t, c.Get(ctx))
}

func TestGetUser_OK(t *testing.T) {
	ctx, rec := setUpContextGetUser()

	name := "test-1"
	addr := "addr-1"
	customer := &Customer{
		&name, &addr,
	}

	r := new(repositoryMock)
	r.On("FindByName", "test-1").Return(customer, nil)

	c := &controller{
		repository: r,
	}

	if assert.NoError(t, c.Get(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"name\":\"test-1\",\"address\":\"addr-1\"}", rec.Body.String())
	}
}

func TestGetUser_NotFound(t *testing.T) {
	ctx, rec := setUpContextGetUser()

	r := new(repositoryMock)
	r.On("FindByName", "test-1").Return(nil, nil)

	c := &controller{
		repository: r,
	}

	if assert.NoError(t, c.Get(ctx)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"customer not found\"", rec.Body.String())
	}
}

func TestCreateUser_OK(t *testing.T) {
	ctx, rec := setUpContextCreateUser("{\"name\":\"test-1\",\"address\":\"addr-1\"}")

	name := "test-1"
	addr := "addr-1"
	customer := &Customer{
		&name, &addr,
	}

	r := new(repositoryMock)
	r.On("Create", customer).Return(nil)

	v := new(validatorMock)
	v.On("validate", mock.Anything).Return(nil, nil)

	c := &controller{
		repository: r,
		validator:  v,
	}

	if assert.NoError(t, c.Create(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "", rec.Body.String())
	}
}

func TestCreateUser_ValidationFailed(t *testing.T) {
	ctx, rec := setUpContextCreateUser("{\"name\":\"test-1\",\"address\":\"addr-1\"}")

	name := "test-1"
	addr := "addr-1"
	customer := &Customer{
		&name, &addr,
	}

	r := new(repositoryMock)
	r.On("Create", customer).Return(nil)

	v := new(validatorMock)
	v.On("validate", *customer).Return(ValidationErrors{
		&ValidationError{Field: "name", Message: "non unique name"},
	}, nil)

	c := &controller{
		repository: r,
		validator:  v,
	}

	if assert.NoError(t, c.Create(ctx)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "[{\"Field\":\"name\",\"Message\":\"non unique name\"}]", rec.Body.String())
	}
}

func TestCreateUser_Error(t *testing.T) {
	ctx, _ := setUpContextCreateUser("{\"name\":\"test-1\",\"address\":\"addr-1\"}")

	name := "test-1"
	addr := "addr-1"
	customer := &Customer{
		&name, &addr,
	}

	r := new(repositoryMock)
	r.On("Create", customer).Return(errors.New("db error"))

	v := new(validatorMock)
	v.On("validate", mock.Anything).Return(nil, nil)

	c := &controller{
		repository: r,
		validator:  v,
	}

	assert.Error(t, c.Create(ctx))
}
