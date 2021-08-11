package smistparse

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_exec"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/pglang/analysis/lex_pgl"
	"github.com/robertkrimen/otto"
)

// Parser file parse.
type Parser interface {
	OpenFile(filePath string) error
	OnRead(*lex_pgl.LexProduct) error
	Close()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// ClikePraser parse C like file.
type ClikePraser struct {
	file        *os.File
	vm          *otto.Otto
	ignoreInput bool
	filePath    string
	jsPath      string
}

func (c *ClikePraser) setIgnoreInput(ignoreInput bool) {
	if c.ignoreInput != ignoreInput && ignoreInput {
		_, err := c.file.WriteString("\n")
		check(err)
	}

	c.ignoreInput = ignoreInput
}

func (c *ClikePraser) set(name string, val interface{}) {
	check(c.vm.Set(name, val))
}

func (c *ClikePraser) OpenFile(filePath, jsPath string, vmIniter func(vm *otto.Otto) error) (err error) {
	wrap := func(err error) error {
		return fmt.Errorf("ClinkeParser::OpenFile Err : %w", err)
	}
	if c.file, err = smn_file.CreateNewFile(filePath); err != nil {
		return wrap(err)
	}

	c.jsPath = jsPath
	c.filePath = filePath

	c.vm = otto.New()
	c.set("set", func(name string, value interface{}) {
		check(c.vm.Set(name, value))
	})
	c.set("setIgnoreInput", func(ignoreInput bool) {
		c.setIgnoreInput(ignoreInput)
	})
	c.set("write", func(str string) {
		_, err := c.file.WriteString(str)
		check(err)
	})
	c.set("readFile", func(filePath string) string {
		data, err := smn_file.FileReadAll(filePath)
		check(err)

		return string(data)
	})
	c.set("exec", func(dir, cmd string, args ...string) string {
		oInfo, oErr, err := smn_exec.DirExecGetOut(dir, cmd, args...)
		check(err)
		log.Print(oErr)

		return oInfo
	})
	c.set("include", func(jsPath string) {
		data, err := smn_file.FileReadAll(c.jsPath + "/" + jsPath)
		if err != nil {
			panic("when parse file " + filePath + " include file: " + err.Error())
		}
		code := string(data)
		_, err = c.vm.Run(code)
		if err != nil {
			panic("include code fail file = " + c.filePath + ", code is \n" + codeAddLine(code) + "\n error is \n" + err.Error())
		}
	})
	c.set("panic", func(reason interface{}) {
		panic(reason)
	})

	return vmIniter(c.vm)
}

// getCommentBody get comment body.
func getCommentBody(str string) string {
	if strings.HasPrefix(str, "//") {
		return str[2:]
	}

	return str[2 : len(str)-2]
}

func (c *ClikePraser) OnRead(lex *lex_pgl.LexProduct) error {
	if !c.ignoreInput {
		if _, err := c.file.Write([]byte(lex.Value)); err != nil {
			return err
		}
	}

	if lex_pgl.IsComment(lex) {
		comm := getCommentBody(lex.Value)
		comm = strings.TrimSpace(comm)
		before := c.ignoreInput
		if strings.HasPrefix(comm, "@SMIST") {
			if _, err := c.vm.Run(comm[6:]); err != nil {
				return err
			}

			if before != c.ignoreInput && !c.ignoreInput {
				if _, err := c.file.Write([]byte(lex.Value)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *ClikePraser) Close() {
	c.file.Close()
	c.file = nil
}

func (c *ClikePraser) DeferClose() {
	if c.file != nil {
		c.file.Close()
		os.Remove(c.filePath)
	}
}
