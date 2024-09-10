# Golang gRPC server

This application accepts function call requests using gRPC.  
c.Fetch(url) - receives CSV document by endpoint, adds/updates data to MongoDB  
c.List(params) - returns sorted and paginated data from MongoDB
    

  ### Tools:
  - go 1.19
  - MongoDB
  - Protobuf 

 ### How to use this
 Run containers with MongoDB:

```cmd
docker run --rm -d --name audit-log-mongo -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=qwerty -p 27017:27017 mongo:latest
```

Build project and run. Waiting gRPC requests from client.
