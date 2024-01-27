# Task 1 -- Counter Implementation using PostgreSQL

## Details:

- `docker-compose` setup with `postgresql` in a separate
  container. [Link](../../db_environment/postgresql/docker_compose.yaml)
- DAO is implemented in Go. [Link](../../db_experiments/postgresql/counter.go). This scripts implement:
    - Lost-update
    - In-place update
    - Row-level locking
    - Optimistic concurrency control
- Tests for correctness and benchmarks are implemented with built-in Go
  tools. [Link](../../db_experiments/postgresql/counter_test.go)
    - Raw results of final test session can be found [here](../raw_data/postgres_test.log).

## Benchmark results

| Update Strategy                | Final value is always 100k | Number of concurrent clients (co-routines) | Operation per client | CPU                | Number of CPU threads | Time of execution |
|--------------------------------|----------------------------|--------------------------------------------|----------------------|--------------------|-----------------------|-------------------|
| Lost-update                    | ❌                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 41,49 sec         |
| In-place update                | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 41,64 sec         |
| Row-level locking              | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 52,59 sec         |
| Optimistic concurrency control | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 360,97 sec        |