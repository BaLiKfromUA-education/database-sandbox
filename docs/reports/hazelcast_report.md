# Task 2 -- Counter Implementation using Hazelcast

- `docker-compose` setup with `postgresql` in a separate
  container. [Link](../../db_environment/hazelcast/docker_compose.yaml)
- DAO is implemented in Go. [Link](../../db_experiments/hazelcast/counter.go). This scripts implement:
    - Counter via Distributed Map **without blocking
      **. [Based on example](https://docs.hazelcast.com/imdg/latest/data-structures/map#locking-maps)
    - Counter via Distributed Map **with pessimistic blocking
      **. [Based on example](https://docs.hazelcast.com/imdg/latest/data-structures/map#pessimistic-locking)
    - Counter via Distributed Map **with optimistic blocking
      **. [Based on example](https://docs.hazelcast.com/imdg/latest/data-structures/map#optimistic-locking)
    - Counter via **Atomic Long
      **. [Based on example](https://docs.hazelcast.com/hazelcast/5.1/data-structures/iatomiclong)
- Tests for correctness and benchmarks are implemented with built-in Go
  tools. [Link](../../db_experiments/hazelcast/counter_test.go)
    - Raw results of final test session can be found [here](../raw_data/hazelcast_test_without_lock_case.log).

## Tests results

| Update Strategy               | Final value is always 100k | Number of concurrent clients (co-routines) | Operation per client | CPU                | Number of CPU threads | Time of execution |
|-------------------------------|----------------------------|--------------------------------------------|----------------------|--------------------|-----------------------|-------------------|
| Map without blocking          | ❌                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 5,74 sec          |
| Map with pessimistic blocking | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | TBD*              |
| Map with optimistic blocking  | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 26,06 sec         |
| Atomic Long                   | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 9,076 sec         |
