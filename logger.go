package onelogin

import (
    "os"
    "github.com/op/go-logging"
)

// Logger
var logger = logging.MustGetLogger("onelogin")

func init() {
    StderrFormatter("")
}

func LogFormat()(logging.Formatter) {
    return logging.MustStringFormatter(
        `%{shortfile} %{shortfunc} ▶ %{level:.8s} %{message}`,
    )
}

func StderrBackend(prefix string)(logging.Backend) {
    return logging.NewLogBackend(os.Stderr, prefix, 0)
}

func StderrFormatter(prefix string) {
    logging.SetFormatter(LogFormat())
}
