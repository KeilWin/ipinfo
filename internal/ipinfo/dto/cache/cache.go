package cache

import (
	"errors"
	"fmt"

	"github.com/KeilWin/ipinfo/internal/common"
)

type CacheType string

const (
	ValkeyCacheType CacheType = "valkey"
	RedisCacheType  CacheType = "redis"
)

type Cache interface {
	common.Storage

	AddIpInfo()

	GetIpInfo(ipAddress string)
}

func NewCache(cacheConfig *CacheConfig) (Cache, error) {
	switch cacheConfig.Type {
	case ValkeyCacheType:
		return NewValkeyCache(cacheConfig), nil
	case RedisCacheType:
		return nil, errors.New("redis not implemented")
	default:
		return nil, fmt.Errorf("unknown cache type: %s", cacheConfig.Type)
	}
}
