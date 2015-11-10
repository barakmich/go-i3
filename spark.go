package goi3

// Len = 8
var spark = []rune("▁▂▃▄▅▆▇█")

type ConstantSparkGenerator struct {
	min, max, history int
	fun               IntTick
	data              []int
}

func NewConstantSparkGenerator(min, max, history int, fun IntTick) *ConstantSparkGenerator {
	return &ConstantSparkGenerator{
		min, max, history, fun, make([]int, history),
	}
}

func (c *ConstantSparkGenerator) MakeBlock(ctx Context) (Block, error) {
	newint := c.fun.Tick(ctx)
	c.data = append(c.data, newint)
	c.data = c.data[1:]
	var out []rune
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
	return Block{
		FullText: string(out),
		Color:    Color{0, 0, 255},
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
