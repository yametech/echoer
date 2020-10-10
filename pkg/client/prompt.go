package client

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Prompt struct {
	mu         sync.Mutex
	prefix     string
	r          *bufio.Reader
	commandBuf []string
}

func NewPrompt() *Prompt {
	p := &Prompt{
		mu:         sync.Mutex{},
		commandBuf: make([]string, 0),
	}
	p.r = bufio.NewReader(os.Stdin)
	return p
}

func (p *Prompt) Handler(h func(s string)) {
	sourcePrefix := p.prefix
	for {
		fmt.Print(p.prefix)
		line, _, err := p.r.ReadLine()
		if err != nil {
			continue
		}

		if len(line) < 1 || (len(line) == 1 && line[0] == byte(' ')) {
			continue
		}
		strLine := string(bytes.TrimSuffix(line, []byte(" ")))
		suffix := string(strLine[len(strLine)-1])

		if strings.ToUpper(strLine) == "EXIT" {
			break
		}
		p.Put(strLine)

		if suffix != "/" {
			p.SetPrefix("")
			continue
		}
		h(p.Clean())
		p.SetPrefix(sourcePrefix)
	}
}

func (p *Prompt) SetPrefix(s string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.prefix = s
}

func (p *Prompt) Put(input string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commandBuf = append(p.commandBuf, input)
}

func (p *Prompt) Clean() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	var res = strings.TrimSuffix(strings.Join(p.commandBuf, ""), "/")
	p.commandBuf = p.commandBuf[:0]
	return res
}
