# gocoin

A blockchain/bitcoin implementation written in Go.

## Setup

go version go1.9beta2

After building the app it will run on port 8000. (but this will be flexible in the future).

[localhost:8000](http://localhost:8000)

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

### Network

_Entire section is a TODO_

[GET] `http://localhost:8000/miners` Get a list of miners  
[POST] `http://localhost:8000/miners` Add a miner to the network

```
{
 "ip": "123.456.78.90",
 "name": "Hello",
 "description": "Hello World"
}
```

## TODO

+ add multiple servers
+ distribute blockchain
+ standarise transactions
+ GET transaction call, based on hash id
+ validate transactions