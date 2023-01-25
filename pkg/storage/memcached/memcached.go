package memcached

import (
	"context"
	"errors"
	memcacheInstance "github.com/bradfitz/gomemcache/memcache"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/appengine/memcache"
	"time"
)

type Environment string

const (
	Localhost Environment = "localhost"
	GAE       Environment = "gae"
)

type instance struct {
	Environment Environment `json:"environment"`
	Client      *memcacheInstance.Client
}

func New(env Environment) fiber.Storage {
	var i = &instance{
		Environment: env,
	}
	if env == Localhost {
		i.Client = memcacheInstance.New("127.0.0.1:11211")
	}
	return i
}

func (m *instance) Get(key string) ([]byte, error) {
	var resp interface{}
	var err error
	if m.Client != nil {
		resp, err = m.Client.Get(key)
	} else {
		resp, err = memcache.Get(context.TODO(), key)
	}
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) || errors.Is(err, memcacheInstance.ErrCacheMiss) {
			return nil, nil
		}
		return nil, err
	}
	return resp.(*memcacheInstance.Item).Value, nil
}

func (m *instance) Set(key string, val []byte, ttl time.Duration) error {
	if m.Client != nil {
		return m.Client.Set(&memcacheInstance.Item{
			Key:        key,
			Value:      val,
			Expiration: int32(ttl.Seconds()),
		})
	}
	return memcache.Set(context.TODO(), &memcache.Item{
		Key:        key,
		Value:      val,
		Expiration: ttl,
	})
}

func (m *instance) Delete(key string) error {
	if m.Client != nil {
		return m.Client.Delete(key)
	}
	return memcache.Delete(context.TODO(), key)
}

func (_ *instance) Reset() error {
	return nil
}
func (_ *instance) Close() error {
	return nil
}
