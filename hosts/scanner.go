// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

// WARNING: If this file is scanner.go then DO NOT EDIT.
// scanner.go is generated by golex from scanner.l (see the Makefile).


package hosts

import (
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
)

type lex struct {
	hosts     File
	startCond int
	buf       []byte
	peek      byte
	line      int
	column    int
	src       *strings.Reader
}

func (l *lex) getc(c byte) byte {
	if c != 0 {
		l.buf = append(l.buf, c)
	}
	if b, err := l.src.ReadByte(); err == nil {
		l.peek = b
		if b == '\n' {
			l.begin(sc_LINE_START)
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		return b
	} else {
		if err == os.EOF {
			l.peek = 0
			return 0
		}
		panic(err)
	}
	panic("unreachable")
}

func newLex(source *strings.Reader) (l *lex) {
	l = &lex{}
	l.line = 1
	l.column = 1
	l.src = source
	l.begin(sc_LINE_START)
	l.getc(0)
	return
}

func (l *lex) Error(e string) {
	panic(os.NewError(e))
}

func (l *lex) begin(sc int) {
	l.startCond = sc
}

const (
	sc_INITIAL = iota
	sc_LINE_START
)

func (l *lex) Lex(lval *yySymType) (ret int) {
	c := l.peek
	ret = -1

yystate0:

	if ret >= 0 {
		lval.str = string(l.buf)
		return
	}
	l.buf = l.buf[:0]

	switch yyt := l.startCond; yyt {
	default:
		panic(fmt.Errorf(`invalid start condition %d`, yyt))
	case 0: // start condition: INITIAL
		goto yystart1
	case 1: // start condition: lineStart
		goto yystart7
	}

	goto yystate1 // silence unused label error
yystate1:
	c = l.getc(c)
yystart1:
	switch {
	default:
		goto yyabort
	case c == '\t' || c == ' ':
		goto yystate5
	case c == '#':
		goto yystate6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate2
	}

yystate2:
	c = l.getc(c)
	switch {
	default:
		goto yyrule3
	case c == '-':
		goto yystate3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate2
	case c == '.':
		goto yystate4
	}

yystate3:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == '-':
		goto yystate3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate2
	}

yystate4:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate2
	}

yystate5:
	c = l.getc(c)
	switch {
	default:
		goto yyrule4
	case c == '\t' || c == ' ':
		goto yystate5
	case c == '#':
		goto yystate6
	}

yystate6:
	c = l.getc(c)
	switch {
	default:
		goto yyrule5
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate6
	}

	goto yystate7 // silence unused label error
yystate7:
	c = l.getc(c)
yystart7:
	switch {
	default:
		goto yyabort
	case c == '#':
		goto yystate6
	case c >= '0' && c <= '9':
		goto yystate8
	case c == ':':
		goto yystate35
	case c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate13
	case c == '\t' || c == ' ':
		goto yystate34
	}

yystate8:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate12
	case c >= '0' && c <= '9':
		goto yystate9
	case c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate14
	case c == '.':
		goto yystate22
	}

yystate9:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate12
	case c >= '0' && c <= '9':
		goto yystate10
	case c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate15
	case c == '.':
		goto yystate22
	}

yystate10:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate12
	case c == '.':
		goto yystate22
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate11
	}

yystate11:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate12
	}

yystate12:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate13
	case c == ':':
		goto yystate16
	}

yystate13:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate14
	}

yystate14:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate15
	}

yystate15:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate11
	}

yystate16:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate17
	}

yystate17:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate21
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate18
	}

yystate18:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate19
	case c == ':':
		goto yystate21
	}

yystate19:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate20
	case c == ':':
		goto yystate21
	}

yystate20:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c == ':':
		goto yystate21
	}

yystate21:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate17
	}

yystate22:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate23
	}

yystate23:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate26
	case c >= '0' && c <= '9':
		goto yystate24
	}

yystate24:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate26
	case c >= '0' && c <= '9':
		goto yystate25
	}

yystate25:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate26
	}

yystate26:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate27
	}

yystate27:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate30
	case c >= '0' && c <= '9':
		goto yystate28
	}

yystate28:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate30
	case c >= '0' && c <= '9':
		goto yystate29
	}

yystate29:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate30
	}

yystate30:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate31
	}

yystate31:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9':
		goto yystate32
	}

yystate32:
	c = l.getc(c)
	switch {
	default:
		goto yyrule2
	case c >= '0' && c <= '9':
		goto yystate33
	}

yystate33:
	c = l.getc(c)
	goto yyrule2

yystate34:
	c = l.getc(c)
	switch {
	default:
		goto yyrule1
	case c == '#':
		goto yystate6
	case c == '\t' || c == ' ':
		goto yystate34
	}

yystate35:
	c = l.getc(c)
	switch {
	default:
		goto yyabort
	case c == ':':
		goto yystate16
	}

yyrule1: // [ \t]+

	goto yystate0
yyrule2: // {ip_address}
	{

		ip := net.ParseIP(string(l.buf))
		if ip == nil {
			panic(fmt.Errorf("invalid IP %q", l.buf))
		}
		lval.ip = ip
		l.begin(0)
		return tIP_ADDRESS
	}
yyrule3: // {hostname}
	{

		ret = tHOST_NAME
		goto yystate0
	}
yyrule4: // [ \t]+
	{

		return ' '
	}
yyrule5: // [ \t]*#.*

		goto yystate0
		panic("unreachable")

		goto yyabort // silence unused label error

yyabort: // no lexem recognized
	// fail
	if len(l.buf) == 0 {
		ret = int(c)
		l.getc(0)
		return
	}

	return unicode.ReplacementChar
}
