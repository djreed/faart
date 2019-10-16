package log

import (
	"io"
	"log"
	"os"
)

const LOG_FLAGS = log.Ldate | log.Ltime | log.Lmicroseconds

var ERR = MakeLogger(os.Stderr)
var OUT = MakeLogger(os.Stdout)

func MakeLogger(out io.Writer) *log.Logger {
	return log.New(out, "", LOG_FLAGS)
}
