package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

func SetLogFile(logFile io.Writer) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	_, err := os.Stdout.Write([]byte("\n"))
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	} else {
		log.SetOutput(logFile)
	}
}

func getPrefix(level string) string {
	return fmt.Sprintf("[%s]", level)
}

func Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("%s %s", getPrefix("DEBUG"), format), args...)
}

func Infof(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("%s %s", getPrefix("INFO"), format), args...)
}

func Warnf(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("%s %s", getPrefix("WARN"), format), args...)
}

func Errorf(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("%s %s", getPrefix("ERROR"), format), args...)
}

func Debugln(args ...interface{}) {
	args = append([]interface{}{getPrefix("DEBUG")}, args...)
	log.Println(args...)
}

func Infoln(args ...interface{}) {
	args = append([]interface{}{getPrefix("INFO")}, args...)
	log.Println(args...)
}

func Warnln(args ...interface{}) {
	args = append([]interface{}{getPrefix("WARN")}, args...)
	log.Println(args...)
}

func Errorln(args ...interface{}) {
	args = append([]interface{}{getPrefix("ERROR")}, args...)
	log.Println(args...)
}

func Debug(args ...interface{}) {
	args = append([]interface{}{getPrefix("DEBUG")}, args...)
	log.Print(args...)
}

func Info(args ...interface{}) {
	args = append([]interface{}{getPrefix("INFO")}, args...)
	log.Print(args...)
}

func Warn(args ...interface{}) {
	args = append([]interface{}{getPrefix("WARN")}, args...)
	log.Print(args...)
}

func Error(args ...interface{}) {
	args = append([]interface{}{getPrefix("ERROR")}, args...)
	log.Print(args...)
}
