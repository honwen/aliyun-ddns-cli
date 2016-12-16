package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

const regxIP = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`

func getIPipcn() (ip string) {
	res, err := http.Get("http://ip.cn")
	if err == nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		ip = regexp.MustCompile(regxIP).FindString(string(body))
	}
	return
}

func getIPipio() (ip string) {
	res, err := http.Get("http://ipinfo.io")
	if err == nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		ip = regexp.MustCompile(regxIP).FindString(string(body))
	}
	return
}
