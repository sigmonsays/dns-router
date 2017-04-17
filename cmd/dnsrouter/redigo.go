package main

import (
	"github.com/garyburd/redigo/redis"
)

func NewRedigo() *Redigo {
	return &Redigo{
		StringMap: redis.StringMap,
	}
}

type Redigo struct {
	StringMap interface{}
}
