package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/grrrben/golog"
	"time"
)

type wallet struct {
	hash   string
	credit float64
}

func createWallet() wallet {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, time.Now().Unix())
	if err != nil {
		golog.Warningf("Could not createWallet. Msg: %s", err)
	}

	w := wallet{
		hash:   fmt.Sprintf("%x", sha256.Sum256(buf.Bytes())),
		credit: 0,
	}
	return w
}
