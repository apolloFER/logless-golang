# LogLess
> LogLess is a centralized logging/event stack using AWS Kinesis/Lambda. The Golang version works as a Logrus hook.

## Usage

Create a new `LoglessHook` with the name of your Kinesis stream as the first parameter. The second parameter is the `logrus` log to use for internal logging. If you don't supply one (`nil`), LogLess will create a new one.

You should defer the `logless.OnClose()` method in your `main()` function since it's used for flushing outstanding entries to Kinesis.

## Example
```go
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/apolloFER/logless-golang"
)

func main() {
	hook, _ := logless.NewLoglessHook("logless-test", log.New())
	defer logless.OnClose()

	log.AddHook(hook)

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously 1!")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 123,
	}).Info("The group's number increased tremendously 2!")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 124,
	}).Error("The group's number increased tremendously 3!")
}
```

### License
Apache License 2.0
