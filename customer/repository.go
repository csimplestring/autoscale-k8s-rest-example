package customer

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

// Repository defines the CRUD functions for Customer.
type Repository interface {
	FindByName(name string) (*Customer, error)
	Create(customer *Customer) error
}

// NewRepository returns the default repository implementation.
func NewRepository(client *redis.Client) Repository {
	return &repository{
		store: client,
	}
}

type repository struct {
	store *redis.Client
}

// GenerateID returns the redis hash key.
func (c *repository) GenerateID(name string) string {
	return "customer:" + name
}

// FindByName find the customer based on name, if not found, a nil Customer is returned.
func (c *repository) FindByName(name string) (*Customer, error) {
	r := c.store.Get(c.GenerateID(name))
	redisErr := r.Err()

	if redisErr != nil && redisErr != redis.Nil {
		return nil, redisErr
	}

	if redisErr == redis.Nil {
		return nil, nil
	}

	data, err := r.Bytes()
	if err != nil {
		return nil, err
	}

	customer := &Customer{}
	err = json.Unmarshal(data, customer)

	return customer, err
}

// Create creates a new customer in redis.
func (c *repository) Create(customer *Customer) error {
	b, err := json.Marshal(customer)
	if err != nil {
		return err
	}

	return c.store.Set(c.GenerateID(*customer.Name), string(b), 0).Err()
}
