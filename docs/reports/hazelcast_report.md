# Task 2 -- Counter Implementation using Hazelcast

- `docker-compose` setup with **3** `hazelcast` nodes as a separate
  containers. [Link](../../db_environment/hazelcast/docker_compose.yaml)
- DAO is implemented in Go. [Link](../../db_experiments/hazelcast/counter.go). This scripts implement:
    - Counter via Distributed Map **without blocking**. [Based on example](https://docs.hazelcast.com/imdg/latest/data-structures/map#locking-maps)
    - Counter via Distributed Map **with pessimistic blocking**. [Based on example](https://docs.hazelcast.com/imdg/latest/data-structures/map#pessimistic-locking)
    - Counter via Distributed Map **with optimistic blocking**. [Based on example](https://docs.hazelcast.com/imdg/latest/data-structures/map#optimistic-locking)
    - Counter via **Atomic Long**. [Based on example](https://docs.hazelcast.com/hazelcast/5.1/data-structures/iatomiclong)
- Tests for correctness and benchmarks are implemented with built-in Go
  tools. [Link](../../db_experiments/hazelcast/counter_test.go)
    - Raw results of final test session can be found [here](../raw_data/hazelcast_test_without_lock_case.log).

## Tests results

| Update Strategy               | Final value is always 100k | Number of concurrent clients (co-routines) | Operation per client | CPU                | Number of CPU threads | Time of execution |
|-------------------------------|----------------------------|--------------------------------------------|----------------------|--------------------|-----------------------|-------------------|
| Map without blocking          | ❌                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 5,74 sec          |
| Map with pessimistic blocking | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | [46,83 min](https://github.com/hazelcast/hazelcast-go-client/issues/992)             |
| Map with optimistic blocking  | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 26,06 sec         |
| Atomic Long                   | ✅                          | 10                                         | 10 000               | AMD Ryzen 9 5900HS | 16                    | 9,076 sec         |


### Proofs that CP subsystem is enabled

Can be found in [logs](../raw_data/hazelcast_test_without_lock_case.log):

```
hazelcast_2_1        | CP Group Members {groupId: default(2082), size:3, term:1, logIndex:0} [
hazelcast_2_1        |  CPMember{uuid=0ee7b71a-2639-49ec-b77d-4a2383e73273, address=[hazelcast_1]:5701} - LEADER
hazelcast_2_1        |  CPMember{uuid=8514f6c9-5e30-4d0b-9b02-aefa00b21b81, address=[hazelcast_3]:5701}
hazelcast_2_1        |  CPMember{uuid=95e22f28-499f-4c04-b5a0-dc4c2f76d461, address=[hazelcast_2]:5701} - FOLLOWER this
hazelcast_2_1        | ]
```
