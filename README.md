Caching Library Documentation

Overview:-
	This caching library provides two types of caches: an in-memory cache and a Memcached cache. Both caches support essential operations such as setting, getting, deleting keys, and retrieving all keys. The in-memory cache uses a Least Recently Used (LRU) eviction policy and supports key expiration, while the Memcached cache interacts with a Memcached server and also supports key expiration.

Key Features:
   In-Memory Cache

	LRU Eviction: Automatically removes the least recently used items when the cache reaches its capacity.
	Expiration: Supports time-based expiration of cache entries.
	High Performance: Optimized for fast access to cache data within a single application instance.

   Memcached Cache

	Distributed Caching: Suitable for caching data across multiple servers.
	Expiration: Supports time-based expiration of cache entries.
	Integration: Utilizes the Memcached server for caching, enabling scalability and distribution of cache data.

----- Key Components ----
API (api/cache_api.go):

Provides interfaces or endpoints for application components to interact with caching functionalities. This abstraction allows easy integration and decouples caching logic from the rest of the application.

Caches (cache/):

inmemory.go: Implements an in-memory cache with support for LRU eviction. This cache is suitable for small to medium-sized datasets that require fast access times.
lru_cache.go: Implements an LRU cache, which prioritizes recently accessed items for efficient data retrieval and eviction.
memcached.go: Implements caching using Memcached, a distributed memory caching system. It provides fast access to cached data across multiple nodes.

Defines interfaces that specify common cache operations (e.g., Set, Get, Delete, GetAllKeys, DeleteAllKeys). By using interfaces, the project ensures consistency and allows easy swapping of different cache implementations.

Main Application (cmd/main.go): Acts as the entry point for the caching library. It initializes configurations, sets up caches based on application needs, and orchestrates caching operations.

Configuration (config/):

config.go: Contains logic for handling configuration settings, such as cache server addresses, timeouts, and other parameters.
configs/: Directory that stores different configuration files tailored for various deployment environments (e.g., development, production).
Documentation (docs/benchmarks.md):

Testing (tests/):
	benchmarks/: Contains benchmark tests (inmemory_benchmark_test.go) that evaluate the performance of caching operations for the in-memory cache implementation.
	integration/: Includes integration tests (integration_test.go) that validate interactions between the caching library and external dependencies, ensuring correct functionality in real-world scenarios.
	unit/: Houses unit tests (inmemory_test.go, memcached_test.go) that verify the correctness of individual cache functionalities, such as Set, Get, Delete, and error handling.
Utilities (utils/):

	Provides utility functions or helper modules (utils/) that assist in various tasks related to caching operations or other functionalities within the caching library.

----- Benchmark Comparison ------

The following benchmark results compare the performance of the in-memory cache and the Memcached cache. The benchmarks measure the time taken for various operations, with results 	in nanoseconds per operation (ns/op).

In-Memory Cache Performance:

Set Operation: 5728 ns/op
Get Operation: 28.61 ns/op
Delete Operation: 126.9 ns/op
GetAllKeys Operation: 22015 ns/op
DeleteAllKeys Operation: 69.10 ns/op

Memcached Cache Performance:

Set Operation: 182761 ns/op
Get Operation: 197102 ns/op
Delete Operation: 189627 ns/op
GetAllKeys Operation: 23683 ns/op
DeleteAllKeys Operation: 80.90 ns/op

Multicached Cache Performance:

BenchmarkMulticacheSet-8             	    3009	    376100 ns/op
BenchmarkMulticacheGet-8             	 4959327	       201.9 ns/op
BenchmarkMulticacheRemove-8          	 3482068	       309.8 ns/op
BenchmarkMulticacheGetAllKeys-8      	     114	   9808498 ns/op
BenchmarkMulticacheDeleteAllKeys-8   	 4154070	       251.0 ns/op
BenchmarkMulticacheClear-8           	 4178968	       249.1 ns/op
PASS
ok  	command-line-arguments	26.384s


Performance Analysis:
Best Performer: In-Memory Cache

Reasons:

Speed: The in-memory cache outperforms the Memcached cache across all benchmarked operations. 
The most notable difference is in the Get operation, where the in-memory cache is more than 6,800 times faster than Memcached.

Reliability: The in-memory cache successfully completes all operations, whereas the Memcached cache encounters issues with the Delete operation due to cache misses.

Efficiency: With operations consistently taking fewer nanoseconds, the in-memory cache is significantly more efficient for frequent and rapid cache access needs.

Conclusion
Based on the benchmark results, the in-memory cache is the superior performer in terms of speed, reliability, and efficiency. It is the best choice for applications requiring high-speed access and manipulation of cached data. However, the Memcached cache may still be preferred for distributed caching scenarios where cache data needs to be shared across multiple servers
