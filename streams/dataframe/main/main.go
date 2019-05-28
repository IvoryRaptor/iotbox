package main

import (
	"fmt"
	"github.com/google/btree"
)

type Int struct {
	Value int
}

func (i *Int) Less(than btree.Item) bool {
	return i.Value < than.(*Int).Value
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func main() {
	tr := btree.New(3)
	for i := 0; i < 10; i++ {
		tr.ReplaceOrInsert(&Int{Value: i})
	}
	fmt.Println("len:       ", tr.Len())
	fmt.Println("get3:      ", tr.Get(&Int{3}))
	fmt.Println("get100:    ", tr.Get(&Int{100}))
	fmt.Println("del4:      ", tr.Delete(&Int{4}))
	fmt.Println("del100:    ", tr.Delete(&Int{100}))
	fmt.Println("replace5:  ", tr.ReplaceOrInsert(&Int{5}))
	fmt.Println("replace100:", tr.ReplaceOrInsert(&Int{100}))
	tr.Max()
	fmt.Println("min:       ", tr.Min())
	fmt.Println("delmin:    ", tr.DeleteMin())
	fmt.Println("max:       ", tr.Max())
	fmt.Println("delmax:    ", tr.DeleteMax())
	fmt.Println("len:       ", tr.Len())
}
