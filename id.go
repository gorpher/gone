package gone

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
)

var n *snowflake.Node

func init() { //nolint
	snowflake.Epoch = time.Now().Unix()
	rand.Seed(rand.Int63n(time.Now().UnixNano())) // nolint
	node := 110 + rand.Int63n(1023-110)           // nolint
	n, _ = snowflake.NewNode(node)                // nolint
}

func IDInt64() int64 {
	return n.Generate().Int64()
}

func IDString() string {
	return n.Generate().String()
}
