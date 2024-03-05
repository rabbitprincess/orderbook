package orderbook

import (
	"fmt"

	"github.com/gokch/snum_sort/snum"
	"github.com/google/btree"
)

//---------------------------------------------------------------------------------------------------//
// OrderBook

type OrderBook struct {
	precise PreciseType

	max Step
	min Step

	precises map[Step]*Precise
}

func (t *OrderBook) Init(pt PreciseType, max Step, min Step) error {
	t.precise = pt
	t.max = max
	t.min = min
	t.precises = make(map[Step]*Precise)

	if pt == dynamic {
		if min <= 0 {
			return fmt.Errorf("minus step is not allow in dynamic type")
		}
	}

	for idx := t.max; idx >= t.min; idx-- {
		t.precises[idx] = &Precise{}
		t.precises[idx].Init(pt, idx)
	}
	return nil
}

func (t *OrderBook) Update(ow OrderWay, p *snum.Snum, a *snum.Snum) error {
	for _, precise := range t.precises {
		err := precise.Update(ow, p, a)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *OrderBook) Get(st Step) *Precise {
	if one, exist := t.precises[st]; exist != true {
		return nil
	} else {
		return one
	}
}

//---------------------------------------------------------------------------------------------------//
// Precise

type Precise struct {
	pt PreciseType
	st Step

	mapBuy  map[string]*Group
	mapSell map[string]*Group

	treeBuy  *btree.BTree
	treeSell *btree.BTree
}

func (t *Precise) Init(pt PreciseType, st Step) error {
	t.pt = pt
	t.st = st

	t.mapBuy = make(map[string]*Group)
	t.mapSell = make(map[string]*Group)
	t.treeBuy = btree.New(3)
	t.treeSell = btree.New(3)
	return nil
}

func (t *Precise) Update(ow OrderWay, price *snum.Snum, amount *snum.Snum) error {
	price = (*price).Copy()
	amount = (*amount).Copy()
	sPrice := price.String()

	var group *Group
	if ow == Buy {
		if t.pt == dynamic {
			(*price).RoundUp(int(t.st))
		} else {
			(*price).GroupUp(int(t.st))
		}
		sPrice = (*price).String()
		group = t.mapBuy[sPrice]
	} else if ow == Sell {
		if t.pt == dynamic {
			(*price).RoundDown(int(t.st))
		} else {
			(*price).GroupDown(int(t.st))
		}
		sPrice = (*price).String()
		group = t.mapSell[sPrice]
	}

	if group == nil { // 새로운 봉 생성
		group = &Group{}
		group.init(price, amount)

		if ow == Buy {
			t.treeBuy.ReplaceOrInsert(group)
			t.mapBuy[sPrice] = group
		} else {
			t.treeSell.ReplaceOrInsert(group)
			t.mapSell[sPrice] = group
		}
	} else { // 기존 봉 수량 업데이트
		group.Change(amount)

		if group.IsEmpty() == true { // 수량이 비었을 경우 봉 삭제
			if ow == Buy {
				t.treeBuy.Delete(group)
				delete(t.mapBuy, sPrice)
			} else {
				t.treeSell.Delete(group)
				delete(t.mapSell, sPrice)
			}
		}
	}
	return nil
}

func (t *Precise) Print() {
	fmt.Printf("buy ( %v )\n", t.treeBuy.Len())
	t.treeBuy.Ascend(func(item btree.Item) bool {
		group := item.(*Group)
		fmt.Printf("%10v : %10v\n", (*group.Price).String(), (*group.Amount).String())
		return true
	})

	fmt.Printf("sell ( %v )\n", t.treeSell.Len())
	t.treeSell.Ascend(func(item btree.Item) bool {
		group := item.(*Group)
		fmt.Printf("%10v : %10v\n", (*group.Price).String(), (*group.Amount).String())
		return true
	})

}

//---------------------------------------------------------------------------------------------------//
// group

type Group struct {
	Price  *snum.Snum
	Amount *snum.Snum
}

func (t *Group) init(price, amount *snum.Snum) {
	t.Price = price
	t.Amount = amount
}

func (t *Group) Change(amount *snum.Snum) {
	t.Amount.Add(amount)
}

func (t *Group) IsEmpty() bool {
	zero := &snum.Snum{}
	zero.SetStr("0")
	if (*t.Amount).Cmp(zero) <= 0 {
		return true
	}
	return false
}

func (t *Group) Less(item btree.Item) bool {
	group := item.(*Group)

	// order way 가 같을 경우 price 비교
	cmp := (*t.Price).Cmp(group.Price)
	if cmp == -1 {
		return true
	}
	return false
}
