package boot

import (
	"CheckIPv6/global"
	"CheckIPv6/until"
)

func Init(manage bool, server string) error {
	if !manage {
		ip, err := until.GetLocalIP()
		if err != nil {
			return err
		}
		if err := until.Add(ip, server); err != nil {
			return err
		}
	} else {
		return global.LoadIPs()
	}
	return nil
}
