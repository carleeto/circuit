// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gocircuit/circuit/client"

	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

// circuit mkproc /X1234/hola/charlie << EOF
// { … }
// EOF
// TODO: Proc element disappears if command misspelled and error condition not obvious.
func mkproc(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("mkproc needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	buf, _ := ioutil.ReadAll(os.Stdin)
	var cmd client.Cmd
	if err := json.Unmarshal(buf, &cmd); err != nil {
		fatalf("command json not parsing: %v", err)
	}
	p, err := c.Walk(w).MakeProc(cmd)
	if err != nil {
		fatalf("mkproc error: %s", err)
	}
	ps := p.Peek()
	if ps.Exit != nil {
		fatalf("%v", ps.Exit)
	}
}

func sgnl(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		fatalf("signal needs an anchor and a signal name arguments")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Proc)
	if !ok {
		fatalf("not a process")
	}
	if err := u.Signal(args[1]); err != nil {
		fatalf("signal error: %v", err)
	}
}

func wait(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("wait needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Proc)
	if !ok {
		fatalf("not a process")
	}
	stat, err := u.Wait()
	if err != nil {
		fatalf("wait error: %v", err)
	}
	buf, _ := json.MarshalIndent(stat, "", "\t")
	fmt.Println(string(buf))
}
