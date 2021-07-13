package parsefile

import (
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
	SetIgnoreInput(ignoreIntput bool)
	Close()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// ClikePrase parse C like file.
type ClikePrase struct {
	file        *os.File
	vm          *otto.Otto
	ignoreInput bool
}

func (c *ClikePrase) SetIgnoreInput(ignoreIntput bool) {
	if c.ignoreInput != ignoreIntput {
		_, err := c.file.WriteString("\n")
		check(err)
	}

	c.ignoreInput = ignoreIntput
}
func (c *ClikePrase) set(name string, val interface{}) {
	check(c.vm.Set(name, val))
}

func (c *ClikePrase) OpenFile(filePath string) (err error) {
	if c.file, err = smn_file.CreateNewFile(filePath + ".out"); err != nil {
		return err
	}

	c.vm = otto.New()
	c.set("set", func(name string, value interface{}) {
		check(c.vm.Set(name, value))
	})
	c.set("setIgnoreInput", func(ignoreInput bool) {
		c.SetIgnoreInput(ignoreInput)
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
	c.set("panic", func(reason ...interface{}) {
		panic(reason)
	})

	return nil
}

// getCommentBody get comment body.
func getCommentBody(str string) string {
	if strings.HasPrefix(str, "//") {
		return str[2:]
	}

	return str[2 : len(str)-2]
}

func (c *ClikePrase) OnRead(lex *lex_pgl.LexProduct) error {
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

func (c *ClikePrase) Close() {
	c.file.Close()
}
