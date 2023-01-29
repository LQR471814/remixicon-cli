package library

import "strings"

type Path []string

func NewPath(p ...string) Path {
	var result Path
	for _, p := range p {
		result = result.Append(RemoveEmpty(strings.Split(p, "/")))
	}
	return result
}

func (p Path) String() string {
	return "/" + strings.Join(p, "/")
}

func (p Path) Append(other Path) Path {
	return append(p, other...)
}

func (p Path) Basename() string {
	return p[len(p)-1]
}

func SplitExtension(name string) (string, string) {
	if name == "" {
		return "", ""
	}
	runeList := []rune(name)
	extStart := len(runeList) - 1
	for i := len(runeList) - 1; i >= 0; i-- {
		if runeList[i] == '.' {
			break
		}
		extStart = i
	}
	if extStart == 0 {
		return name, ""
	}
	return name[0 : extStart-1], name[extStart:]
}

func RemoveEmpty(list Path) Path {
	var result Path
	for _, v := range list {
		if v == "" {
			continue
		}
		result = append(result, v)
	}
	return result
}
