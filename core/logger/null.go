package logger

import "github.com/sirupsen/logrus"

type (
	null struct{}
)

func (n null) WithField(key string, value interface{}) *logrus.Entry { return nil }
func (n null) WithFields(fields logrus.Fields) *logrus.Entry         { return nil }
func (n null) WithError(err error) *logrus.Entry                     { return nil }
func (n null) Debugf(format string, args ...interface{})             { return }
func (n null) Infof(format string, args ...interface{})              { return }
func (n null) Printf(format string, args ...interface{})             { return }
func (n null) Warnf(format string, args ...interface{})              { return }
func (n null) Warningf(format string, args ...interface{})           { return }
func (n null) Errorf(format string, args ...interface{})             { return }
func (n null) Fatalf(format string, args ...interface{})             { return }
func (n null) Panicf(format string, args ...interface{})             { return }
func (n null) Debug(args ...interface{})                             { return }
func (n null) Info(args ...interface{})                              { return }
func (n null) Print(args ...interface{})                             { return }
func (n null) Warn(args ...interface{})                              { return }
func (n null) Warning(args ...interface{})                           { return }
func (n null) Error(args ...interface{})                             { return }
func (n null) Fatal(args ...interface{})                             { return }
func (n null) Panic(args ...interface{})                             { return }
func (n null) Debugln(args ...interface{})                           { return }
func (n null) Infoln(args ...interface{})                            { return }
func (n null) Println(args ...interface{})                           { return }
func (n null) Warnln(args ...interface{})                            { return }
func (n null) Warningln(args ...interface{})                         { return }
func (n null) Errorln(args ...interface{})                           { return }
func (n null) Fatalln(args ...interface{})                           { return }
func (n null) Panicln(args ...interface{})                           { return }

func NewNull() logrus.FieldLogger {
	return &null{}
}
