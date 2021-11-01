package main

import (
	"net"
	"regexp"
	"strings"
)

func CheckFileSystemTypeValid(fs string) bool {
	if fs == "" {
		return false
	}
	for _, fst := range FileSystemTypes {
		if fs == fst {
			return true
		}
	}
	return false
}

func CheckMountPointValid(mp string) bool {
	if len(mp) == 0 {
		return false
	}
	if mp == "swap" {
		return true
	}
	fistCharacter := mp[0:1]

	if fistCharacter != "/" {
		return false
	}
	return true
}

func CheckIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}

}

func CheckHost(host string) bool {
	host = strings.Trim(host, " ")
	//re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if re.MatchString(host) {
		return true
	}
	return false
}
