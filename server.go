package main

import (
	"github.com/labstack/echo"
	"github.com/go-redis/redis"
	"github.com/csimplestring/echo-rest-example/customer"
	"github.com/labstack/gommon/log"
)


// add tests, message queue increse ops, server stop(close connection), validation, flags command, separate files, store cluster

func main() {

	redisCli := redis.NewClient(&redis.Options{
		Addr:     "192.168.99.100:6379",
		DB:       0,
	})

	customerController := customer.NewController(
		customer.NewRepository(redisCli))

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.GET("/customer/:name", customerController.Get)
	e.POST("/customer", customerController.Create)

	e.Logger.Fatal(e.Start(":8088"))
}



