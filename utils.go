package main

import (
	"errors"
	"net"
	"reflect"

	"github.com/honwen/golibs/cip"
)

var ipBlacklist = []string{
	"127.0.0.0/8", "0.0.0.0/24",
}

var funcs = map[string]interface{}{
	"myip":    cip.MyIPv4,
	"reslove": cip.ResloveIPv4,
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func containsCIDR(a, b *net.IPNet) bool {
	ones1, _ := a.Mask.Size()
	ones2, _ := b.Mask.Size()
	return ones1 <= ones2 && a.Contains(b.IP)
}

func containsCIDRString(a, b string) (bool, error) {
	_, net1, err := net.ParseCIDR(a)
	if err != nil {
		return false, err
	}
	_, net2, err := net.ParseCIDR(b)
	if err != nil {
		return false, err
	}
	result := containsCIDR(net1, net2)
	return result, err
}

func myip() (ip string) {
	if result, err := Call(funcs, "myip"); err == nil {
		for _, r := range result {
			ip = r.String()
			for _, it := range ipBlacklist {
				ok, _ := containsCIDRString(it, ip)
				if ok {
					ip = ""
					break
				}
			}
		}
	}
	return
}

func reslove(domain string) (ip string) {
	if result, err := Call(funcs, "reslove", domain); err == nil {
		for _, r := range result {
			return r.String()
		}
	}
	return
}
