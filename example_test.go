package logless

import log "github.com/Sirupsen/logrus"

func Example() {
	hook1, _ := NewLoglessHook("logless-test", log.New())
	hook2, _ := NewLoglessHook("logless-test", log.New())
	defer OnClose()

	log.AddHook(hook1)
	log.AddHook(hook2)

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

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 125,
	}).Debug("The group's number increased tremendously 4!")
}