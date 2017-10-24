# gocoin

A blockchain/bitcoin implementation written in Go.

## Setup

go version go1.9beta2

After building the app it will run on port 8000 unless a -p flag is set.

[localhost:8000](http://localhost:8000)

## Flags

`-p` Port number on which the client will run.
Usage: `-p=8001`

`--verbose` Verbose will let the app show debug messages in the terminal.
Usage: `--verbose=1`

## API calls

Postman [collection](https://www.getpostman.com/collections/ca46387e102621040d2c) of the call's.

### Transactions

[POST] `http://localhost:8000/transaction`

Add a new transaction
```
{
 "sender": "my address",
 "recipient": "someone else's address",
 "amount": 5
}
```

### Chain

[GET] `http://localhost:8000/mine`   
Mine the next block.  
Response e.g.  

```
{
    "Block": {
        "Index": 3,
        "Timestamp": 1507534014669759993,
        "Transactions": [
            {
                "Sender": "my address",
                "Recipient": "someone else's address",
                "Amount": 5
            },
            {
                "Sender": "0",
                "Recipient": "recipient",
                "Amount": 1
            }
        ],
        "Proof": 27562,
        "PreviousHash": "484dbea2061eb70559cba363897d6c6e63383b233e00fca9a403165a31d5689b"
    },
    "length": 3,
    "message": "New block mined.",
    "transactions": 2
}
```

[GET] `http://localhost:8000/validate`  
Validate the chain.  
```
{
   "length": 3,
   "valid": true
}
```

[GET] `http://localhost:8000/chain`  
Fetch the entire chain.  

```
{
    "chain": [
        {
            "Index": 1,
            "Timestamp": 1507533982542409663,
            "Transactions": [],
            "Proof": 100,
            "PreviousHash": "_"
        },
        {
            "Index": 2,
            "Timestamp": 1507533994836490217,
            "Transactions": [
                {
                    "Sender": "my address",
                    "Recipient": "someone else's address",
                    "Amount": 5
                },
                {
                    "Sender": "0",
                    "Recipient": "recipient",
                    "Amount": 1
                }
            ],
            "Proof": 52838,
            "PreviousHash": "c3b09e9d4930e8af16eb0892d8629572f694741bf046b596cf05c8ca1553799b"
        }
    ],
    "length": 2,
    "transactions": null
}
```

[GET] `http://localhost:8000/resolve`

Resolve conflicts in the chain.
The client checks the list of other nodes in the network and replaces it's blockchain if a larger one is found.
Responses with true if the chain is replaced, otherwise false.

### Network

[GET] `http://localhost:8000/client` Get a list of clients

The response exists of a `length`, representing the total number oof clients, and a `list` of all clients.

```
{
    "length": 3,
    "list": [
        {
            "Ip": "127.0.0.1",
            "Protocol": "http://",
            "Port": 8000,
            "Name": "client1",
            "Hash": "f1c13a0c8292fa5c9dfe565a19f79c2993619e9b6c5da0669b5c886043224673"
        },
        {
            ...
        }
    ]
}
```

[POST] `http://localhost:8000/client` Add a client to the network

The POSTed data should be consistent with a Client.

```
{
 "ip": "123.456.78.90", // string
 "protocol": "http://", // string
 "port": 8080, // int
 "name": "This is me" // string
}
```

## TODO

+ distribute blockchain
+ standarise transactions
+ GET transaction call, based on hash id
+ validate transactions
+ Remove the messenger for (debug) feedback, just build or `go get` another logger.