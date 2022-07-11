package example

import (
	"testing"

	"github.com/hikitani/anyhash"
)

func TestExample(t *testing.T) {
	type Bar struct {
		a int
		b [2]bool
	}

	type Foo struct {
		str    string
		number int
		slice  []Bar
	}

	fooHasher, err := anyhash.New[Foo](0)
	if err != nil {
		panic(err)
	}

	f1 := Foo{
		str:    "str",
		number: 1234,
		slice: []Bar{
			{
				a: 4321,
				b: [2]bool{true, false},
			},
		},
	}

	println(fooHasher.GetHash(f1))
}
