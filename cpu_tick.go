package goi3

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type CPUTick struct {
	total    int
	idle     int
	lastTime time.Time
}

func NewCPUTick() *CPUTick {
	c := &CPUTick{}
	t, i, at := c.getCounts()
	c.total = t
	c.idle = i
	c.lastTime = at
	return c
}

func (c *CPUTick) getCounts() (int, int, time.Time) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, time.Now()
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	ok := s.Scan()
	if !ok {
		err := s.Err()
		if err != nil {
			return 0, 0, time.Now()
		}
	}
	fields := strings.Split(s.Text(), " ")
	var total int
	var idle int
	for i, x := range fields[1:5] {
		v, _ := strconv.Atoi(x)
		if i == 3 {
			idle = v
		} else {
			total += v
		}
	}
	return total, idle, time.Now()
}

func (c *CPUTick) Tick(ctx Context) int {
	t, i, at := c.getCounts()
	secs := at.Sub(c.lastTime).Seconds()
	tpercent := int((float64(t-c.total) / float64(runtime.NumCPU())) / secs)
	c.total = t
	c.idle = i
	c.lastTime = at
	return tpercent
}
