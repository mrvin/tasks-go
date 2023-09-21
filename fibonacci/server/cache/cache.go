package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/mrvin/tasks-go/004-fibonacci/server/config"
)

var ctx = context.Background()

type CacheRDB struct {
	rdb          *redis.Client
	maxCachedNum uint64
}

type Cache interface {
	Connect(conf *config.Config) error
	GetFromCache(from, to uint64) ([]string, error)
	SetToCache(slValFib []string, from, to uint64) error
	GetMaxCachedNum() uint64
	Close() error
}

func (c *CacheRDB) Connect(conf *config.Config) error {
	c.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.DB.Host, conf.DB.Port),
		Password: conf.DB.Password,
		DB:       conf.DB.NameDB,
	})
	//c.rdb = rdb
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	maxCachedNumStr, err := c.rdb.Get(ctx, "maxСachedNum").Result()
	if err != redis.Nil {
		if err != nil {
			return err
		}
		c.maxCachedNum, err = strconv.ParseUint(maxCachedNumStr, 10, 64)
		if err != nil {
			return err
		}
	} else {
		if err := c.SetToCache([]string{"0", "1"}, 0, 1); err != nil {
			return err
		}
	}

	return nil
}

func (c *CacheRDB) GetFromCache(from, to uint64) ([]string, error) {
	slValFib := make([]string, 0, to-from+1)
	for i := from; i <= to; i++ {
		val, err := c.rdb.Get(ctx, strconv.FormatUint(i, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("can't get from cache num %d: %v", i, err)
		}
		slValFib = append(slValFib, val)
	}

	return slValFib, nil
}

func (c *CacheRDB) SetToCache(slValFib []string, from, to uint64) error {
	for i := from; i <= to; i++ {
		if err := c.rdb.Set(ctx, strconv.FormatUint(i, 10), slValFib[i-from], 0).Err(); err != nil {
			return fmt.Errorf("can't set to cash num %d: %v", i, err)
		}
	}
	if err := c.rdb.Set(ctx, "maxСachedNum", strconv.FormatUint(to, 10), 0).Err(); err != nil {
		return fmt.Errorf("can't set to cash maxСachedNum: %v", err)
	}
	c.maxCachedNum = to

	return nil
}

func (c *CacheRDB) Close() error {
	return c.rdb.Close()
}

func (c *CacheRDB) GetMaxCachedNum() uint64 {
	return c.maxCachedNum
}
