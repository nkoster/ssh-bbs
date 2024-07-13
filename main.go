package main

import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"strings"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		clearScreen(s)
		items := []string{"Item1", "Item2", "Item3"}
		selectedIndex := 0

		drawMenu(s, items, selectedIndex)

		input := make([]byte, 1)
		for {
			_, err := s.Read(input)
			if err != nil {
				io.WriteString(s, "\nFailed to read input: "+err.Error()+"\n")
				return
			}

			switch input[0] {
			case 65: // up arrow (ANSI escape code)
				if selectedIndex > 0 {
					selectedIndex--
					drawMenu(s, items, selectedIndex)
				}
			case 66: // down arrow (ANSI escape code)
				if selectedIndex < len(items)-1 {
					selectedIndex++
					drawMenu(s, items, selectedIndex)
				}
			case 10, 13: // enter key (newline and carriage return)
				io.WriteString(s, "\nYou selected "+items[selectedIndex]+"!\n")
				return
			}
		}
	})

	serverOptions := []ssh.Option{
		ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
			return true
		}),
		ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		ssh.HostKeyFile("./ssh_host_rsa_key"),
	}

	log.Println("Starting SSH server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil, serverOptions...))
}

func clearScreen(s ssh.Session) {
	io.WriteString(s, "\033[H\033[2J")
}

func drawMenu(s ssh.Session, items []string, selectedIndex int) {
	clearScreen(s)

	// Determine the width of the box
	maxWidth := 0
	for _, item := range items {
		if len(item) > maxWidth {
			maxWidth = len(item)
		}
	}
	maxWidth += 2 // Add padding and marker space

	// Draw the top of the box
	io.WriteString(s, "┌"+strings.Repeat("─", maxWidth)+"┐\n")

	// Draw each item
	for i, item := range items {
		if i == selectedIndex {
			io.WriteString(s, "│\033[42m\033[1m\033[30m "+item+" \033[0m"+strings.Repeat("", maxWidth-len(item))+"│\n")
		} else {
			io.WriteString(s, "│ "+item+strings.Repeat(" ", maxWidth-len(item)-1)+"│\n")
		}
	}

	// Draw the bottom of the box
	io.WriteString(s, "└"+strings.Repeat("─", maxWidth)+"┘\n")
}
