# MongoDB basics

My environment:

- Run MongoDB instance [via docker-compose](../../db_environment/mongo/docker_compose.yaml)
- Execute queries via MongoDB console from my JetBrains IDE.

### Items collection

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