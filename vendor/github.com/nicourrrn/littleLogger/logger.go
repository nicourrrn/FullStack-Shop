package littleLogger

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

//TODO перерассмотреть реализацию с запуском гоуротин
//TODO обдумать реализацию под микросервисы (возможно создание lock файлов или подобное)
type Logger struct {
	target    io.Writer
	mu        sync.Mutex
	levels    [4]bool
	formatter func() string
}

//NewLogger lvls = {
//	1 - debug (false)
//	2 - warning (true) TODO Рассмотреть порядок еще раз
//	3 - info (true)
//	4 - error (true)
//}
//	formatter: func()string
//
//
func NewLogger(target io.Writer, lvls ...bool) (*Logger, error) {
	if len(lvls) > 4 {
		return nil, errors.New("Count of lvls less than 5")
	}
	levels := [4]bool{false, true, true, true}
	for i, l := range lvls {
		levels[i] = l
	}
	_, err := fmt.Fprintln(target, "Starting logs")
	if err != nil {
		return nil, err
	}
	return &Logger{target: target, levels: levels, formatter: FormatterMinimal}, nil
}

func (l *Logger) log(typeLog, msg string) {
	defer l.mu.Unlock()
	fmt.Fprint(l.target,
		strings.Replace(l.formatter(), "$msg",
			fmt.Sprintf("%s%s\n", typeLog, msg), -1))
}

func (l *Logger) Error(msg string) {
	if l.levels[3] {
		l.mu.Lock()
		go l.log("Error message: ", msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.levels[2] {
		l.mu.Lock()
		go l.log("Info message: ", msg)
	}
}

func (l *Logger) Warning(msg string) {
	if l.levels[1] {
		l.mu.Lock()
		go l.log("Warning message: ", msg)
	}
}

func (l *Logger) Debug(msg string) {
	if l.levels[0] {
		l.mu.Lock()
		go l.log("Debug message: ", msg)
	}
}

//SetFormatter must contain $msg substring
//Example: fmt.Sprintf("%s: $msg", time.Now().String())
func (l *Logger) SetFormatter(newFormatter func() string) {
	l.mu.Lock()
	l.formatter = newFormatter
	l.mu.Unlock()
}

func (l *Logger) Wait() {
	time.Sleep(time.Millisecond * 100)
	waiter := make(chan int) // TODO переделать тип канала
	go func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		waiter <- 1
	}()
	select {
	case <-waiter:
		return
	case <-time.After(time.Second * 5):
		fmt.Println("Log file is blocked")
	}
}
