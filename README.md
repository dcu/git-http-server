# Git HTTP Server

Server and library to handle git requests via HTTP.


## Installation

To install it just type:

```shell
go get github.com/dcu/git-http-server
```

## Usage as library

```go
import (
    "github.com/dcu/git-http-server"
    "fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "nothing to see here\n")
}

func main() {
    config := &gitserver.Config{
        ReposRoot:     "/tmp/repos",
        AutoInitRepos: true,
    }

    err := gitserver.Init(config)
	if err != nil {
		panic(err)
	}

    http.HandleFunc("/", gitserver.MiddlewareFunc(handler))
	http.ListenAndServe(*listenAddressFlag, nil)
}
```

## Usage as binary

Start the server using:

```shell
git-http-server -repos.root "/tmp/repos" -repos.autoinit -web.listen-address ":5000"
```

Then try pushing any code to it:

```shell
git push http://localhost:5000/foo.git master
```

Since the `autoinit` option was given it should create the repository
automatically.


