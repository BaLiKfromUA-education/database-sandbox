# Cassandra basics

My environment:

- Run Cassandra instance [via docker-compose](../../db_environment/cassandra/basics/docker_compose.yaml)
- Execute queries via Cassandra console from my JetBrains IDE

- Queries to insert records and relationships can be found [here](../../db_experiments/cassandra/basics/insert_data.cql)
- Queries to analyze data, according to given task, can be found [here](../../db_experiments/cassandra/basics/analyze_data.cql)

## Items

It is necessary to read fast based on **product category**.

1) Write a query that shows the structure of the created table (DESCRIBE command)
```sql
DESCRIBE TABLE shop.items;
```

```json
[
  {
    "keyspace_name": "shop",
    "type": "table",
    "name": "items",
    "create_statement": "CREATE TABLE shop.items (\n    category text,\n    price int,\n    id int,\n    model text,\n    producer text,\n    properties map<text, text>,\n    PRIMARY KEY (category, price, id)\n) WITH CLUSTERING ORDER BY (price DESC, id ASC)\n    AND additional_write_policy = '99p'\n    AND allow_auto_snapshot = true\n    AND bloom_filter_fp_chance = 0.01\n    AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'}\n    AND cdc = false\n    AND comment = ''\n    AND compaction = {'class': 'org.apache.cassandra.db.compaction.SizeTieredCompactionStrategy', 'max_threshold': '32', 'min_threshold': '4'}\n    AND compression = {'chunk_length_in_kb': '16', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor'}\n    AND memtable = 'default'\n    AND crc_check_chance = 1.0\n    AND default_time_to_live = 0\n    AND extensions = {}\n    AND gc_grace_seconds = 864000\n    AND incremental_backups = true\n    AND max_index_interval = 2048\n    AND memtable_flush_period_in_ms = 0\n    AND min_index_interval = 128\n    AND read_repair = 'BLOCKING'\n    AND speculative_retry = '99p';"
  },
  {
    "keyspace_name": "shop",
    "type": "index",
    "name": "properties_entry_idx",
    "create_statement": "CREATE INDEX properties_entry_idx ON shop.items (entries(properties));"
  },
  {
    "keyspace_name": "shop",
    "type": "index",
    "name": "properties_key_idx",
    "create_statement": "CREATE INDEX properties_key_idx ON shop.items (keys(properties));"
  }
]
```

2) Write a query that displays all items in a certain category sorted by price
```sql
SELECT * FROM shop.items WHERE category = 'Phone';
```

```text
+--------+-----+--+------------------+--------+----------------------+
|category|price|id|model             |producer|properties            |
+--------+-----+--+------------------+--------+----------------------+
|Phone   |1200 |2 |iPhone 15         |Apple   |{'isAndroid': 'false'}|
|Phone   |800  |3 |Samsung Galaxy A22|Samsung |{'isAndroid': 'true'} |
|Phone   |666  |1 |iPhone 6          |Apple   |{'isAndroid': 'false'}|
+--------+-----+--+------------------+--------+----------------------+
```

3) Write queries that select items according to various criteria **within a certain category** (if necessary, use Matirialized view instead of an index):
- model name
```sql
CREATE MATERIALIZED VIEW items_by_model_and_category AS
SELECT *
FROM items
WHERE model IS NOT NULL
AND id IS NOT NULL
AND price IS NOT NULL
AND category IS NOT NULL
PRIMARY KEY ((category, model), price, id);

SELECT * FROM shop.items_by_model_and_category WHERE category='TV' AND model='LG TV 2';
```

```text
+--------+-------+-----+--+--------+------------------------------------------------+
|category|model  |price|id|producer|properties                                      |
+--------+-------+-----+--+--------+------------------------------------------------+
|TV      |LG TV 2|699  |5 |LG      |{'isAndroid': 'true', 'screenResolution': '42"'}|
+--------+-------+-----+--+--------+------------------------------------------------+
```

- price (in range)
```sql
SELECT * FROM shop.items WHERE category='Phone' AND price >= 300 AND price <= 1000;
```

```text
+--------+-----+--+------------------+--------+----------------------+
|category|price|id|model             |producer|properties            |
+--------+-----+--+------------------+--------+----------------------+
|Phone   |800  |3 |Samsung Galaxy A22|Samsung |{'isAndroid': 'true'} |
|Phone   |666  |1 |iPhone 6          |Apple   |{'isAndroid': 'false'}|
+--------+-----+--+------------------+--------+----------------------+
```

- price and producer
```sql
CREATE MATERIALIZED VIEW items_by_producer_and_category AS
SELECT *
FROM items
WHERE producer IS NOT NULL
  AND id IS NOT NULL
  AND price IS NOT NULL
  AND category IS NOT NULL
PRIMARY KEY ((category, producer), price, id);

SELECT * FROM shop.items_by_producer_and_category WHERE category='Phone' AND producer='Apple' AND price >= 300;
```

```text
+--------+--------+-----+--+---------+----------------------+
|category|producer|price|id|model    |properties            |
+--------+--------+-----+--+---------+----------------------+
|Phone   |Apple   |1200 |2 |iPhone 15|{'isAndroid': 'false'}|
|Phone   |Apple   |666  |1 |iPhone 6 |{'isAndroid': 'false'}|
+--------+--------+-----+--+---------+----------------------+
```

4) Write queries that select items based on:
- existence of a property
```sql
SELECT * FROM shop.items WHERE properties CONTAINS KEY 'isAndroid';
```

```text
+--------+-----+--+-------------------------+--------+------------------------------------------------+
|category|price|id|model                    |producer|properties                                      |
+--------+-----+--+-------------------------+--------+------------------------------------------------+
|Phone   |1200 |2 |iPhone 15                |Apple   |{'isAndroid': 'false'}                          |
|Phone   |800  |3 |Samsung Galaxy A22       |Samsung |{'isAndroid': 'true'}                           |
|Phone   |666  |1 |iPhone 6                 |Apple   |{'isAndroid': 'false'}                          |
|TV      |700  |4 |Samsung Smart TV model 22|Samsung |{'isAndroid': 'true', 'screenResolution': '45"'}|
|TV      |699  |5 |LG TV 2                  |LG      |{'isAndroid': 'true', 'screenResolution': '42"'}|
+--------+-----+--+-------------------------+--------+------------------------------------------------+
```

- property value
```sql
SELECT * FROM shop.items WHERE properties['isAndroid']='true';
```

```text
+--------+-----+--+-------------------------+--------+------------------------------------------------+
|category|price|id|model                    |producer|properties                                      |
+--------+-----+--+-------------------------+--------+------------------------------------------------+
|Phone   |800  |3 |Samsung Galaxy A22       |Samsung |{'isAndroid': 'true'}                           |
|TV      |700  |4 |Samsung Smart TV model 22|Samsung |{'isAndroid': 'true', 'screenResolution': '45"'}|
|TV      |699  |5 |LG TV 2                  |LG      |{'isAndroid': 'true', 'screenResolution': '42"'}|
+--------+-----+--+-------------------------+--------+------------------------------------------------+
```

5) Update properties of item:
- update existing property
```sql
UPDATE shop.items SET properties['screenResolution'] = 'UNKNOWN' WHERE id=5 and category='TV' and price=699;
SELECT * FROM shop.items WHERE properties['screenResolution'] = 'UNKNOWN';
```

```text
+--------+-----+--+-------+--------+----------------------------------------------------+
|category|price|id|model  |producer|properties                                          |
+--------+-----+--+-------+--------+----------------------------------------------------+
|TV      |699  |5 |LG TV 2|LG      |{'isAndroid': 'true', 'screenResolution': 'UNKNOWN'}|
+--------+-----+--+-------+--------+----------------------------------------------------+
```

- add new property
```sql
UPDATE shop.items SET properties = properties + {'onSale':'true'} WHERE id=5 and category='TV' and price=699;
SELECT * FROM shop.items WHERE properties CONTAINS KEY 'onSale';
```

```text
+--------+-----+--+-------+--------+----------------------------------------------------------------------+
|category|price|id|model  |producer|properties                                                            |
+--------+-----+--+-------+--------+----------------------------------------------------------------------+
|TV      |699  |5 |LG TV 2|LG      |{'isAndroid': 'true', 'onSale': 'true', 'screenResolution': 'UNKNOWN'}|
+--------+-----+--+-------+--------+----------------------------------------------------------------------+

```

- delete property
```sql
UPDATE shop.items SET properties = properties - {'onSale'} WHERE id=5 and category='TV' and price=699;
SELECT * FROM shop.items WHERE properties CONTAINS KEY 'onSale';
```

```text
shop> UPDATE shop.items SET properties = properties - {'onSale'} WHERE id=5 and category='TV' and price=699
[2024-03-10 18:59:06] completed in 9 ms
shop> SELECT * FROM shop.items WHERE properties CONTAINS KEY 'onSale'
[2024-03-10 18:59:06] 0 rows retrieved in 28 ms (execution: 8 ms, fetching: 20 ms)
```

## Orders

