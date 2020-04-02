package log

import (
  "os"
  
  "github.com/sirupsen/logrus"
)

var App  *logrus.Entry

func Setup() {
  // app log
  log := logrus.New()
  log.Out = os.Stdout
  App = log.WithFields(logrus.Fields{"channel":"app"})

  // ...

}

