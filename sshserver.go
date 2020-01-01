package loud

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	gossh "golang.org/x/crypto/ssh"

	"github.com/gliderlabs/ssh"
)

const loudPubkey = "LOUD-pubkey"

func handleConnection(world World, session ssh.Session) {
	user := world.GetUser(session.User())
	screen := NewSSHScreen(session, world, user)
	pubKey, _ := session.Context().Value(loudPubkey).(string)
	userSSH, ok := user.(UserSSHAuthentication)

	if len(session.Command()) > 0 {
		session.Write([]byte("Commands are not supported.\n"))
		session.Close()
	}

	if ok {
		if userSSH.SSHKeysEmpty() {
			userSSH.AddSSHKey(pubKey)
			log.Printf("Saving SSH key for %s", "eugen")
		} else if !userSSH.ValidateSSHKey(pubKey) {
			session.Write([]byte("This is not the SSH key verified for this user. Try another username.\n"))
			log.Printf("User %s doesn't use this key.", "eugen")
			return
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	logMessage := fmt.Sprintf("Logged in as %s via %s at %s", "eugen", session.RemoteAddr(), time.Now().UTC().Format(time.RFC3339))
	log.Println(logMessage)

	done := session.Context().Done()
	tick := time.Tick(500 * time.Millisecond)
	stringInput := make(chan inputEvent, 1)
	reader := bufio.NewReader(session)

	go handleKeys(reader, stringInput, cancel)

	for {
		select {
		case inputString := <-stringInput:
			if inputString.err != nil {
				screen.Reset()
				session.Close()
				continue
			}
			switch inputString.inputString {
			case "UP":
			case "DOWN":
			case "LEFT":
			case "RIGHT":
			case "TAB":
			case "/":
			case "BACKSPACE":
			case "ENTER":
			default:
				screen.HandleInputKey(inputString.inputString)
			}
		case <-ctx.Done():
			cancel()
		case <-tick:
			user.Reload()
			screen.Render()
			continue
		case <-done:
			log.Printf("Disconnected %v@%v", "eugen", session.RemoteAddr())
			screen.Reset()
			session.Close()
			return
		}
	}
}

// ServeSSH runs the main SSH server loop.
func ServeSSH(listen string) {
	rand.Seed(time.Now().Unix())

	world := LoadWorldFromDB("./world.db")
	defer world.Close()

	privateKey := makeKeyFiles()

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		marshal := gossh.MarshalAuthorizedKey(key)
		ctx.SetValue(loudPubkey, string(marshal))
		return true
	})

	log.Printf("Starting SSH server on %v", listen)
	log.Fatal(ssh.ListenAndServe(listen, func(s ssh.Session) {
		handleConnection(world, s)
	}, publicKeyOption, ssh.HostKeyFile(privateKey)))
}
