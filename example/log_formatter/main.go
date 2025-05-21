package main

import (
	"Taurus/pkg/logx"
	"fmt"
	"time"
)

// if you want to use a custom formatter, you can implement the Formatter interface, and register it to the logx.RegisterFormatter
// then you must set the Formatter name to the config/logger/logger.yaml
type DemoFormatter struct {
}

func (d *DemoFormatter) Format(level logx.LogLevel, message string) string {
	return fmt.Sprintf("[%s] [%s]: %s", time.Now().Format("2006-01-02 15:04:05"), logx.GetLevelSTR(level), message)
}

func init() {
	logx.RegisterFormatter("demo", &DemoFormatter{})
}

func main() {
}
