package main

import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		_, _ = io.WriteString(s, "Welcome "+s.User()+"!\n")
		//cmd := exec.Command("/usr/bin/env")
		//cmdOutput, err := cmd.CombinedOutput()
		//if err != nil {
		//	_, _ = io.WriteString(s, "Error: "+err.Error()+"\n")
		//} else {
		//	_, _ = io.WriteString(s, string(cmdOutput))
		//}
	})

	serverOptions := []ssh.Option{
		ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
			return true
		}),
		ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		// Create a host key:
		// ssh-keygen -t rsa -b 4096 -f ssh_host_rsa_key -N ""
		ssh.HostKeyFile("./ssh_host_rsa_key"),
	}

	log.Println("Starting SSH server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil, serverOptions...))
}
