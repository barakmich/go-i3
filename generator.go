package goi3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
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

type PowerGenerator struct {
	ID string
}

func stringFile(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	b = bytes.TrimRight(b, "\n")
	return string(b), nil
}

func (p PowerGenerator) MakeBlock(_ Context) (Block, error) {
	path := filepath.Join("/sys/class/power_supply", p.ID)
	out, err := stringFile(filepath.Join(path, "status"))
	if err != nil {
		return Block{}, err
	}
	cur, err := stringFile(filepath.Join(path, "charge_now"))
	if err != nil {
		return Block{}, err
	}
	full, err := stringFile(filepath.Join(path, "charge_full"))
	if err != nil {
		return Block{}, err
	}
	c, err := strconv.Atoi(cur)
	if err != nil {
		return Block{}, err
	}
	f, err := strconv.Atoi(full)
	if err != nil {
		return Block{}, err
	}
	percent := (c * 100) / f
	out = fmt.Sprintf("%s %d%%", out, percent)
	return Block{
		FullText: out,
		Color:    Color{255, 255, 0},
	}, nil
}
