package chinaip

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	
	"github.com/cyfdecyf/bufio"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

var chinaIPs = [][]uint32{}

// IP2Int converts ip from string format to int format
func IP2Int(ip string) (int, error) {
	strs := strings.Split(ip, ".")
	if len(strs) != 4 {
		return -1, Error("isn't ipv4 addr")
	}
	ret := 0
	mul := 1
	for i := 3; i >= 0; i-- {
		a, err := strconv.Atoi(strs[i])
		if err != nil {
			return -1, err
		}
		ret += a * mul
		mul *= 256
	}
	return ret, nil
}

// IsChinaIP returns whether a IPv4 address belong to China
func IsChinaIP(ip string) bool {
	var k, _ = IP2Int(ip)
	var i = uint32(k)
	var l = 0
	var r = len(chinaIPs) - 1
	for l <= r {
		var mid = int((l + r) / 2)
		if i < chinaIPs[mid][0] {
			r = mid - 1
		} else if i > chinaIPs[mid][1] {
			l = mid + 1
		} else {
			return true
		}
	}
	return false
}

func LoadChnRoute(chnRouteFile string) (bool, error) {
	f, err := os.Open(chnRouteFile)
	if err != nil {
		return false, Error("Error opening file: " + chnRouteFile)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "/")
		if len(parts) != 2 {
			panic(errors.New("Invalid CIDR Format"))
		}
		ip := parts[0]
		mask := parts[1]
		count, err := cidrCalc(mask)
		if err != nil {
			panic(err)
		}

		ipLong, err := ipToUint32(ip)
		if err != nil {
			panic(err)
		}
		a := [][]uint32{{ipLong, ipLong+count-1}}
		chinaIPs = append(chinaIPs, a...)
	}
	return true, nil
}

func cidrCalc(mask string) (uint32, error) {
	i, err := strconv.Atoi(mask)
	if err != nil || i > 32 {
		return 0, errors.New("Invalid Mask")
	}
	p := 32 - i
	res := uint32(intPow2(p))
	return res, nil
}

func intPow2(p int) int {
	r := 1
	for i := 0; i < p; i++ {
		r *= 2
	}
	return r
}

func ipToUint32(ipstr string) (uint32, error) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0, errors.New("Invalid IP")
	}
	ip = ip.To4()
	if ip == nil {
		return 0, errors.New("Not IPv4")
	}
	return binary.BigEndian.Uint32(ip), nil
}
