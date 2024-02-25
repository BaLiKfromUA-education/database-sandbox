```js
db.createCollection("items")
```

```js
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

```js
db.items.find({});
```

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
])
```

```js
db.items.distinct("category").length
```

```js
db.items.distinct("producer")
```

```js
db.items.find({category: "Phone", price: {$lt: 1000, $gt: 24}});
```

```js
db.items.find({$or: [{model: "Samsung Galaxy 42"}, {model: "G-Shock 42"}]});
```


```js
db.items.find({producer: {$in: ["LG", "Samsung"]}});
```

```js
db.items.updateMany(
    {category: "Smart Watch"},
    {
        $set: {sale: true, model: "Sale model 42"},
    });

db.items.find({sale: {$exists: true}});
```

```js
db.items.updateMany(
    {sale: {$exists: true}},
    {$inc: {price: 100}}
)

db.items.find({sale: {$exists: true}});
```

```js
db.createCollection("orders");
```

```js
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
        "items_id": ["65db112fa820ad44b33c01be", "65db112fa820ad44b33c01c0"]
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
        "items_id": ["65db112fa820ad44b33c01c2", "65db112fa820ad44b33c01c0"]
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
        "items_id": ["65db112fa820ad44b33c01bc", "65db112fa820ad44b33c01c0"]
    },
])
```

```js
db.orders.find({});
```

```js
db.orders.find({total_sum: {$gt: 1923.4}});
```

```js
db.orders.find({"customer.name": "Andrii", "customer.surname": "Zahuta"});
```

```js
db.orders.find({items_id: {$in: ["65db112fa820ad44b33c01c0"]}});
```

```js
db.orders.updateMany(
    {items_id: {$in: ["65db112fa820ad44b33c01bc"]}},
    {$push: {"items_id": {$each: ["65db112fa820ad44b33c01be"]}}, $inc: {total_sum: 600}}
);
```

```js
db.orders.findOne({order_number: 201515}).items_id.length;
```

```js
db.orders.find({total_sum: {$gt: 1923.4}}, {
    _id: 0,
    customer: 1,
    "payment.cardId": 1
});


db.orders.updateMany(
    {date: {$lt: ISODate("2023-04-14"), $gt: ISODate("2015-04-12")}},
    {$pull: {"items_id": {$in: ["65db112fa820ad44b33c01be"]}}});

db.orders.updateMany(
    {},
    {$set: {"customer.name": "Valentyn"}}
    );
```
