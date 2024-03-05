package orderbook

type OrderWay uint8

const (
	Buy OrderWay = iota + 1
	Sell
)

type PreciseType uint8

const (
	dynamic PreciseType = iota
	fixed
)

type Step int8
