package sparse

import (
	"fmt"
	"github.com/bits-and-blooms/bitset"
	"testing"
)

func Test(t *testing.T) {
	var x = bitset.New(3)
	x.Set(0)
	x.Set(1)
	x.Set(2)
	index := 0
	for index, ok := x.NextSet(uint(index)); ok; index, ok = x.NextSet(uint(index + 1)) {
		fmt.Println(index)
	}

}
