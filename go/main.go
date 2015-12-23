package main

// build: ignore

var buf string

func main() {
	buf = "xyzsdlskdlskdlk flkjsldfkhjskgfdjhsl kjlskjf"
	println(buf)
	buf = buf[20:25]
	println(buf)
}
