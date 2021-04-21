/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package devicestate

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/warthog618/gpiod"
)

const (
	// baseTimeMilliSec is the minimum time the led stays in one state
	baseTimeMilliSec = 500
)

// LedService implements an led service indicating the device connection state by
// applying a blink pattern on the device led selected.
//
// LED on: Device is connected.
// LED blinking: Device tries to connect.
// LED off: Service terminated.
type LedService struct {
	closed       chan interface{}
	stateChan    chan bool
	blinkPattern blinkPattern
	chip         *gpiod.Chip
	line         *gpiod.Line
	wg           sync.WaitGroup
}

// NewLedService intialize led service
// GPIO used can be configured
func NewLedService(connectionStateChannel chan bool, gpioChip string, lineNr int) (*LedService, error) {

	chip, err := gpiod.NewChip(gpioChip)
	if err != nil {
		return nil, err
	}

	line, err := chip.RequestLine(lineNr, gpiod.AsOutput(0))
	if err != nil {
		return nil, err
	}

	return &LedService{
		closed:    make(chan interface{}),
		chip:      chip,
		line:      line,
		stateChan: connectionStateChannel,
	}, nil
}

// Close Cleaup function for LedService
func (l *LedService) Close() {

	// terminate all goroutines
	close(l.closed)
	// wait for goroutines to finish
	l.wg.Wait()

	// close gpio ressources
	l.chip.Close()
	err := l.line.Reconfigure(gpiod.AsInput)
	if err != nil {
		log.Println(err)
	}
	l.line.Close()
}

// Run runs the led servie
func (l *LedService) Run() {
	l.wg.Add(1)
	defer l.wg.Done()

	// start with blinkpattern
	go l.controlLed()

	// wait for new data from channel from device state goroutine
	for {
		select {
		case <-l.closed: // close function was called
			l.blinkPattern.changePattern(exit) // set blink pattern to terminate
			return
		case connectionState := <-l.stateChan:
			// depending on the connection state change the blink pattern
			if connectionState {
				l.blinkPattern.changePattern(on)

			} else if !connectionState {
				l.blinkPattern.changePattern(blink)
			}
		}
	}
}

// controlLed goroutine which executes the currently selected blink pattern
func (l *LedService) controlLed() {
	l.wg.Add(1)
	defer l.wg.Done()

	// numberSteps number of blink pattern steps defined
	const numberSteps = 2

	steps := map[blinkPatterns][numberSteps]int{
		off:   {0, 0},
		blink: {0, 1},
		on:    {1, 1},
		exit:  {},
	}
	stepIdx := 0
	curPattern := off
	for {
		pattern := l.blinkPattern.getPattern()

		if pattern == exit {
			err := l.line.SetValue(0)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			break
		}
		if pattern != curPattern {
			stepIdx = 0
			curPattern = pattern
		}

		ledVal := steps[curPattern][stepIdx]
		err := l.line.SetValue(ledVal)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		time.Sleep(baseTimeMilliSec * time.Millisecond)
		stepIdx++
		if stepIdx == numberSteps {
			stepIdx = 0
		}
	}
}
