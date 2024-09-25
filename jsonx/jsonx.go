package jsonx

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrCheck = errors.New("jsonx: check failed")

type Checker interface {
	Check() error
}

func DecodeCheck[T Checker](r io.Reader) (item T, err error) {
	if err = json.NewDecoder(r).Decode(&item); err != nil {
		return item, err
	}
	if err = item.Check(); err != nil {
		return item, errors.Join(ErrCheck, err)
	}
	return item, nil
}
func DecodeCheckFunc[T any](r io.Reader, checker func(T) error) (item T, err error) {
	if err = json.NewDecoder(r).Decode(&item); err != nil {
		return item, err
	}
	if err = checker(item); err != nil {
		return item, errors.Join(ErrCheck, err)
	}
	return item, nil
}
