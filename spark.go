package goi3

import "fmt"

// Len = 8
var spark = []rune("▁▂▃▄▅▆▇█")

type SparkGenerator struct {
	min, max, history int
	fun               IntTick
	data              []int
	dynamic           bool
	PrintBounds       bool
	PrintCurrent      bool
}

func NewConstantSparkGenerator(min, max, history int, fun IntTick) *SparkGenerator {
	return &SparkGenerator{
		min:          min,
		max:          max,
		history:      history,
		fun:          fun,
		data:         make([]int, history),
		PrintCurrent: true,
	}
}

func NewDynamicSparkGenerator(history int, fun IntTick) *SparkGenerator {
	return &SparkGenerator{
		min:     0,
		max:     0,
		history: history,
		fun:     fun,
		data:    make([]int, history),
		dynamic: true,
	}
}

func (c *SparkGenerator) MakeBlock(ctx Context) (Block, error) {
	newint := c.fun.Tick(ctx)
	c.data = append(c.data, newint)
	c.data = c.data[1:]
	var out []rune
	if c.dynamic {
		c.max = c.data[0]
		c.min = c.data[0]
		for _, x := range c.data {
			if x > c.max {
				c.max = x
			}
			if x < c.min {
				c.min = x
			}
		}
	}
	segment := ((c.max - c.min) << 8) / 7
	for _, i := range c.data {
		if i < c.min {
			out = append(out, spark[0])

		} else if i >= c.max {
			out = append(out, spark[7])
		} else {
			f := i - c.min
			f = f << 8
			f = f / segment
			out = append(out, spark[f])
		}
	}
	sparkstring := string(out)
	if c.PrintBounds {
		sparkstring = fmt.Sprintf("%s %d %d", sparkstring, c.min, c.max)
	}
	if c.PrintCurrent {
		sparkstring = fmt.Sprintf("%s %d", sparkstring, newint)
	}
	return Block{
		FullText: sparkstring,
		Color:    Color{64, 128, 255},
	}, nil
}

type CountTick struct {
	i int
}

func (t *CountTick) Tick(ctx Context) int {
	t.i++
	if t.i > 100 {
		t.i = 0
	}
	return t.i
}
