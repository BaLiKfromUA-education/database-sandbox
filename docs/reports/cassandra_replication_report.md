1. Configure cassandra cluster with 3
   nodes --> [docker-compose configuration](../../db_environment/cassandra/replication/docker_compose.yaml)
2. Validate configuration via `nodetool status`
```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/cassandra/replication$ docker exec cassandra-1 nodetool status
Datacenter: my-datacenter-1
===========================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load       Tokens  Owns (effective)  Host ID                               Rack 
UN  172.29.0.3  80.06 KiB  16      59.3%             b0b2d43d-021c-4f7d-ada9-8ceb72044e15  rack1
UN  172.29.0.4  80.06 KiB  16      76.0%             53dfef6a-82f8-418e-9386-559ab777fc15  rack1
UN  172.29.0.2  119.8 KiB  16      64.7%             40561bf6-89bd-4182-91a6-b7962e8698f9  rack1
```

3. Create 3 namespaces with replication factor `1`, `2`, and `3`.
```sql
CREATE KEYSPACE IF NOT EXISTS namespace_1
    WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : 1
};

CREATE KEYSPACE IF NOT EXISTS namespace_2
    WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : 2
};

CREATE KEYSPACE IF NOT EXISTS namespace_3
    WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : 3
};
```

```shell
[2024-03-23 13:36:54] completed in 878 ms
[2024-03-23 13:37:14] completed in 886 ms
[2024-03-23 13:37:30] completed in 650 ms
```

4. Create table in each namespace
```sql
CREATE TABLE namespace_1.items
(
    id         int,
    category   text,
    model      text,
    price      int,
    PRIMARY KEY ((category), price, id)
);

CREATE TABLE namespace_2.items
(
    id         int,
    category   text,
    model      text,
    price      int,
    PRIMARY KEY ((category), price, id)
);

CREATE TABLE namespace_3.items
(
    id         int,
    category   text,
    model      text,
    price      int,
    PRIMARY KEY ((category), price, id)
);
```

```shell
[2024-03-23 13:52:03] completed in 1 s 523 ms
[2024-03-23 13:52:12] completed in 891 ms
[2024-03-23 13:52:21] completed in 509 ms
```

5. Try to write/read from different nodes

- `node 1` writes
```sql
INSERT INTO namespace_1.items (id, category, model, price)
VALUES (1, 'Phone', 'iPhone 6', 666);
INSERT INTO namespace_2.items (id, category, model, price)
VALUES (1, 'Phone', 'iPhone 6', 666);
INSERT INTO namespace_3.items (id, category, model, price)
VALUES (1, 'Phone', 'iPhone 6', 666);
```

```shell
[2024-03-23 14:51:57] completed in 17 ms
[2024-03-23 14:51:57] completed in 15 ms
[2024-03-23 14:51:57] completed in 9 ms
```

- `node 2` writes
```sql
INSERT INTO namespace_1.items (id, category, model, price)
VALUES (2, 'Phone', 'iPhone 15', 1200);
INSERT INTO namespace_2.items (id, category, model, price)
VALUES (2, 'Phone', 'iPhone 15', 1200);
INSERT INTO namespace_3.items (id, category, model, price)
VALUES (2, 'Phone', 'iPhone 15', 1200);
```

```shell
[2024-03-23 14:54:37] completed in 5 ms
[2024-03-23 14:54:37] completed in 5 ms
[2024-03-23 14:54:37] completed in 4 ms
```

- `node 3` writes
```sql
INSERT INTO namespace_1.items (id, category, model, price)
VALUES (3, 'Phone', 'Samsung Galaxy A22', 900);
INSERT INTO namespace_2.items (id, category, model, price)
VALUES (3, 'Phone', 'Samsung Galaxy A22', 900);
INSERT INTO namespace_3.items (id, category, model, price)
VALUES (3, 'Phone', 'Samsung Galaxy A22', 900);
```

```shell
[2024-03-23 21:05:06] completed in 12 ms
[2024-03-23 21:05:07] completed in 18 ms
[2024-03-23 21:05:07] completed in 7 ms
```

- `node 1`/`node 2`/`node 3` reads (results are the same)

```sql
SELECT * FROM namespace_1.items WHERE category = 'Phone';
SELECT * FROM namespace_2.items WHERE category = 'Phone';
SELECT * FROM namespace_3.items WHERE category = 'Phone';
```

```shell
> SELECT * FROM namespace_1.items WHERE category = 'Phone'
[2024-03-23 21:07:39] 3 rows retrieved starting from 1 in 78 ms (execution: 17 ms, fetching: 61 ms)
> SELECT * FROM namespace_2.items WHERE category = 'Phone'
[2024-03-23 21:08:19] 3 rows retrieved starting from 1 in 32 ms (execution: 10 ms, fetching: 22 ms)
> SELECT * FROM namespace_3.items WHERE category = 'Phone'
[2024-03-23 21:08:26] 3 rows retrieved starting from 1 in 49 ms (execution: 28 ms, fetching: 21 ms)
```

```text
+--------+-----+--+------------------+
|category|price|id|model             |
+--------+-----+--+------------------+
|Phone   |666  |1 |iPhone 6          |
|Phone   |900  |3 |Samsung Galaxy A22|
|Phone   |1200 |2 |iPhone 15         |
+--------+-----+--+------------------+
```

6. Check distribution of data across all nodes for each namespace

```shell
docker exec cassandra-1 nodetool status namespace_1
Datacenter: my-datacenter-1
===========================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load        Tokens  Owns (effective)  Host ID                               Rack 
UN  172.29.0.3  152.19 KiB  16      35.7%             53dfef6a-82f8-418e-9386-559ab777fc15  rack1
UN  172.29.0.4  191.64 KiB  16      32.7%             40561bf6-89bd-4182-91a6-b7962e8698f9  rack1
UN  172.29.0.2  167.79 KiB  16      31.6%             b0b2d43d-021c-4f7d-ada9-8ceb72044e15  rack1

```

```shell
docker exec cassandra-1 nodetool status namespace_2
Datacenter: my-datacenter-1
===========================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load        Tokens  Owns (effective)  Host ID                               Rack 
UN  172.29.0.3  152.19 KiB  16      76.0%             53dfef6a-82f8-418e-9386-559ab777fc15  rack1
UN  172.29.0.4  191.64 KiB  16      64.7%             40561bf6-89bd-4182-91a6-b7962e8698f9  rack1
UN  172.29.0.2  167.79 KiB  16      59.3%             b0b2d43d-021c-4f7d-ada9-8ceb72044e15  rack1


```

```shell
docker exec cassandra-1 nodetool status namespace_3
Datacenter: my-datacenter-1
===========================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load        Tokens  Owns (effective)  Host ID                               Rack 
UN  172.29.0.3  152.19 KiB  16      100.0%            53dfef6a-82f8-418e-9386-559ab777fc15  rack1
UN  172.29.0.4  191.64 KiB  16      100.0%            40561bf6-89bd-4182-91a6-b7962e8698f9  rack1
UN  172.29.0.2  167.79 KiB  16      100.0%            b0b2d43d-021c-4f7d-ada9-8ceb72044e15  rack1

```

7. For each record from each namespace, print nodes where record is stored

** This lab will be done in couple of hours**