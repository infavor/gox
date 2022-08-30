package disk_test

import (
	"bytes"
	"container/list"
	"github.com/gorilla/mux"
	"github.com/infavor/gox"
	"github.com/infavor/gox/disk"
	"github.com/infavor/gox/file"
	"github.com/infavor/gox/logger"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

var out io.WriteCloser
var writeLock *sync.Mutex
var g = sync.WaitGroup{}
var switchBuffer *disk.SwitchBuffer

func init() {
	logger.Init(nil)
	out, _ = file.AppendFile("test_append.txt")
	writeLock = new(sync.Mutex)
	g.Add(5)
	switchBuffer = disk.NewSwitchBuffer(func(items *list.List) error {
		var buffer bytes.Buffer
		gox.WalkList(items, func(item interface{}) bool {
			entry := item.(*disk.Entry)
			buffer.WriteString(entry.Source.(string))
			return false
		})
		_, err := out.Write(buffer.Bytes())
		return err
	})
	go switchBuffer.Schedule()
}

// Test A
func TestSwitchBufferWrite(t *testing.T) {
	r := mux.NewRouter()
	srv := &http.Server{
		Handler: r,
		Addr:    ":8088",
		// Good practice: enforce timeouts for servers you create!
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      0,
		ReadTimeout:       0,
		MaxHeaderBytes:    1 << 20, // 1MB
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()
	logger.Info("write test started")
	for i := 0; i < 1; i++ {
		go SwitchBufferWriteWork()
	}
	logger.Info("waiting write to end")
	g.Wait()
	logger.Info("write end")
}

// Test B 2277ms
func TestPlainWrite(t *testing.T) {
	logger.Info("write test started")
	start := gox.GetTimestamp(time.Now())
	for i := 0; i < 5; i++ {
		go PlainWriteWork()
	}
	logger.Info("waiting write to end")
	g.Wait()
	logger.Info("write end")
	end := gox.GetTimestamp(time.Now())
	logger.Info("time duration: ", end-start, "ms")
}

func PlainWriteWork() {
	for i := 0; i < 100000; i++ {
		PlainWriteTask("ppqwevqwevq9wven129u3-1931-2n9wur9we9rvn9-1n23719vm0x\n")
	}
	g.Done()
}

func PlainWriteTask(line string) error {
	writeLock.Lock()
	defer writeLock.Unlock()
	_, err := out.Write([]byte(line))
	return err
}

func SwitchBufferWriteWork() {
	for i := 0; i < 100000; i++ {
		err := <-switchBuffer.Push("ppqwevqwevq9wven129u3-1931-2n9wur9we9rvn9-1n23719vm0x\n")
		if err != nil {
			logger.Fatal("write failed:", err)
		}
		logger.Error(err)
	}
	g.Done()
}
