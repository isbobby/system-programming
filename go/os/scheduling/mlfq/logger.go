package main

import "fmt"

type actor string
type action string

var (
	P  actor = "CPU"
	IO actor = "IO"
	S  actor = "Scheduler"

	EXEC     action = "executing job instruction"
	EXPIRE   action = "expire job with no more time allotment"
	SWAP     action = "swap job out for IO"
	COMPLETE action = "complete"

	JobIDKey = "Job ID"
)

type logLevel int

var (
	infoLevel logLevel = 1
	warnLevel logLevel = 2
	errLevel  logLevel = 3
)

type AuditLogger struct {
	SystemTime   *Clock
	Verbose      bool
	SystemOutput []AuditLog
}

type AuditLog struct {
	Action action
	Time   int
	JobID  string
}

func (al AuditLog) String() string {
	return fmt.Sprintf("%v %v", al.JobID, al.Action)
}

func (l *AuditLogger) CPUAuditLog(action action, keyValuePairs ...interface{}) {
	l.printLog(P, infoLevel, string(action), keyValuePairs...)
}

func (l *AuditLogger) CPUWarnLog(message string, keyValuePairs ...interface{}) {
	l.printLog(P, warnLevel, message, keyValuePairs...)
}

func (l *AuditLogger) CPUErrLog(message string, keyValuePairs ...interface{}) {
	l.printLog(P, errLevel, message, keyValuePairs...)
}

func (l *AuditLogger) IOLog(message string, keyValuePairs ...interface{}) {
	l.printLog(IO, infoLevel, message, keyValuePairs...)
}

func (l *AuditLogger) MLFQLog(message string, keyValuePairs ...interface{}) {
	l.printLog(S, infoLevel, message, keyValuePairs...)
}

func (l *AuditLogger) printLog(actor actor, level logLevel, message string, keyValuePairs ...interface{}) {
	kvString := []byte{}

	logTime := l.SystemTime.Time.Load()

	jobID := ""
	if len(keyValuePairs)%2 == 0 {
		for i := 0; i < len(keyValuePairs); i += 2 {
			key, value := keyValuePairs[i], keyValuePairs[i+1]

			kvString = append(kvString, []byte(fmt.Sprintf("%v:%v ", key, value))...)

			if key == JobIDKey {
				jobID = fmt.Sprintf("%v", value)
			}
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

	if actor == IO && l.Verbose {
		fmt.Printf("\033[0;33m[T:%v][%v] %v %v\n", logTime, actor, message, string(kvString))
	} else if actor == S && l.Verbose {
		fmt.Printf("\033[0;34m[T:%v][%v] %v %v\n", logTime, actor, message, string(kvString))
	} else {
		if l.Verbose {
			fmt.Printf("\033[0;32m[T:%v][%v] %v %v\n", logTime, actor, message, string(kvString))
		}

		if jobID != "" {
			l.SystemOutput = append(l.SystemOutput, AuditLog{
				Action: action(message),
				Time:   int(logTime),
				JobID:  jobID,
			})
		}
	}
}
