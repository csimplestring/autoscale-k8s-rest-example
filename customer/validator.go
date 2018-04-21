package customer

// Validator defines the functions to validate Customer data.
type Validator interface {
	validate(customer Customer) (ValidationErrors, error)
}

// ValidationErrors is a slice of ValidationError.
type ValidationErrors []*ValidationError

// ValidationError represents the detailed error message.
type ValidationError struct {
	Field   string
	Message string
}

// newValidator returns the default validator.
func newValidator(repository Repository) Validator {
	return &validator{
		repository: repository,
	}
}

// validator is the default implementation of Validator.
type validator struct {
	repository Repository
}

// validate validates the customer data. If the input data fail to validate,
// ValidationErrors will be returned.
func (v *validator) validate(customer Customer) (ValidationErrors, error) {
	var ve ValidationErrors

	if customer.Name == nil {
		ve = append(ve, fieldMissingError("name"))
	}
	if customer.Address == nil {
		ve = append(ve, fieldMissingError("address"))
	}

	isUnique, err := v.checkUniqueName(*customer.Name)
	if err != nil {
		return nil, err
	}
	if !isUnique {
		ve = append(ve, nonUniqueValueError("name"))
	}

	return ve, nil
}

// checkUniqueName checks if the name already exists in redis.
func (v *validator) checkUniqueName(name string) (isUnique bool, err error) {
	c, err := v.repository.FindByName(name)

	if err != nil {
		return
	}

	isUnique = c == nil
	return
}

func fieldMissingError(field string) *ValidationError {
	return &ValidationError{field, "this field is required."}
}

func nonUniqueValueError(field string) *ValidationError {
	return &ValidationError{field, "the value already exists."}
}
