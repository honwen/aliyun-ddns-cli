package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitDomain001(t *testing.T) {
	rr, domain := splitDomain("a.domain.tld")

	assert.Equal(t, rr, "a")
	assert.Equal(t, domain, "domain.tld")
}

func TestSplitDomain002(t *testing.T) {
	rr, domain := splitDomain("domain.tld")

	assert.Equal(t, rr, "@")
	assert.Equal(t, domain, "domain.tld")
}

func TestSplitDomain003(t *testing.T) {
	rr, domain := splitDomain("*.domain.tld")

	assert.Equal(t, rr, "*")
	assert.Equal(t, domain, "domain.tld")
}

func TestSplitDomain004(t *testing.T) {
	rr, domain := splitDomain("*.a.domain.tld")

	assert.Equal(t, rr, "*.a")
	assert.Equal(t, domain, "domain.tld")
}

func TestSplitDomain005(t *testing.T) {
	rr, domain := splitDomain("*.b.a.domain.tld")

	assert.Equal(t, rr, "*.b.a")
	assert.Equal(t, domain, "domain.tld")
}
