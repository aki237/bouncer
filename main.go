package main

import (
	"flag"
	"fmt"
	"os"
)

func warn(msg string) {
	fmt.Fprintf(os.Stderr, "Warning : %s\n", msg)
}

func log(msg string) {
	fmt.Fprintf(os.Stdout, "Log : %s\n", msg)
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "Fatal : %s\n", msg)
	os.Exit(-1)
}

func DirExists(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

func FileExists(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if !fi.IsDir() {
		return true
	}
	return false
}

func main() {
	configFile := flag.String("c", "/etc/bouncer.conf", "Configuration file to use")
	certKey := flag.String("cert", "", "Full chain certificate to use for TLS connections")
	privKey := flag.String("key", "", "Private Key to use for TLS connections")

	flag.Parse()

	if !FileExists(*configFile) {
		fatal("config file passed as a command line argument doesn't exist or not readable by this user. Peace out!!")
	}

	b := NewBouncer(":8080")
	err := b.ReadConfig(*configFile)
	if err != nil {
		fatal(err.Error())
	}

	fmt.Println(b)
	b.Serve(*certKey, *privKey)
}
