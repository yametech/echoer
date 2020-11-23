// Copyright 2020 <laik.lj@me.com>. All rights reserved.
// Use of this source code is governed by a Apache
// license that can be found in the LICENSE file.

package fss

//#include "token.h"
//#include "fss.lex.h"
// extern int yylex();
// extern int yylineno;
// extern char *yytext;
import "C"

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var _ fssLexer = (*fssLex)(nil)

type fssLex struct {
	yylineno       int
	yytext         string
	lastErr        error
	yylineposition int
}

func NewFssLexer(data []byte) *fssLex {
	p := new(fssLex)
	p.yylineno = 1

	C.yy_scan_bytes(
		(*C.char)(C.CBytes(data)),
		C.yy_size_t(len(data)),
	)
	return p
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (p *fssLex) Lex(lval *fssSymType) int {
	p.lastErr = nil
	var ntoken = C.yylex()
	p.yylineposition += int(C.yylineno)
	p.yytext = C.GoString(C.yytext)
	var err error

	switch ntoken {
	case C.IDENTIFIER:
		lval._identifier = p.yytext
		return IDENTIFIER
	case C.NUMBER_VALUE:
		lval._number, err = strconv.ParseInt(p.yytext, 10, 64)
		if err != nil {
			fmt.Printf("parser number error %s tttext %s \n", err, p.yytext)
			goto END
		}
		return NUMBER_VALUE
	case C.LIST:
		result := make([]interface{}, 0)
		if strings.HasPrefix(p.yytext, "list(") && strings.HasSuffix(p.yytext, ")") {
			newYYText := p.yytext[5 : len(p.yytext)-1]
			for _, item := range strings.Split(newYYText, ",") {
				result = append(result, item)
			}
			lval._list = result
			return LIST
		}

		if strings.HasPrefix(p.yytext, "[") && strings.HasSuffix(p.yytext, "]") {
			newYYText := p.yytext[1 : len(p.yytext)-1]
			for _, item := range strings.Split(newYYText, ",") {
				result = append(result, item)
			}
			lval._list = result
			return LIST
		}
	case C.DICT:
		if !strings.HasPrefix(p.yytext, "dict(") || !strings.HasSuffix(p.yytext, ")") {
			return 0
		}
		result := make(map[string]interface{})
		newYYText := p.yytext[5 : len(p.yytext)-1]
		items := strings.Split(newYYText, ",")
		for _, item := range items {
			kvItem := strings.Split(item, "=")
			if len(kvItem) != 2 {
				continue
			}
			result[kvItem[0]] = kvItem[1]
		}
		lval._dict = result
		return DICT
	case C.STRING_VALUE:
		if len(p.yytext) < 1 {
			return STRING_VALUE
		}
		switch p.yytext[0] {
		case '"':
			lval._string = strings.TrimSuffix(strings.TrimPrefix(p.yytext, `"`), `"`)
		case '`':
			lval._string = strings.TrimSuffix(strings.TrimPrefix(p.yytext, "`"), "`")
		}
		return STRING_VALUE
	case C.FLOW:
		return FLOW
	case C.FLOW_END:
		return FLOW_END
	case C.DECI:
		return DECI
	case C.STEP:
		return STEP
	case C.ACTION:
		return ACTION
	case C.ARGS:
		return ARGS
	case C.LPAREN:
		return LPAREN
	case C.RPAREN:
		return RPAREN
	case C.LSQUARE:
		return LSQUARE
	case C.RSQUARE:
		return RSQUARE
	case C.LCURLY:
		return LCURLY
	case C.RCURLY:
		return RCURLY
	case C.ASSIGN:
		return ASSIGN
	case C.SEMICOLON:
		return SEMICOLON
	case C.OR:
		return OR
	case C.AND:
		return AND
	case C.TO:
		return TO
	case C.COMMA:
		return COMMA
	case C.COLON:
		return COLON
	case C.DEST:
		return DEST
	case C.HTTP:
		lval._http = ActionHTTPMethod
		return HTTP
	case C.GRPC:
		lval._grpc = ActionGRPCMethod
		return GRPC
	case C.INT:
		return INT
	case C.STR:
		return STR
	case C.ACTION_END:
		return ACTION_END
	case C.ADDR:
		return ADDR
	case C.METHOD:
		return METHOD
	case C.FLOW_RUN:
		return FLOW_RUN
	case C.FLOW_RUN_END:
		return FLOW_RUN_END
	case C.RETURN:
		return RETURN
	case C.EOL:
		p.yylineno += 1
		return EOL
	case C.ILLEGAL:
		fmt.Printf("lex: ILLEGAL token, yytext = %q, yylineno = %d, position = %d \n", p.yytext, p.yylineno, p.yylineposition)
	}
END:
	return 0
}

func (p *fssLex) Error(e string) {
	p.lastErr = errors.New("yacc: " + e)
	if err := p.lastErr; err != nil {
		fmt.Printf("lex: lastErr = %s, lineno = %d, position = %d, text = %s \n", p.lastErr, p.yylineno, p.yylineposition, p.yytext)
	}
}
