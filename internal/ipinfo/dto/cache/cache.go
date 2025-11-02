package cache

import (
	"errors"
	"fmt"
)

type CacheType string

const (
	ValkeyCacheType CacheType = "valkey"
	RedisCacheType  CacheType = "redis"
)

type Cache interface {
	StartUp()
	ShutDown()

	AddIpInfo()

	GetIpInfo(ipAddress string)
}

func NewCache(cacheType CacheType) (Cache, error) {
	switch cacheType {
	case ValkeyCacheType:
		return NewValkeyCache(), nil
	case RedisCacheType:
		return nil, errors.New("redis not implemented")
	default:
		return nil, fmt.Errorf("unknown cache type: %s", cacheType)
	}
}
