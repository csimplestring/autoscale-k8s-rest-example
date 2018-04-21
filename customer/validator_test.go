package customer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckUniqueName_OK(t *testing.T) {
	r := new(repositoryMock)
	r.On("FindByName", "test-1").Return(nil, nil)
	r.On("FindByName", "test-2").Return(&Customer{}, nil)

	v := &validator{
		repository: r,
	}

	isUnique, err := v.checkUniqueName("test-1")
	assert.NoError(t, err)
	assert.True(t, isUnique)

	isUnique, err = v.checkUniqueName("test-2")
	assert.NoError(t, err)
	assert.False(t, isUnique)
}

func TestValidate_OK(t *testing.T) {
	r := new(repositoryMock)
	r.On("FindByName", "test-1").Return(nil, nil)

	v := &validator{
		repository: r,
	}

	name := "test-1"
	addr := "addr-1"
	customer := Customer{
		&name, &addr,
	}

	violations, err := v.validate(customer)
	assert.NoError(t, err)
	assert.Nil(t, violations)
}

func TestValidate_ReturnValidationErrs(t *testing.T) {
	r := new(repositoryMock)
	r.On("FindByName", "test-1").Return(&Customer{}, nil)

	v := &validator{
		repository: r,
	}

	name := "test-1"
	customer := Customer{
		Name: &name,
	}

	violations, err := v.validate(customer)
	assert.NoError(t, err)
	assert.ObjectsAreEqualValues(ValidationErrors{
		&ValidationError{"name", "the value already exists."},
		&ValidationError{"address", "this field is required."},
	}, violations)
}
