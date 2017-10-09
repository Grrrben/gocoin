# gocoin

A blockchain/bitcoin implementation written in Go.

## Setup

go version go1.9beta2

After building the app it will run on port 8000. (but this will be flexible in the future).

[localhost:8000](http://localhost:8000)

## API calls

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

[GET] `http://localhost:8000/validate`  
Validate the chain.  

[GET] `http://localhost:8000/chain`  
Fetch the entire chain.

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