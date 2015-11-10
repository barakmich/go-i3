package goi3

type StaticGenerator struct {
	Text  string
	Color Color
}

func (s StaticGenerator) MakeBlock(ctx Context) (Block, error) {
	return Block{
		FullText: s.Text,
		Color:    s.Color,
	}, nil
}
