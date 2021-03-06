%{
// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz


// WARNING: If this file is scanner.go then DO NOT EDIT.
// scanner.go is generated by golex from scanner.l (see the Makefile).


package hosts


import (
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode"

	"github.com/cznic/fileutil"
)


type lex struct {
	hosts     File
	startCond int
	buf       []byte
	peek	  byte
	line      int
	column    int
	src		  *strings.Reader
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
		if fileutil.IsEOF(err) {
			l.peek = 0
			return 0
		}
		panic(err)
	}
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
	panic(errors.New(e))
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
%}

%x lineStart

%yyt l.startCond
%yyc c
%yyn c = l.getc(c)

ip_address {ipv6}|{ipv4}
byte       {digit}{digit}?{digit}?
digit      [0-9]
hex        [0-9a-fA-F]
hex4       {hex}{hex}?{hex}?{hex}?
hexpart    {hexseq}|{hexseq}"::"{hexseq}?|"::"{hexseq}?
hexseq     {hex4}(":"{hex4})*
ipv4       {byte}"."{byte}"."{byte}"."{byte}
ipv6       {hexpart}

hostname       {name}("."{name})*
let_or_digit   [a-zA-Z0-9]
name           {let_or_digit}(({let_or_digit}|"-")*{let_or_digit})?

%%
	if ret >= 0 {
		lval.str = string(l.buf)
		return
	}

	l.buf = l.buf[:0]

<lineStart>[ \t]+

<lineStart>{ip_address}
	ip := net.ParseIP(string(l.buf))
	if ip == nil {
		panic(fmt.Errorf("invalid IP %q", l.buf))
	}
	lval.ip = ip
	l.begin(0)
	return tIP_ADDRESS

{hostname}
	ret = tHOST_NAME

[ \t]+
	return ' '

<*>[ \t]*#.*

%%
	// fail
	if len(l.buf) == 0 {
		ret = int(c)
		l.getc(0)
		return
	}

	return unicode.ReplacementChar
}
