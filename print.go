package print

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"sync"
)

var (
	printDebugStatements    bool
	printDebugStatementsMux sync.RWMutex
	printGoroutines         bool
	printGoroutinesMux      sync.RWMutex
	gids                    []uint64
	gidsMux                 sync.Mutex
)

func init() {
	DebugStatements(true)
	Goroutines(true)
	gids = []uint64{}
}

// DebugStatements sets whether or not to print debug statements.
func DebugStatements(b bool) {
	printDebugStatementsMux.Lock()
	printDebugStatements = b
	printDebugStatementsMux.Unlock()
}

// Goroutines turns on/off goroutine printing markup.
func Goroutines(b bool) {
	printGoroutinesMux.Lock()
	printGoroutines = b
	printGoroutinesMux.Unlock()
}

func getPrintGoroutines() bool {
	printGoroutinesMux.RLock()
	defer printGoroutinesMux.RUnlock()
	return printGoroutines
}

func getPrintDebugStatements() bool {
	printDebugStatementsMux.RLock()
	defer printDebugStatementsMux.RUnlock()
	return printDebugStatements
}

// Msg prints the arguments as they are.
func Msg(args ...interface{}) {
	gidsMux.Lock()
	defer gidsMux.Unlock()
	fmt.Print(tabs())
	for _, a := range args {
		fmt.Print(a)
	}
	fmt.Println()
}

// Debug adds comment-like togglable debugging statements at runtime.
func Debug(stmt string, stmts ...string) {
	if !getPrintDebugStatements() {
		return
	}
	gidsMux.Lock()
	defer gidsMux.Unlock()
	tabs := tabs()
	p(tabs)
	p(tabs, "  | ", stmt)
	for _, stmt := range stmts {
		p(tabs, "  | ", stmt)
	}
	p(tabs, "  v")
}

func p(args ...interface{}) {
	for _, a := range args {
		fmt.Print(a)
	}
	fmt.Println()
}

func tabs() string {
	tabs := ""

	if !getPrintGoroutines() {
		return tabs
	}

	// find gid (CAUTION; anti-pattern)
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	gid, _ := strconv.ParseUint(string(b), 10, 64)

	index := -1
	for i, id := range gids {
		if id == gid {
			index = i
		}
	}
	if index >= 0 { // found
		tabs = strings.Repeat("\t", index)
	} else { // not found; new channel
		index = len(gids)
		gids = append(gids, gid)
		tabs = strings.Repeat("\t", index)
		p()
		p(tabs, "Go| ", index+1)
		p(tabs, "|||")
	}

	return tabs
}
