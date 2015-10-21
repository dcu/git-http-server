package gitserver

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os/exec"
)

// GitCommand is a command to be executed by git
type GitCommand struct {
	procInput *bytes.Reader
	args      []string
}

// WriteGitToHTTP copies the output of the git command to the http socket.
func WriteGitToHTTP(w http.ResponseWriter, gitCommand GitCommand) {
	cmd := exec.Command("git", gitCommand.args...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		w.WriteHeader(404)
		log.Fatal("Error:", err)
		return
	}

	if gitCommand.procInput != nil {
		cmd.Stdin = gitCommand.procInput
	}

	if err := cmd.Start(); err != nil {
		w.WriteHeader(404)
		log.Fatal("Error:", err)
		return
	}

	nbytes, err := io.Copy(w, stdout)
	if err != nil {
		log.Fatal("Error writing to socket", err)
	} else {
		log.Printf("Bytes written: %d", nbytes)
	}
}
