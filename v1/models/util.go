package models

func cutLastChar(s string) string {
	c := []byte(s)
	
	Len := len(s)
	if c[Len-1] == '_' {
		c = c[:Len-1]
	}
	
	return string(c)
}