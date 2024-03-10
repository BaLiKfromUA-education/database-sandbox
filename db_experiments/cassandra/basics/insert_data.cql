// https://docs.datastax.com/en/cql/dse/docs/developing/keyspace-create.html#simple-keyspace
CREATE KEYSPACE IF NOT EXISTS shop
WITH REPLICATION = {
'class':'SimpleStrategy',
'replication_factor':1
};

USE shop;

// items table
CREATE TABLE shop.items
(
id         int,
category   text,
model      text,
producer   text,
price      int,
properties map<text,text>,
PRIMARY KEY ((category), price, id)
)
WITH CLUSTERING ORDER BY (price DESC);

CREATE INDEX IF NOT EXISTS properties_idx
ON shop.items (KEYS(properties));

INSERT INTO shop.items (id, category, model, producer, price, properties) VALUES (1, 'Phone', 'iPhone 6','Apple', 666, {'isAndroid': 'false'});
INSERT INTO shop.items (id, category, model, producer, price, properties) VALUES (2, 'Phone', 'iPhone 15','Apple', 1200, {'isAndroid': 'false'});
INSERT INTO shop.items (id, category, model, producer, price, properties) VALUES (3, 'Phone', 'Samsung Galaxy A22','Samsung', 800, {'isAndroid': 'true'});
INSERT INTO shop.items (id, category, model, producer, price, properties) VALUES (4, 'TV', 'Samsung Smart TV model 22','Samsung', 700, {'isAndroid': 'true', 'screenResolution': '45"'});
INSERT INTO shop.items (id, category, model, producer, price, properties) VALUES (5, 'TV', 'LG TV 2','LG', 699, {'isAndroid': 'true', 'screenResolution': '42"'});
INSERT INTO shop.items (id, category, model, producer, price, properties) VALUES (6, 'Smart Watch', 'Casio G-Shock 42', 'Casio', 300, {});

CREATE MATERIALIZED VIEW items_by_model_and_category AS
SELECT *
FROM items
WHERE model IS NOT NULL
AND id IS NOT NULL
AND price IS NOT NULL
AND category IS NOT NULL
PRIMARY KEY ((category, model), price, id);

CREATE MATERIALIZED VIEW items_by_producer_and_category AS
SELECT *
FROM items
WHERE producer IS NOT NULL
AND id IS NOT NULL
AND price IS NOT NULL
AND category IS NOT NULL
PRIMARY KEY ((category, producer), price, id);

// orders