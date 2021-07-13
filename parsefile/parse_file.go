package parsefile

import (
	"fmt"
	"log"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/pglang/analysis/lex_pgl"
	"github.com/ProtossGenius/pglang/snreader"
)

const (
	ErrCantReadFileBeforeParse = "ErrCantReadFileBeforeParse: %v"
)

// Parse parse a file.
func Parse(filePath string) error {
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

	parseFile(filePath, sm.GetResultChan())

	return nil
}

// parseFile do parse an rewrite to file.
func parseFile(filePath string, ch <-chan snreader.ProductItf) {
	parser := new(ClikePrase)
	parser.OpenFile(filePath)
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
