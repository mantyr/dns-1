// Copyright (c) 2010 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

// Package zone handles master files which may be textual data formated as per RFC 1035 chapter 5
// or binary data produced by Compiler.
package zone

import (
	"bufio"
	"fmt"
	"github.com/cznic/dns/rr"
	"github.com/cznic/fileutil"
	"io"
	"net"
	"os"
	"strings"
)

var errIP = net.ParseIP("0.0.0.0")

var typex = map[rr.Type]int{
	rr.TYPE_A:          tA,
	rr.TYPE_A6:         tA6,
	rr.TYPE_AAAA:       tAAAA,
	rr.TYPE_AFSDB:      tAFSDB,
	rr.TYPE_APL:        tAPL,
	rr.TYPE_ATMA:       tATMA,
	rr.TYPE_CDS:        tCDS,
	rr.TYPE_CERT:       tCERT,
	rr.TYPE_CNAME:      tCNAME,
	rr.TYPE_DHCID:      tDHCID,
	rr.TYPE_DNAME:      tDNAME,
	rr.TYPE_DNSKEY:     tDNSKEY,
	rr.TYPE_DS:         tDS,
	rr.TYPE_EID:        tEID,
	rr.TYPE_GID:        tGID,
	rr.TYPE_GPOS:       tGPOS,
	rr.TYPE_HINFO:      tHINFO,
	rr.TYPE_HIP:        tHIP,
	rr.TYPE_IPSECKEY:   tIPSECKEY,
	rr.TYPE_ISDN:       tISDN,
	rr.TYPE_KEY:        tKEY,
	rr.TYPE_KX:         tKX,
	rr.TYPE_LOC:        tLOC,
	rr.TYPE_MB:         tMB,
	rr.TYPE_MD:         tMD,
	rr.TYPE_MF:         tMF,
	rr.TYPE_MG:         tMG,
	rr.TYPE_MINFO:      tMINFO,
	rr.TYPE_MR:         tMR,
	rr.TYPE_MX:         tMX,
	rr.TYPE_NAPTR:      tNAPTR,
	rr.TYPE_NIMLOC:     tNIMLOC,
	rr.TYPE_NINFO:      tNINFO,
	rr.TYPE_NS:         tNS,
	rr.TYPE_NSAP:       tNSAP,
	rr.TYPE_NSAP_PTR:   tNSAP_PTR,
	rr.TYPE_NSEC:       tNSEC,
	rr.TYPE_NSEC3:      tNSEC3,
	rr.TYPE_NSEC3PARAM: tNSEC3PARAM,
	rr.TYPE_NULL:       tNULL,
	rr.TYPE_NXT:        tNXT,
	rr.TYPE_PTR:        tPTR,
	rr.TYPE_PX:         tPX,
	rr.TYPE_RKEY:       tRKEY,
	rr.TYPE_RP:         tRP,
	rr.TYPE_RRSIG:      tRRSIG,
	rr.TYPE_RT:         tRT,
	rr.TYPE_SIG:        tSIG,
	rr.TYPE_SOA:        tSOA,
	rr.TYPE_SPF:        tSPF,
	rr.TYPE_SRV:        tSRV,
	rr.TYPE_SSHFP:      tSSHFP,
	rr.TYPE_TALINK:     tTALINK,
	rr.TYPE_TKEY:       tTKEY,
	rr.TYPE_TLSA:       tTLSA,
	rr.TYPE_TSIG:       tTSIG,
	rr.TYPE_TXT:        tTXT,
	rr.TYPE_UID:        tUID,
	rr.TYPE_UINFO:      tUINFO,
	rr.TYPE_UNSPEC:     tUNSPEC,
	rr.TYPE_WKS:        tWKS,
	rr.TYPE_X25:        tX25,
}

// Compiler is an object for compiling (large amounts of) DNS resource records
// to binary data for faster loading (no tokenizing/parsing).
type Compiler struct {
	w     io.Writer
	rrs   rr.RRs
	owner string
}

// NewCompiler returns a newly created Compiler writing to w.
// The produced binary data may be embeded in another data as their
// end (a logical EOF) is marked on creation and detected on loading.
// Compiled RRs are also smaller than the original textual data.
// The data compression ratio on some measured real world sample of cca 3.5 million RRs
// was about 60% (i.e. space savings were about 40%).
func NewCompiler(w io.Writer) (c *Compiler) {
	c = &Compiler{w: w}
	return
}

func (c *Compiler) flush() (err error) {
	if len(c.rrs) == 0 {
		return
	}

	var b rr.Bytes
	b.Pack(c.rrs)
	c.rrs = c.rrs[:0]

	n := len(b)
	if err = c.write([]byte{byte(n >> 8), byte(n)}); err != nil {
		return
	}

	return c.write(b)
}

func (c *Compiler) write(b []byte) (err error) {
	//fmt.Printf(".write %04x: % x\n", len(b), b)
	n, m := len(b), 0
	if m, err = c.w.Write(b); err != nil {
		return
	}

	if m != n {
		err = fmt.Errorf("zone.Compiler.write() - short write %d/%d", m, n)
	}
	return
}

var eof = []byte{0, 0}

// Done marks the logical EOF of the compiled RRs block.
// After calling Done the Compiler is no more usable and will panic on attempts
// to use it. If the io.Writer passed to NewCompiler is an bufio.Writer then
// the user of Compiler is responsible for calling Writer.Flush and possibly Writer.Close
// after calling Compiler.Done.
// The semantics of Done are like os.File.Close except that no io.Writer closing
// is performed and the closed entity is the Compiler per se.
//
// Warning: Failure to invoke Done as a last task of a compilation causes loss and/or
// corruption of the produced binary data.
func (c *Compiler) Done() (err error) {
	defer func() { c.w = nil }()
	if err = c.flush(); err != nil {
		return
	}

	return c.write(eof) // logical EOF marker
}

// Write appends r to the "compilation".
func (c *Compiler) Write(r *rr.RR) (err error) {
	owner := strings.ToLower(r.Name)
	if owner != c.owner {
		if err = c.flush(); err != nil {
			return
		}

		c.owner = owner
	}
	c.rrs = append(c.rrs, r)
	return
}

// Load attempts to load a zone/master (RFC1034/5.1) file from fname.
// On syntax error the errHandler is invoked if it's not nil, otherwise
// the loading is aborted and Error returned.
// If errHandler is not nil and returns false the loading is also aborted.
// On unrecoverable errors like file not found the load is aborted
// and Error returned.
// rrHandler is invoked for every resource record found in the zone file.
// If rrHandler returns false the loading is aborted and returns nil Error.
func Load(fname string, errHandler func(e string) bool, rrHandler func(rr *rr.RR) bool) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var file *os.File
	if file, err = os.Open(fname); err != nil {
		panic(fmt.Errorf("zone load: %s", err))
	}

	defer file.Close()

	lx := newLex(fname, bufio.NewReader(file), errHandler, rrHandler)
	if yyParse(lx) != 0 {
		panic(fmt.Errorf("%s:%d:%d - synatx error", fname, lx.line, lx.column))
	}

	return
}

// Load attempts to load compiled RRs from r.
// rrHandler is invoked for every RRs pack found in the data.
// If rrHandler returns false the loading is aborted without error.
func LoadBinary(r io.Reader, rrHandler func(rr.Bytes) bool) (err error) {
	lbuf := []byte{0, 0}
	for {
		if err = fileutil.Read(r, lbuf); err != nil {
			return
		}

		msglen := int(lbuf[0])<<8 | int(lbuf[1])
		if msglen == 0 {
			break // Done
		}

		rbuf := make([]byte, msglen)
		if err = fileutil.Read(r, rbuf); err != nil {
			return
		}

		if !rrHandler(rr.Bytes(rbuf)) {
			return
		}

	}
	return
}
