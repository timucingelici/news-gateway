package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

func New() (*config, error) {
	c := NewWithDefaults()
	err := c.Parse()
	return c, err
}

func NewWithDefaults() *config {
	return &config{
		60,
		"tcp",
		"redis:6379",
		1,
		1,
	}
}

type config struct {
	FetcherInterval   time.Duration
	RedisProtocol     string
	RedisAddr         string
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration
}

func (c *config) Parse() error {
	var e string

	// FetcherInterval
	e = os.Getenv("FETCHER_INTERVAL")

	v, err := strconv.Atoi(e)

	if e == "" || err != nil {
		return errors.New("FETCHER_INTERVAL needed and must be an integer")
	}

	c.FetcherInterval = time.Duration(v) * time.Second

	// RedisAddr
	e = os.Getenv("REDIS_ADDR")

	if e == "" {
		return errors.New("REDIS_ADDR needed and must must not be empty")
	}

	c.RedisAddr = e

	// RedisReadTimeout
	e = os.Getenv("REDIS_READ_TIMEOUT")

	v, err = strconv.Atoi(e)

	if e == "" || err != nil {
		return errors.New("REDIS_READ_TIMEOUT needed and must be an integer")
	}

	c.RedisReadTimeout = time.Duration(v) * time.Second

	// RedisReadTimeout
	e = os.Getenv("REDIS_WRITE_TIMEOUT")

	v, err = strconv.Atoi(e)

	if e == "" || err != nil {
		return errors.New("REDIS_WRITE_TIMEOUT needed and must be an integer")
	}

	c.RedisWriteTimeout = time.Duration(v) * time.Second

	return nil
}
