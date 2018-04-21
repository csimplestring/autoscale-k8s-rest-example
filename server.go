package main

import (
	"github.com/labstack/echo"
	"github.com/go-redis/redis"
	"github.com/csimplestring/echo-rest-example/customer"
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
		Addr:     redisAddr,
		DB:       redisDB,
		DialTimeout: time.Second,
	})

	customerController := customer.NewController(
		customer.NewRepository(redisCli))

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.GET("/customer/:name", customerController.Get)
	e.POST("/customer", customerController.Create)

	e.Logger.Fatal(e.Start(addr))
}



