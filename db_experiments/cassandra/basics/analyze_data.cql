// items
DESCRIBE TABLE shop.items;
SELECT * FROM shop.items WHERE category = 'Phone';

SELECT * FROM shop.items_by_model_and_category WHERE category='TV' AND model='LG TV 2';
SELECT * FROM shop.items WHERE category='Phone' AND price >= 300 AND price <= 1000;
SELECT * FROM shop.items_by_producer_and_category WHERE category='Phone' AND producer='Apple' AND price >= 300;

SELECT * FROM shop.items WHERE properties CONTAINS KEY 'isAndroid';
SELECT * FROM shop.items WHERE properties['isAndroid']='true';

UPDATE shop.items SET properties['screenResolution'] = 'UNKNOWN' WHERE id=5 and category='TV' and price=699;
SELECT * FROM shop.items WHERE properties['screenResolution'] = 'UNKNOWN';

UPDATE shop.items SET properties = properties + {'onSale':'true'} WHERE id=5 and category='TV' and price=699;
SELECT * FROM shop.items WHERE properties CONTAINS KEY 'onSale';

UPDATE shop.items SET properties = properties - {'onSale'} WHERE id=5 and category='TV' and price=699;
SELECT * FROM shop.items WHERE properties CONTAINS KEY 'onSale';

// orders
DESCRIBE TABLE shop.orders;

SELECT * FROM shop.orders WHERE customer_id='Olga';

SELECT * FROM shop.orders WHERE customer_id='Olga' AND item_ids CONTAINS 4;

SELECT customer_id, count(*) FROM shop.orders WHERE customer_id='Olga' AND order_date > '2024-02-11 03:03:40.000' AND order_date < '2024-03-10 03:03:40.000';

SELECT customer_id, sum(price) FROM shop.orders GROUP BY customer_id;

SELECT customer_id, order_date, item_ids, max(price) FROM shop.orders GROUP BY customer_id;

UPDATE shop.orders SET item_ids = item_ids + [1], price = 2566 WHERE customer_id='Olga' AND order_date='2024-03-11 03:03:40+0000';
SELECT customer_id, order_date, item_ids, max(price) FROM shop.orders WHERE customer_id='Olga';

SELECT customer_id, WRITETIME(price) FROM shop.orders;

INSERT INTO shop.orders (customer_id, order_date, item_ids, price) VALUES ('Andrii', '2024-03-11 03:03:40+0000', [2, 4], 1900) USING TTL 10;
SELECT * FROM shop.orders WHERE customer_id='Andrii';

INSERT INTO shop.orders JSON '{"customer_id": "Andrii", "order_date": "2024-03-01 13:01:00.000Z", "item_ids": [2,4], "price": 1900}';
SELECT JSON * FROM shop.orders WHERE customer_id='Andrii';
