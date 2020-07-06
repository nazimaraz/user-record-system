# user-record-system


In this project 9997 users was saved to Redis database using gRPC. On the server side Go language was used and the client side was used Python language. On the client side, user information is read from json files and this information is sent to the server via gRPC.

INSTALL
----------
```
$ git clone https://github.com/nazimaraz/user-record-system.git
$ cd user-record-system
$ go run server/main.go & python3 client/client.py
```
