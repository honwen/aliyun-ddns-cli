package main

import (
	"errors"
	"reflect"

	"github.com/honwen/golibs/cip"
)

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

func myip() (ip string) {
	if result, err := Call(funcs, "myip"); err == nil {
		for _, r := range result {
			if ip := r.String(); ip != "127.0.0.1" {
				return ip
			}
		}
	}
	return
}

func reslove(domain string) (ip string) {
	if result, err := Call(funcs, "reslove", domain); err == nil {
		for _, r := range result {
			if ip := r.String(); ip != "127.0.0.1" {
				return ip
			}
		}
	}
	return
}
