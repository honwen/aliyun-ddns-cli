package main

import (
	"regexp"
	"testing"

	"github.com/honwen/golibs/cip"
	"github.com/honwen/golibs/domain"
	"github.com/ysmood/got"
)

func TestGetIPv4(t *testing.T) {
	funcs["myip"] = cip.MyIPv4
	ip4 := myip()
	got.T(t).True(regexp.MustCompile(cip.RegxIPv4).MatchString(ip4) || len(ip4) == 0)
}

func TestGetIPv6(t *testing.T) {
	funcs["myip"] = cip.MyIPv6
	ip6 := myip()
	got.T(t).True(regexp.MustCompile(cip.RegxIPv6).MatchString(ip6) || len(ip6) == 0)
}

func TestResloveIPv4(t *testing.T) {
	funcs["reslove"] = cip.ResloveIPv4
	got.T(t).Has([]string{"8.8.8.8", "8.8.4.4"}, reslove("dns.google"))
	got.T(t).Has([]string{"223.6.6.6", "223.5.5.5"}, reslove("dns.alidns.com"))
}

func TestResloveIPv6(t *testing.T) {
	funcs["reslove"] = cip.ResloveIPv6
	got.T(t).Has([]string{"2001:4860:4860::8844", "2001:4860:4860::8888"}, reslove("dns.google"))
	got.T(t).Has([]string{"2400:3200::1", "2400:3200:baba::1"}, reslove("dns.alidns.com"))
}

func TestSplitDomain001(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("a.example.com")

	got.T(t).Eq(rr, "a")
	got.T(t).Eq(domain, "example.com")
}

func TestSplitDomain002(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("example.com")

	got.T(t).Eq(rr, "@")
	got.T(t).Eq(domain, "example.com")
}

func TestSplitDomain003(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("*.example.com")

	got.T(t).Eq(rr, "*")
	got.T(t).Eq(domain, "example.com")
}

func TestSplitDomain004(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("*.a.example.com")

	got.T(t).Eq(rr, "*.a")
	got.T(t).Eq(domain, "example.com")
}

func TestSplitDomain005(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("*.b.a.example.com")

	got.T(t).Eq(rr, "*.b.a")
	got.T(t).Eq(domain, "example.com")
}
func TestSplitDomain006(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("a.example.co.kr")

	got.T(t).Eq(rr, "a")
	got.T(t).Eq(domain, "example.co.kr")
}

func TestSplitDomain007(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("*.a.example.co.kr")

	got.T(t).Eq(rr, "*.a")
	got.T(t).Eq(domain, "example.co.kr")
}

func TestSplitDomain008(t *testing.T) {
	rr, domain := domain.SplitDomainToRR("example.co.kr")

	got.T(t).Eq(rr, "@")
	got.T(t).Eq(domain, "example.co.kr")
}
