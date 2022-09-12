package util

import (
	"log"
	"os"
)

func GetLogger(name string) *log.Logger {
	return log.New(os.Stdout, name+":  ", log.Lshortfile|log.Lmicroseconds|log.Ldate)
}
