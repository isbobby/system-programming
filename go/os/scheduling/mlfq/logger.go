package main

import "fmt"

type actor string

var (
	P  actor = "CPU"
	IO actor = "IO"
	S  actor = "Scheduler"
)

type logLevel int

var (
	infoLevel logLevel = 1
	warnLevel logLevel = 2
	errLevel  logLevel = 3
)

type Logger struct {
	SystemTime *Clock
}

func (l Logger) CPULog(message string, keyValuePairs ...interface{}) {
	l.printLog(P, infoLevel, message, keyValuePairs...)
}

func (l Logger) CPUWarnLog(message string, keyValuePairs ...interface{}) {
	l.printLog(P, warnLevel, message, keyValuePairs...)
}

func (l Logger) CPUErrLog(message string, keyValuePairs ...interface{}) {
	l.printLog(P, errLevel, message, keyValuePairs...)
}

func (l Logger) IOLog(message string, keyValuePairs ...interface{}) {
	l.printLog(IO, infoLevel, message, keyValuePairs...)
}

func (l Logger) MLFQLog(message string, keyValuePairs ...interface{}) {
	l.printLog(S, infoLevel, message, keyValuePairs...)
}

func (l Logger) printLog(actor actor, level logLevel, message string, keyValuePairs ...interface{}) {
	kvString := []byte{}

	if len(keyValuePairs)%2 == 0 {
		for i := 0; i < len(keyValuePairs); i += 2 {
			key, value := keyValuePairs[i], keyValuePairs[i+1]

			kvString = append(kvString, []byte(fmt.Sprintf("%v:%v ", key, value))...)
		}
	}

	// trim space
	if len(kvString) > 0 {
		kvString = kvString[:len(kvString)-1]
		head := []byte{byte('[')}
		head = append(head, kvString...)
		head = append(head, ']')
		kvString = head
	}

	if actor == IO {
		fmt.Printf("\033[0;33m[T:%v][%v] %v %v\n", l.SystemTime.Time.Load(), actor, message, string(kvString))
	} else if actor == S {
		fmt.Printf("\033[0;34m[T:%v][%v] %v %v\n", l.SystemTime.Time.Load(), actor, message, string(kvString))
	} else {
		fmt.Printf("\033[0;32m[T:%v][%v] %v %v\n", l.SystemTime.Time.Load(), actor, message, string(kvString))
	}
}
