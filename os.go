package gone

import (
	"fmt"
	"net"
	"regexp"
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

// binary units (IEC 60027)
const (
	_   = 1.0 << (10 * iota) // ignore first value by assigning to blank identifier
	KiB                      // 1 KiB = 1024 Bytes
	MiB                      // 1 MiB = 1048576 Bytes
	GiB                      // 1 GiB = 1073741824 Bytes
	TiB                      // 1 TiB = 1099511627776 Bytes
	PiB                      // 1 PiB = 1125899906842624 Bytes
	EiB                      // 1 EiB = 1152921504606846976 Bytes
)

const (
	KB = 1000
	MB = 1_000_000
	GB = 1_000_000_000
	TB = 1_000_000_000_000
	PB = 1_000_000_000_000_000
	EB = 1_000_000_000_000_000_000
)

var (
	patternBinary  = regexp.MustCompile(`(?i)^(-?\d+(?:\.\d+)?)\s?([KMGTPE]iB?)$`)
	patternDecimal = regexp.MustCompile(`(?i)^(-?\d+(?:\.\d+)?)\s?([KMGTPE]B?|B?)$`)
)

// FormatBinary formats bytes integer to human readable string according to IEC 60027.
// For example, 31323 bytes will return 30.59KB.
func FormatBinary(value int64) string {
	multiple := ""
	val := float64(value)

	switch {
	case value >= EiB:
		val /= EiB
		multiple = "EiB"
	case value >= PiB:
		val /= PiB
		multiple = "PiB"
	case value >= TiB:
		val /= TiB
		multiple = "TiB"
	case value >= GiB:
		val /= GiB
		multiple = "GiB"
	case value >= MiB:
		val /= MiB
		multiple = "MiB"
	case value >= KiB:
		val /= KiB
		multiple = "KiB"
	case value == 0:
		return "0"
	default:
		return strconv.FormatInt(value, 10) + "B"
	}

	return fmt.Sprintf("%.2f%s", val, multiple)
}

// FormatBinaryDecimal formats bytes integer to human readable string according to SI international system of units.
// For example, 31323 bytes will return 31.32KB.
func FormatBinaryDecimal(value int64) string {
	multiple := ""
	val := float64(value)

	switch {
	case value >= EB:
		val /= EB
		multiple = " EB"
	case value >= PB:
		val /= PB
		multiple = " PB"
	case value >= TB:
		val /= TB
		multiple = " TB"
	case value >= GB:
		val /= GB
		multiple = " GB"
	case value >= MB:
		val /= MB
		multiple = " MB"
	case value >= KB:
		val /= KB
		multiple = " KB"
	case value == 0:
		return "0 B"
	default:
		return strconv.FormatInt(value, 10) + " B"
	}

	return fmt.Sprintf("%.2f%s", val, multiple)
}

// ParseBytes parses human readable bytes string to bytes integer.
// For example, 6GiB (6Gi is also valid) will return 6442450944, and
// 6GB (6G is also valid) will return 6000000000.
func ParseBytes(value string) (int64, error) {
	i, err := ParseBinaryString(value)
	if err == nil {
		return i, err
	}
	return ParseStringDecimal(value)
}

// ParseStringDecimal parses human readable bytes string to bytes integer.
// For example, 6GB (6G is also valid) will return 6000000000.
func ParseStringDecimal(value string) (i int64, err error) {
	parts := patternDecimal.FindStringSubmatch(value)
	if len(parts) < 3 {
		return 0, fmt.Errorf("error parsing value=%s", value)
	}
	bytesString := parts[1]
	multiple := strings.ToUpper(parts[2])
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil {
		return
	}

	switch multiple {
	case "K", "KB":
		return int64(bytes * KB), nil
	case "M", "MB":
		return int64(bytes * MB), nil
	case "G", "GB":
		return int64(bytes * GB), nil
	case "T", "TB":
		return int64(bytes * TB), nil
	case "P", "PB":
		return int64(bytes * PB), nil
	case "E", "EB":
		return int64(bytes * EB), nil
	default:
		return int64(bytes), nil
	}
}

// ParseBinaryString parses human readable bytes string to bytes integer.
// For example, 6GiB (6Gi is also valid) will return 6442450944.
func ParseBinaryString(value string) (i int64, err error) {
	parts := patternBinary.FindStringSubmatch(value)
	if len(parts) < 3 {
		return 0, fmt.Errorf("error parsing value=%s", value)
	}
	bytesString := parts[1]
	multiple := strings.ToUpper(parts[2])
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil {
		return
	}

	switch multiple {
	case "KI", "KIB":
		return int64(bytes * KiB), nil
	case "MI", "MIB":
		return int64(bytes * MiB), nil
	case "GI", "GIB":
		return int64(bytes * GiB), nil
	case "TI", "TIB":
		return int64(bytes * TiB), nil
	case "PI", "PIB":
		return int64(bytes * PiB), nil
	case "EI", "EIB":
		return int64(bytes * EiB), nil
	default:
		return int64(bytes), nil
	}
}

// FormatBytesString 格式化bytes单位成可阅读单位形式,由于电脑制造商使用的是1000为单位计算磁盘大小
// 所以基本上使用该函数格式化存储大小.
func FormatBytesString(b int64) string { //nolint
	if b < KB {
		return strconv.FormatInt(b, 10) + " B"
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

// FormatBytesStringOhMyGod 格式化存储大小，理论上应该使用该方式格式化存储大小，但是实际上不是这样的，呜呜呜呜呜呜.
func FormatBytesStringOhMyGod(b int64) string { //nolint
	if b < KiB {
		return strconv.FormatInt(b, 10) + " B"
	}
	if b == KiB {
		return "1 KiB"
	}
	if b < MiB {
		return formatBytesUnit(b, KiB, " KiB")
	}
	if b == MiB {
		return "1 MiB"
	}
	if b < GiB {
		return formatBytesUnit(b, MiB, " MiB")
	}
	if b == GiB {
		return "1 GiB"
	}
	if b < TiB {
		return formatBytesUnit(b, GiB, " GiB")
	}
	if b == TiB {
		return "1 TiB"
	}
	if b < PiB {
		return formatBytesUnit(b, TiB, " TiB")
	}
	if b == PiB {
		return "1 PiB"
	}
	if b < EiB {
		return formatBytesUnit(b, PiB, " PiB")
	}
	if b == EiB {
		return "1 EiB"
	}
	return formatBytesUnit(b, EiB, " EiB")
}

func formatBytesUnit(a, b int64, suffix string) string {
	return decimal.NewFromInt(a).DivRound(decimal.NewFromInt(b), 2).Truncate(2).String() + suffix
}
