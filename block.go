package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grrrben/glog"
)

type Block struct {
	Index        int64         `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int64         `json:"proof"`
	PreviousHash string        `json:"previousHash"`
}

// announceMinedBlock shares the block with other nodes. It is done in a goroutine.
// Other nodes should check the validity of the new block on their chain and add it.
func announceMinedBlock(cl Node, bl Block) {
	url := fmt.Sprintf("%s/block/distributed", cl.getAddress())

	blockAndSender := map[string]interface{}{"block": bl, "sender": me.getAddress()}
	payload, err := json.Marshal(blockAndSender)
	if err != nil {
		glog.Errorf("Could not marshall block or node. Msg: %s", err)
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		glog.Errorf("Request setup error: %s", err)
	} else {
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			glog.Warningf("POST request error: %s", err)
			// I don't want to panic here, but it might be a good idea to
			// remove the node from the list
		} else {
			resp.Body.Close()
		}
	}
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
