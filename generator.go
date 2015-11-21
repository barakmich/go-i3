package goi3

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
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
			metadata: map[string]string{"down": "true"},
		}, nil
	}
	ip := strings.Split(addrs[0].String(), "/")[0]
	return Block{
		FullText: fmt.Sprintf("%s: %s", i.IFace, ip),
		Color:    Color{0, 255, 0},
		metadata: map[string]string{"down": "false"},
	}, nil
}

type WifiGenerator struct {
	iface *IFaceGenerator
}

func NewWifiGenerator(name string) *WifiGenerator {
	return &WifiGenerator{
		iface: &IFaceGenerator{
			IFace: name,
		},
	}
}

func (w WifiGenerator) MakeBlock(ctx Context) (Block, error) {
	blk, err := w.iface.MakeBlock(ctx)
	if err != nil {
		return blk, err
	}
	var down bool
	if blk.metadata["down"] == "true" {
		down = true
	}
	if down {
		return blk, nil
	}
	cmd := exec.Command("iwconfig", w.iface.IFace)
	data, err := cmd.Output()
	if err != nil {
		return blk, err
	}
	re, err := regexp.Compile("ESSID:\"(.*?)\"")
	if err != nil {
		return blk, err
	}
	s := re.FindAllStringSubmatch(string(data), -1)
	ssid := "(no SSID)"
	if s != nil && len(s) > 0 && len(s[0]) > 0 {
		ssid = s[0][1]
	}
	blk.FullText = fmt.Sprintf("%s @ %s", blk.FullText, ssid)
	return blk, nil
}
