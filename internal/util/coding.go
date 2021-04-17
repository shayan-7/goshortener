package util

import (
	"time"

	"github.com/speps/go-hashids"
)

func GetHashID() string {
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	id, _ := h.Encode([]int{int(now.Unix())})
	return id
}
