package pg

import (
	"bytes"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/timer"
	"math"
	"os"
	"sync"
	"time"
)

var (
	a, b, c, d, e rune
)

type Pos byte

type pattern struct {
	title  string
	left   string
	finish string
	cur    string
	blank  string
	right  string
}

const (
	Top  Pos = 1
	Left Pos = 2
)

var (
	defaultPattern = pattern{
		title:  "",
		left:   "[",
		finish: "-",
		cur:    ">",
		blank:  " ",
		right:  "]",
	}
)

type progress struct {
	Title    string
	TitlePos Pos
	Value    int
	MaxValue int
	Width    int
	buffer   bytes.Buffer
	last     int
	pat      pattern
	lock     *sync.Mutex
	shine    bool
	timeLeft float64
	lastTime int64
}

func New(maxValue int, width int, title string, titlePos Pos) *progress {
	p := &progress{
		MaxValue: maxValue,
		Width:    width,
		Title:    title,
		TitlePos: titlePos,
		buffer:   bytes.Buffer{},
		pat:      defaultPattern,
		lock:     new(sync.Mutex),
		lastTime: gox.GetTimestamp(time.Now()),
	}
	if titlePos == Top {
		fmt.Println(title)
	}

	timer.Start(0, 0, time.Millisecond*500, func(t *timer.Timer) {
		p.render(t)
	})
	return p
}

func (p *progress) Update(value int) {
	defer p.render(nil)
	p.Value += value
}

func (p *progress) Increase() {
	defer p.render(nil)
	if p.Value >= p.MaxValue {
		return
	}
	p.Value += 1
}

func (p *progress) render(t *timer.Timer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	defer func() {
		p.buffer.Reset()
		p.last = p.Value
		p.lastTime = gox.GetTimestamp(time.Now())
		if t != nil {
			p.shine = !p.shine
		}
		if p.Value >= p.MaxValue && t != nil {
			t.Destroy()
			fmt.Println("\n\n结束")
		}
	}()
	finish := int(math.Floor(float64(p.Value) / float64(p.MaxValue) * float64(p.Width-2)))
	for i := 0; i < finish; i++ {
		p.buffer.WriteString(p.pat.finish)
	}
	fs := p.buffer.String()
	p.buffer.Reset()
	for i := 0; i < p.Width-2-finish-1; i++ {
		p.buffer.WriteString(p.pat.blank)
	}
	bs := p.buffer.String()
	percent := float64(p.Value) / float64(p.MaxValue) * 100
	if p.Value > p.last {
		p.timeLeft = math.Ceil(float64(p.MaxValue-p.Value) / float64(p.Value-p.last) * float64(time.Millisecond*500) / float64(time.Second))
	}

	fmt.Fprintf(os.Stdout, "\r%s%s%s%s%s%s%s%.2f%%%s",
		gox.TValue(p.TitlePos == Left, p.Title+" ", ""),
		p.pat.left,
		fs,
		// gox.TValue(p.Value == p.MaxValue, "", gox.TValue(!p.shine, " ", p.pat.cur)),
		gox.TValue(p.Value == p.MaxValue, "", p.pat.cur),
		bs,
		p.pat.right,
		gox.TValue(percent < 10, "  ", gox.TValue(percent < 100, " ", "")),
		percent,
		gox.TValue(p.Value == p.MaxValue, "", fmt.Sprintf(" | %.fs", p.timeLeft)))
}
