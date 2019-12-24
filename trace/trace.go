package trace

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"sync"
	"context"
	"strconv"
	"runtime"
	"io/ioutil"
)

const (
	_traceHead string = "zddgo-http-header-tid"	//trace id
	_spandHead string = "zddgo-http-header-sid"	//span id
	TraceContextKey string = "zddgo-trace" //context key
)

func GetTraceHeadKey() string{
	return _traceHead
}
func GetSpandHeadKey() string{
	return _spandHead
}

var (
	_processId int = 0
)


func init(){
	_processId = os.Getpid()
}

func NewTrace() (t *Trace){
	t = &Trace{}
	t.SetPid(_processId)
	t.remoteAddr = "127.0.0.1"
	t.host = "127.0.0.1"
	t.traceId = generateTraceId(time.Now(), t.remoteAddr, t.pid)
	t.spanId = "0"
	return
}

//
//inherit trace info from http.Request
//
func InheritHttpTrace(r *http.Request) (t *Trace){
	t = &Trace{}
	t.SetPid(_processId)
	if r.RemoteAddr != "" {
		t.remoteAddr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	if t.remoteAddr == "::1" {
		t.remoteAddr = "127.0.0.1"
	}
	if len(r.Host) > 0 {
		t.host = strings.Split(r.Host, ":")[0]
	}
	t.traceId = r.Header.Get(_traceHead)
	if t.traceId == "" {
		t.traceId = generateTraceId(time.Now(), t.remoteAddr, t.pid)
	}
	t.spanId = r.Header.Get(_spandHead)
	if t.spanId == "" {
		t.spanId = "0"
	}
	return
}

func InheritContextTrace(ctx context.Context) *Trace{
	if v := ctx.Value(TraceContextKey); v != nil {
        if trace, ok := v.(*Trace); ok {
			return trace
		}
    }
    return nil
}

type Trace struct {
	mu sync.Mutex
	traceId string
	spanId string
	remoteAddr string
	host string
	pid int
	rpcId int
}
func(t *Trace) SetPid(pid int){
	t.mu.Lock()
	defer t.mu.Unlock()
	t.pid = pid
}
func(t *Trace) IncrRpc(){
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rpcId = t.rpcId + 1
	t.spanId = t.spanId + "." + strconv.Itoa(t.rpcId)
}
func (t *Trace) Get() (string, string, string, string, int, int){
	return t.traceId, t.spanId, t.remoteAddr, t.host, t.pid, t.rpcId
}

//
//route生成规则: 0~7 ip地址 8~15 id生成时间 16~19 生成id的nginx启动时间 20~23 生成的nginx进程号 24~29 循环自增序列 30~31 02 固定
//
func generateTraceId(t time.Time, ip string, pid int) string{
	b := bytes.Buffer{}
	b.WriteString(hex.EncodeToString(net.ParseIP(ip).To4()))
	b.WriteString(fmt.Sprintf("%x", uint32(t.Unix())&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", t.UnixNano()&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0")
	return b.String()
}



var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack returns a nicely formatted stack frame, skipping skip frames.
func Stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}