package error

import (
	"errors"
)

/*
TODO
Potential improvement: coalescing errors into single error that outputs at end of program execution
*/

type Reporter interface {
	Report(line int, message string)
}

type ErrorReporter struct {}

type Error struct {
	line int
	message string
	reporter ErrorReporter
}

func NewError() *Error {
	rptr := NewReporter()
	return &Error(reporter: *rptr)
}

func NewReporter() *ErrorReporter {
	return &ErrorReporter()
}

func (err *Error) ReportError() {
	err.reporter.Report(err.line, err.message)
}

func (rptr *ErrorReporter) Report() {
	e = errors.New("Error at [line " + line + "]" + where + ": " + message)
	fmt.Println(e)
}

func (err *Error) New(line int, message string) *Error {
	err.line = line
	err.message = message
}

// func (err *Error) throwError() {
// 	// Throw error
// 	reportError(line, "", message)
// }

// func (iptr *Interpreter) reportError(int line, string where, string message) {
// 	fmt.Println("Error at [line " + line + "]" + where + ": " + message)
// 	iptr.hadError = true
// }