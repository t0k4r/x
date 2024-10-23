package jsonx

import (
	"encoding/json"
	"io"
)

type Item interface {
	Validate() error
}

func Unmarshall[T Item](bytes []byte, item *T) error {
	if err := json.Unmarshal(bytes, item); err != nil {
		return err
	}
	return (*item).Validate()
}

func UnmarshallRead[T Item](r io.Reader, item *T) error {
	if err := json.NewDecoder(r).Decode(item); err != nil {
		return err
	}
	return (*item).Validate()
}

func Read[T Item](r io.Reader) (item T, err error) {
	err = UnmarshallRead(r, &item)
	return item, err
}
