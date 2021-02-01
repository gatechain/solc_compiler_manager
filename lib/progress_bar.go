package lib

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Bar struct {
	mu       sync.Mutex
	line     int
	prefix   string
	total    int
	width    int
	advance  chan bool
	done     chan bool
	currents map[string]int
	current  int
	before   int
	rate     int
	speed    int
	cost     int
	estimate int
	fast     int
	slow     int
}

var (
	bar1 string
	bar2 string
)

const (
	defaultFast = 20
	defaultSlow = 5
)

func initBar(width int) {
	for i := 0; i < width; i++ {
		bar1 += "="
		bar2 += "-"
	}
}

func NewBar(line int, prefix string, total int) *Bar {
	if total <= 0 {
		return nil
	}

	if line <= 0 {
		gMaxLine++
		line = gMaxLine
	}

	bar := &Bar{
		line:     line,
		prefix:   prefix,
		total:    total,
		fast:     defaultFast,
		slow:     defaultSlow,
		width:    100,
		advance:  make(chan bool),
		done:     make(chan bool),
		currents: make(map[string]int),
	}

	initBar(bar.width)
	go bar.updateCost()
	go bar.run()

	return bar
}

func (b *Bar) SetSpeedSection(fast, slow int) {
	if fast > slow {
		b.fast, b.slow = fast, slow
	} else {
		b.fast, b.slow = slow, fast
	}
}

func (b *Bar) Add(n ...int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	step := 1
	if len(n) > 0 {
		step = n[0]
	}

	b.current += step

	lastRate := b.rate
	lastSpeed := b.speed

	b.count()

	if lastRate != b.rate || lastSpeed != b.speed {
		b.advance <- true
	}

	if b.rate >= 100 {
		close(b.done)
		close(b.advance)
	}
}

func (b *Bar) Set(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.current = n

	lastRate := b.rate
	lastSpeed := b.speed

	b.count()

	if lastRate != b.rate || lastSpeed != b.speed {
		b.advance <- true
	}

	if b.rate >= 100 {
		close(b.done)
		close(b.advance)
	}
}

func (b *Bar) count() {
	now := time.Now()
	nowKey := now.Format("20060102150405")
	befKey := now.Add(time.Minute * -1).Format("20060102150405")
	b.currents[nowKey] = b.current
	if v, ok := b.currents[befKey]; ok {
		b.before = v
	}
	delete(b.currents, befKey)

	b.rate = b.current * 100 / b.total
	if b.cost == 0 {
		b.speed = b.current * 100
	} else if b.before == 0 {
		b.speed = b.current * 100 / b.cost
	} else {
		b.speed = (b.current - b.before) * 100 / 60
	}

	if b.speed != 0 {
		b.estimate = (b.total - b.current) * 100 / b.speed
	}
}

func (b *Bar) updateCost() {
	for {
		select {
		case <-time.After(time.Second):
			b.cost++
			b.mu.Lock()
			b.count()
			b.mu.Unlock()
			b.advance <- true
		case <-b.done:
			return
		}
	}
}

func (b *Bar) run() {
	for range b.advance {
		printf(b.line, "\r%s", b.barMsg())
	}
}

func (b *Bar) barMsg() string {
	prefix := fmt.Sprintf("%s", b.prefix)
	rate := fmt.Sprintf("%3d%%", b.rate)
	speed := fmt.Sprintf("%3.2fps", 0.01*float64(b.speed))
	cost := b.timeFmt(b.cost)
	estimate := b.timeFmt(b.estimate)
	ct := fmt.Sprintf(" (%d/%d)", b.current, b.total)
	barLen := b.width - len(prefix) - len(rate) - len(speed) - len(cost) - len(estimate) - len(ct) - 10
	bar1Len := barLen * b.rate / 100
	bar2Len := barLen - bar1Len

	realBar1 := bar1[:bar1Len]
	var realBar2 string
	if bar2Len > 0 {
		realBar2 = ">" + bar2[:bar2Len-1]
	}

	msg := fmt.Sprintf(`%s %s%s [%s%s] %s %s in: %s`, prefix, rate, ct, realBar1, realBar2, speed, cost, estimate)
	switch {
	case b.speed <= b.slow*100:
		return "\033[0;31m" + msg + "\033[0m"
	case b.speed > b.slow*100 && b.speed < b.fast*100:
		return "\033[0;33m" + msg + "\033[0m"
	default:
		return "\033[0;32m" + msg + "\033[0m"
	}
}

func (b *Bar) timeFmt(cost int) string {
	var h, m, s int
	h = cost / 3600
	m = (cost - h*3600) / 60
	s = cost - h*3600 - m*60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

type Pgbar struct {
	title   string
	maxLine int
	bars    []*Bar
}

func NewMultiProgressBar(title string) *Pgbar {

	p := &Pgbar{
		title:   title,
		maxLine: gMaxLine,
	}

	gMaxLine++
	return p
}

func (p *Pgbar) NewBar(prefix string, total int) *Bar {
	gMaxLine++
	return NewBar(gMaxLine, prefix, total)
}

var (
	mu           sync.Mutex
	gSrcLine     = 0 //起点行
	gCurrentLine = 0 //当前行
	gMaxLine     = 0 //最大行
)

func move(line int) {
	//fmt.Println("\n\n\n\n", gCurrentLine, line)
	fmt.Printf("\033[%dA\033[%dB", gCurrentLine, line)
	gCurrentLine = line
}

func print(line int, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	move(line)
	var realArgs []interface{}
	realArgs = append(realArgs, "\r")
	realArgs = append(realArgs, args...)
	fmt.Print(realArgs...)
	move(gMaxLine)
}

func printf(line int, format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	move(line)
	fmt.Printf("\r"+format, args...)
	move(gMaxLine)
}

func println(line int, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	move(line)
	var realArgs []interface{}
	realArgs = append(realArgs, "\r")
	realArgs = append(realArgs, args...)
	fmt.Print(realArgs...)
	move(gMaxLine)
}

func Print(args ...interface{}) {
	mu.Lock()
	lf := countLF("", args...)
	if gMaxLine == 0 {
		gMaxLine += lf + 1
	} else {
		gMaxLine += lf
	}
	mu.Unlock()

	print(gMaxLine, args...)
}

func Printf(format string, args ...interface{}) {
	mu.Lock()

	lf := countLF(format, args...)
	if gMaxLine == 0 {
		gMaxLine += lf + 1
	} else {
		gMaxLine += lf
	}
	mu.Unlock()

	printf(gMaxLine, format, args...)
}

func Println(args ...interface{}) {
	mu.Lock()

	lf := countLF("", args...)
	lf++
	if gMaxLine == 0 {
		gMaxLine += lf + 1
	} else {
		gMaxLine += lf
	}
	mu.Unlock()

	println(gMaxLine, args...)
}

func countLF(format string, args ...interface{}) int {
	var count int
	count = strings.Count(format, "\n")
	for _, arg := range args {
		count += strings.Count(String(arg), "\n")
	}

	return count
}

//func ShowProgress(path string, i chan int64, end chan bool) {
//	b := NewBar(path, "", 100)
//	bar.NewOption(0, 100)
//	for {
//		if <-end {
//			break
//		}
//		bar.Play(<-i)
//	}
//	bar.Finish()
//}