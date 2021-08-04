package gone

import (
	"net"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// MacAddr 获取机器mac地址，返回mac字串数组.
func MacAddr() (upMac []string, err error) {
	var interfaces []net.Interface
	// 获取本机的MAC地址
	interfaces, err = net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr // 获取本机MAC地址
		if len(mac.String()) > 0 && strings.Contains(inter.Flags.String(), "up") {
			upMac = append(upMac, mac.String())
		}
	}
	return upMac, nil
}

const (
	Bytes = 1
	KB    = Bytes << 10 // 1 KB = 1024 Bytes
	MB    = Bytes << 20 // 1 MB = 1048576 Bytes
	GB    = Bytes << 30 // 1 GB = 1073741824 Bytes
	TB    = Bytes << 40 // 1 TB = 1099511627776 Bytes
	PB    = Bytes << 50 // 1 PB = 1125899906842624 Bytes
	EB    = Bytes << 60 // 1 EB = 1152921504606846976 Bytes
)

// FormatBytesStringOhMyGod 格式化存储大小，理论上应该使用该方式格式化存储大小，但是实际上不是这样的，呜呜呜呜呜呜.
func FormatBytesStringOhMyGod(b int64) string { //nolint
	if b < KB {
		return strconv.FormatInt(b, 10) + " bytes"
	}
	if b == KB {
		return "1 KB"
	}
	if b < MB {
		return formatBytesUnit(b, KB, " KB")
	}
	if b == MB {
		return "1 MB"
	}
	if b < GB {
		return formatBytesUnit(b, MB, " MB")
	}
	if b == GB {
		return "1 GB"
	}
	if b < TB {
		return formatBytesUnit(b, GB, " GB")
	}
	if b == TB {
		return "1 TB"
	}
	if b < PB {
		return formatBytesUnit(b, TB, " TB")
	}
	if b == PB {
		return "1 PB"
	}
	if b < EB {
		return formatBytesUnit(b, PB, " PB")
	}
	if b == EB {
		return "1 EB"
	}
	return formatBytesUnit(b, EB, " EB")
}

const (
	kb = 1000
	mb = 1_000_000
	gb = 1_000_000_000
	tb = 1_000_000_000_000
	pb = 1_000_000_000_000_000
	eb = 1_000_000_000_000_000_000
)

// FormatBytesString 格式化bytes单位成可阅读单位形式,由于电脑制造商使用的是1000为单位计算磁盘大小
// 所以基本上使用该函数格式化存储大小.
func FormatBytesString(b int64) string { //nolint
	if b < kb {
		return strconv.FormatInt(b, 10) + " bytes"
	}
	if b == kb {
		return "1 KB"
	}
	if b < mb {
		return formatBytesUnit(b, kb, " KB")
	}
	if b == mb {
		return "1 MB"
	}
	if b < gb {
		return formatBytesUnit(b, mb, " MB")
	}
	if b == gb {
		return "1 GB"
	}
	if b < tb {
		return formatBytesUnit(b, gb, " GB")
	}
	if b == tb {
		return "1 TB"
	}
	if b < pb {
		return formatBytesUnit(b, tb, " TB")
	}
	if b == pb {
		return "1 PB"
	}
	if b < eb {
		return formatBytesUnit(b, pb, " PB")
	}
	if b == eb {
		return "1 EB"
	}
	return formatBytesUnit(b, eb, " EB")
}

func formatBytesUnit(a, b int64, suffix string) string {
	return decimal.NewFromInt(a).DivRound(decimal.NewFromInt(b), 2).Truncate(2).String() + suffix
}
