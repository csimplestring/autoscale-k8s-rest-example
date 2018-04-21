package main

import (
	"github.com/csimplestring/echo-rest-example/customer"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"os"
	"strconv"
	"time"
)

func main() {
	addr := os.Getenv("API_SERVER_ADDR")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisCli := redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		DB:          redisDB,
		DialTimeout: time.Second,
	})

	s := newServer(redisCli)

	s.Logger.Fatal(s.Start(addr))
}

func newServer(redisCli *redis.Client) *echo.Echo {

	customerController := customer.NewController(
		customer.NewRepository(redisCli))

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.GET("/customer/:name", customerController.Get)
	e.POST("/customer", customerController.Create)

	return e
}
