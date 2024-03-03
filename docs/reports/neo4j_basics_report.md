# neo4j basics

1. Find all items in an order by order ID
```cypher
MATCH (order:Order {id: 1001})-[:CONTAINS]->(item:Item)
RETURN order,item;
```

![img.png](../img/neo4j_1.png)

2. Calculate the total price of a specific Order
```cypher
MATCH (order:Order {id: 1001})-[:CONTAINS]->(item:Item)
RETURN order.id AS orderID, sum(item.price) AS totalOrderPrice;
```
```text
╒═══════╤═══════════════╕
│orderID│totalOrderPrice│
╞═══════╪═══════════════╡
│1001   │1966           │
└───────┴───────────────┘
```

3. Find all Orders of a specific Customer
```cypher
MATCH (customer:Customer {id: 101})-[:BOUGHT]->(order:Order)
RETURN customer, order;
```
![img.png](../img/neo4j_2.png)
