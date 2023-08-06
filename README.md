## MONGOSH 

- Create db 
```bash=

use mydb

db.collection.insert({name:"john", age:30})


db.createCollection("mycollection")

show dbs
show collections

```

### Replicats 

> Note: Make sure to have data1 folder before only 
- replicat 1 2 3 :
```bash=
mongod --replSet myReplicaSet --dbpath data1 --port 27017 --oplogSize 128 --logpath log1.txt
```


- replicat 2:
```bash=
mongod --replSet myReplicaSet --dbpath data1 --port 27017 --oplogSize 128 

mongod --replSet myReplicaSet --dbpath data2 --port 27018 --oplogSize 128 

mongod --replSet myReplicaSet --dbpath data3 --port 27019 --oplogSize 128   
```
- Add replica sets

### Connect Replicas Each other
```bash=
mongo --port 27017

rs.initiate({
  _id: "myReplicaSet",
  members: [
    {_id: 0, host: "localhost:27017"},
    {_id: 1, host: "localhost:27018"},
    {_id: 2, host: "localhost:27019"}
  ]
})
```
# mongodb-examples
# mongodb-examples
