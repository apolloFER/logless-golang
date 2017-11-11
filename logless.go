package logless

import (
	"github.com/Sirupsen/logrus"
	"github.com/a8m/kinesis-producer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/vmihailenco/msgpack"
	"github.com/satori/go.uuid"

	"os"
	"time"
)

var hooks []*LoglessHook
var hostname string

type LoglessHook struct {
	producer *producer.Producer
}

type LogEntry struct {
	Message string `msgpack:"message"`
	Level string `msgpack:"level"`
	Timestamp string `msgpack:"timestamp"`
	Source string `msgpack:"source"`
	Fields map[string]interface{} `msgpack:"fields"`
}

func init() {
	hst, err := os.Hostname()

	if err == nil {
		hostname = hst
	}
}

func NewLoglessHook(streamName string, log *logrus.Logger) (*LoglessHook, error) {
	s, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, err
	}

	if log == nil {
		log = logrus.New()
	}

	client := kinesis.New(s)
	pr := producer.New(&producer.Config{
		StreamName:   streamName,
		FlushInterval: time.Second,
		BacklogCount: 128,
		Client:       client,
		Logger:       log,
	})

	go processFailure(pr)

	pr.Start()

	hook :=  &LoglessHook{producer:pr}

	logrus.RegisterExitHandler(func() {
		hook.Stop()
	})

	hooks = append(hooks, hook)

	return hook, nil
}

func processFailure(pr *producer.Producer) {
	for msg := range pr.NotifyFailures() {
		pr.Logger.Error(msg.Error())
	}
}

func levelConverter(level logrus.Level) (string) {
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return "error"
	case logrus.WarnLevel:
		return "warning"
	case logrus.InfoLevel:
		return "info"
	case logrus.DebugLevel:
		return "debug"
	default:
		return ""
	}
}

func (hook *LoglessHook) Fire(entry *logrus.Entry) error {
	packed, err := msgpack.Marshal(&LogEntry{Message:entry.Message,
		Level:levelConverter(entry.Level),
		Timestamp:entry.Time.UTC().Format("2006-01-02T15:04:05.999999999Z"),
		Source:hostname,
		Fields:entry.Data})

	if err != nil {
		return err
	}

	if err = hook.producer.Put(packed, uuid.NewV4().String()); err != nil {
		return err
	}

	return nil
}

func (hook *LoglessHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *LoglessHook) Stop() {
	hook.producer.Stop()

}

func OnClose() {
	for _, hook := range hooks {
		hook.Stop()
	}
}