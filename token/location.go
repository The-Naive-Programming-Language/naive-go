package token

import "strconv"

type Location struct {
	FileName string
	Line     int
	Column   int
}

func (loc Location) String() (s string) {
	if len(loc.FileName) == 0 {
		s = "<unknown>"
	} else {
		s = loc.FileName
	}
	s += ":" + strconv.Itoa(loc.Line) + ":" + strconv.Itoa(loc.Column)
	return s
}
