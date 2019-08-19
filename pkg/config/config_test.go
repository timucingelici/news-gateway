package config_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/timucingelici/news-gateway/pkg/config"
	"testing"
	"time"
)

var mock = struct {
	FetcherInterval   time.Duration
	RedisProtocol     string
	RedisAddr         string
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration
}{
	60,
	"tcp",
	"redis:6379",
	1,
	1,
}

func TestConfig_NewShouldReturnError(t *testing.T) {
	_, err := config.New()

	assert.NotNil(t, err)

}
func TestConfig_NewWithDefaults(t *testing.T) {
	conf := config.NewWithDefaults()

	assert.Equal(t, mock.FetcherInterval, conf.FetcherInterval)
	assert.Equal(t, mock.RedisProtocol, conf.RedisProtocol)
	assert.Equal(t, mock.RedisAddr, conf.RedisAddr)
	assert.Equal(t, mock.RedisReadTimeout, conf.RedisReadTimeout)
	assert.Equal(t, mock.RedisWriteTimeout, conf.RedisWriteTimeout)

}
