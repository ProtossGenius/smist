package parsefile

import (
	"fmt"
	"log"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/pglang/analysis/lex_pgl"
	"github.com/ProtossGenius/pglang/snreader"
	"github.com/robertkrimen/otto"
)

const (
	ErrCantReadFileBeforeParse = "ErrCantReadFileBeforeParse: %v"
)

// Parse parse a file.
func Parse(filePath string, vmIniter func(vm *otto.Otto) error) error {
	sm := lex_pgl.NewLexAnalysiser()
	str, err := smn_file.FileReadAll(filePath)

	if err != nil {
		return fmt.Errorf(ErrCantReadFileBeforeParse, err)
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

	return nil
}

func initNothing(vm *otto.Otto) error {
	return nil
}

// parseFile do parse an rewrite to file.
func parseFile(filePath string, ch <-chan snreader.ProductItf, vmIniter func(vm *otto.Otto) error) {
	parser := new(ClikePraser)
	if vmIniter == nil {
		vmIniter = initNothing
	}
	parser.OpenFile(filePath, vmIniter)
	defer parser.Close()

	for {
		p := <-ch
		lex := lex_pgl.ToLexProduct(p)
		if lex.ProductType() < 0 {
			if lex.ProductType() != snreader.ResultEnd {
				log.Fatalf("parse file error, ProductType is %v ,reason is %v", lex.ProductType(), lex.Value)
			}
			break
		}
		fmt.Print(lex.Value)
		parser.OnRead(lex)
	}
}
