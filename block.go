package main

type Block struct {
	Index        int64
	Timestamp    int64
	Transactions []Transaction
	Proof        int64
	PreviousHash string
}

// Block = {
// 	'Index': 1,
// 	'Timestamp': 1506057125.900785,
// 	'Transactions': [
// 	{
// 		'Sender': "8527147fe1f5426f9dd545de4b27ee00",
// 		'Recipient': "a77f5cdfa2934df3954a5c7c7da5df1f",
// 		'Amount': 5,
// 	}
// 	],
// 	'Proof': 324984774000,
// 	'previous_hash': "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
// }
