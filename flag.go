package main

import (
	"fmt"
	"net/http"
	"strings"
)

type stringsVar struct {
	Values []string
}

// Set is flag.Value.Set
func (fv *stringsVar) Set(v string) error {
	fv.Values = append(fv.Values, v)
	return nil
}

func (fv *stringsVar) String() string {
	return fmt.Sprint(fv.Values)
}

type headersVar struct {
	Values map[string]string
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *headersVar) Help() string {
	separator := ":"
	return fmt.Sprintf("a HTTP header KEY%sVALUE", separator)
}

// Set is flag.Value.Set
func (fv *headersVar) Set(v string) error {
	separator := ":"
	i := strings.Index(v, separator)
	if i < 0 {
		return fmt.Errorf(`"%s" must have the form KEY%sVALUE`, v, separator)
	}
	fv.Texts = append(fv.Texts, v)
	if fv.Values == nil {
		fv.Values = make(map[string]string)
	}
	key, value := v[:i], v[i+len(separator):]
	key = strings.TrimSpace(key)
	key = http.CanonicalHeaderKey(key)
	value = strings.TrimLeft(value, " ")
	fv.Values[key] = value
	return nil
}

func (fv *headersVar) String() string {
	return fmt.Sprint(fv.Texts)
}
