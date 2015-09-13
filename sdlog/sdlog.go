package sdlog

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
)

//standard log
var logpath string

func SdLogSetup(dirpath string) {
	logpath = dirpath
}

func Err(err error, innerException string) {
	var errmsg string
	errmsg = fmt.Sprintf("ERROR:%s", err.Error())
	if innerException != "" {
		errmsg = errmsg + fmt.Sprintf("\ninner exception:%s", innerException)
	}
	WriteLogFile(errmsg, true)
}

func Errf(err error, innerExFormat string, params ...interface{}) {
	var errmsg string
	errmsg = fmt.Sprintf("ERROR:%s", err.Error())
	if innerExFormat != "" {
		errmsg = errmsg + fmt.Sprintf("\ninner exception:")
		errmsg = errmsg + fmt.Sprintf(innerExFormat, params...)
	}
	WriteLogFile(errmsg, true)
}

func Info(msg string) {
	//var errmsg string
	msg = fmt.Sprintf("INFO:%s", msg)
	WriteLogFile(msg, true)
}

func Infof(format string, params ...interface{}) {
	var msg string
	msg = msg + fmt.Sprint("INFO:")
	msg = msg + fmt.Sprintf(format, params...)
	WriteLogFile(msg, true)
}

func Debugf(format string, params ...interface{}) {
	var msg string
	msg = msg + fmt.Sprint("DEBUG:")
	msg = msg + fmt.Sprintf(format, params...)
	log.Print(msg)
}

func WriteLogFile(msg string, isPrintStack bool) {
	if logpath != "" {
		now := time.Now()
		logfile := logpath + "/" + now.Format("2006Jan02") + ".log"
		f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "log file error")
			return
		}
		defer f.Close()
		log.SetOutput(f)
	}
	logmsg := fmt.Sprintln(msg)
	if isPrintStack {
		logmsg = logmsg + fmt.Sprintln(string(debug.Stack()))
	}
	log.Printf(logmsg)
}
