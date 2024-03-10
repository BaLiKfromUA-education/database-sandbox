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