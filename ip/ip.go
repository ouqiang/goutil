package ip

import (
        "net"
)

// InternalIP 获取内网IP
func InternalIP() string {
        inters, err := net.Interfaces()
        if err != nil {
                return ""
        }
        for _, inter := range inters {
                addrs, err := inter.Addrs()
                if err != nil {
                        continue
                }
                for _, addr := range addrs {
                        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                                if ipnet.IP.To4() != nil {
                                        return ipnet.IP.String()
                                }
                        }
                }
        }
        return ""
}
