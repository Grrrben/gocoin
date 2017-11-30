package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/grrrben/golog"
	"net/http"
)

type Block struct {
	Index        int64
	Timestamp    int64
	Transactions []Transaction
	Proof        int64
	PreviousHash string
}

// announceMinedBlock shares the block with other clients. It is done in a goroutine.
// Other clients should check the validity of the new block on their chain and add it.
func announceMinedBlock(cl Client, bl Block) {
	url := fmt.Sprintf("%s/block/distributed", cls.getAddress(cl))

	blockAndSender := map[string]interface{}{"block": bl, "sender": cls.getAddress(me)}
	payload, err := json.Marshal(blockAndSender)
	if err != nil {
		golog.Errorf("Could not marshall block or client. Msg: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		golog.Warningf("Request setup error: %s", err)
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Warningf("POST request error: %s", err)
		// I don't want to panic here, but it might be a good idea to
		// remove the client from the list
	}
	defer resp.Body.Close()
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
