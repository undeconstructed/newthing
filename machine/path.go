package machine

import "strings"

func splitPath(p string) []string {
	path := []string{}
	for _, s := range strings.Split(p, "/") {
		if s == "" {
			continue
		}
		path = append(path, s)
	}
	return path
}

type pathPiece struct {
	t int
	v string
}

type pathex struct {
	pieces []pathPiece
}

func parsePathex(s string) pathex {
	pieces := []pathPiece{}
	for _, p := range strings.Split(s, "/") {
		if p == "" {
			continue
		}
		t := 0
		if p[0] == '{' && p[len(p)-1] == '}' {
			t = 1
			p = p[1 : len(p)-1]
		}
		pieces = append(pieces, pathPiece{
			t: t,
			v: p,
		})
	}
	return pathex{pieces}
}

type matchType int

func (p *pathex) match(path []string) (map[string]string, []string, bool) {
	if len(p.pieces) > len(path) {
		return nil, path, false
	}
	vars := map[string]string{}
	for i, p := range p.pieces {
		s := path[i]
		if p.t == 0 {
			if p.v != s {
				return nil, path, false
			}
		} else {
			vars[p.v] = s
		}
	}
	return vars, path[len(p.pieces):], true
}
