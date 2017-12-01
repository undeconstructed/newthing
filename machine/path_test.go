package machine

import (
	"testing"
)

func TestPathExRoot(t *testing.T) {
	p := parsePathex("/")
	_, r, m := p.match(splitPath("/"))
	if !m {
		t.Errorf("not match")
	}
	if len(r) != 0 {
		t.Errorf("bad rest: %v", r)
	}
}

func TestPathExPart(t *testing.T) {
	p := parsePathex("/foo/{bar}")
	_, r, m := p.match(splitPath("/foo/bat/charm"))
	if !m {
		t.Errorf("not match")
	}
	if len(r) != 1 {
		t.Errorf("bad rest: %v", r)
	}
}

func TestPathExMiss(t *testing.T) {
	p := parsePathex("/foo/bar")
	_, r, m := p.match(splitPath("/foo/bat"))
	if m {
		t.Errorf("match")
	}
	if len(r) != 2 {
		t.Errorf("bad rest: %v", r)
	}
}

func TestPathExFull(t *testing.T) {
	p := parsePathex("/foo/bar")
	_, r, m := p.match(splitPath("/foo/bar"))
	if !m {
		t.Errorf("not match")
	}
	if len(r) != 0 {
		t.Errorf("bad rest: %v", r)
	}
}
