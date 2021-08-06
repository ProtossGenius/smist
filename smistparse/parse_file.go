package smistparse

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/pglang/analysis/lex_pgl"
	"github.com/ProtossGenius/pglang/snreader"
	"github.com/robertkrimen/otto"
)

const (
	ErrCantReadFileBeforeParse = "ErrCantReadFileBeforeParse: %v"
)

// logInfo print logInfo message.
func logInfo(objs ...interface{}) {
	log.Println(objs...)
}

// logWarn print info message.
func logWarn(objs ...interface{}) {
	log.Println(objs...)
}

// err print info message.
func logErr(objs ...interface{}) {
	log.Println(objs...)
}

// Parse parse a file.
func Parse(filePath string, vmIniter func(vm *otto.Otto) error, workGroup *sync.WaitGroup) {
	workGroup.Add(1)
	defer workGroup.Done()

	sm := lex_pgl.NewLexAnalysiser()

	var str []byte

	var err error

	if str, err = smn_file.FileReadAll(filePath); err != nil {
		log.Printf("error happened %v\n", fmt.Sprintf(ErrCantReadFileBeforeParse, err))

		return
	}

	go func() {
		for _, char := range string(str) {
			err := sm.Read(&lex_pgl.PglaInput{Char: char})
			if err != nil {
				log.Fatalf("When parse lex find error, error is %v", err)

				break
			}
		}

		sm.End()
	}()

	parseFile(filePath, sm.GetResultChan(), vmIniter)
}

func initNothing(vm *otto.Otto) error {
	return nil
}

func codeAddLine(code string) string {
	lines := strings.Split(code, "\n")
	for i := range lines {
		lines[i] = strconv.Itoa(i+1) + "|" + lines[i]
	}

	return strings.Join(lines, "\n")
}

// parseFile do parse an rewrite to file.
func parseFile(filePath string, ch <-chan snreader.ProductItf, vmIniter func(vm *otto.Otto) error) {
	parser := new(ClikePraser)
	newPath := filePath + ".smist_temp"

	if vmIniter == nil {
		vmIniter = initNothing
	}

	err := parser.OpenFile(newPath, vmIniter)
	if err != nil {
		logErr("open file error", err)

		return
	}

	defer parser.DeferClose()

	for {
		p := <-ch
		lex := lex_pgl.ToLexProduct(p)

		if lex.ProductType() < 0 {
			if lex.ProductType() != snreader.ResultEnd {
				log.Printf("parse file error, ProductType is %v ,reason is %v\n", lex.ProductType(), codeAddLine(lex.Value))
			}

			break
		}

		err = parser.OnRead(lex)

		if err != nil {
			logErr("when parsing error happened, code = ", lex, "error is ", err)
		}
	}
	parser.Close()

	err = os.Rename(newPath, filePath)
	if err != nil {
		logErr("when rename temp file, error happened, error is : ", err)
	}
}
