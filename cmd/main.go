package main

import (
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	memcachedAddr := "localhost:11211" //default
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	//for redis initialisation
	redisAddr = "localhost:6379" //default
	redisCache = cache.NewRedisCache(redisAddr)
	//for in memoryCache initialisation
	inMemorycache := cache.NewInMemorycache() //inmemory is basically type of caching mehcanism stored in RAM rather than disk or in dbs.
	key := "User:111"
	user := User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragav@gmail.com",
		CreatedAt: time.Now(),
	}
	//memcacehe pushing
	expiration := 10 * time.Minute
	err := memcachedCache.Set(key, user, expiration)
	if err != nil {
		fmt.Printf("Error in setting the values in Memcached: %v\n", err)

	}
	//memcache retrieving
	var cachedUser User
	err := memcachedCache.Get(key, &cachedUser)
	if err != nil {
		fmt.printf("error in getting the value from memecachedCache: %v\n", err)
	} else {
		fmt.printf("Successfully retrievd the data from memcachedCache")
	}
	//for redis
	expiration = 5 * time.Minute
	err := redisCache.Set(key, user, expiration)
	if err != nil {
		fmt.printf("Error in setting the values in redis%: v\n", err)
	}
	//retreival
	err := redisCache.Get(key, user, expiration)
	if err != nil {
		fmt.printf("Error in fetchig the values from redis: %v\n", err)
	} else {
		fmt.printf("Successfully retrieved the values")
	}
	//in memorycache value setting
	expiration := 1 * time.Hour
	err := inMemorycache.Set(key, user, expiration)
	if err != nil {
		fmt.printf("Error in Setting the Values in inMemorycache: %v\n ", err)
	}
	//in memorycache value getting
	err := inMemorycache.Get(key, user, expiration)
	if err != nil {
		fmt.printf("Error in setting the values to inmemory cache: %v\n", err)
	} else {
		fmt.printf("MISSION SUCCESS!")
	}

}
