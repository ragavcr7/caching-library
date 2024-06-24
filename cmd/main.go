package main

/*
import (
	"fmt"
	"time"
	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	memcachedAddr := "localhost:11211" //default
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	//for redis initialisation
	redisAddr := "localhost:6379" //default
	redisPassword := "123"
	redisdB := 0
	redisCache := cache.NewRedisCache(redisAddr, redisPassword, redisdB)
	//for in memoryCache initialisation
	capacity := 100 // Example capacity, adjust as per your requirement
	expiration := 24 * time.Hour
	inMemorycache := cache.NewInMemoryCache(capacity, expiration) //inmemory is basically type of caching mehcanism stored in RAM rather than disk or in dbs.

	key := "User:111"
	user := User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragav@gmail.com",
		CreatedAt: time.Now(),
	}
	//memcacehe pushing
	expiration = 10 * time.Minute
	err := memcachedCache.Set(key, user, expiration)
	if err != nil {
		fmt.Printf("Error in setting the values in Memcached: %v\n", err)

	}
	//memcache retrieving
	var cachedUser User
	err = memcachedCache.Get(key, &cachedUser)
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

type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
}
*/
/* somewhat okish ....
import (
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost:6379"
	//redisPassword := "itsragav"
	redisDB := 0 //default value
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB) // redisPassword
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Example data to cache
	key := "user:111"
	user := cache.User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Cache with Memcached
	expirationMem := 10 * time.Minute
	err := memcachedCache.Set(key, user, expirationMem)
	if err != nil {
		fmt.Printf("Error setting value in Memcached: %v\n", err)
	}

	// Retrieve from Memcached
	var cachedUser cache.User
	cachedValue, found := memcachedCache.Get(key)
	if found != nil {
		fmt.Printf("Key %s not found in Memcached cache\n", key)
	} else {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from Memcached: %+v\n", cachedUser)
	}

	// Cache with Redis
	expirationRedis := 5 * time.Minute
	err = redisCache.Set(key, user, expirationRedis)
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
	}

	// Retrieve from Redis
	cachedValue, err = redisCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else if cachedValue == nil {
		fmt.Printf("Key %s not found in Redis cache\n", key)
	} else {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from Redis: %+v\n", cachedUser)
	}

	// Cache with InMemory
	expirationInMem := 1 * time.Hour
	inMemoryCache.Set(key, user, expirationInMem) // InMemoryCache's Set method only takes key and value
	cachedValue = inMemoryCache.Get(key)
	if found != nil {
		fmt.Printf("Key %s not found in InMemory cache\n", key)
	} else {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from InMemory cache: %+v\n", cachedUser)
	}
}
*/
/*
import ( //this main working for redis cache
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Use default DB
	})

	// Example data to cache
	key := "user:111"
	user := User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Example context usage
	ctx := context.Background()

	// Marshal user to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshalling user: %v\n", err)
		return
	}

	// Set user in Redis with expiration
	expirationRedis := 5 * time.Minute
	err = rdb.Set(ctx, key, userJSON, expirationRedis).Err()
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
		return
	}
	fmt.Println("User cached in Redis successfully")

	// Retrieve from Redis
	cachedValue, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else {
		var cachedUser User
		// Unmarshal from JSON
		err := json.Unmarshal(cachedValue, &cachedUser)
		if err != nil {
			fmt.Printf("Error unmarshalling user: %v\n", err)
		} else {
			fmt.Printf("User retrieved from Redis: %+v\n", cachedUser)
		}
	}
}
*/
//this main working for inmemory cache
/*
import (
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost:6379"
	redisDB := 0 //default value
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB)
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Example data to cache
	key := "user:111"
	user := cache.User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Cache with Memcached
	expirationMem := 10 * time.Minute
	err := memcachedCache.Set(key, user, expirationMem)
	if err != nil {
		fmt.Printf("Error setting value in Memcached: %v\n", err)
	}

	// Retrieve from Memcached
	var cachedUser cache.User
	cachedValue, found := memcachedCache.Get(key)
	if found != nil {
		fmt.Printf("Key %s not found in Memcached cache\n", key)
	} else {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from Memcached: %+v\n", cachedUser)
	}

	// Cache with Redis
	expirationRedis := 5 * time.Minute
	err = redisCache.Set(key, user, expirationRedis)
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
	}

	// Retrieve from Redis
	cachedValue, err = redisCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else if cachedValue == nil {
		fmt.Printf("Key %s not found in Redis cache\n", key)
	} else {
		// Since redisCache.Get() returns []byte, we need to unmarshal it to cache.User
		err := user.UnmarshalBinary(cachedValue.([]byte))
		if err != nil {
			fmt.Printf("Error unmarshalling user from Redis: %v\n", err)
		} else {
			fmt.Printf("User retrieved from Redis: %+v\n", user)
		}
	}

	// Cache with InMemory
	expirationInMem := 1 * time.Hour
	inMemoryCache.Set(key, user, expirationInMem)
	cachedValue = inMemoryCache.Get(key)
	if cachedValue == nil {
		fmt.Printf("Key %s not found in InMemory cache\n", key)
	} else {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from InMemory cache: %+v\n", cachedUser)
	}
}
*/

/*
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost"
	redisDB := 0 //default value
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB)
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Example data to cache
	key := "user:111"
	user := cache.User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Cache with Memcached
	expirationMem := 10 * time.Minute
	err := memcachedCache.Set(key, user, expirationMem)
	if err != nil {
		fmt.Printf("Error setting value in Memcached: %v\n", err)
	} else {
		fmt.Printf("User cached in Memcached: %+v\n", user)
	}

	// Retrieve from Memcached
	var cachedUser cache.User
	cachedValue, found := memcachedCache.Get(key)
	if found != nil {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from Memcached: %+v\n", cachedUser)
	} else {
		fmt.Printf("Key %s not found in Memcached cache\n", key)
	}

	// Cache with Redis
	expirationRedis := 5 * time.Minute

	// Serialize user to JSON for Redis
	jsonValue, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshalling user to JSON: %v\n", err)
		return
	}

	err = redisCache.Set(key, string(jsonValue), expirationRedis)
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
	} else {
		fmt.Printf("User cached in Redis: %+v\n", user)
	}

	// Retrieve from Redis
	cachedValueBytes, err := redisCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else if cachedValueBytes == nil {
		fmt.Printf("Key %s not found in Redis cache\n", key)
	} else {
		// Deserialize JSON from Redis to user
		var cachedUser cache.User
		err := json.Unmarshal([]byte(cachedValueBytes), &cachedUser)
		if err != nil {
			fmt.Printf("Error unmarshalling user from Redis: %v\n", err)
		} else {
			fmt.Printf("User retrieved from Redis: %+v\n", cachedUser)
		}
	}

	// Cache with InMemory
	expirationInMem := 1 * time.Hour
	inMemoryCache.Set(key, user, expirationInMem)
	cachedValue2 := inMemoryCache.Get(key)
	if cachedValue2 == nil {
		fmt.Printf("Key %s not found in InMemory cache\n", key)
	} else {
		cachedUser = cachedValue2.(cache.User)
		fmt.Printf("User retrieved from InMemory cache: %+v\n", cachedUser)
	}
}
*/
/*
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost:6379"
	redisDB := 0 // Default Redis DB
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB)
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Example data to cache
	key := "user:111"
	user := cache.User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Cache with Memcached
	expirationMem := 10 * time.Minute
	err := memcachedCache.Set(key, user, expirationMem)
	if err != nil {
		fmt.Printf("Error setting value in Memcached: %v\n", err)
	} else {
		fmt.Printf("User cached in Memcached: %+v\n", user)
	}

	// Retrieve from Memcached
	var cachedUser cache.User
	cachedValue, found := memcachedCache.Get(key)
	if found == nil {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from Memcached: %+v\n", cachedUser)
	} else {
		fmt.Printf("Key %s not found in Memcached cache\n", key)
	}

	// Cache with Redis
	expirationRedis := 5 * time.Minute

	// Serialize user to JSON for Redis
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshalling user to JSON: %v\n", err)
		return
	}

	err = redisCache.Set(key, string(userJSON), expirationRedis)
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
	} else {
		fmt.Printf("User cached in Redis: %+v\n", user)
	}

	// Retrieve from Redis
	cachedValueBytes, err := redisCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else if cachedValueBytes != nil {
		// Assert that cachedValueBytes is a string (JSON data)
		cachedJSON, ok := cachedValueBytes.(string)
		if !ok {
			fmt.Println("Unexpected type in Redis cache")
			return
		}

		// Deserialize JSON from Redis to user
		var cachedUser cache.User
		err := json.Unmarshal([]byte(cachedJSON), &cachedUser)
		if err != nil {
			fmt.Printf("Error unmarshalling user from Redis: %v\n", err)
		} else {
			fmt.Printf("User retrieved from Redis: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in Redis cache\n", key)
	}

	// Cache with InMemory
	expirationInMem := 1 * time.Hour
	inMemoryCache.Set(key, user, expirationInMem)
	cachedValue = inMemoryCache.Get(key)
	if cachedValue != nil {
		cachedUser = cachedValue.(cache.User)
		fmt.Printf("User retrieved from InMemory cache: %+v\n", cachedUser)
	} else {
		fmt.Printf("Key %s not found in InMemory cache\n", key)
	}
}
*/
/*previously used
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost:6379"
	redisDB := 0 // Default Redis DB
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB)
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Example data to cache
	key := "user:111"
	user := cache.User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Cache with Memcached
	expirationMem := 10 * time.Minute
	err := memcachedCache.Set(key, user, expirationMem)
	if err != nil {
		fmt.Printf("Error setting value in Memcached: %v\n", err)
	} else {
		fmt.Printf("User cached in Memcached: %+v\n", user)
	}

	// Retrieve from Memcached
	var cachedUser cache.User
	fmt.Println(cachedUser)
	cachedValue, found := memcachedCache.Get(key)
	if found != nil {
		cachedUser, ok := cachedValue.(cache.User)
		if !ok {
			fmt.Printf("Failed to type assert cached value to User type from Memcached\n")
		} else {
			fmt.Printf("User retrieved from Memcached: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in Memcached cache\n", key)
	}

	// Cache with Redis
	expirationRedis := 5 * time.Minute

	// Serialize user to JSON for Redis
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshalling user to JSON: %v\n", err)
		return
	}

	err = redisCache.Set(key, string(userJSON), expirationRedis)
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
	} else {
		fmt.Printf("User cached in Redis: %+v\n", user)
	}

	// Retrieve from Redis
	cachedValueBytes, err := redisCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else if cachedValueBytes != nil {
		// Ensure cachedValueBytes is of type []byte
		cachedJSON, ok := cachedValueBytes.([]byte)
		if !ok {
			fmt.Printf("Failed to type assert cached value to []byte from Redis\n")
		} else {
			var cachedUser cache.User
			err := json.Unmarshal(cachedJSON, &cachedUser)
			if err != nil {
				fmt.Printf("Error unmarshalling user from Redis: %v\n", err)
			} else {
				fmt.Printf("User retrieved from Redis: %+v\n", cachedUser)
			}
		}
	} else {
		fmt.Printf("Key %s not found in Redis cache\n", key)
	}

	// Cache with InMemory
	expirationInMem := 1 * time.Hour
	inMemoryCache.Set(key, user, expirationInMem)
	cachedValue = inMemoryCache.Get(key)
	if cachedValue != nil {
		cachedUser, ok := cachedValue.(cache.User)
		if !ok {
			fmt.Printf("Failed to type assert cached value to User type from InMemory\n")
		} else {
			fmt.Printf("User retrieved from InMemory cache: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in InMemory cache\n", key)
	}
}
*/
/* working except memcache fetching
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost:6379"
	redisDB := 0 // Default Redis DB
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB)
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Example data to cache
	key := "user:111"
	user := cache.User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Cache with Memcached
	expirationMem := 10 * time.Minute
	err := memcachedCache.Set(key, user, expirationMem)
	if err != nil {
		fmt.Printf("Error setting value in Memcached: %v\n", err)
	} else {
		fmt.Printf("User cached in Memcached: %+v\n", user)
	}

	// Retrieve from Memcached
	cachedValue, found := memcachedCache.Get(key)
	if found != nil {
		cachedUser, ok := cachedValue.(cache.User)
		if !ok {
			fmt.Printf("Failed to type assert cached value to User type from Memcached\n")
		} else {
			fmt.Printf("User retrieved from Memcached: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in Memcached cache\n", key)
	}

	// Cache with Redis
	expirationRedis := 5 * time.Minute

	// Serialize user to JSON for Redis
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshalling user to JSON: %v\n", err)
		return
	}

	err = redisCache.Set(key, string(userJSON), expirationRedis)
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
	} else {
		fmt.Printf("User cached in Redis: %+v\n", user)
	}

	// Retrieve from Redis
	cachedValueStr, err := redisCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else if cachedValueStr != "" {
		var cachedUser cache.User
		err := json.Unmarshal([]byte(cachedValueStr), &cachedUser)
		if err != nil {
			fmt.Printf("Error unmarshalling user from Redis: %v\n", err)
		} else {
			fmt.Printf("User retrieved from Redis: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in Redis cache\n", key)
	}

	// Cache with InMemory
	expirationInMem := 1 * time.Hour
	inMemoryCache.Set(key, user, expirationInMem)
	cachedValue = inMemoryCache.Get(key)
	if cachedValue != nil {
		cachedUser, ok := cachedValue.(cache.User)
		if !ok {
			fmt.Printf("Failed to type assert cached value to User type from InMemory\n")
		} else {
			fmt.Printf("User retrieved from InMemory cache: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in InMemory cache\n", key)
	}
}
*/
/***********working one***************/

/*
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost:6379"
	redisDB := 0 // Default Redis DB
	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB)
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Example data to cache
	key := "user:111"
	user := cache.User{
		ID:        111,
		Username:  "ragav",
		Email:     "ragavcr7@yahoo.com",
		CreatedAt: time.Now(),
	}

	// Cache with Memcached
	expirationMem := 10 * time.Minute
	err := memcachedCache.Set(key, user, expirationMem)
	if err != nil {
		fmt.Printf("Error setting value in Memcached: %v\n", err)
	} else {
		fmt.Printf("User cached in Memcached: %+v\n", user)
	}

	// Retrieve from Memcached
	var cachedUser cache.User
	cachedValue, err := memcachedCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Memcached: %v\n", err)
	} else {
		err = json.Unmarshal(cachedValue.([]byte), &cachedUser)
		if err != nil {
			fmt.Printf("Error unmarshalling user from Memcached: %v\n", err)
		} else {
			fmt.Printf("User retrieved from Memcached: %+v\n", cachedUser)
		}
	}

	// Cache with Redis
	expirationRedis := 5 * time.Minute

	// Serialize user to JSON for Redis
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshalling user to JSON: %v\n", err)
		return
	}

	err = redisCache.Set(key, string(userJSON), expirationRedis)
	if err != nil {
		fmt.Printf("Error setting value in Redis: %v\n", err)
	} else {
		fmt.Printf("User cached in Redis: %+v\n", user)
	}

	// Retrieve from Redis
	cachedValueStr, err := redisCache.Get(key)
	if err != nil {
		fmt.Printf("Error getting value from Redis: %v\n", err)
	} else if cachedValueStr != "" {
		var cachedUser cache.User
		err := json.Unmarshal([]byte(cachedValueStr), &cachedUser)
		if err != nil {
			fmt.Printf("Error unmarshalling user from Redis: %v\n", err)
		} else {
			fmt.Printf("User retrieved from Redis: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in Redis cache\n", key)
	}

	// Cache with InMemory
	expirationInMem := 1 * time.Hour
	inMemoryCache.Set(key, user, expirationInMem)
	cachedValue = inMemoryCache.Get(key)
	if cachedValue != nil {
		cachedUser, ok := cachedValue.(cache.User)
		if !ok {
			fmt.Printf("Failed to type assert cached value to User type from InMemory\n")
		} else {
			fmt.Printf("User retrieved from InMemory cache: %+v\n", cachedUser)
		}
	} else {
		fmt.Printf("Key %s not found in InMemory cache\n", key)
	}
}
*/

// trying with manual inputs
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	cache "github.com/ragavcr7/caching-library/api"
	"github.com/ragavcr7/caching-library/cache"
)

func main() {
	// Initialize caches
	memcachedAddr := "localhost:11211"
	redisAddr := "localhost:6379"
	redisDB := 0 // Default Redis DB

	memcachedCache := cache.NewMemcachedCache(memcachedAddr)
	redisCache := cache.NewRedisCache(redisAddr, redisDB)

	// Capacity and expiration for InMemoryCache
	capacity := 100
	expiration := 24 * time.Hour
	inMemoryCache := cache.NewInMemoryCache(capacity, expiration)

	// Create the Gin router
	router := gin.Default()

	// Initialize API handlers
	cacheHandler := api.NewCacheHandler(inMemoryCache, memcachedCache, redisCache)
	userHandler := api.NewUserHandler()

	// Routes for caching endpoints
	router.POST("/cache", cacheHandler.Set)
	router.GET("/cache/:key", cacheHandler.Get)
	router.DELETE("/cache/:key", cacheHandler.Delete)
	router.GET("/cache", cacheHandler.GetAll)
	router.DELETE("/cache", cacheHandler.DeleteAll)

	// Routes for user endpoints
	router.POST("/user", userHandler.CreateUser)
	router.GET("/user/:id", userHandler.GetUser)
	router.PUT("/user/:id", userHandler.UpdateUser)
	router.DELETE("/user/:id", userHandler.DeleteUser)

	// Start the server
	port := 8080
	address := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Example user structure for demonstration
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// Example handler for caching operations
func handleCacheSet(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, ok := data["key"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key"})
		return
	}

	value := data["value"]
	expiration := data["expiration"]

	// Example: store in Redis
	jsonValue, err := json.Marshal(value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal value"})
		return
	}

	err = redisCache.Set(key, string(jsonValue), time.Duration(expiration.(int))*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to set key in Redis: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "value set in cache"})
}

func handleCacheGet(c *gin.Context) {
	key := c.Param("key")

	// Example: retrieve from Redis
	cachedValue, err := redisCache.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get key from Redis: %v", err)})
		return
	}

	var value interface{}
	err = json.Unmarshal([]byte(cachedValue), &value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to unmarshal value: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func handleCacheDelete(c *gin.Context) {
	key := c.Param("key")

	// Example: delete from Redis
	err := redisCache.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete key from Redis: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "key deleted from cache"})
}

func handleCacheGetAll(c *gin.Context) {
	// Example: retrieve all keys from Redis
	keys, err := redisCache.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get all keys from Redis: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func handleCacheDeleteAll(c *gin.Context) {
	// Example: delete all keys from Redis
	err := redisCache.DeleteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete all keys from Redis: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all keys deleted from cache"})
}

// Example handler for user operations
func handleUserCreate(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Example: save user data to database or cache
	// In this case, we just return the received user data as a response
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func handleUserGet(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Example: retrieve user data from database or cache
	user := User{
		ID:        id,
		Username:  "example_user",
		Email:     "example@example.com",
		CreatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func handleUserUpdate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Example: update user data in database or cache
	// In this case, we just return the updated user data as a response
	updatedUser.ID = id
	c.JSON(http.StatusOK, gin.H{"updatedUser": updatedUser})
}

func handleUserDelete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Example: delete user data from database or cache
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with ID %d deleted", id)})
}
