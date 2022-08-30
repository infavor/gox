// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// This file contains some common functions.
package gox

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/infavor/gox/convert"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

type IP struct {
	Address string
	Name    string
}

// List2Array converts list to array.
func List2Array(ls *list.List) []interface{} {
	if ls == nil {
		return nil
	}
	arr := make([]interface{}, ls.Len())
	index := 0
	for ele := ls.Front(); ele != nil; ele = ele.Next() {
		arr[index] = ele.Value
		index++
	}
	return arr
}

// ParseHostPortFromConnStr parses host and port from connection string.
func ParseHostPortFromConnStr(connStr string) (string, int) {
	host := strings.Split(connStr, ":")[0]
	port, _ := strconv.Atoi(strings.Split(connStr, ":")[1])
	return host, port
}

// TOperation simulates ternary operation.
func TOperation(condition bool, trueOperation func() interface{}, falseOperation func() interface{}) interface{} {
	if condition {
		if trueOperation == nil {
			return nil
		}
		return trueOperation()
	}
	if falseOperation == nil {
		return nil
	}
	return falseOperation()
}

// TValue ternary operation
func TValue(condition bool, trueValue interface{}, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// WalkList walk a list.
// walker return value as break signal,
// if it is true, break walking
func WalkList(ls *list.List, walker func(item interface{}) bool) {
	if ls == nil {
		return
	}
	for ele := ls.Front(); ele != nil; ele = ele.Next() {
		breakWalk := walker(ele.Value)
		if breakWalk {
			break
		}
	}
}

// Md5Sum calculates md5 value of some strings.
func Md5Sum(input ...string) string {
	h := md5.New()
	if input != nil {
		for _, v := range input {
			io.WriteString(h, v)
		}
	}
	sliceCipherStr := h.Sum(nil)
	sMd5 := hex.EncodeToString(sliceCipherStr)
	return sMd5
}

// LimitRange limits a variable's value in a value range.
// returns defaultValue if the range not provided or the value is not in the range,
func LimitRange(value interface{}, defaultValue interface{}, rangeValues ...interface{}) interface{} {
	if rangeValues == nil || len(rangeValues) == 0 {
		return defaultValue
	}
	in := false
	for _, v := range rangeValues {
		if v == value {
			in = true
			break
		}
	}
	if in {
		return value
	}
	return defaultValue
}

// GetPreferredIPAddress get self ip address by preferred interface
func GetMyAddress(preferredNetworks ...string) string {
	addresses, _ := net.Interfaces()
	var ret list.List
	for i := range addresses {
		info := scan(addresses[i])
		if info == nil {
			continue
		}
		ret.PushBack(info)
	}
	if ret.Len() == 0 {
		return "127.0.0.1"
	}
	selected := ""
	for _, p := range preferredNetworks {
		WalkList(&ret, func(item interface{}) bool {
			if p == item.(*IP).Address {
				selected = p
				return true
			}
			if strings.HasPrefix(item.(*IP).Address, p) || item.(*IP).Name == p {
				if selected == "" {
					selected = item.(*IP).Address
				}
			}
			return false
		})
	}
	return TValue(selected == "", ret.Front().Value.(*IP).Address, selected).(string)
}

func scan(itf net.Interface) *IP {
	var (
		addr      *net.IPNet
		addresses []net.Addr
		err       error
	)
	if addresses, err = itf.Addrs(); err != nil {
		return nil
	}
	if !strings.Contains(itf.Flags.String(), "up") {
		return nil
	}
	for _, a := range addresses {
		if ipNet, ok := a.(*net.IPNet); ok {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				addr = &net.IPNet{
					IP:   ip4,
					Mask: ipNet.Mask[len(ipNet.Mask)-4:],
				}
				break
			}
		}
	}
	if addr == nil {
		return nil
	}
	if addr.IP[0] == 127 {
		return nil
	}
	if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return nil
	}
	return &IP{
		Address: addr.IP.String(),
		Name:    itf.Name,
	}
}

// BlockTest blocks test methods.
func BlockTest() {
	listener, err := net.Listen("tcp", "127.0.0.1:"+convert.IntToStr(10000+rand.Intn(10000)))
	if err != nil {
		os.Exit(3)
	}
	for {
		_, err := listener.Accept()
		fmt.Println(err)
	}
}
