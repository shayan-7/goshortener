package models

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis"
)

type Item struct {
	ID  string `json:"id"  validate:"required"`
	URL string `json:"url" validate:"required"`
}

var ErrItemNotFound = errors.New("Item not found")

func FindItem(id string, r *redis.Client) (*Item, error) {
	keys := r.Keys(id).Val()
	if len(keys) == 0 {
		return nil, ErrItemNotFound
	}
	url := r.Get(id).Val()
	return &Item{ID: id, URL: url}, nil
}

func AddItem(i *Item, r *redis.Client) {
	r.Set(i.ID, i.URL, 0)
	// itemList = append(itemList, i)
}

func (i *Item) Validate() error {
	validate := validator.New()
	return validate.Struct(i)
}

func (i *Item) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(i)
}

// FromJSON takes an io.Reader and converts the content of Reader to Item
// value if it's convertable
func (i *Item) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(i)
}
