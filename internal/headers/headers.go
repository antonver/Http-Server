package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

var appropriate_characters string = "abcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"


func create_validate_map() map[rune]bool{
	valid_map := make(map[rune]bool)
	for _, i := range appropriate_characters{
		valid_map[i] = true
}
return valid_map
}



func validate(data string) bool {
	valid_map := create_validate_map()
	for _, i := range data{
		if !valid_map[i]{
			return false
		}
	}
	return true
}


func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	sData := string(data)
	

	crlfIdx := strings.Index(sData, "\r\n")
	if crlfIdx == -1 {
		return 0, false, nil
	}

	if crlfIdx == 0 {
		return 2, true, nil
	}

	line := sData[:crlfIdx]
	totalConsumed := crlfIdx + 2

	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		return 0, false, errors.New("invalid header format: missing colon")
	}

	key := strings.ToLower(line[:colonIdx])
	if !validate(key){
			return 0, false, errors.New("All the characters must be from this list: abcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~")
		}
	if strings.Contains(key, " ") || strings.HasPrefix(key, " ") || strings.HasSuffix(key, " ") {
		return 0, false, errors.New("invalid spacing in header key")
	}

	value := strings.TrimSpace(line[colonIdx+1:])
	
	_, ok := h[key]
	if ok{
		h[key] = h[key] + "," + value
 	}else{
	h[key] = value
	}
	return totalConsumed, false, nil
}

func (h Headers) Get(key string) string{
	v, ok := h[strings.ToLower(key)]
	if ok{
		return v
	}
	return ""
}
