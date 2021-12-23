package littleLogger

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

//TODO обдумать возможность рефакторинга под DRY
//TODO перерассмотреть реализацию с запуском гоуротин
//TODO обдумать реализацию под микросервисы (возможно создание lock файлов или подобное)
type Logger struct {
	target io.Writer
	mu sync.Mutex
	levels    [4]bool
	formatter func()string
}

//NewLogger lvls = {
//	1 - debug (false)
//	2 - warning (true) TODO Рассмотреть порядок еще раз
//	3 - info (true)
//	4 - error (true)
//}
//	formatter: func()string TODO проверить на деле удобность
//
//
func NewLogger(target io.Writer, lvls ...int) (*Logger, error) {
	if len(lvls) > 4 {
		return nil, errors.New("Count of lvls less than 5")
	}
	levels := [4]bool{false,true,true,true}
	for i, l := range lvls{
		levels[i] = l == 1
	}
	_, err := fmt.Fprint(target, "Starting logs")
	if err != nil {
		return nil, err
	}
	return &Logger{target: target, levels: levels, formatter: FormatterMinimal}, nil
}

func (l *Logger) Error(msg string) {
	if !l.levels[3]{
		return
	}
	go func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprint(l.target,
			strings.Replace(l.formatter(), "$msg",
				fmt.Sprintf("Error message: %s", msg), 1))
	}()

}

func (l *Logger) Info(msg string) {
	if !l.levels[2]{
		return
	}
	go func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprint(l.target,
			strings.Replace(l.formatter(), "$msg",
				fmt.Sprintf("Info message: %s", msg), 1))
	}()
}

func (l *Logger) Warning(msg string) {
	if !l.levels[1]{
		return
	}
	go func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprint(l.target,
			strings.Replace(l.formatter(), "$msg",
				fmt.Sprintf("Warning message: %s", msg), 1))
	}()
}

func (l *Logger) Debug(msg string) {
	if !l.levels[0]{
		return
	}
	go func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprint(l.target,
			strings.Replace(l.formatter(), "$msg",
				fmt.Sprintf("Debug message: %s", msg), 1))
	}()
}

func (l *Logger) SetFormatter(newFormatter func() string) {
	l.formatter = newFormatter
}