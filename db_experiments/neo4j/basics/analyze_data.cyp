// 1. Get all items from 1 order by order id
MATCH (order:Order {id: 1001})-[:CONTAINS]->(item:Item)
RETURN order, item;
// 2. Calculate the total price of a specific Order
MATCH (order:Order {id: 1001})-[:CONTAINS]->(item:Item)
RETURN order.id AS orderID, sum(item.price) AS totalOrderPrice;
// 3. Find all Orders of a specific Customer
MATCH (customer:Customer {id: 101})-[:BOUGHT]->(order:Order)
RETURN customer, order;
// 4. Find all Items bought by a specific Customer
MATCH (customer:Customer {id: 101})-[:BOUGHT]->(order:Order)-[:CONTAINS]->(item:Item)
RETURN customer, order, item;
// 5. Find for the Customer the total amount of items he purchased (through his orders)
MATCH (customer:Customer {id: 103})-[:BOUGHT]->(order:Order)-[:CONTAINS]->(item:Item)
RETURN customer.name, count(item) AS itemCount;
// 6. Find how many times each item was purchased and sort by purchase count
MATCH (order:Order)-[:CONTAINS]->(item:Item)
RETURN item.name, count(order) AS purchaseCount
  ORDER BY purchaseCount DESC;
// 7. Find all Items viewed by a specific Customer
MATCH (customer:Customer {id: 104})-[:VIEWED]->(item:Item)
RETURN customer, item;
// 8. Find other Items purchased together with a specific Item
MATCH (targetItem:Item {id: 1})<-[:CONTAINS]-(order:Order)-[:CONTAINS]->(relatedItem:Item)
  WHERE targetItem <> relatedItem
RETURN order, relatedItem;
// 9. Find all Customers who bought a specific Item
MATCH (customer:Customer)-[:BOUGHT]->(order:Order)-[:CONTAINS]->(item:Item {id: 3})
RETURN customer, order, item;
// 10. Find products for a specific Customer that he viewed but did not buy
MATCH (customer:Customer {id: 101})-[:VIEWED]->(viewedItem:Item)
  WHERE NOT (customer)-[:BOUGHT]->(:Order)-[:CONTAINS]->(viewedItem)
RETURN customer, viewedItem;