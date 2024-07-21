package main

import (
	"github.com/gliderlabs/ssh"
	"io"
)

func w(s ssh.Session, data string) {
	_, _ = io.WriteString(s, data)
}
