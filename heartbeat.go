// Copyright 2022 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/lextoumbourou/idle"
	"github.com/micmonay/keybd_event"
)

const (
	key2press         = 0x10 + 0xFFF // keybd_event._VK_SHIFT
	maxSleepInMinutes = 120
)

var (
	wait   = time.Minute * 5
	verbal = false
)

func main() {

	for _, arg := range os.Args[1:] {
		if arg == "verbal" {
			verbal = true
		} else if input, err := strconv.ParseInt(arg, 10, 32); err == nil && input > 0 && input <= maxSleepInMinutes {
			wait = time.Minute * time.Duration(input)
		}
	}

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Select keys to be pressed
	kb.SetKeys(key2press) // keybd_event.VK_A, keybd_event.VK_B)

	// Set shift to be pressed
	kb.HasSHIFT(true)

	display("check frequency:", wait)
	for {
		idleTime, err := idle.Get()
		if err != nil {
			panic(err)
		}

		display("    idle time:", idleTime)

		sleep := wait
		if idleTime < wait {
			sleep -= idleTime
		} else {
			display("pressing")
			// Press the selected keys
			err = kb.Launching()
			if err != nil {
				panic(err)
			}
		}

		display("    sleeping for:", sleep)
		time.Sleep(sleep)
	}

	// Or you can use Press and Release
	kb.Press()
	time.Sleep(10 * time.Millisecond)
	kb.Release()
}

func display(a ...interface{}) {
	if verbal {
		_, _ = fmt.Println(a...)
	}
}
