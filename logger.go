package onelogin

import (
    "os"
    "github.com/op/go-logging"
)

const LoggerModule = "onelogin"

// Logger
var logger = logging.MustGetLogger(LoggerModule)

func init() {
    SetLogLevel(logging.WARNING)
}

func SetLogLevel(level logging.Level) {
    logging.SetFormatter(LogFormat())
    backend := logging.AddModuleLevel(StderrBackend(""))
    backend.SetLevel(level, "")
    logger.SetBackend(backend)
}

func LogFormat()(logging.Formatter) {
    return logging.MustStringFormatter(
        `%{level:.8s} %{shortfile} %{shortfunc} ▶  %{message}`,
    )
}

func StderrBackend(prefix string)(logging.Backend) {
    return logging.NewLogBackend(os.Stderr, prefix, 0)
}

func StderrFormatter(prefix string) {
    logging.SetFormatter(LogFormat())
}
