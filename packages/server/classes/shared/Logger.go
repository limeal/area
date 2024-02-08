package shared

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/uuid"
)

// `Logger` is a type that has a `id` field of type `uuid.UUID`, a `File` field of type `*os.File`, and
// a `Writer` field of type `*bufio.Writer`.
// @property id - A unique identifier for the logger.
// @property File - This is the file that the logger will write to.
// @property Writer - A buffered writer that writes to the file.
type Logger struct {
	id     uuid.UUID
	File   *os.File
	Writer *bufio.Writer
}

// It creates a new file in the logs directory, and returns a pointer to a Logger struct that contains
// a file handle and a buffered writer
func NewLogger(id uuid.UUID) *Logger {

	derr := os.Mkdir("logs", 0755)
	if derr != nil && !os.IsExist(derr) {
		return nil
	}
	path := "logs/" + id.String() + ".log"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	datawriter := bufio.NewWriter(file)
	if datawriter == nil {
		file.Close()
		return nil
	}

	return &Logger{
		id:     id,
		File:   file,
		Writer: datawriter,
	}
}

// Closing the file that the logger is writing to.
func (l *Logger) Close() error {
	return l.File.Close()
}

// Writing an error message to the log file.
func (l *Logger) WriteError(message string) error {
	formatMsg := "[ERROR - " + l.id.String() + "]: " + message + "\n"
	_, err := l.Writer.WriteString(formatMsg)
	if err != nil {
		return err
	}
	fmt.Println(formatMsg)
	l.Writer.Flush()
	return nil
}

// Writing an info message to the log file.
func (l *Logger) WriteInfo(message string, t bool) error {
	formatMsg := "[INFO - " + l.id.String() + "]:" + message + "\n"
	_, err := l.Writer.WriteString(formatMsg)
	if err != nil {
		return err
	}
	if t {
		fmt.Println(formatMsg)
	}
	l.Writer.Flush()
	return nil
}
