package jsonx

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrJson = errors.New("jsonx: error with json")
var ErrCheck = errors.New("jsonx: check failed")

type Checker interface {
	Check() error
}

func DecodeCheck[T Checker](r io.Reader) (T, error) {
	var item T
	err := json.NewDecoder(r).Decode(&item)
	if err != nil {
		return item, errors.Join(ErrJson, err)
	}
	err = item.Check()
	if err != nil {
		return item, errors.Join(ErrCheck, err)
	}
	return item, nil
}
func DecodeCheckFunc[T any](r io.Reader, checker func(T) error) (T, error) {
	var item T
	err := json.NewDecoder(r).Decode(&item)
	if err != nil {
		return item, errors.Join(ErrJson, err)
	}
	err = checker(item)
	if err != nil {
		return item, errors.Join(ErrCheck, err)
	}
	return item, nil
}
