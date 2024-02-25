# MongoDB basics

My environment:

- Run MongoDB instance [via docker-compose](../../db_environment/mongo/docker_compose.yaml)
- Execute queries via MongoDB console from my JetBrains IDE.

### Items

1. Create multiple items with different set of properties Phone/TV/Smart Watch/ ....

```js
db.createCollection("items");

db.items.insertMany([
    {category: "Phone", model: "Samsung Galaxy 42", producer: "Samsung", price: 42},
    {category: "Phone", model: "iPhone 6", producer: "Apple", price: 600},
    {category: "Phone", model: "iPhone 15", producer: "Apple", price: 1600},
    {category: "Phone", model: "iPhone 10", producer: "Apple", price: 1000},
    {category: "TV", model: "LG model 42", producer: "LG", price: 999},
    {category: "TV", model: "Samsung TV 42", producer: "Samsung", price: 777},
    {category: "Smart Watch", model: "Garmin RUN 42", producer: "Garmin", price: 200},
    {category: "Smart Watch", model: "G-Shock 42", producer: "Casio", price: 400},
]);
```

```shell
test> db.createCollection("items");
[2024-02-25 11:24:13] completed in 157 ms
test> db.items.insertMany([
          {category: "Phone", model: "Samsung Galaxy 42", producer: "Samsung", price: 42},
          {category: "Phone", model: "iPhone 6", producer: "Apple", price: 600},
          {category: "Phone", model: "iPhone 15", producer: "Apple", price: 1600},
          {category: "Phone", model: "iPhone 10", producer: "Apple", price: 1000},
          {category: "TV", model: "LG model 42", producer: "LG", price: 999},
          {category: "TV", model: "Samsung TV 42", producer: "Samsung", price: 777},
          {category: "Smart Watch", model: "Garmin RUN 42", producer: "Garmin", price: 200},
          {category: "Smart Watch", model: "G-Shock 42", producer: "Casio", price: 400},
      ]);
[2024-02-25 11:24:13] completed in 215 ms
```

2. Write a query that displays all items (display in JSON)

```js
db.items.find({});
```

```json
[
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01cf"
    },
    "category": "Phone",
    "model": "Samsung Galaxy 42",
    "price": 42,
    "producer": "Samsung"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d0"
    },
    "category": "Phone",
    "model": "iPhone 6",
    "price": 600,
    "producer": "Apple"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d1"
    },
    "category": "Phone",
    "model": "iPhone 15",
    "price": 1600,
    "producer": "Apple"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d2"
    },
    "category": "Phone",
    "model": "iPhone 10",
    "price": 1000,
    "producer": "Apple"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d3"
    },
    "category": "TV",
    "model": "LG model 42",
    "price": 999,
    "producer": "LG"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d4"
    },
    "category": "TV",
    "model": "Samsung TV 42",
    "price": 777,
    "producer": "Samsung"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d5"
    },
    "category": "Smart Watch",
    "model": "Garmin RUN 42",
    "price": 200,
    "producer": "Garmin"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d6"
    },
    "category": "Smart Watch",
    "model": "G-Shock 42",
    "price": 400,
    "producer": "Casio"
  }
]

```

3. Count how many items are in a certain category

```js
db.items.aggregate([
    {
        $group: {
            _id: "$category",
            count: {$sum: 1}
        }
    },
    {
        $project: {
            category: "$_id",
            count: 1,
            _id: 0
        }
    }
]);
```

```json
[
  {
    "category": "TV",
    "count": 2
  },
  {
    "category": "Phone",
    "count": 4
  },
  {
    "category": "Smart Watch",
    "count": 2
  }
]
```

4. Count how many different product categories there are

```js
db.items.distinct("category").length;
```

```json
3
```

5. Display a list of all producers without repetitions

```js
db.items.distinct("producer");
```

```json
[
  {
    "result": "Apple"
  },
  {
    "result": "Casio"
  },
  {
    "result": "Garmin"
  },
  {
    "result": "LG"
  },
  {
    "result": "Samsung"
  }
]
```

6. Write queries that select items based on various criteria and their combinations:
    - category and price (in the interval) - construction `$and`:

```js
db.items.find({category: "Phone", price: {$lt: 1000, $gt: 24}});
```

```json
[
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01cf"
    },
    "category": "Phone",
    "model": "Samsung Galaxy 42",
    "price": 42,
    "producer": "Samsung"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d0"
    },
    "category": "Phone",
    "model": "iPhone 6",
    "price": 600,
    "producer": "Apple"
  }
]
```

- model either one or the other - construction `$or`

```js
db.items.find({$or: [{model: "Samsung Galaxy 42"}, {model: "G-Shock 42"}]});
```

```json
[
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01cf"
    },
    "category": "Phone",
    "model": "Samsung Galaxy 42",
    "price": 42,
    "producer": "Samsung"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d6"
    },
    "category": "Smart Watch",
    "model": "G-Shock 42",
    "price": 400,
    "producer": "Casio"
  }
]
```

- producer from the list - design `$in`

```js
db.items.find({producer: {$in: ["LG", "Samsung"]}});
```

```json
[
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01cf"
    },
    "category": "Phone",
    "model": "Samsung Galaxy 42",
    "price": 42,
    "producer": "Samsung"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d3"
    },
    "category": "TV",
    "model": "LG model 42",
    "price": 999,
    "producer": "LG"
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d4"
    },
    "category": "TV",
    "model": "Samsung TV 42",
    "price": 777,
    "producer": "Samsung"
  }
]
```

7. Update certain products by changing existing items and
   add new properties (characteristics) to all items according to a certain criterion.

```js
db.items.updateMany(
    {category: "Smart Watch"},
    {
        $set: {sale: true, model: "Sale model 42"},
    });

```

```shell
test> db.items.updateMany(
          {category: "Smart Watch"},
          {
              $set: {sale: true, model: "Sale model 42"},
          });
[2024-02-25 15:58:59] completed in 233 ms
```

8. Find items that have (present field) certain properties

```js
db.items.find({sale: {$exists: true}});
```

```json
[
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d5"
    },
    "category": "Smart Watch",
    "model": "Sale model 42",
    "price": 200,
    "producer": "Garmin",
    "sale": true
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d6"
    },
    "category": "Smart Watch",
    "model": "Sale model 42",
    "price": 400,
    "producer": "Casio",
    "sale": true
  }
]
```

9. For found items, increase their price by a certain amount

```js
db.items.updateMany(
    {
        sale: {
            $exists: true
        }
    },
    {$inc: {price: 100}}
);

db.items.find({
    sale: {
        $exists: true
    }
});
```

```json
[
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d5"
    },
    "category": "Smart Watch",
    "model": "Sale model 42",
    "price": 300,
    "producer": "Garmin",
    "sale": true
  },
  {
    "_id": {
      "$oid": "65db235da820ad44b33c01d6"
    },
    "category": "Smart Watch",
    "model": "Sale model 42",
    "price": 500,
    "producer": "Casio",
    "sale": true
  }
]
```

### Orders

1. Create several orders with different sets of items, but so that one of the items is in several orders

```js
db.createCollection("orders");

db.orders.insertMany([
   {
      "order_number": 201513,
      "date": ISODate("2015-04-14"),
      "total_sum": 1923.4,
      "customer": {
         "name": "Andrii",
         "surname": "Rodinov",
         "phones": ["+9876543", "+1234567"],
         "address": "PTI, Peremohy 37, Kyiv, UA"
      },
      "payment": {
         "card_owner": "Andrii Rodionov",
         "cardId": "12345678"
      },
      "items_id": [ObjectId("65db235da820ad44b33c01cf"), ObjectId("65db235da820ad44b33c01d3")]
   },
   {
      "order_number": 201514,
      "date": ISODate("2023-04-14"),
      "total_sum": 3333,
      "customer": {
         "name": "Valentyn",
         "surname": "Yukhymenko",
         "phones": ["+33333", "+44444"],
         "address": "London, UK"
      },
      "payment": {
         "card_owner": "Valentyn Yukhymenko",
         "cardId": "987654321"
      },
      "items_id": [ObjectId("65db235da820ad44b33c01d5"), ObjectId("65db235da820ad44b33c01d3")]
   },
   {
      "order_number": 201515,
      "date": ISODate("2023-09-14"),
      "total_sum": 2222.333,
      "customer": {
         "name": "Andrii",
         "surname": "Zahuta",
         "phones": ["+12345678"],
         "address": "Prague"
      },
      "payment": {
         "card_owner": "Valentyn Yukhymenko",
         "cardId": "987654321"
      },
      "items_id": [ObjectId("65db235da820ad44b33c01d1"), ObjectId("65db235da820ad44b33c01d3")]
   },
]);
```

```shell
[2024-02-25 16:02:17] completed in 199 ms
[2024-02-25 16:04:41] completed in 285 ms
```

2. Display all orders

```js
db.orders.find({});
```

```json
[
   {
      "_id": {"$oid": "65db6fbe907e240b7da62f7f"},
      "customer": {
         "name": "Andrii",
         "surname": "Rodinov",
         "phones": ["+9876543", "+1234567"],
         "address": "PTI, Peremohy 37, Kyiv, UA"
      },
      "date": {"$date": "2015-04-14T00:00:00.000Z"},
      "items_id": [
         {"$oid": "65db235da820ad44b33c01cf"},
         {"$oid": "65db235da820ad44b33c01d3"}
      ],
      "order_number": 201513,
      "payment": {
         "card_owner": "Andrii Rodionov",
         "cardId": "12345678"
      },
      "total_sum": 1923.4
   },
   {
      "_id": {"$oid": "65db6fbe907e240b7da62f80"},
      "customer": {
         "name": "Valentyn",
         "surname": "Yukhymenko",
         "phones": ["+33333", "+44444"],
         "address": "London, UK"
      },
      "date": {"$date": "2023-04-14T00:00:00.000Z"},
      "items_id": [
         {"$oid": "65db235da820ad44b33c01d5"},
         {"$oid": "65db235da820ad44b33c01d3"}
      ],
      "order_number": 201514,
      "payment": {
         "card_owner": "Valentyn Yukhymenko",
         "cardId": "987654321"
      },
      "total_sum": 3333
   },
   {
      "_id": {"$oid": "65db6fbe907e240b7da62f81"},
      "customer": {
         "name": "Andrii",
         "surname": "Zahuta",
         "phones": ["+12345678"],
         "address": "Prague"
      },
      "date": {"$date": "2023-09-14T00:00:00.000Z"},
      "items_id": [
         {"$oid": "65db235da820ad44b33c01d1"},
         {"$oid": "65db235da820ad44b33c01d3"}
      ],
      "order_number": 201515,
      "payment": {
         "card_owner": "Valentyn Yukhymenko",
         "cardId": "987654321"
      },
      "total_sum": 2222.333
   }
]
```

3. Display orders with a value greater than a certain value

```js
db.orders.find({total_sum: {$gt: 1923.4}});
```

```json
[
   {
      "_id": {"$oid": "65db6fbe907e240b7da62f80"},
      "customer": {
         "name": "Valentyn",
         "surname": "Yukhymenko",
         "phones": ["+33333", "+44444"],
         "address": "London, UK"
      },
      "date": {"$date": "2023-04-14T00:00:00.000Z"},
      "items_id": [
         {"$oid": "65db235da820ad44b33c01d5"},
         {"$oid": "65db235da820ad44b33c01d3"}
      ],
      "order_number": 201514,
      "payment": {
         "card_owner": "Valentyn Yukhymenko",
         "cardId": "987654321"
      },
      "total_sum": 3333
   },
   {
      "_id": {"$oid": "65db6fbe907e240b7da62f81"},
      "customer": {
         "name": "Andrii",
         "surname": "Zahuta",
         "phones": ["+12345678"],
         "address": "Prague"
      },
      "date": {"$date": "2023-09-14T00:00:00.000Z"},
      "items_id": [
         {"$oid": "65db235da820ad44b33c01d1"},
         {"$oid": "65db235da820ad44b33c01d3"}
      ],
      "order_number": 201515,
      "payment": {
         "card_owner": "Valentyn Yukhymenko",
         "cardId": "987654321"
      },
      "total_sum": 2222.333
   }
]
```

4. Find orders made by one customer

```js
db.orders.find({"customer.name": "Andrii", "customer.surname": "Zahuta"});
```

```json
[
   {
      "_id": {"$oid": "65db6fbe907e240b7da62f81"},
      "customer": {
         "name": "Andrii",
         "surname": "Zahuta",
         "phones": ["+12345678"],
         "address": "Prague"
      },
      "date": {"$date": "2023-09-14T00:00:00.000Z"},
      "items_id": [
         {"$oid": "65db235da820ad44b33c01d1"},
         {"$oid": "65db235da820ad44b33c01d3"}
      ],
      "order_number": 201515,
      "payment": {
         "card_owner": "Valentyn Yukhymenko",
         "cardId": "987654321"
      },
      "total_sum": 2222.333
   }
]
```

5. Find all orders with a certain item(s) (you can search by ObjectId)

```js
db.orders.find({items_id: {$in: [ObjectId("65db235da820ad44b33c01d3")]}});
```

```json
[
  {
    "_id": {
      "$oid": "65db6fbe907e240b7da62f7f"
    },
    "customer": {
      "name": "Andrii",
      "surname": "Rodinov",
      "phones": [
        "+9876543",
        "+1234567"
      ],
      "address": "PTI, Peremohy 37, Kyiv, UA"
    },
    "date": {
      "$date": "2015-04-14T00:00:00.000Z"
    },
    "items_id": [
      {
        "$oid": "65db235da820ad44b33c01cf"
      },
      {
        "$oid": "65db235da820ad44b33c01d3"
      }
    ],
    "order_number": 201513,
    "payment": {
      "card_owner": "Andrii Rodionov",
      "cardId": "12345678"
    },
    "total_sum": 1923.4
  },
  {
    "_id": {
      "$oid": "65db6fbe907e240b7da62f80"
    },
    "customer": {
      "name": "Valentyn",
      "surname": "Yukhymenko",
      "phones": [
        "+33333",
        "+44444"
      ],
      "address": "London, UK"
    },
    "date": {
      "$date": "2023-04-14T00:00:00.000Z"
    },
    "items_id": [
      {
        "$oid": "65db235da820ad44b33c01d5"
      },
      {
        "$oid": "65db235da820ad44b33c01d3"
      }
    ],
    "order_number": 201514,
    "payment": {
      "card_owner": "Valentyn Yukhymenko",
      "cardId": "987654321"
    },
    "total_sum": 3333
  },
  {
    "_id": {
      "$oid": "65db6fbe907e240b7da62f81"
    },
    "customer": {
      "name": "Andrii",
      "surname": "Zahuta",
      "phones": [
        "+12345678"
      ],
      "address": "Prague"
    },
    "date": {
      "$date": "2023-09-14T00:00:00.000Z"
    },
    "items_id": [
      {
        "$oid": "65db235da820ad44b33c01d1"
      },
      {
        "$oid": "65db235da820ad44b33c01d3"
      }
    ],
    "order_number": 201515,
    "payment": {
      "card_owner": "Valentyn Yukhymenko",
      "cardId": "987654321"
    },
    "total_sum": 2222.333
  }
]
```

6. Add one more item to all orders with a certain item and increase the existing order value by some X value

```js
db.orders.updateMany(
        {items_id: {$in: [ObjectId("65db235da820ad44b33c01d1")]}},
        {$push: {"items_id": {$each: [ObjectId("65db235da820ad44b33c01d0")]}}, $inc: {total_sum: 600}}
);
```

```shell
[2024-02-25 16:15:47] completed in 188 ms
```

7. Display the number of items in a certain order

```js
db.orders.findOne({order_number: 201515}).items_id.length;
```

```json
3
```

8. Display only customer information and credit card numbers for orders over a certain price

```js
db.orders.find({total_sum: {$gt: 1923.4}}, {
    _id: 0,
    customer: 1,
    "payment.cardId": 1
});

```

```json
[
  {
    "customer": {
      "name": "Valentyn",
      "surname": "Yukhymenko",
      "phones": [
        "+33333",
        "+44444"
      ],
      "address": "London, UK"
    },
    "payment": {
      "cardId": "987654321"
    }
  },
  {
    "customer": {
      "name": "Andrii",
      "surname": "Zahuta",
      "phones": [
        "+12345678"
      ],
      "address": "Prague"
    },
    "payment": {
      "cardId": "987654321"
    }
  }
]
```

9. Remove an item from orders made within a specific date range

```js
db.orders.updateMany(
        {date: {$lt: ISODate("2023-04-14"), $gt: ISODate("2015-04-12")}},
        {$pull: {"items_id": {$in: [ObjectId("65db235da820ad44b33c01d3")]}}});
```

````shell
[2024-02-25 16:15:47] completed in 188 ms
````

10. Rename the name (surname) of the customer in all orders

```js
db.orders.updateMany(
    {},
    {$set: {"customer.name": "Valentyn"}}
    );
```

```shell
[2024-02-25 16:22:54] completed in 178 ms
```

11. Find the orders made by one customer, and display only information about the customer and the ordered items by substituting the names of the products and their cost instead of ObjectId("***") (similar to a join between the orders and items tables).

```js
db.orders.aggregate([
   {
      $match: { "customer.surname": "Rodinov" }
   },
   {
      $lookup: {
         from: "items",
         localField: "items_id",
         foreignField: "_id",
         as: "ordered_items"
      }
   },
   {
      $project: {
         "_id": 0,
         "customer": 1,
         "ordered_items.model": 1,
         "ordered_items.price": 1
      }
   }
])
```

```json
[
  {
    "customer": {
      "name": "Valentyn",
      "surname": "Rodinov",
      "phones": [
        "+9876543",
        "+1234567"
      ],
      "address": "PTI, Peremohy 37, Kyiv, UA"
    },
    "ordered_items": [
      {
        "model": "Samsung Galaxy 42",
        "price": 42
      }
    ]
  }
]
```


### Reviews (capped collection)

1. Create collection
```js
db.createCollection( "reviews", { capped: true, size: 100000, max: 5 } )
```

```shell
test> db.createCollection( "reviews", { capped: true, size: 100000, max: 5 } )
[2024-02-25 16:58:51] completed in 171 ms
```

2. Add 5 reviews

```js
db.reviews.insertMany(
    [
        {date: ISODate("2015-04-14"), message: "First great review"},
        {date: ISODate("2015-04-15"), message: "Second great review"},
        {date: ISODate("2015-04-16"), message: "Third great review"},
        {date: ISODate("2015-04-17"), message: "Fourth great review"},
        {date: ISODate("2015-04-18"), message: "Fifth great review"},
    ]);
```

```shell
test> db.reviews.insertMany(
          [
              {date: ISODate("2015-04-14"), message: "First great review"},
              {date: ISODate("2015-04-15"), message: "Second great review"},
              {date: ISODate("2015-04-16"), message: "Third great review"},
              {date: ISODate("2015-04-17"), message: "Fourth great review"},
              {date: ISODate("2015-04-18"), message: "Fifth great review"},
          ]);
[2024-02-25 17:01:19] completed in 208 ms
```

3. Get all reviews

```js
db.reviews.find({});
```

```json
[
  {
    "_id": {"$oid": "65db725f907e240b7da62f84"},
    "date": {"$date": "2015-04-14T00:00:00.000Z"},
    "message": "First great review"
  },
  {
    "_id": {"$oid": "65db725f907e240b7da62f85"},
    "date": {"$date": "2015-04-15T00:00:00.000Z"},
    "message": "Second great review"
  },
  {
    "_id": {"$oid": "65db725f907e240b7da62f86"},
    "date": {"$date": "2015-04-16T00:00:00.000Z"},
    "message": "Third great review"
  },
  {
    "_id": {"$oid": "65db725f907e240b7da62f87"},
    "date": {"$date": "2015-04-17T00:00:00.000Z"},
    "message": "Fourth great review"
  },
  {
    "_id": {"$oid": "65db725f907e240b7da62f88"},
    "date": {"$date": "2015-04-18T00:00:00.000Z"},
    "message": "Fifth great review"
  }
]
```

4. Add one more review

```js
db.reviews.insertOne({date: ISODate("2015-04-19"), message: "First BAD review"});
```

```shell
test> db.reviews.insertOne({date: ISODate("2015-04-19"), message: "First BAD review"});
[2024-02-25 17:02:58] completed in 174 ms
```

5. Check reviews again

```js
db.reviews.find({});
```


```json
[
  {
    "_id": {
      "$oid": "65db725f907e240b7da62f85"
    },
    "date": {
      "$date": "2015-04-15T00:00:00.000Z"
    },
    "message": "Second great review"
  },
  {
    "_id": {
      "$oid": "65db725f907e240b7da62f86"
    },
    "date": {
      "$date": "2015-04-16T00:00:00.000Z"
    },
    "message": "Third great review"
  },
  {
    "_id": {
      "$oid": "65db725f907e240b7da62f87"
    },
    "date": {
      "$date": "2015-04-17T00:00:00.000Z"
    },
    "message": "Fourth great review"
  },
  {
    "_id": {
      "$oid": "65db725f907e240b7da62f88"
    },
    "date": {
      "$date": "2015-04-18T00:00:00.000Z"
    },
    "message": "Fifth great review"
  },
  {
    "_id": {
      "$oid": "65db72c2907e240b7da62f8a"
    },
    "date": {
      "$date": "2015-04-19T00:00:00.000Z"
    },
    "message": "First BAD review"
  }
]
```