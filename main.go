package main

import (
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/smist/smistparse"
)

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
/*
# smist
Through formatted comments, use tools to automatically parse and generate code
*/
/*@SMIST setIgnoreInput(false)//.*/
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var exts string

	var path string

	flag.StringVar(&path, "path", ".", "smist parse path")
	flag.StringVar(&exts, "exts", ".go", "smist parse extra-name; multi split with ','")
	flag.Parse()

	extList := strings.Split(exts, ",")

	workGroup := &sync.WaitGroup{}

	_, err := smn_file.DeepTraversalDir(path, func(p string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}

		hasExt := false

		for _, ext := range extList {
			if strings.HasSuffix(info.Name(), ext) {
				hasExt = true

				break
			}
		}

		if !hasExt {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}

		smistparse.Parse(p, nil, workGroup)

		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})

	check(err)
}
