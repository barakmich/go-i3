package goi3

func BuildBlocks(gens []Generator, ctx Context) []Block {
	var out []Block
	for _, g := range gens {
		b, err := g.MakeBlock(ctx)
		if err != nil {
			out = append(out, ErrorBlock(err))
			continue
		}
		out = append(out, b)
	}
	return out
}

func ErrorBlock(err error) Block {
	return Block{
		FullText: err.Error(),
		Color:    RED,
	}
}

var RED Color = Color{255, 0, 0}
