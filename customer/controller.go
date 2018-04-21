package customer

import (
	"github.com/labstack/echo"
	"net/http"
)

// Controller defines the CRUD handlers for Customer.
type Controller interface {
	Create(ctx echo.Context) error
	Get(ctx echo.Context) error
}

// NewController returns the default controller implementation.
func NewController(repository Repository) Controller {
	return &controller{
		repository: repository,
		validator: newValidator(repository),
	}
}

type controller struct {
	repository Repository
	validator  Validator
}

// Create creates a new customer.
func (c *controller) Create(ctx echo.Context) error {
	customer := &Customer{}
	if err := ctx.Bind(customer); err != nil {
		return err
	}

	validationErrs, err := c.validator.validate(*customer)
	if err != nil {
		return err
	}
	if validationErrs != nil {
		return ctx.JSON(http.StatusBadRequest, validationErrs)
	}

	if err := c.repository.Create(customer); err != nil {
		return err
	}

	return ctx.String(http.StatusCreated, "")
}

// Get returns a customer.
func (c *controller) Get(ctx echo.Context) error {
	name := ctx.Param("name")
	customer, err := c.repository.FindByName(name)
	if err != nil {
		return err
	}

	if customer == nil {
		return ctx.JSON(http.StatusNotFound, "customer not found")
	}

	return ctx.JSON(http.StatusOK, customer)
}

