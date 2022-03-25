package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/transerver/commons/logger"
)

var (
	// GenerateId examples
	// Int64  ID: 1318819864034414592
	// Bytes  ID: [49 51 49 56 56 49 57 56 54 52 48 51 52 52 49 52 53 57 50]
	// String ID: 1318819864034414592
	// Base2  ID: 1001001001101011000111010110111001110100000000001000000000000
	// Base32 ID: brumdiz8eyryy
	// Base36 ID: a0pmmhku9kw0
	// Base58 ID: 44ycTvdggum
	// Base64 ID: MTMxODgxOTg2NDAzNDQxNDU5Mg==
	// ID Time  : 1603266132731
	// ID Node  : 1
	// ID Step  : 0
	uids map[string]*snowflake.Node
)

const defaultNodeId = ""

func FetchUidNode(nodeId string) (*snowflake.Node, bool) {
	if uid, ok := uids[nodeId]; ok {
		return uid, ok
	}

	var err error
	uid, err := snowflake.NewNode(1)
	if err != nil {
		logger.Errorf("Snowflakes generator init error: %+v", err)
		return nil, false
	}
	if uids == nil {
		uids = make(map[string]*snowflake.Node)
	}
	uids[nodeId] = uid
	return uid, true
}

func GenerateIdWithDefaultNode() snowflake.ID {
	id, _ := GenerateId(defaultNodeId)
	return id
}

func GenerateId(nodeId string) (snowflake.ID, bool) {
	if uid, ok := FetchUidNode(nodeId); ok {
		return uid.Generate(), true
	}
	logger.Errorf("generate snowflake ID fail")
	return snowflake.ID(-1), false
}
