package gitsrv

import (
    "net/http"
    "log"
    "os/exec"
    "io"
    "bytes"
)

type GitCommand struct {
    procInput *bytes.Reader
    args []string
}

func WriteGitToHttp(w http.ResponseWriter, gitCommand GitCommand) {
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


