package main

import "github.com/ProtossGenius/smist/parsefile"

/*@SMIST
setIgnoreInput(true)
console.log("here ???")
str = readFile("./README.md")
write("/*\n")
write(str)
write("*\/")
*/

// aaaaaaaaaaaaaaaaaaaa.

/*@SMIST
setIgnoreInput(false)
*/

func main() {
	err := parsefile.Parse("./main.go")
	if err != nil {
		panic(err)
	}
}
