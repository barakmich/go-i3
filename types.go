package goi3

import "fmt"

type Block struct {
	FullText  string    `json:"full_text,omitempty"`
	ShortText string    `json:"short_text,omitempty"`
	Color     Color     `json:"color,omitempty"`
	MinWidth  int       `json:"min_width,omitempty"`
	Align     Alignment `json:"align,omitempty"`
	Name      string    `json:"name,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	metadata  map[string]string
}

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
)

func (a Alignment) String() string {
	switch a {
	case AlignLeft:
		return "left"
	case AlignRight:
		return "right"
	case AlignCenter:
		return "center"
	}
	panic("unknown alignment")
}

func (a Alignment) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", a, "\"")), nil
}

type IntTick interface {
	Tick(ctx Context) int
}

type Context struct {
	NetworkIsConnected bool
}

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
