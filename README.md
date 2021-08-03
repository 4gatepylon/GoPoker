# Building
Make sure you have Golang 1.16 installed.

Run `go get google.golang.org/grpc`. We are using `v1.39.0` and protoc `v3.17.3`. I installed with homebrew on my mac.

You'll also need to install `protoc-gen-go` v1.26 and `protoc-gen-go-grpc` v1.1. I followed the tutorial [here](https://grpc.io/docs/languages/go/quickstart/). I ran the two commands below.

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
export PATH="$PATH:$(go env GOPATH)/bin" # Add to bashrc or zshrc
```

**First** you'll need to build any proto files into go. `cd net/proto` then
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ping.proto
```

**Second** (and this is important if you make any modifications) run `go mod tidy` from the root directory. This will ensure that any implicit dependencies of the generated protobuf golang code will be required by the GoPoker module.

**Last** To build the client/server stack simply `cd main` then `go build main`. This should create an executable that works on your system.

PS, In the future, as I add more protos and whatnot, it's possible that I may create a builder of some kind. It would be a good experience. I was originally using Bazel, but it was a massive pain in the butt and I never got it to run with grpc.

# Running
`cd main` after building and then with a `tmux` or pair of tabs/windows, run `./main -client=0` for the server and `./main` for your clients.

# What's left
1. Finishing up the game itself (check todos inside `game.go`; also, make sure to test `utils/naming.go` as well)
2. Creating a CLI interface for clients
3. Creating a flags CLI interface for server launch
4. Creating a server/client protocol
5. Creating a server/client abstraction/object for the game server

Ideally draw an ascii schematic or something of the sort too if you can. Below are helpful docs/resources.

- Remember to test grpc as well: `https://stackoverflow.com/questions/42102496/testing-a-grpc-service`.
- Protobuf Docs: `https://developers.google.com/protocol-buffers/docs/overview`.
- Example Protobuf + GRPC: `https://www.youtube.com/watch?v=SBPjEbZcgf8`.
- GRPC API In go: `https://grpc.io/docs/languages/go/`.