package cpst

import (
	"fmt"
	"math"
	"testing"
)

func testSingleNumber(number uint64) {
	s := NumberToChar(number)
	n := CharToNumber(s)
	fmt.Printf("testSingleNumber: %d->%s->%d\n", number, s, n)
}

func testSingleChar(char string) {
	n := CharToNumber(char)
	c := NumberToChar(n)
	fmt.Printf("testSingleChar: %s->%d->%s\n", char, n, c)
}

func TestNumberToChar(t *testing.T) {
	initCharArray()
	testSingleNumber(uint64(0))
	testSingleNumber(uint64(2333))
	testSingleNumber(uint64(114840))
	testSingleNumber(uint64(math.MaxInt32))
	testSingleChar("233333")
	testSingleChar("123456")
	testSingleChar("AAAAAA")
	testSingleChar("ZZZZZZ")
	testSingleChar("AAA3F3")
}
