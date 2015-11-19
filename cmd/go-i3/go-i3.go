package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/barakmich/go-i3"
)

var config []goi3.Generator = []goi3.Generator{
	//goi3.StaticGenerator{
	//"Baba Booey 😁",
	//goi3.Color{0, 255, 0},
	//},
	goi3.IFaceGenerator{"wlan0"},
	goi3.IFaceGenerator{"lo"},
	goi3.NewCPUGenerator(goi3.Color{64, 128, 255}),
	goi3.TimeGenerator{
		"2006-01-02 15:04:05",
		goi3.Color{255, 255, 255},
	},
}

func main() {
	fmt.Println("{\"version\": 1}")
	fmt.Println("[")
	// Read config
	for {
		blocks := goi3.BuildBlocks(config, goi3.Context{})
		bytes, err := json.Marshal(blocks)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))
		fmt.Println(",")
		time.Sleep(1000 * time.Millisecond)
	}
}
