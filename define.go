package order_book

type TD_OrderWay uint8

const (
	TD_orderWay_buy TD_OrderWay = iota + 1
	TD_orderWay_sell
)

type TD_PreciseType uint8

const (
	TD_preciseType_dynamic TD_PreciseType = iota
	TD_preciseType_fixed
)

type TD_Precise int8
