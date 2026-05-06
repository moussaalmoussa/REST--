// logger/logger.go
func New() *logrus.Logger {
    l := logrus.New()
    l.SetFormatter(&logrus.JSONFormatter{})
    l.SetLevel(logrus.InfoLevel)
    return l
}
