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

const minTimeout = 333 * time.Millisecond
const regxIP = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`

var ipAPI = []string{
	"http://ip.cn", "http://ipinfo.io", "http://ifconfig.co", "http://myip.ipip.net",
	"http://cnc.synology.cn:81", "http://jpc.synology.com:81", "http://usc.synology.com:81",
	"http://ip.6655.com/ip.aspx", "http://pv.sohu.com/cityjson?ie=utf-8", "http://whois.pconline.com.cn/ipJson.jsp",
	"http://ddns.oray.com/checkip",
}

func getIP() (ip string) {
	var (
		length   = len(ipAPI)
		ipMap    = make(map[string]int, length)
		cchan    = make(chan string, length)
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
		if ipMap[v] >= length/2 {
			return v
		}
		ipMap[v]++
	}
	for k, v := range ipMap {
		if v > maxCount {
			maxCount = v
			ip = k
		}
	}

	// Use First ipAPI as failsafe
	if 0 == len(ip) {
		ip = regexp.MustCompile(regxIP).FindString(wGet(ipAPI[0], 20*minTimeout))
	}
	return
}

func wGet(url string, timeout time.Duration) (str string) {
	request, err := http.NewRequest("GET", url, nil)
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
