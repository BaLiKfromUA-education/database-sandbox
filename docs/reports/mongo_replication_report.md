# MongoDB replication

My environment:

- Run MongoDB instance [via docker-compose](../../db_environment/mongo/replication/docker_compose.yaml)
- Half of the MongoDB queries were done via JetBrains IDE, others were done directly via installed `mongosh` due to IDE bugs that I faced during execution of some queries.


## Experiments

1. Configure replication for 3 nodes: 1 primary and 2 secondaries.
   docker-compose config

```js
// connected to 'replica-1'
config = {
    "_id": "rs0",
    "members": [{"_id": 0, "host": "replica-1:27017"}, {"_id": 1, "host": "replica-2:27018"}, {
        "_id": 2,
        "host": "replica-3:27019"
    }]
}

rs.initiate(config)
```

```shell
[2024-03-21 21:09:56] Connected
test> config = {
          "_id": "rs0",
          "members": [{"_id": 0, "host": "replica-1:27017"}, {"_id": 1, "host": "replica-2:27018"}, {
              "_id": 2,
              "host": "replica-3:27019"
          }]
      }
[2024-03-21 21:09:57] 1 row retrieved starting from 1 in 443 ms (execution: 219 ms, fetching: 224 ms)
test> rs.initiate(config)
[2024-03-21 21:09:57] 1 row retrieved starting from 1 in 311 ms (execution: 278 ms, fetching: 33 ms)
```

```js
rs.conf();
```

```json
[
  {
    "_id": "rs0",
    "members": [
      {
        "_id": 0,
        "host": "replica-1:27017",
        "arbiterOnly": false,
        "buildIndexes": true,
        "hidden": false,
        "priority": 1,
        "tags": {
        },
        "secondaryDelaySecs": 0,
        "votes": 1
      },
      {
        "_id": 1,
        "host": "replica-2:27018",
        "arbiterOnly": false,
        "buildIndexes": true,
        "hidden": false,
        "priority": 1,
        "tags": {
        },
        "secondaryDelaySecs": 0,
        "votes": 1
      },
      {
        "_id": 2,
        "host": "replica-3:27019",
        "arbiterOnly": false,
        "buildIndexes": true,
        "hidden": false,
        "priority": 1,
        "tags": {
        },
        "secondaryDelaySecs": 0,
        "votes": 1
      }
    ],
    "protocolVersion": 1,
    "settings": {
      "chainingAllowed": true,
      "heartbeatIntervalMillis": 2000,
      "heartbeatTimeoutSecs": 10,
      "electionTimeoutMillis": 10000,
      "catchUpTimeoutMillis": -1,
      "catchUpTakeoverDelayMillis": 30000,
      "getLastErrorModes": {
      },
      "getLastErrorDefaults": {
        "w": 1,
        "wtimeout": 0
      },
      "replicaSetId": {
        "$oid": "65fca22545e6ee4a412bd18c"
      }
    },
    "term": 1,
    "version": 1,
    "writeConcernMajorityJournalDefault": true
  }
]
```

2. Test different **Read Preference Modes**: reading from primary or from secondary.

- Add data via primary connection

```js
// connected to primary node
db.createCollection("messages")

db.messages.insertOne({text: "message 1 to primary"})
```

- Read preference -- `secondary`

```js
db.getMongo().setReadPref('secondary');
```

```shell
rs0 [direct: primary] test> db.getMongo().setReadPref('secondary');

rs0 [direct: primary] test> db.messages.find({});
[
  {
    _id: ObjectId('65fcb130b2c2d775c6c00e7e'),
    text: 'message 1 to primary'
  }
]
```

- Read preference -- `primary`

```js
db.getMongo().setReadPref('primary');
```

```shell
rs0 [direct: primary] test> db.getMongo().setReadPref('primary');

rs0 [direct: primary] test> db.messages.find({});
[
  {
    _id: ObjectId('65fcb130b2c2d775c6c00e7e'),
    text: 'message 1 to primary'
  }
]
```

3. Try to write with **one disabled node** + write concern level 3 and infinite timeout. Try to turn on the disabled
   node during the timeout

```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command               State                            Ports                         
--------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Up      0.0.0.0:27017->27017/tcp,:::27017->27017/tcp           
replication_replica-2_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27018->27018/tcp,:::27018->27018/tcp
replication_replica-3_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27019->27019/tcp,:::27019->27019/tcp

(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker stop replication_replica-1_1 
replication_replica-1_1

(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command                State                              Ports                         
-----------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Exit 137                                                          
replication_replica-2_1   docker-entrypoint.sh --rep ...   Up         27017/tcp, 0.0.0.0:27018->27018/tcp,:::27018->27018/tcp
replication_replica-3_1   docker-entrypoint.sh --rep ...   Up         27017/tcp, 0.0.0.0:27019->27019/tcp,:::27019->27019/tcp

(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker start replication_replica-1_1 
replication_replica-1_1
```

```js
db.messages.insertOne(
    {text: "test infinite timeout with w=3"},
    {writeConcern: {w: 3}})

db.find({})
```

```shell
test> db.messages.insertOne(
          {text: "test infinite timeout with w=3"},
          {writeConcern: {w: 3}})
[2024-03-21 22:30:10] completed in 46 s 418 ms

test> db.messages.find({})
[2024-03-21 22:46:37] 2 rows retrieved starting from 1 in 114 ms (execution: 98 ms, fetching: 16 ms)
```

```js
[
    {
        "_id": {"$oid": "65fcb130b2c2d775c6c00e7e"},
        "text": "message 1 to primary"
    },
    {
        "_id": {"$oid": "65fcb4c4c64ba84a442237ce"},
        "text": "test infinite timeout with w=3"
    }
]
```

4. Similar to the previous point, but set a finite timeout and wait for it to end. Check whether the data has been
   written and is available for reading with `readConcern` level: `majority`

```js
db.messages.insertOne(
    {text: "test finite timeout with w=3"},
    {writeConcern: {w: 3, wtimeout: 5000}})
```

```shell
test> db.messages.insertOne(
          {text: "test finite timeout with w=3"},
          {writeConcern: {w: 3, wtimeout: 5000}})
[2024-03-21 23:02:26] waiting for replication timed out
```

```js
db.messages.find({})
db.messages.find({}).readConcern("majority")
```

```shell
test> db.messages.find({})
[2024-03-21 23:02:31] 4 rows retrieved starting from 1 in 113 ms (execution: 95 ms, fetching: 18 ms)
test> db.messages.find({}).readConcern("majority")
[2024-03-21 23:03:54] 4 rows retrieved starting from 1 in 123 ms (execution: 108 ms, fetching: 15 ms)
```

```js
[
        ...
   {
      "_id": {"$oid": "65fcbc7dc64ba84a442237d2"},
      "text": "test finite timeout with w=3"
   }
]
```

**Result:** message was stored to database

5. Demonstrate primary node re-elections by disabling the current primary (Replica Set Elections),
   and that after the old primary is restored, new data that appeared during its downtime is replicated to it


- **replication_replica-1_1** is a primary node
```shell
balik@balik:~$ mongosh "mongodb://localhost:27017"
...
rs0 [direct: primary] test> 
```
```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command               State                            Ports                         
--------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Up      0.0.0.0:27017->27017/tcp,:::27017->27017/tcp           
replication_replica-2_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27018->27018/tcp,:::27018->27018/tcp
replication_replica-3_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27019->27019/tcp,:::27019->27019/tcp
```

- stop primary
```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker stop replication_replica-1_1 
replication_replica-1_1
```

```shell
balik@balik:~$ mongosh "mongodb://localhost:27017"
Current Mongosh Log ID:	65fcc4e21d7c81bfa7c00e7d
Connecting to:		mongodb://localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.2.1
MongoNetworkError: connect ECONNREFUSED 127.0.0.1:27017
```

- connect to new primary and write new messages
```shell
balik@balik:~$ mongosh "mongodb://localhost:27018"
...
rs0 [direct: primary] test> db.messages.insertOne({text: "this is a new message during old primary out -- task 5"})
{
  acknowledged: true,
  insertedId: ObjectId('65fcc5877bc1a272bac00e7e')
}
```

- enable old primary back and check message
```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker start replication_replica-1_1 
replication_replica-1_1
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command               State                            Ports                         
--------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Up      0.0.0.0:27017->27017/tcp,:::27017->27017/tcp           
replication_replica-2_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27018->27018/tcp,:::27018->27018/tcp
replication_replica-3_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27019->27019/tcp,:::27019->27019/tcp
```

```shell
balik@balik:~$ mongosh "mongodb://localhost:27017"
...
rs0 [direct: secondary] test> db.getMongo().setReadPref('secondary')

rs0 [direct: secondary] test> db.messages.find({}).readConcern('local')
[
  ...
  {
    _id: ObjectId('65fcc5877bc1a272bac00e7e'),
    text: 'this is a new message during old primary out -- task 5'
  }
]
```

- check that `term` value increased
```shell
rs0 [direct: secondary] test> rs.conf()
{
  _id: 'rs0',
  version: 1,
  term: 3, // it is not 2 since I made one experiment before making report
....
}
```

6. Bring the cluster to an inconsistent state using the moment when the primary node does not immediately notice the absence of the secondary node

- after disconnecting two secondary nodes, within 5 seconds, write the value (with `w:1`) to the primary and check that it has been written
```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command               State                            Ports                         
--------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Up      0.0.0.0:27017->27017/tcp,:::27017->27017/tcp           
replication_replica-2_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27018->27018/tcp,:::27018->27018/tcp
replication_replica-3_1   docker-entrypoint.sh --rep ...   Up      27017/tcp, 0.0.0.0:27019->27019/tcp,:::27019->27019/tcp
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker stop replication_replica-3_1 
replication_replica-3_1
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker stop replication_replica-1_1 
replication_replica-1_1
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command                State                              Ports                         
-----------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Exit 137                                                          
replication_replica-2_1   docker-entrypoint.sh --rep ...   Up         27017/tcp, 0.0.0.0:27018->27018/tcp,:::27018->27018/tcp
replication_replica-3_1   docker-entrypoint.sh --rep ...   Exit 137
```

```shell
test> db.messages.insertOne({text: "only primary message -- task 6"}, {writeConcern: {w: 1}})
{
  acknowledged: true,
  insertedId: ObjectId('65fccce23b19c9c20fc00e7f')
}
```

- try to read this value with different levels of read concern - `readConcern: {level: <"majority"|"local"| "linearizable">}` - the value should be available only when reading with the “local” level

```shell
rs0 [direct: secondary] test> db.getMongo().setReadPref('secondary');

rs0 [direct: secondary] test> db.messages.find({}).readConcern('local')
[
  ...
  {
    _id: ObjectId('65fccce23b19c9c20fc00e7f'),
    text: 'only primary message -- task 6'
  }
]

rs0 [direct: secondary] test> db.messages.find({}).readConcern('linearizable')
MongoServerError[NotWritablePrimary]: cannot satisfy linearizable read concern on non-primary node

rs0 [direct: secondary] test> db.messages.find({}).readConcern('majority')
[
  {
    _id: ObjectId('65fcb130b2c2d775c6c00e7e'),
    text: 'message 1 to primary'
  },
  ...
  {
    _id: ObjectId('65fccb793b19c9c20fc00e7e'),
    text: 'only primary message -- task 6'
  }
]

// Majority works for me :(((

```

- turn on the other two nodes so that they do not see the previous primary (it can be disabled) and wait for them to choose a new primary

```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker stop replication_replica-2_1 
replication_replica-2_1

(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker start replication_replica-1_1 
replication_replica-1_1

(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker start replication_replica-3_1 
replication_replica-3_1

(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command                State                              Ports                         
-----------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Up         0.0.0.0:27017->27017/tcp,:::27017->27017/tcp           
replication_replica-2_1   docker-entrypoint.sh --rep ...   Exit 137                                                          
replication_replica-3_1   docker-entrypoint.sh --rep ...   Up         27017/tcp, 0.0.0.0:27019->27019/tcp,:::27019->27019/tcp

```

- connect (turn on) the previous primary node to the cluster and see what happened to the value that was written to it - it should disappear

```shell
balik@balik:~$ mongosh "mongodb://localhost:27018"
...
rs0 [direct: secondary] test> db.getMongo().setReadPref('secondary');

rs0 [direct: secondary] test> db.messages.find({}).readConcern('local')
[
  {
    _id: ObjectId('65fcb130b2c2d775c6c00e7e'),
    text: 'message 1 to primary'
  },
  {
    _id: ObjectId('65fcb4c4c64ba84a442237ce'),
    text: 'test infinite timeout with w=3'
  },
  {
    _id: ObjectId('65fcbc51c64ba84a442237d0'),
    text: 'test infinite timeout with w=3'
  },
  {
    _id: ObjectId('65fcbc7dc64ba84a442237d2'),
    text: 'test finite timeout with w=3'
  },
  {
    _id: ObjectId('65fcc5877bc1a272bac00e7e'),
    text: 'this is a new message during old primary out -- task 5'
  }
]

```


7. Simulate eventual consistency by setting the replication delay for the replica

```shell
rs0 [direct: primary] test> conf.members[1].priority = 0
0
rs0 [direct: primary] test> conf.members[1].hidden = true
true
rs0 [direct: primary] test> conf.members[1].secondaryDelaySecs = 70
70
rs0 [direct: primary] test> rs.reconfig(conf)
rs0 [direct: primary] test> rs.conf()
```

```json
[
  {
    "_id": "rs0",
    "members": [
      {
        "_id": 0,
        "host": "replica-1:27017",
        "arbiterOnly": false,
        "buildIndexes": true,
        "hidden": false,
        "priority": 1,
        "tags": {
        },
        "secondaryDelaySecs": 0,
        "votes": 1
      },
      {
        "_id": 1,
        "host": "replica-2:27018",
        "arbiterOnly": false,
        "buildIndexes": true,
        "hidden": true,
        "priority": 0,
        "tags": {
        },
        "secondaryDelaySecs": 70,
        "votes": 1
      },
      {
        "_id": 2,
        "host": "replica-3:27019",
        "arbiterOnly": false,
        "buildIndexes": true,
        "hidden": false,
        "priority": 1,
        "tags": {
        },
        "secondaryDelaySecs": 0,
        "votes": 1
      }
    ],
    "protocolVersion": 1,
    "settings": {
      "chainingAllowed": true,
      "heartbeatIntervalMillis": 2000,
      "heartbeatTimeoutSecs": 10,
      "electionTimeoutMillis": 10000,
      "catchUpTimeoutMillis": -1,
      "catchUpTakeoverDelayMillis": 30000,
      "getLastErrorModes": {
      },
      "getLastErrorDefaults": {
        "w": 1,
        "wtimeout": 0
      },
      "replicaSetId": {"$oid": "65fcafcc31d513f6ab5c9a29"}
    },
    "term": 6,
    "version": 2,
    "writeConcernMajorityJournalDefault": true
  }
]
```

8. Leave the primary and secondary for which the replication delay is configured (disable the other secondary). Record several values. Try to read value from `readConcern: {level: "linearizable"}`
   There should be a delay until the values are replicated to most nodes

```shell
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker stop replication_replica-3_1 
replication_replica-3_1
(venv) balik@balik:~/Desktop/database-sandbox/db_environment/mongo/replication$ docker-compose -f docker_compose.yaml ps
         Name                        Command                State                              Ports                         
-----------------------------------------------------------------------------------------------------------------------------
replication_replica-1_1   docker-entrypoint.sh --rep ...   Up         0.0.0.0:27017->27017/tcp,:::27017->27017/tcp           
replication_replica-2_1   docker-entrypoint.sh --rep ...   Up         27017/tcp, 0.0.0.0:27018->27018/tcp,:::27018->27018/tcp
replication_replica-3_1   docker-entrypoint.sh --rep ...   Exit 137     
```

```js
db.messages.insertOne({text: "last task final exp -- msg 1"}, {writeConcern: {w: 1}})
db.messages.insertOne({text: "last task final exp -- msg 2"}, {writeConcern: {w: 1}})
db.messages.find({}).readConcern('linearizable')
```

```shell
test> db.messages.insertOne({text: "last task final exp -- msg 1"}, {writeConcern: {w: 1}})
[2024-03-22 00:48:48] completed in 317 ms
test> db.messages.insertOne({text: "last task final exp -- msg 2"}, {writeConcern: {w: 1}})
[2024-03-22 00:48:48] completed in 229 ms
test> db.messages.find({}).readConcern('linearizable')
[2024-03-22 00:49:58] 9 rows retrieved starting from 1 in 1 m 10 s 155 ms (execution: 131 ms, fetching: 1 m 10 s 24 ms)
```


```json
[
  {
    "_id": {
      "$oid": "65fcb130b2c2d775c6c00e7e"
    },
    "text": "message 1 to primary"
  },
  ...
  {
    "_id": {
      "$oid": "65fcd5705c2c22055957b747"
    },
    "text": "last task final exp -- msg 1"
  },
  {
    "_id": {
      "$oid": "65fcd5705c2c22055957b748"
    },
    "text": "last task final exp -- msg 2"
  }
]
```