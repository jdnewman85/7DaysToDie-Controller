package main

import (
	"regexp"
)

type Point struct {
	x float64
	y float64
	z float64
}

type stateFunc func()

type triggerFunc func(map[string]string)

type trigger struct {
	reString string
	re       *regexp.Regexp
	callback triggerFunc
}

type Player struct {
	id      string
	name    string
	steamId string

	lastUpdate int64
	online     bool
	remote     bool
	ip         string
	ping       uint64

	x, y, z float64 //TODO Point types?
	u, v, w float64

	//TEMP OPT Stop using 64 bits, and do the conversions where needed
	health uint64
	deaths uint64
	zKills uint64
	pKills uint64
	score  uint64
	level  uint64

	//Added stats
	blinkLocations []*Point
}
