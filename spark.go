package goi3

import "fmt"

// Len = 8
var spark = []rune("▁▂▃▄▅▆▇█")

type SparkGenerator struct {
	min, max, history int
	fun               IntTick
	data              []int
	dynamic           bool
}

func NewConstantSparkGenerator(min, max, history int, fun IntTick) *SparkGenerator {
	return &SparkGenerator{
		min:     min,
		max:     max,
		history: history,
		fun:     fun,
		data:    make([]int, history),
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
	return Block{
		FullText: sparkstring,
	}, nil
}

type CPUGenerator struct {
	spark *SparkGenerator
	color Color
}

func NewCPUGenerator(color Color) *CPUGenerator {
	return &CPUGenerator{
		spark: NewConstantSparkGenerator(0, 99, 10, NewCPUTick()),
		color: color,
	}
}

func (cpu *CPUGenerator) MakeBlock(ctx Context) (Block, error) {
	blk, err := cpu.spark.MakeBlock(ctx)
	if err != nil {
		return blk, err
	}
	blk.Color = cpu.color
	blk.FullText = fmt.Sprintf("%s %02d", blk.FullText, cpu.spark.data[len(cpu.spark.data)-1])
	return blk, nil
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
