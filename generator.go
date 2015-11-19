package goi3

import (
	"fmt"
	"net"
	"strings"
	"time"
)

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

type TimeGenerator struct {
	Format string
	Color  Color
}

func (t TimeGenerator) MakeBlock(ctx Context) (Block, error) {
	now := time.Now()
	return Block{
		FullText: now.Format(t.Format),
		Color:    t.Color,
	}, nil
}

type IFaceGenerator struct {
	IFace string
}

func (i IFaceGenerator) MakeBlock(ctx Context) (Block, error) {
	iface, err := net.InterfaceByName(i.IFace)
	if err != nil {
		return Block{}, err
	}
	addrs, err := iface.Addrs()
	if len(addrs) == 0 {
		return Block{
			FullText: fmt.Sprintf("%s: down", i.IFace),
			Color:    Color{255, 0, 0},
		}, nil
	}
	ip := strings.Split(addrs[0].String(), "/")[0]
	return Block{
		FullText: fmt.Sprintf("%s: %s", i.IFace, ip),
		Color:    Color{0, 255, 0},
	}, nil
}
