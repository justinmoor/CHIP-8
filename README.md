# CHIP-8
A CHIP-8 emulator implemented in Go for learning purposes which runs in the browser using WebAssembly or in a local terminal.

### Requirements to run
* You have a computer
* You have Go installed on that computer

### How to run
* Clone this repository

#### Web Assembly
The folder ```static``` contains all the static files needed to run, including the WASM build (```main.wasm```). 
* To run you need to serve the HTTP server using ```go run cmd\server\main.go``` (or any other HTTP server you'd like to use)
* If you want to modify the emulator code, make sure you compile ```cmd\wasm\main.go``` again. This will update the WASM build ```main.wasm``` in the ```static``` folder. For compiling to WASM use ```GOOS=js;GOARCH=wasm```.

More info on compiling Go to WASM can be found [here](https://github.com/golang/go/wiki/WebAssembly).

#### Locally in your terminal
* Just run ```go run cmd\terminal\main.go [path to any ROM]```

### Known issues
* Still somewhat buggy. For example TIC TAC TOE does not work properly.
* I did not take the time to add the beep :(.
* WASM build sometimes runs a bit slow.

Enjoy, I hope you can learn as much from it as I did myself. :)
