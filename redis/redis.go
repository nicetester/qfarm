// Package redis provides common methods for getting data from redis.
package redis

import (
	"fmt"
	"log"
	"time"

	"errors"
	"github.com/garyburd/redigo/redis"
)

type Config struct {
	Connection string
	Password   string
}

func NewConfig() *Config {
	return &Config{Connection: "127.0.0.1:6379"}
}

func (c *Config) WithConnection(conn string) *Config {
	c.Connection = conn
	return c
}

func (c *Config) WithPassword(pass string) *Config {
	c.Password = pass
	return c
}

var ErrNotFound = errors.New("Error not found")

// NewService creates new Service using the given redis config.
func NewService(conf *Config) (*Service, error) {
	redisPool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Connection)
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				_, err = c.Do("AUTH", conf.Password)
				if err != nil {
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &Service{rdb: redisPool}, nil
}

// Service provides common methods for getting data from Redis.
type Service struct {
	rdb *redis.Pool
}

// Get returns value from redis for given key.
func (s *Service) Get(key string) ([]byte, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, fmt.Errorf("error while fetching data from redis: %v", err)
	}
	log.Printf("Getting from redis: %s\n", key)

	return reply, nil
}

// Set stores given data under given key with specified ttl.
func (s *Service) Set(key string, ttl int, data interface{}) error {
	conn := s.rdb.Get()
	defer conn.Close()

	var reply interface{}
	var err error
	if ttl < 0 {
		reply, err = conn.Do("SET", key, data)
	} else {
		reply, err = conn.Do("SETEX", key, ttl, data)
	}

	if err != nil {
		return fmt.Errorf("can't insert encoded elements, key: %s, err: %v", key, err)
	}

	result, ok := reply.(string)
	if !ok {
		return fmt.Errorf("can't decode redis response, key: %s", key)
	}
	if result != "OK" {
		return fmt.Errorf("error! redis response with status(%s), key: %s", result, key)
	}

	log.Printf("Set redis done: %s", key)
	return nil
}

// Del delete specified key and associated data.
func (s *Service) Del(key string) error {
	conn := s.rdb.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		return fmt.Errorf("can't delete element, key: %s, err: %v", key, err)
	}

	return nil
}

// Keys returns all redis keys which match the pattern.
func (s *Service) Keys(pattern string) ([]string, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return nil, fmt.Errorf("can't find keys with pattern: %s, err: %v", pattern, err)
	}

	return reply, nil
}

// AddToSet stores given data under given key inside the set.
func (s *Service) AddToSet(setKey string, data []interface{}) (int64, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	n, err := conn.Do("SADD", redis.Args{}.Add(setKey).AddFlat(data)...)
	if err != nil {
		return -1, fmt.Errorf("can't insert encoded elements, key: %s, err: %v", setKey, err)
	}

	result, ok := n.(int64)
	if !ok {
		return -1, fmt.Errorf("can't decode redis response, key: %s", setKey)
	}

	log.Printf("SADD inserted: %d elements, under key %s", result, setKey)
	return result, nil
}

// GetSet returns set of values from redis for given key.
func (s *Service) GetSet(key string) ([]interface{}, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, fmt.Errorf("error while fetching data from redis: %v", err)
	}
	log.Printf("Getting from redis: %s", key)

	return reply, nil
}

// DelKeys delete all keys which match the pattern.
func (s *Service) DelKeys(pattern string) error {
	conn := s.rdb.Get()
	defer conn.Close()

	matchedKeys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return fmt.Errorf("can't find keys with pattern: %s, err: %v", pattern, err)
	}

	for _, s := range matchedKeys {
		_, err := redis.Int(conn.Do("DEL", s))
		if err != nil {
			return fmt.Errorf("can't delete element, key: %s, err: %v", s, err)
		}
	}

	return nil
}

// ListPush push an element to the list.
func (s *Service) ListPush(key string, data interface{}) error {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := conn.Do("RPUSH", key, data)
	if err != nil {
		return fmt.Errorf("can't push data, key: %s, err: %v", key, err)
	}

	result, ok := reply.(int64)
	if !ok {
		return fmt.Errorf("can't decode redis response, key: %s", key)
	}
	if result < 1 {
		return fmt.Errorf("error! redis should add at least one element, result: %d, key: %s", result, key)
	}

	log.Printf("Push redis done: %s", key)
	return nil
}

// ListPop pops element from the list.
func (s *Service) ListPop(key string) (interface{}, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("LPOP", key))
	if err != nil {
		return nil, fmt.Errorf("error while poping data from redis: %v", err)
	}
	log.Printf("Getting from redis: %s\n", key)

	return reply, nil
}

// ListGetLast get last element from the list.
func (s *Service) ListGetLast(key string) (interface{}, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("LINDEX", key, -1))
	if err != nil {
		if err == redis.ErrNil {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error while poping data from redis: %v", err)
	}

	log.Printf("Getting from redis: %s\n", key)
	return reply, nil
}

// ListGetLast get last element from the list.
func (s *Service) ListGetLastElements(key string, elems int) ([][]byte, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	entries, err := redis.ByteSlices(conn.Do("LRANGE", key, -elems, elems-1))
	if err != nil {
		if err == redis.ErrNil {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error while poping data from redis: %v", err)
	}

	log.Printf("Getting from redis: %s\n", key)
	return entries, nil
}

// ListGetAllElements get all elements from the list.
func (s *Service) ListGetAllElements(key string) ([][]byte, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	entries, err := redis.ByteSlices(conn.Do("LRANGE", key, 0, -1))
	if err != nil {
		if err == redis.ErrNil {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error while poping data from redis: %v", err)
	}

	log.Printf("Getting from redis: %s\n", key)
	return entries, nil
}

// Subscribe subscribes to redis pubsub topic and run job func when got any notification from topic.
func (s *Service) Subscribe(topic string, job func(interface{}) error) error {
	conn := s.rdb.Get()
	defer conn.Close()

	if err := conn.Send("SUBSCRIBE", topic); err != nil {
		return fmt.Errorf("error while subscribing to channel: %v", err)
	}

	if err := conn.Flush(); err != nil {
		return fmt.Errorf("error while flush: %v", err)
	}

	for {
		reply, err := conn.Receive()
		if err != nil {
			return err
		}

		if err := job(reply); err != nil {
			return err
		}
	}
}

// Publish publishes data to redis pubsub topic.
func (s *Service) Publish(topic string, data interface{}) error {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := conn.Do("PUBLISH", topic, data)
	if err != nil {
		return fmt.Errorf("can't publish data, topic: %s, err: %v", topic, err)
	}

	result, ok := reply.(int64)
	if !ok {
		return fmt.Errorf("can't decode redis response, key: %s", topic)
	}
	if result < 1 {
		return fmt.Errorf("error! redis should publish to at least one subcriber, result: %d, key: %s", result, topic)
	}

	return nil
}

// SortedSetRank adds value to sorted set and increasing its rank.
// Returns new ranking score.
func (s *Service) SortedSetRank(key string, value string) (string, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("ZINCRBY", key, 1, value))
	if err != nil {
		return "", fmt.Errorf("can't increase score of %s, key: %s, err: %v", value, key, err)
	}

	return string(reply), nil
}

// SortedSetAdd adds value to sorted set with the rank.
// Returns ranking score.
func (s *Service) SortedSetAdd(key string, value []byte, rank int) (int64, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("ZADD", key, rank, value))
	if err != nil {
		return -1, fmt.Errorf("can't add value to set %s, key: %s, err: %v", value, key, err)
	}

	return reply, nil
}


// SortedSetGetAllRev gets all items from the sorted set in reverse order
// (from highest to lowest rank).
func (s *Service) SortedSetGetAllRev(key string) ([][]byte, error) {
	conn := s.rdb.Get()
	defer conn.Close()

	entries, err := redis.ByteSlices(conn.Do("ZREVRANGE", key, 0, -1))
	if err != nil {
		if err == redis.ErrNil {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error while poping data from redis: %v", err)
	}

	log.Printf("Getting from redis: %s\n", key)
	return entries, nil
}
