package main

import "fmt"

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
