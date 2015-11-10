package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/barakmich/go-i3"
)

var config []goi3.Generator = []goi3.Generator{
	goi3.StaticGenerator{
		"Baba Booey",
		goi3.Color{0, 255, 0},
	},
	goi3.NewConstantSparkGenerator(
		0, 20, 10, &goi3.CountTick{},
	),
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
