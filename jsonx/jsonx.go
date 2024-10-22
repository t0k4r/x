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

type Decoder[T Item] struct {
	*json.Decoder
}

func NewDecoder[T Item](r io.Reader) *Decoder[T] {
	return &Decoder[T]{Decoder: json.NewDecoder(r)}
}
func (d *Decoder[T]) Decode(item *T) error {
	if err := d.Decoder.Decode(item); err != nil {
		return err
	}
	return (*item).Validate()
}

func Read[T Item](r io.Reader) (T, error) {
	var item T
	err := NewDecoder[T](r).Decode(&item)
	return item, err
}
