# wasm-go-pomodoro
Pomodoro build in web assembly using go (golang)

![image](https://github.com/sskanishk/wasm-go-pomodoro/assets/29313203/9a3f9c7a-ab8b-47ad-a898-5b9afd04dbe0)


### Requirements
go should be installed and have stable version prefered go1.20.5


### Setup (Windows)
```
git clone https://github.com/sskanishk/wasm-go-pomodoro.git

// copy wasm_exec.js from where Go binaries is stored
cd web/
cp "$(go env GOROOT)\misc\wasm\wasm_exec.js" .

// go to wasm folder to create build
cd cmd/wasm/
GOOS=js GOARCH=wasm go build -o ../../web/main.wasm

// go to server folder to run server
cd cmd/server/
go run main.go
```

### Ref
https://youtu.be/10Mz3z-W1BE
https://www.sitepen.com/blog/compiling-go-to-webassembly
