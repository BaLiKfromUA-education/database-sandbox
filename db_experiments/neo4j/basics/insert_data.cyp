MATCH(shop) DETACH DELETE shop; // DROP * FROM shop_db

// Add items
CREATE (shop: Item {id: 1, name: "iPhone 6", price: 666});
CREATE (shop: Item {id: 2, name: "iPhone 15", price: 1200});
CREATE (shop: Item {id: 3, name: "Fishing rod", price: 100});
CREATE (shop: Item {id: 4, name: "Samsung Smart TV model 42", price: 800});
CREATE (shop: Item {id: 5, name: "Casio G-Shock 42", price: 300});


// Add customers
CREATE (shop: Customer {id: 101, name: "Valentyn"});
CREATE (shop: Customer {id: 102, name: "Olga"});
CREATE (shop: Customer {id: 103, name: "Andrii"});
CREATE (shop: Customer {id: 104, name: "Mykyta"});
CREATE (shop: Customer {id: 105, name: "Ivan"});

// Add order
CREATE (shop: Order {id: 1001, date: date("2015-04-14")});
CREATE (shop: Order {id: 1002, date: date("2015-04-15")});
CREATE (shop: Order {id: 1003, date: date("2015-04-16")});
CREATE (shop: Order {id: 1004, date: date("2015-04-17")});
CREATE (shop: Order {id: 1005, date: date("2015-04-18")});
CREATE (shop: Order {id: 1006, date: date("2015-04-18")});

// BOUGHT relationships (CUSTOMER -> ORDER)
// Valentyn has 2 orders
MATCH (customer:Customer {id: 101}), (order:Order {id: 1001})
CREATE (customer)-[:BOUGHT]->(order);
MATCH (customer:Customer {id: 101}), (order:Order {id: 1002})
CREATE (customer)-[:BOUGHT]->(order);
// Olga has 1 order
MATCH (customer:Customer {id: 102}), (order:Order {id: 1003})
CREATE (customer)-[:BOUGHT]->(order);
// Andrii has 3 orders
MATCH (customer:Customer {id: 103}), (order:Order {id: 1004})
CREATE (customer)-[:BOUGHT]->(order);
MATCH (customer:Customer {id: 103}), (order:Order {id: 1005})
CREATE (customer)-[:BOUGHT]->(order);
MATCH (customer:Customer {id: 103}), (order:Order {id: 1006})
CREATE (customer)-[:BOUGHT]->(order);
// Mykyta and Ivan have no orders:(

// CONTAINS relationship (ORDER -> ITEM)
// First order contains 3 items
MATCH (order:Order {id: 1001}), (item:Item {id: 1})
CREATE (order)-[:CONTAINS]->(item);
MATCH (order:Order {id: 1001}), (item:Item {id: 2})
CREATE (order)-[:CONTAINS]->(item);
MATCH (order:Order {id: 1001}), (item:Item {id: 3})
CREATE (order)-[:CONTAINS]->(item);
// Second order is the same as first by content
MATCH (order:Order {id: 1002}), (item:Item {id: 1})
CREATE (order)-[:CONTAINS]->(item);
MATCH (order:Order {id: 1002}), (item:Item {id: 2})
CREATE (order)-[:CONTAINS]->(item);
MATCH (order:Order {id: 1002}), (item:Item {id: 3})
CREATE (order)-[:CONTAINS]->(item);
// Third order contains 2 items
MATCH (order:Order {id: 1003}), (item:Item {id: 3})
CREATE (order)-[:CONTAINS]->(item);
MATCH (order:Order {id: 1003}), (item:Item {id: 4})
CREATE (order)-[:CONTAINS]->(item);
// Fourth order contains 1 item
MATCH (order:Order {id: 1004}), (item:Item {id: 5})
CREATE (order)-[:CONTAINS]->(item);
// Fifth order contains 1 item
MATCH (order:Order {id: 1005}), (item:Item {id: 3})
CREATE (order)-[:CONTAINS]->(item);
// Sixth order contains 2 items
MATCH (order:Order {id: 1006}), (item:Item {id: 1})
CREATE (order)-[:CONTAINS]->(item);
MATCH (order:Order {id: 1006}), (item:Item {id: 2})
CREATE (order)-[:CONTAINS]->(item);

// VIEW relationship (CUSTOMER -> ITEM)
// All users viewed each item in our shop
MATCH (customer:Customer), (item:Item)
CREATE (customer)-[:VIEWED]->(item);