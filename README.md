# wasm-go-pomodoro
Pomodoro build in web assembly using go (golang)

### Requirements
go should be installed and have stable version prefered go1.20.5


### Setup
```
git clone https://github.com/sskanishk/wasm-go-pomodoro.git

// copy wasm_exec.js from where Go binaries is stored (for now repo contain this file)
cd web/
cp "$(go env GROOT)\misc\wasm\wasm_exec.js" .

// go to wasm folder to create build
cd cmd/wasm/
GOOS=js GOARCH=wasm go build -o ../../web/main.wasm

// go to server folder to run server
cd cmd/server/
go run main.go
```

### Ref
https://youtu.be/10Mz3z-W1BE
