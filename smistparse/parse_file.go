package smistparse

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/pglang/analysis/lex_pgl"
	"github.com/ProtossGenius/pglang/snreader"
	"github.com/robertkrimen/otto"
)

const (
	ErrCantReadFileBeforeParse = "ErrCantReadFileBeforeParse: %v"
)

// Parse parse a file.
func Parse(filePath string, vmIniter func(vm *otto.Otto) error, workGroup *sync.WaitGroup) {
	defer workGroup.Done()

	sm := lex_pgl.NewLexAnalysiser()
	str, err := smn_file.FileReadAll(filePath)

	if err != nil {
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

	err = parseFile(filePath+".smist_temp", sm.GetResultChan(), vmIniter)
	if err != nil {
		log.Printf("error happened %v", err)
	}
}

func initNothing(vm *otto.Otto) error {
	return nil
}

// parseFile do parse an rewrite to file.
func parseFile(filePath string, ch <-chan snreader.ProductItf, vmIniter func(vm *otto.Otto) error) (err error) {
	parser := new(ClikePraser)
	newPath := filePath + ".smist_temp"

	if vmIniter == nil {
		vmIniter = initNothing
	}

	err = parser.OpenFile(newPath, vmIniter)
	if err != nil {
		return err
	}

	defer parser.DeferClose()

	for {
		p := <-ch
		lex := lex_pgl.ToLexProduct(p)

		if lex.ProductType() < 0 {
			if lex.ProductType() != snreader.ResultEnd {
				log.Printf("parse file error, ProductType is %v ,reason is %v\n", lex.ProductType(), lex.Value)
			}

			break
		}

		err = parser.OnRead(lex)

		if err != nil {
			return err
		}
	}
	parser.Close()

	return os.Rename(newPath, filePath)
}
