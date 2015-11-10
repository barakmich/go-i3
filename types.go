package goi3

import "fmt"

type Block struct {
	FullText  string `json:"full_text,omitempty"`
	ShortText string `json:"short_text,omitempty"`
	Color     Color  `json:"color,omitempty"`
}

type IntTick interface {
	Tick(ctx Context) int
}

type Context struct{}

type Generator interface {
	MakeBlock(ctx Context) (Block, error)
}

type Color struct {
	R, G, B uint8
}

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func (c Color) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", c, "\"")), nil
}
