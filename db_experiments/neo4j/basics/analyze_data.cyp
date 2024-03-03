// Get all items from 1 order by order id
MATCH (order:Order {id: 1001})-[:CONTAINS]->(item:Item)
RETURN order,item;
// Calculate the total price of a specific Order
MATCH (order:Order {id: 1001})-[:CONTAINS]->(item:Item)
RETURN order.id AS orderID, sum(item.price) AS totalOrderPrice;
// Find all Orders of a specific Customer
MATCH (customer:Customer {id: 101})-[:BOUGHT]->(order:Order)
RETURN customer, order;