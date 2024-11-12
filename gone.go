package gone

import (
	"github.com/gorpher/gone/conv"
	"github.com/gorpher/gone/osutil"
)

// 导出常用方法

var (
	BytesToStr = conv.BytesToStr
	XID        = osutil.XID
	NumberID   = osutil.NumberID
	UUID       = osutil.UUID
)
