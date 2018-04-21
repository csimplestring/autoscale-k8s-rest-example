package main

import (
	"testing"
	"os"
	"strconv"
	"time"
	"context"
	"github.com/go-redis/redis"
	"net/http"
	"github.com/labstack/echo"
	"strings"
	"github.com/stretchr/testify/assert"
	"fmt"
	"io/ioutil"
)

func TestServer(t *testing.T) {

	addr := os.Getenv("TEST_API_SERVER_ADDR")
	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	redisDB, _ := strconv.Atoi(os.Getenv("TEST_REDIS_DB"))

	redisCli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       redisDB,
		DialTimeout: time.Second,
	})

	s := newServer(redisCli)

	go func() {
		if err := s.Start(addr); err != nil {
			s.Logger.Info("shutting down the server")
		}
	}()

	redisCli.Del("customer:test-1")
	redisCli.Del("customer:test-2")

	POSTTests := []struct {
		payload string
		expectedBody string
		expectedStatus int
	}{
		{
			"{\"name\":\"test-1\",\"address\":\"addr-1\"}",
			"",
			201,
		},
		{
			"{\"name\":\"test-2\",\"address\":\"addr-2\"}",
			"",
			201,
		},
		{
			"{\"name\":\"test-1\"}",
			"[{\"Field\":\"address\",\"Message\":\"this field is required.\"},{\"Field\":\"name\",\"Message\":\"the value already exists.\"}]",
			400,
		},
	}

	for i, test := range POSTTests {
		t.Run(fmt.Sprintf("post-test-%d", i), func(t *testing.T) {
			res, err := http.Post("http://"+addr+"/customer", echo.MIMEApplicationJSON, strings.NewReader(test.payload))
			assert.NoError(t, err)

			actual, err := ioutil.ReadAll(res.Body)
			res.Body.Close()

			assert.NoError(t, err)
			assert.Equal(t, test.expectedStatus, res.StatusCode)
			assert.Equal(t, test.expectedBody, string(actual))
		})
	}

	GETTests := []struct {
		expectedBody string
		name string
		expectedStatus int
	}{
		{
			"{\"name\":\"test-1\",\"address\":\"addr-1\"}",
			"test-1",
			200,
		},
		{
			"{\"name\":\"test-2\",\"address\":\"addr-2\"}",
			"test-2",
			200,
		},
		{
			"\"customer not found\"",
			"test-3",
			404,
		},
	}

	for i, test := range GETTests {
		t.Run(fmt.Sprintf("get-test-%d", i), func(t *testing.T) {
			res, err := http.Get("http://"+addr+"/customer/"+test.name)
			assert.NoError(t, err)

			actual, err := ioutil.ReadAll(res.Body)
			res.Body.Close()

			assert.NoError(t, err)
			assert.Equal(t, test.expectedStatus, res.StatusCode)
			assert.Equal(t, test.expectedBody, string(actual))
		})
	}


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Fatal(err)
	}
}




