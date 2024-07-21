package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"log"
	"strings"
	"time"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		fmt.Printf("SSH: %s@%s", s.User(), s.RemoteAddr())
		waitForBanner(s)
		items := []string{"Item 1", "Item 2", "Item 3"}
		selectedIndex := 0

		drawMenu(s, items, selectedIndex)

		input := make([]byte, 1)
		for {
			_, err := s.Read(input)
			if err != nil {
				w(s, "\nFailed to read input: "+err.Error()+"\n")
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
				w(s, "\nYou selected "+items[selectedIndex]+"!\n")
				fmt.Printf(" (%s)\n", items[selectedIndex])
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

func waitForBanner(s ssh.Session) {
	// Wait a second...
	time.Sleep(time.Second)
	// ...then clear the screen
	w(s, "\033[H\033[2J")
}

func clearScreen(s ssh.Session) {
	w(s, "\033[H\033[2J")
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
	w(s, "┌"+strings.Repeat("─", maxWidth)+"┐\n")

	// Draw each item
	for i, item := range items {
		if i == selectedIndex {
			w(s, "│\033[0;102m\033[1m\033[30m "+item+" \033[0m"+strings.Repeat("", maxWidth-len(item))+"│\n")
		} else {
			w(s, "│ "+item+strings.Repeat(" ", maxWidth-len(item)-1)+"│\n")
		}
	}

	// Draw the bottom of the box
	w(s, "└"+strings.Repeat("─", maxWidth)+"┘\n")
}
