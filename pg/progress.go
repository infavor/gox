package pg

import (
	"bytes"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/timer"
	"io"
	"math"
	"os"
	"strings"
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
		finish: "=",
		cur:    ">",
		blank:  " ",
		right:  "]",
	}
)

type progress struct {
	Title    string
	TitlePos Pos
	Value    int64
	MaxValue int64
	Width    int
	buffer   bytes.Buffer
	last     int64
	pat      pattern
	lock     *sync.Mutex
	shine    bool
	timeLeft int64
	lastTime int64
	reader   *WrappedReader
	writer   *WrappedWriter
	timer    *timer.Timer
}

type WrappedWriter struct {
	Writer io.Writer
	p      *progress
}

func (w *WrappedWriter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	w.p.Update(int64(n))
	return
}

type WrappedReader struct {
	Reader io.Reader
	p      *progress
}

func (w *WrappedReader) Read(p []byte) (n int, err error) {
	n, err = w.Reader.Read(p)
	w.p.Update(int64(n))
	return
}

func New(maxValue int64, width int, title string, titlePos Pos) *progress {
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

	p.timer = timer.Start(0, 0, time.Millisecond*500, func(t *timer.Timer) {
		p.render()
	})
	return p
}

func NewWrappedReaderProgress(maxValue int64, width int, title string, titlePos Pos, reader *WrappedReader) *progress {
	p := &progress{
		MaxValue: maxValue,
		Width:    width,
		Title:    title,
		TitlePos: titlePos,
		buffer:   bytes.Buffer{},
		pat:      defaultPattern,
		lock:     new(sync.Mutex),
		lastTime: gox.GetTimestamp(time.Now()),
		reader:   reader,
	}
	reader.p = p
	if titlePos == Top {
		fmt.Println(title)
	}

	p.timer = timer.Start(0, 0, time.Millisecond*500, func(t *timer.Timer) {
		p.render()
	})
	return p
}

func NewWrappedWriterProgress(maxValue int64, width int, title string, titlePos Pos, writer *WrappedWriter) *progress {
	p := &progress{
		MaxValue: maxValue,
		Width:    width,
		Title:    title,
		TitlePos: titlePos,
		buffer:   bytes.Buffer{},
		pat:      defaultPattern,
		lock:     new(sync.Mutex),
		lastTime: gox.GetTimestamp(time.Now()),
		writer:   writer,
	}
	writer.p = p
	if titlePos == Top {
		fmt.Println(title)
	}

	p.timer = timer.Start(0, 0, time.Millisecond*500, func(t *timer.Timer) {
		p.render()
	})
	return p
}

func (p *progress) Update(value int64) {
	p.lock.Lock()
	defer func() {
		p.lock.Unlock()
		if p.Value >= p.MaxValue {
			p.render()
		}
	}()
	p.Value += value
}

func (p *progress) Increase() {
	p.lock.Lock()
	defer func() {
		p.lock.Unlock()
		if p.Value >= p.MaxValue {
			p.render()
		}
	}()
	if p.Value >= p.MaxValue {
		return
	}
	p.Value += 1
}

func (p *progress) render() {
	p.lock.Lock()
	defer p.lock.Unlock()
	defer func() {
		p.buffer.Reset()
		p.last = p.Value
		p.lastTime = gox.GetTimestamp(time.Now())
		if p.Value >= p.MaxValue {
			logger.Debug("stop progress")
			p.timer.Destroy()
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
		p.timeLeft = int64(math.Ceil(float64(p.MaxValue-p.Value) / float64(p.Value-p.last) * float64(time.Millisecond*500) / float64(time.Second)))
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
		gox.TValue(p.Value == p.MaxValue, "\n", fmt.Sprintf(" | %s", HumanReadableTime(p.timeLeft))))
}

func HumanReadableTime(second int64) string {
	if second < 60 {
		return strings.Join([]string{"        ", convertBlank(second), "s"}, "")
	}
	if second < 3600 {
		return strings.Join([]string{"    ", convertBlank(second / 60), "m", convertBlank(second % 60), "s"}, "")
	}
	if second < 86400 {
		return strings.Join([]string{convertBlank(second / 3600), "h", convertBlank(second % 3600 / 60), "m", convertBlank(second % 3600 % 60), "s"}, "")
	}
	return "1d+"
}

func convertBlank(v int64) string {
	return gox.TValue(v < 10, " "+convert.Int64ToStr(v), convert.Int64ToStr(v)).(string)
}
