package order_book

import (
	"fmt"

	"github.com/gokch/snum"
	"github.com/google/btree"
)

//---------------------------------------------------------------------------------------------------//
// OrderBook

type OrderBook struct {
	preciseType TD_PreciseType

	precise_max TD_Precise
	precise_min TD_Precise

	map_precise map[TD_Precise]*Precise
}

func (t *OrderBook) Init(_stepType TD_PreciseType, _step_max TD_Precise, _step_min TD_Precise) error {
	t.preciseType = _stepType
	t.precise_max = _step_max
	t.precise_min = _step_min
	t.map_precise = make(map[TD_Precise]*Precise)

	if _stepType == TD_preciseType_dynamic {
		if _step_min <= 0 {
			return fmt.Errorf("minus td_step is not allow in dynamic type")
		}
	}

	for idx := t.precise_max; idx >= t.precise_min; idx-- {
		t.map_precise[idx] = &Precise{}
		t.map_precise[idx].Init(_stepType, idx)
	}
	return nil
}

func (t *OrderBook) Update(_orderWay TD_OrderWay, _price *snum.Snum, _amount *snum.Snum) error {
	for _, precise := range t.map_precise {
		err := precise.Update(_orderWay, _price, _amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *OrderBook) Get_step(_precise TD_Precise) *Precise {
	precise_one, is_exist := t.map_precise[_precise]
	if is_exist != true {
		return nil
	}
	return precise_one
}

//---------------------------------------------------------------------------------------------------//
// Precise

type Precise struct {
	preciseType TD_PreciseType
	precise     TD_Precise

	map_group_buy  map[string]*Group
	map_group_sell map[string]*Group

	tree_group_buy  *btree.BTree
	tree_group_sell *btree.BTree
}

func (t *Precise) Init(_preciseType TD_PreciseType, _precise TD_Precise) error {
	t.preciseType = _preciseType
	t.precise = _precise

	t.map_group_buy = make(map[string]*Group)
	t.map_group_sell = make(map[string]*Group)
	t.tree_group_buy = btree.New(3)
	t.tree_group_sell = btree.New(3)
	return nil
}

func (t *Precise) Update(_orderWay TD_OrderWay, _price *snum.Snum, _amount *snum.Snum) error {
	price := (*_price).Copy()
	amount := (*_amount).Copy()
	sn_price := price.String()

	var group *Group
	if _orderWay == TD_orderWay_buy {
		if t.preciseType == TD_preciseType_dynamic {
			(*price).Round_up(int(t.precise))
		} else {
			(*price).Group_up(int(t.precise))
		}
		sn_price = (*price).String()
		group = t.map_group_buy[sn_price]
	} else if _orderWay == TD_orderWay_sell {
		if t.preciseType == TD_preciseType_dynamic {
			(*price).Round_down(int(t.precise))
		} else {
			(*price).Group_down(int(t.precise))
		}
		sn_price = (*price).String()
		group = t.map_group_sell[sn_price]
	}

	if group == nil { // 새로운 봉 생성
		group = &Group{}
		group.init(price, amount)

		if _orderWay == TD_orderWay_buy {
			t.tree_group_buy.ReplaceOrInsert(group)
			t.map_group_buy[sn_price] = group
		} else {
			t.tree_group_sell.ReplaceOrInsert(group)
			t.map_group_sell[sn_price] = group
		}
	} else { // 기존 봉 수량 업데이트
		group.change_amount_left(amount)

		if group.is_amount_empty() == true { // 수량이 비었을 경우 봉 삭제
			if _orderWay == TD_orderWay_buy {
				t.tree_group_buy.Delete(group)
				delete(t.map_group_buy, sn_price)
			} else {
				t.tree_group_sell.Delete(group)
				delete(t.map_group_sell, sn_price)
			}
		}
	}
	return nil
}

func (t *Precise) Print() {
	fmt.Printf("buy ( %v )\n", t.tree_group_buy.Len())
	t.tree_group_buy.Ascend(func(_i_item btree.Item) bool {
		pt_group := _i_item.(*Group)
		fmt.Printf("%10v : %10v\n", (*pt_group.Price).String(), (*pt_group.Amount).String())
		return true
	})

	fmt.Printf("sell ( %v )\n", t.tree_group_sell.Len())
	t.tree_group_sell.Ascend(func(_i_item btree.Item) bool {
		pt_group := _i_item.(*Group)
		fmt.Printf("%10v : %10v\n", (*pt_group.Price).String(), (*pt_group.Amount).String())
		return true
	})

}

//---------------------------------------------------------------------------------------------------//
// group

type Group struct {
	Price  *snum.Snum
	Amount *snum.Snum
}

func (t *Group) init(_price, _amount *snum.Snum) {
	t.Price = _price
	t.Amount = _amount
}

func (t *Group) change_amount_left(_amount *snum.Snum) {
	t.Amount.Add(_amount)
}

func (t *Group) is_amount_empty() bool {
	zero := &snum.Snum{}
	zero.Set__str("0")
	if (*t.Amount).Cmp(zero) <= 0 {
		return true
	}
	return false
}

func (t *Group) Less(_item btree.Item) bool {
	group := _item.(*Group)

	// order way 가 같을 경우 price 비교
	cmp := (*t.Price).Cmp(group.Price)
	if cmp == -1 {
		return true
	}
	return false
}
