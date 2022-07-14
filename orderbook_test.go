package order_book

import (
	"fmt"
	"testing"

	"github.com/gokch/snum"
)

func TestOrderBook(_t *testing.T) {
	type T_order struct {
		order_way TD_OrderWay
		price     string
		amount    string
	}

	fn_test := func(_preciseType TD_PreciseType, _step TD_Precise, _order_input, _order_expect []*T_order) error {
		orderBook := &OrderBook{}
		orderBook.Init(_preciseType, _step, _step)

		for _, pt_input := range _order_input {
			price := &snum.Snum{}
			amount := &snum.Snum{}
			price.Set__str(pt_input.price)
			amount.Set__str(pt_input.amount)

			err := orderBook.Update(pt_input.order_way, price, amount)
			if err != nil {
				_t.Fatal(err)
			}
		}

		// 결과물
		orderBook_one := orderBook.Get_step(_step)
		for _, order_expect := range _order_expect {
			var pt_order_result *Group
			if order_expect.order_way == TD_orderWay_buy {
				pt_order_result = orderBook_one.map_group_buy[order_expect.price]
			} else {
				pt_order_result = orderBook_one.map_group_sell[order_expect.price]
			}
			if pt_order_result == nil {
				_t.Errorf("not exist price | way : %v | price : %s", order_expect.order_way, order_expect.price)
				return fmt.Errorf("error")
			}

			// amount 비교
			if (*pt_order_result.Amount).Cmp__str(order_expect.amount) != 0 {
				_t.Errorf("amount not equal | way : %v | price : %s | amount : %s | result amount : %s",
					order_expect.order_way, order_expect.price, order_expect.amount, (*pt_order_result.Amount).String())
				return fmt.Errorf("error")
			}
		}

		return nil
	}

	// fixed
	{
		fn_test(
			TD_preciseType_fixed,
			-6,
			[]*T_order{
				// buy 는 자리 올림 ( 0.000002 )
				{TD_orderWay_buy, "0.000001", "10000"},
				{TD_orderWay_buy, "0.0000011", "10000"},
				{TD_orderWay_buy, "0.0000015", "10000"},
				{TD_orderWay_buy, "0.0000019", "10000"},
				{TD_orderWay_buy, "0.000002", "10000"},
				// sell 은 자리 내림 ( 0.000001 )
				{TD_orderWay_sell, "0.000001", "10000"},
				{TD_orderWay_sell, "0.0000011", "10000"},
				{TD_orderWay_sell, "0.0000015", "10000"},
				{TD_orderWay_sell, "0.0000019", "10000"},
				{TD_orderWay_sell, "0.000002", "10000"},
			},
			[]*T_order{
				{TD_orderWay_buy, "0.000001", "10000"},
				{TD_orderWay_buy, "0.000002", "40000"},
				{TD_orderWay_sell, "0.000001", "40000"},
				{TD_orderWay_sell, "0.000002", "10000"},
			},
		)

		fn_test(
			TD_preciseType_fixed,
			-4,
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.000987", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.00987", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.11111", "10000"},
				{TD_orderWay_sell, "15.1111", "10000"},
				{TD_orderWay_sell, "16.111", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.001", "20000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.0099", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.1111", "20000"},
				{TD_orderWay_sell, "16.111", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
		)

		fn_test(
			TD_preciseType_fixed,
			-1,
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.000987", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.00987", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.11111", "10000"},
				{TD_orderWay_sell, "15.1111", "10000"},
				{TD_orderWay_sell, "16.111", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
			[]*T_order{
				{TD_orderWay_buy, "0.1", "90000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.3", "20000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "30000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.1", "20000"},
				{TD_orderWay_sell, "16.1", "10000"},
				{TD_orderWay_sell, "17.1", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
		)

		fn_test(
			TD_preciseType_fixed,
			1,
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.000987", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.00987", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.11111", "10000"},
				{TD_orderWay_sell, "15.1111", "10000"},
				{TD_orderWay_sell, "16.111", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
			[]*T_order{
				{TD_orderWay_buy, "10", "140000"},
				{TD_orderWay_buy, "20", "10000"},

				{TD_orderWay_sell, "0", "30000"},
				{TD_orderWay_sell, "10", "80000"},
				{TD_orderWay_sell, "80", "10000"},
				{TD_orderWay_sell, "880", "10000"},
				{TD_orderWay_sell, "8870", "20000"},
			},
		)

		fn_test(
			TD_preciseType_fixed,
			4,
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.000987", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.00987", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.11111", "10000"},
				{TD_orderWay_sell, "15.1111", "10000"},
				{TD_orderWay_sell, "16.111", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
			[]*T_order{
				{TD_orderWay_buy, "10000", "150000"},

				{TD_orderWay_sell, "0", "150000"},
			},
		)
	}

	// dynamic
	{
		fn_test(
			TD_preciseType_dynamic,
			4,
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.000987", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.00987", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.11111", "10000"},
				{TD_orderWay_sell, "15.1111", "10000"},
				{TD_orderWay_sell, "16.111", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.000987", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.00987", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.11", "20000"},
				{TD_orderWay_sell, "16.11", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "20000"},
			},
		)

		fn_test(
			TD_preciseType_dynamic,
			2,
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.000987", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.00987", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.0987", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.23", "10000"},
				{TD_orderWay_buy, "1.24", "10000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "10.1", "10000"},

				{TD_orderWay_sell, "1.2", "10000"},
				{TD_orderWay_sell, "1.23", "10000"},
				{TD_orderWay_sell, "1.24", "10000"},
				{TD_orderWay_sell, "10", "10000"},
				{TD_orderWay_sell, "10.1", "10000"},
				{TD_orderWay_sell, "15.11111", "10000"},
				{TD_orderWay_sell, "15.1111", "10000"},
				{TD_orderWay_sell, "16.111", "10000"},
				{TD_orderWay_sell, "17.11", "10000"},
				{TD_orderWay_sell, "18.1", "10000"},
				{TD_orderWay_sell, "18.9", "10000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "887", "10000"},
				{TD_orderWay_sell, "8876", "10000"},
				{TD_orderWay_sell, "8876.5", "10000"},
			},
			[]*T_order{
				{TD_orderWay_buy, "0.0009", "10000"},
				{TD_orderWay_buy, "0.00098", "10000"},
				{TD_orderWay_buy, "0.00099", "10000"},
				{TD_orderWay_buy, "0.009", "10000"},
				{TD_orderWay_buy, "0.0098", "10000"},
				{TD_orderWay_buy, "0.0099", "10000"},
				{TD_orderWay_buy, "0.098", "10000"},
				{TD_orderWay_buy, "0.099", "10000"},
				{TD_orderWay_buy, "0.09", "10000"},
				{TD_orderWay_buy, "0.9", "10000"},
				{TD_orderWay_buy, "1.2", "10000"},
				{TD_orderWay_buy, "1.3", "20000"},
				{TD_orderWay_buy, "10", "10000"},
				{TD_orderWay_buy, "11", "10000"},

				{TD_orderWay_sell, "1.2", "30000"},
				{TD_orderWay_sell, "10", "20000"},
				{TD_orderWay_sell, "15", "20000"},
				{TD_orderWay_sell, "16", "10000"},
				{TD_orderWay_sell, "17", "10000"},
				{TD_orderWay_sell, "18", "20000"},
				{TD_orderWay_sell, "88", "10000"},
				{TD_orderWay_sell, "880", "10000"},
				{TD_orderWay_sell, "8800", "20000"},
			},
		)
	}
}
