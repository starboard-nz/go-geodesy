package utils_test

import (
	"fmt"

	"github.com/starboard-nz/orb"
	"github.com/starboard-nz/go-geodesy/utils"
)

func ExampleArea() {
	// +
	// |\
	// | \
	// |  \
	// +---+

	r := orb.Ring{{0, 0}, {3, 0}, {0, 4}, {0, 0}}
	a := utils.Area(r)

	fmt.Println(a)
	// Output:
	// 6
}

func ExampleDistance() {
	d := utils.Distance(orb.Point{0, 0}, orb.Point{3, 4})

	fmt.Println(d)
	// Output:
	// 5
}

