package gone

import (
	"math/rand"
	"strings"
	"time"

	uuid "github.com/google/uuid"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/xid"
)

// UUID	16 bytes	36 chars	configuration free, not sortable
// shortuuid	16 bytes	22 chars	configuration free, not sortable
// Snowflake	8 bytes	up to 20 chars	needs machine/DC configuration, needs central server, sortable
// MongoID	12 bytes	24 chars	configuration free, sortable
// xid	12 bytes	20 chars	configuration free, sortable

type IDGenerator interface {
	Snowflake() snowflake.ID
	XID() xid.ID
	UUID4() uuid.UUID
	SInt64() int64
	SString() string
	XString() string
	UString() string
	RandString(int) string
}

type id struct {
	n *snowflake.Node
}

var ID id

func (i *id) Snowflake() snowflake.ID {
	return i.n.Generate()
}

func (i *id) XID() xid.ID {
	return xid.New()
}

func (i *id) UUID4() uuid.UUID {
	return uuid.New()
}

func (i *id) SInt64() int64 {
	return i.n.Generate().Int64()
}

func (i *id) SString() string {
	return i.n.Generate().String()
}

func (i *id) XString() string {
	return xid.New().String()
}

func (i *id) UString() string {
	return uuid.New().String()
}

// RandString 最大64个字母.
func (i *id) RandString(count int) string {
	if count > 62 {
		count = 62
	}
	a := []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}
	var r strings.Builder
	for j := 0; j < count; j++ {
		k := RandInt32(0, int32(len(a)))
		r.WriteRune(rune(a[int(k)]))
		a = append(a[:k], a[k+1:]...)
	}

	return r.String()
}

func init() { //nolint
	snowflake.Epoch = time.Now().Unix()
	rand.Seed(rand.Int63n(time.Now().UnixNano())) // nolint
	node := 110 + rand.Int63n(1023-110)           // nolint
	n, _ := snowflake.NewNode(node)               // nolint
	ID = id{
		n: n,
	}
}
