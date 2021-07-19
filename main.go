package main

import "github.com/ProtossGenius/smist/parsefile"

/*@SMIST
setIgnoreInput(true)
console.log("here ???")
str = readFile("./README.md")
write("/*\n")
write(str)
write("*\/")
include("./meta_datas/split.js")
split("hello")
*/

// aaaaaaaaaaaaaaaaaaaa.

/*@SMIST
setIgnoreInput(false)
*/

func main() {
	err := parsefile.Parse("./main.go", nil)
	if err != nil {
		panic(err)
	}
}
