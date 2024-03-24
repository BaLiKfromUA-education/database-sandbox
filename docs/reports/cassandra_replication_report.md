# Cassandra replication

1. Configure cassandra cluster with 3
   nodes --> [docker-compose configuration](../../db_environment/cassandra/replication/docker_compose.yaml)
2. Validate configuration via `nodetool status`

```shell
docker exec cassandra-1 nodetool status
Datacenter: my-datacenter-1
===========================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address     Load       Tokens  Owns (effective)  Host ID                               Rack 
UN  172.29.0.3  80.06 KiB  16      59.3%             b0b2d43d-021c-4f7d-ada9-8ceb72044e15  rack1
UN  172.29.0.4  80.06 KiB  16      76.0%             53dfef6a-82f8-418e-9386-559ab777fc15  rack1
UN  172.29.0.2  119.8 KiB  16      64.7%             40561bf6-89bd-4182-91a6-b7962e8698f9  rack1
```

3. Create 3 keyspaces with replication factor `1`, `2`, and `3`.

```sql
CREATE
KEYSPACE IF NOT EXISTS namespace_1
WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : 1
};

CREATE
KEYSPACE IF NOT EXISTS namespace_2
WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : 2
};

CREATE
KEYSPACE IF NOT EXISTS namespace_3
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

4. Create table in each keyspace

```sql
CREATE TABLE namespace_1.items
(
    id       int,
    category text,
    model    text,
    price    int,
    PRIMARY KEY ((category),
    price,
    id
)
    );

CREATE TABLE namespace_2.items
(
    id       int,
    category text,
    model    text,
    price    int,
    PRIMARY KEY ((category),
    price,
    id
)
    );

CREATE TABLE namespace_3.items
(
    id       int,
    category text,
    model    text,
    price    int,
    PRIMARY KEY ((category),
    price,
    id
)
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
SELECT *
FROM namespace_1.items
WHERE category = 'Phone';
SELECT *
FROM namespace_2.items
WHERE category = 'Phone';
SELECT *
FROM namespace_3.items
WHERE category = 'Phone';
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

6. Check distribution of data across all nodes for each keyspace

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

7. For each record from each keyspace, print nodes where record is stored

- `namespace-1`

```shell
docker exec cassandra-1 nodetool getendpoints -- namespace_1 items 1
172.29.0.3

docker exec cassandra-1 nodetool getendpoints -- namespace_1 items 2
172.29.0.3

docker exec cassandra-1 nodetool getendpoints -- namespace_1 items 3
172.29.0.3
```

- `namespace-2`

```shell
docker exec cassandra-1 nodetool getendpoints -- namespace_2 items 1
172.29.0.3
172.29.0.2

docker exec cassandra-1 nodetool getendpoints -- namespace_2 items 2
172.29.0.3
172.29.0.4

docker exec cassandra-1 nodetool getendpoints -- namespace_2 items 3
172.29.0.3
172.29.0.2
```

- `namespace-3`

```shell
docker exec cassandra-1 nodetool getendpoints -- namespace_3 items 1
172.29.0.3
172.29.0.2
172.29.0.4
```

```shell
docker exec cassandra-1 nodetool getendpoints -- namespace_3 items 2
172.29.0.3
172.29.0.4
172.29.0.2
```

```shell
docker exec cassandra-1 nodetool getendpoints -- namespace_3 items 3
172.29.0.3
172.29.0.2
172.29.0.4
```

8. Stop one of the nodes. For each keyspace, find out with which consistency level we can write or read,
   and which of them guarantee strong consistency

```shell
docker-compose -f docker_compose.yaml ps
   Name                  Command                  State                                                            Ports                                                      
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
cassandra-1   docker-entrypoint.sh cassa ...   Up (healthy)   0.0.0.0:7000->7000/tcp,:::7000->7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9042->9042/tcp,:::9042->9042/tcp, 9160/tcp
cassandra-2   docker-entrypoint.sh cassa ...   Up (healthy)   7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9043->9042/tcp,:::9043->9042/tcp, 9160/tcp                                
cassandra-3   docker-entrypoint.sh cassa ...   Up (healthy)   7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9044->9042/tcp,:::9044->9042/tcp, 9160/tcp                                

docker stop cassandra-3
cassandra-3

docker-compose -f docker_compose.yaml ps
   Name                  Command                  State                                                            Ports                                                      
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
cassandra-1   docker-entrypoint.sh cassa ...   Up (healthy)   0.0.0.0:7000->7000/tcp,:::7000->7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9042->9042/tcp,:::9042->9042/tcp, 9160/tcp
cassandra-2   docker-entrypoint.sh cassa ...   Up (healthy)   7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9043->9042/tcp,:::9043->9042/tcp, 9160/tcp                                
cassandra-3   docker-entrypoint.sh cassa ...   Exit 143   
```

```sql
CONSISTENCY <LEVEL>;
INSERT INTO namespace_1.items (id, category, model, price)
VALUES (404, 'TV', 'LG Model 42', 555);
INSERT INTO namespace_2.items (id, category, model, price)
VALUES (404, 'TV', 'LG Model 42', 555);
INSERT INTO namespace_3.items (id, category, model, price)
VALUES (404, 'TV', 'LG Model 42', 555);

SELECT * FROM namespace_1.items WHERE category = 'TV' AND price=555;
SELECT * FROM namespace_2.items WHERE category = 'TV' AND price=555;
SELECT * FROM namespace_3.items WHERE category = 'TV' AND price=555;
```

| Consistency Level | Write `namespace_1` | Write `namespace_2` | Write `namespace_3` | Read `namespace_1` | Read `namespace_2` | Read `namespace_3` | Strong consistency                |
|-------------------|---------------------|---------------------|---------------------|--------------------|--------------------|--------------------|-----------------------------------|
| `ANY`             | ✅                   | ✅                   | ✅                   | ❌                  | ❌                  | ❌                  | ❌                                 |
| `ONE`             | ✅                   | ✅                   | ✅                   | ✅                  | ✅                  | ✅                  | ❌                                 |
| `TWO`             | ❌                   | ❌                   | ✅                   | ❌                  | ❌                  | ✅                  | ❌                                 |
| `THREE`           | ❌                   | ❌                   | ❌                   | ❌                  | ❌                  | ❌                  | ❌                                 |
| `QUORUM`          | ✅                   | ❌                   | ✅                   | ✅                  | ❌                  | ✅                  | ✅                                 |
| `ALL`             | ✅                   | ❌                   | ❌                   | ✅                  | ❌                  | ❌                  | ✅                                 |
| `LOCAL_QUORUM`    | ✅                   | ❌                   | ✅                   | ✅                  | ❌                  | ✅                  | ✅                                 |
| `EACH_QUORUM`     | ✅                   | ❌                   | ✅                   | ✅                  | ❌                  | ✅                  | ✅                                 |
| `SERIAL`          | ❌                   | ❌                   | ❌                   | ✅                  | ❌                  | ✅                  | ✅ (lightweight transactions only) |
| `LOCAL_SERIAL`    | ❌                   | ❌                   | ❌                   | ✅                  | ❌                  | ✅                  | ✅ (lightweight transactions only) |

9. Disable connection between all nodes
```shell
docker-compose -f docker_compose.yaml ps
   Name                  Command                  State                                                            Ports                                                      
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
cassandra-1   docker-entrypoint.sh cassa ...   Up (healthy)   0.0.0.0:7000->7000/tcp,:::7000->7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9042->9042/tcp,:::9042->9042/tcp, 9160/tcp
cassandra-2   docker-entrypoint.sh cassa ...   Up (healthy)   7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9043->9042/tcp,:::9043->9042/tcp, 9160/tcp                                
cassandra-3   docker-entrypoint.sh cassa ...   Up (healthy)   7000/tcp, 7001/tcp, 7199/tcp, 0.0.0.0:9044->9042/tcp,:::9044->9042/tcp, 9160/tcp   

docker network disconnect replication_cassandra-net cassandra-1
docker network disconnect replication_cassandra-net cassandra-2
docker network disconnect replication_cassandra-net cassandra-3
```

10. For a keyspace with replication factor 3, set the consistency level to 1. Write the same value, with the same primary key, but different other values to each of the nodes ( create a conflict)
```shell
docker exec -it cassandra-1 cqlsh
WARNING: cqlsh was built against 5.0-beta1, but this server is 5.0.  All features may not work!
Connected to my-cluster at 127.0.0.1:9042
[cqlsh 6.2.0 | Cassandra 5.0-beta1 | CQL spec 3.4.7 | Native protocol v5]
Use HELP for help.
cqlsh> CONSISTENCY ONE;
Consistency level set to ONE.
cqlsh> INSERT INTO namespace_3.items (id, category, model, price) VALUES (42, 'TV', 'LG Model node 1', 555);
cqlsh> exit

docker exec -it cassandra-2 cqlsh
WARNING: cqlsh was built against 5.0-beta1, but this server is 5.0.  All features may not work!
Connected to my-cluster at 127.0.0.1:9042
[cqlsh 6.2.0 | Cassandra 5.0-beta1 | CQL spec 3.4.7 | Native protocol v5]
Use HELP for help.
cqlsh> CONSISTENCY ONE;
Consistency level set to ONE.
cqlsh>  INSERT INTO namespace_3.items (id, category, model, price) VALUES (42, 'TV', 'LG Model node 2', 555);
cqlsh> exit

docker exec -it cassandra-3 cqlsh
WARNING: cqlsh was built against 5.0-beta1, but this server is 5.0.  All features may not work!
Connected to my-cluster at 127.0.0.1:9042
[cqlsh 6.2.0 | Cassandra 5.0-beta1 | CQL spec 3.4.7 | Native protocol v5]
Use HELP for help.
cqlsh> CONSISTENCY ONE;
Consistency level set to ONE.
cqlsh> INSERT INTO namespace_3.items (id, category, model, price) VALUES (42, 'TV', 'LG Model node 3', 555);
cqlsh> exit
```

11. Combine the nodes into a cluster and determine which value was accepted by the cluster and according to which principle -- **LAST WRITE WINS**

```shell
docker network connect replication_cassandra-net cassandra-1
docker network connect replication_cassandra-net cassandra-2
docker network connect replication_cassandra-net cassandra-3
```

```sql
CONSISTENCY ALL;
SELECT * FROM namespace_3.items WHERE category='TV' and id=42 and price=555;
```

```text
+--------+-----+--+---------------+
|category|price|id|model          |
+--------+-----+--+---------------+
|TV      |555  |42|LG Model node 3|
+--------+-----+--+---------------+
```

12. Check the behavior of lightweight transactions for previous points in a split and non-split cluster
- split state
```shell
docker network disconnect replication_cassandra-net cassandra-1
docker network disconnect replication_cassandra-net cassandra-2
docker network disconnect replication_cassandra-net cassandra-3

docker exec -it cassandra-1 cqlsh
WARNING: cqlsh was built against 5.0-beta1, but this server is 5.0.  All features may not work!
Connected to my-cluster at 127.0.0.1:9042
[cqlsh 6.2.0 | Cassandra 5.0-beta1 | CQL spec 3.4.7 | Native protocol v5]
Use HELP for help.
cqlsh> CONSISTENCY ONE;
Consistency level set to ONE.
cqlsh> INSERT INTO namespace_3.items (id, category, model, price) VALUES (777, 'node', ' node 1', 555) IF NOT EXISTS;
NoHostAvailable: ('Unable to complete the operation against any hosts', {<Host: 127.0.0.1:9042 my-datacenter-1>: Unavailable('Error from server: code=1000 [Unavailable exception] message="Cannot achieve consistency level SERIAL" info={\'consistency\': \'SERIAL\', \'required_replicas\': 2, \'alive_replicas\': 1}')})
```

- non-split state
```shell
docker network connect replication_cassandra-net cassandra-1
docker network connect replication_cassandra-net cassandra-2
docker network connect replication_cassandra-net cassandra-3
```

```sql
CONSISTENCY ONE;

INSERT INTO namespace_3.items (id, category, model, price) VALUES (777, 'node', 'try 1', 555) IF NOT EXISTS;
INSERT INTO namespace_3.items (id, category, model, price) VALUES (777, 'node', 'try 2', 555) IF NOT EXISTS;
INSERT INTO namespace_3.items (id, category, model, price) VALUES (777, 'node', 'try 3', 555) IF NOT EXISTS;

SELECT * FROM namespace_3.items WHERE category='node' AND id=777 and price=555;
```

```text
+--------+-----+---+-----+
|category|price|id |model|
+--------+-----+---+-----+
|node    |555  |777|try 1|
+--------+-----+---+-----+

```
