package store

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"net"
	"time"
)

func New(protocol, address string, readTimeout time.Duration, writeTimeout time.Duration) (DataStore, error) {
	netConn, err := net.Dial(protocol, address)

	if err != nil {
		return nil, err
	}

	rConn := redis.NewConn(
		netConn, readTimeout, writeTimeout,
	)

	return &store{
		protocol:     protocol,
		address:      address,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		Conn:         rConn,
	}, nil
}

type DataStore interface {
	GetKeys(pattern string) (reply []string, err error)
	Get(key string) (reply string, err error)
	Set(key string, value []byte) (reply interface{}, err error)
	GetAllWithValues(pattern string) (map[string]string, error)
}

type store struct {
	protocol     string
	address      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	Conn         redis.Conn
}

func (s *store) Set(key string, value []byte) (interface{}, error) {
	return s.Conn.Do("SET", key, value)
}

func (s *store) Get(key string) (string, error) {
	resp, err := s.Conn.Do("GET", key)

	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", nil
	}

	return redis.String(s.Conn.Do("GET", key))
}

func (s *store) GetKeys(pattern string) ([]string, error) {
	return redis.Strings(s.Conn.Do("KEYS", pattern))
}

func (s *store) GetAllWithValues(pattern string) (map[string]string, error) {

	var r = make(map[string]string)

	resp, err := s.GetKeys(pattern)

	if err != nil {
		return nil, err
	}

	for _, key := range resp {

		v, err := s.Get(key)

		if err != nil {
			log.Printf("Failed to fetch the news for key %s. Err : %s\n", key, err)
			continue
		}

		r[key] = v

	}

	return r, nil
}
