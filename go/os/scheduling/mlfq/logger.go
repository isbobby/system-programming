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

func CPULog(message string, keyValuePairs ...interface{}) {
	printLog(P, infoLevel, message, keyValuePairs...)
}

func CPUWarnLog(message string, keyValuePairs ...interface{}) {
	printLog(P, warnLevel, message, keyValuePairs...)
}

func CPUErrLog(message string, keyValuePairs ...interface{}) {
	printLog(P, errLevel, message, keyValuePairs...)
}

func IOLog(message string, keyValuePairs ...interface{}) {
	printLog(IO, infoLevel, message, keyValuePairs...)
}

func MLFQLog(message string, keyValuePairs ...interface{}) {
	printLog(S, infoLevel, message, keyValuePairs...)
}

func printLog(actor actor, level logLevel, message string, keyValuePairs ...interface{}) {
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
		fmt.Printf("\033[0;33m[%v] %v %v\n", actor, message, string(kvString))
	} else if actor == S {
		fmt.Printf("\033[0;34m[%v] %v %v\n", actor, message, string(kvString))
	} else {
		fmt.Printf("\033[0;32m[%v] %v %v\n", actor, message, string(kvString))
	}
}
