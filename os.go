package gone

import (
	"net"
	"strings"
)

//MacAddr 获取机器mac地址，返回mac字串数组
func MacAddr() (upMac []string, err error) {
	var interfaces []net.Interface
	// 获取本机的MAC地址
	interfaces, err = net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		if len(mac.String()) > 0 && strings.Contains(inter.Flags.String(), "up") {
			upMac = append(upMac, mac.String())
		}
	}
	return upMac, nil
}
