package log

import (
	"os"

	"github.com/kovetskiy/lorg"
	structured "github.com/reconquest/cog"
)

type Logger struct {
	*structured.Logger
}

func New(debug bool, trace bool, traceFile string) *Logger {
	stderr := lorg.NewLog()
	stderr.SetIndentLines(true)
	stderr.SetFormat(
		lorg.NewFormat("${time} ${level:[%s]:right:short} ${prefix}%s"),
	)

	stderr.Infof("trace log file: %s", traceFile)

	if traceFile != "" {
		logfile, err := os.OpenFile(
			traceFile,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND,
			0666,
		)
		if err != nil {
			stderr.Fatalf(
				"unable to create log file: %s", err,
			)
		}

		if logfile.Name() != os.Stderr.Name() {
			output := lorg.NewOutput(logfile)

			output.SetLevelWriterCondition(lorg.LevelTrace, logfile)
			output.SetLevelWriterCondition(lorg.LevelDebug, logfile)
			output.SetLevelWriterCondition(lorg.LevelFatal, logfile, os.Stderr)
			output.SetLevelWriterCondition(lorg.LevelError, logfile, os.Stderr)
			output.SetLevelWriterCondition(lorg.LevelWarning, logfile, os.Stderr)
			output.SetLevelWriterCondition(lorg.LevelInfo, logfile, os.Stderr)

			stderr.SetOutput(output)
		}
	}

	if debug {
		stderr.SetLevel(lorg.LevelDebug)
	}

	if trace {
		stderr.SetLevel(lorg.LevelTrace)
	}

	return &Logger{structured.NewLogger(stderr)}
}

func (logger *Logger) TraceJSON(obj interface{}) (encoded string) {
	return logger.Logger.TraceJSON(obj)
}

func (logger *Logger) NewChild() *Logger {
	return &Logger{
		logger.Logger.NewChild(),
	}
}

func (logger *Logger) NewChildWithPrefix(prefix string) *Logger {
	return &Logger{
		logger.Logger.NewChildWithPrefix(prefix),
	}
}

func (logger *Logger) Println(v ...interface{}) {
	logger.Print(v...)
}