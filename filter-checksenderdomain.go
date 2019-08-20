//
// Copyright (c) 2019 Gilles Chehade <gilles@poolp.org>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var filters = map[string]func(string, []string) {
	"mail-from": filterMailFrom,
}

func filterMailFrom(sessionId string, params[] string) {
	token := params[0]
	sender := params[1]

	// mailer daemon
	if sender == "" {
		fmt.Printf("filter-result|%s|%s|proceed\n", token, sessionId)
		return
	}

	// local user
	parts := strings.Split(sender, "@")
	if len(parts) == 1 {
		fmt.Printf("filter-result|%s|%s|proceed\n", token, sessionId)
		return
	}

	go resolveDomain(sessionId, token, parts[1])
}

func resolveDomain(sessionId string, token string, domain string) {
	_, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("filter-result|%s|%s|reject|550 unknown sender domain\n", token, sessionId)
		return
	}

	fmt.Printf("filter-result|%s|%s|proceed\n", token, sessionId)
}

func filterInit() {
	for k := range filters {
		fmt.Printf("register|filter|smtp-in|%s\n", k)
	}
	fmt.Println("register|ready")	
}

func trigger(currentSlice map[string]func(string, []string), atoms []string) {
	found := false
	for k, v := range currentSlice {
		if k == atoms[4] {
			v(atoms[5], atoms[6:])
			found = true
			break
		}
	}
	if !found {
		os.Exit(1)
	}
}

func skipConfig(scanner *bufio.Scanner) {
	for {
		if !scanner.Scan() {
			os.Exit(0)
		}
		line := scanner.Text()
		if line == "config|ready" {
			return
		}
	}
}

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	skipConfig(scanner)
	filterInit()

	for {
		if !scanner.Scan() {
			os.Exit(0)
		}
		
		atoms := strings.Split(scanner.Text(), "|")
		if len(atoms) < 6 {
			os.Exit(1)
		}

		if atoms[0] != "filter" {
			os.Exit(1)
		}

		trigger(filters, atoms)
	}
}
