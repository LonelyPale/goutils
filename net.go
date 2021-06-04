package goutils

import (
	"net"
)

// 获取本机第一个不为空的MAC地址
func Mac() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}

	mac := ""
	for _, inter := range interfaces {
		mac = inter.HardwareAddr.String() //获取本机MAC地址
		if len(mac) > 0 {
			//fmt.Printf("0x%x\n", binary.LittleEndian.Uint32(inter.HardwareAddr))
			//fmt.Printf("0x%x\n", binary.BigEndian.Uint32(inter.HardwareAddr))
			break
		}
	}

	return mac
}
