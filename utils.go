package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const minTimeout = 2000 * time.Millisecond
const regxIP = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`
const regxIP6 = `([0-9A-Fa-f]{0,4}:){2,7}([0-9A-Fa-f]{1,4}$|((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`

var ipAPI = []string{
	"http://v4.myip.la", "http://www.net.cn/static/customercare/yourip.asp", "http://ddns.oray.com/checkip", "http://haoip.cn",
	"http://members.3322.org/dyndns/getip", "http://ns1.dnspod.net:6666", "http://14.215.150.17:6666",
	"http://pv.sohu.com/cityjson?ie=utf-8", "http://whois.pconline.com.cn/ipJson.jsp",
	"http://api-ipv4.ip.sb/ip", "http://ip-api.com/", "http://whatismyip.akamai.com/",
}

var ip6API = []string{
	"http://speed.neu6.edu.cn/getIP.php", "http://v6.myip.la", "http://api-ipv6.ip.sb/ip", "http://ip6only.me/api/", "http://v6.ipv6-test.com/api/myip.php", "https://v6.ident.me",
}

var curlVer = []string{
	"7.70.0", "7.69.1", "7.69.0", "7.68.0", "7.67.0", "7.66.0", "7.65.3", "7.65.2", "7.65.1", "7.65.0", "7.64.1", "7.64.0",
	"7.63.0", "7.62.0", "7.61.1", "7.61.0", "7.60.0", "7.59.0", "7.58.0", "7.57.0", "7.56.1", "7.56.0", "7.55.1", "7.55.0",
	"7.54.1", "7.54.0", "7.53.1", "7.53.0", "7.52.1", "7.52.0", "7.51.0", "7.50.3", "7.50.2", "7.50.1", "7.50.0", "7.49.1",
	"7.49.0", "7.48.0", "7.47.1", "7.47.0", "7.46.0", "7.45.0", "7.44.0", "7.43.0", "7.42.1", "7.42.0", "7.41.0", "7.40.0",
	"7.39.0", "7.38.0", "7.37.1", "7.37.0", "7.36.0",
}

func getIP() (ip string) {
	return apiGetIP(ipAPI, regxIP)
}

func getIP6() (ip string) {
	return apiGetIP(ip6API, regxIP6)
}

func apiGetIP(ipAPI []string, regxIP string) (ip string) {
	var (
		length   = len(ipAPI)
		ipMap    = make(map[string]int, length/5)
		cchan    = make(chan string, length/2)
		regx     = regexp.MustCompile(regxIP)
		maxCount = -1
	)
	for _, url := range ipAPI {
		go func(url string) {
			cchan <- regx.FindString(wGet(url, minTimeout))
		}(url)
	}
	for i := 0; i < length; i++ {
		v := <-cchan
		if 0 == len(v) {
			continue
		}
		if ipMap[v]++; ipMap[v] >= length/2 {
			return v
		}
	}
	for k, v := range ipMap {
		if v > maxCount {
			maxCount = v
			ip = k
		}
	}

	// Use First ipAPI as failsafe
	if 0 == len(ip) {
		ip = regexp.MustCompile(regxIP).FindString(wGet(ipAPI[0], 5*minTimeout))
	}
	return
}

func wGet(url string, timeout time.Duration) (str string) {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "curl/"+curlVer[rand.Intn(len(curlVer))])
	if err != nil {
		return
	}
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	str = string(body)
	// fmt.Println(url, regexp.MustCompile(regxIP).FindString(str))
	return
}

func getDNS(domain string) (ip string) {
	var (
		maxServerCount = 5
		dnsMap         = make(map[string]int, maxServerCount)
		cchan          = make(chan string, maxServerCount)
		maxCount       = -1
		udpClient      = &dns.Client{Net: "udp", Timeout: time.Second}
	)

	for i := 0; i < maxServerCount; i++ {
		go func(dns string) {
			cchan <- getFisrtARecord(udpClient, dns, domain)
		}(fmt.Sprintf("dns%d.hichina.com:53", rand.Intn(30)+1))
	}

	for i := 0; i < maxServerCount; i++ {
		v := <-cchan
		if len(v) == 0 {
			continue
		}
		if dnsMap[v] >= maxServerCount/2 {
			return v
		}
		dnsMap[v]++
	}

	for k, v := range dnsMap {
		if v > maxCount {
			maxCount = v
			ip = k
		}
	}
	return
}

func getFisrtARecord(client *dns.Client, dnsServer, targetDomain string) (ip string) {
	if !strings.HasSuffix(targetDomain, ".") {
		targetDomain += "."
	}
	msg := new(dns.Msg)
	msg.SetQuestion(targetDomain, dns.TypeA)
	r, _, err := client.Exchange(msg, dnsServer)
	if err != nil && (r == nil || r.Rcode != dns.RcodeSuccess) {
		return
	}
	for _, rr := range r.Answer {
		if a, ok := rr.(*dns.A); ok {
			ip = a.A.String()
			break
		}
	}
	return
}
