package color

import (
	"fmt"
	"log"
	"strings"
)

const infoLogTemplate = "\x1b[33;1m%s\x1b[0m\n"
const errorLogTemplate = "\x1b[31;1m%s\x1b[0m\n"

// Println prints the specified message in a colored fashion.
func Println(msg string) {
	log.Printf(infoLogTemplate, format(msg))
}

// Printf prints the specified formatted message in a colored fashion.
func Printf(f string, v ...interface{}) {
	log.Printf(fmt.Sprintf(infoLogTemplate, format(f)), v...)
}

// Fatal prints the specified fatal message in a colored fashion.
func Fatal(msg string) {
	log.Fatal(fmt.Sprintf(errorLogTemplate, format(msg)))
}

// Fatalf prints the specified formatted fatal message in a colored fashion.
func Fatalf(f string, v ...interface{}) {
	log.Fatalf(fmt.Sprintf(errorLogTemplate, format(f)), v...)
}

// Warnln prints the specified warning message in a colored fashion.
func Warnln(msg string) {
	log.Printf(errorLogTemplate, format(msg))
}

// Warnf prints the specified formatted warning message in a colored fashion.
func Warnf(f string, v ...interface{}) {
	log.Printf(fmt.Sprintf(errorLogTemplate, format(f)), v...)
}

func format(msg string) string {
	m := msg
	if strings.HasSuffix(m, "\n") {
		m = m[:len(m)-1]
	}
	return m
}

