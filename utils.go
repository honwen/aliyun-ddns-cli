package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

const regxIP = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`

var ipAPI = []string{"http://ip.cn", "http://ipinfo.io", "http://ifconfig.co", "http://pv.sohu.com/cityjson?ie=utf-8", "http://whois.pconline.com.cn/ipJson.jsp"}

func getIP() (ip string) {
	ipMap := make(map[string]int, len(ipAPI))
	for _, url := range ipAPI {
		ip = regexp.MustCompile(regxIP).FindString(wGet(url))
		// log.Println(ip, url)
		if len(ip) > 0 {
			ipMap[ip]++
		}
	}
	max := 0
	for k, v := range ipMap {
		if v > len(ipAPI)/2 {
			return k
		} else if v > max {
			max = v
			ip = k
		}
	}
	return
}

func wGet(url string) (str string) {
	if res, err := http.Get(url); err == nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		str = string(body)
	}
	return
}
