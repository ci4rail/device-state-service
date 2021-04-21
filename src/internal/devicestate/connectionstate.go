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
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	nsLookupURI = "portal.azure.com"
)

// ConnectionState implements a service which continously checks the connection state
// and provides the information via channel connectionStateChannel
type ConnectionState struct {
	closed          chan interface{}
	stateChan       chan bool
	wg              sync.WaitGroup
	checkIntervalMs int
}

// NewConnectionState intialize connection state
func NewConnectionState(connectionStateChannel chan bool, checkIntervalMs int) *ConnectionState {
	return &ConnectionState{
		stateChan:       connectionStateChannel,
		closed:          make(chan interface{}),
		checkIntervalMs: checkIntervalMs,
	}
}

// Run runs the connection state service
func (c *ConnectionState) Run() {

	// add this goroutine to waitgroups
	c.wg.Add(1)
	defer c.wg.Done()

	// pre-initialize connection state
	connState := false

	for {
		select {
		case <-c.closed: // close function was called
			// terminate goroutine
			close(c.stateChan)
			return
		default:
			newConnState := getNetworkStatus()
			if newConnState != connState {
				connState = newConnState
				// Log connection state change
				if connState {
					fmt.Println("Device switched to connected state.")
				} else {
					fmt.Println("Device switched to disconnected state.")
				}
				// Provide connection change via channel
				c.stateChan <- connState
			}
		}
		time.Sleep(time.Duration(c.checkIntervalMs) * time.Millisecond)
	}
}

// Close Cleaup function for LedService
func (c *ConnectionState) Close() {

	// terminate goroutine
	close(c.closed)
	// wait for goroutine to finish
	c.wg.Wait()
}

// getNetworkStatus checks if a valid internet connection is established.
// This is checked by executing a dns lookup.
func getNetworkStatus() (status bool) {
	_, err := net.LookupIP(nsLookupURI)
	if err != nil {
		status = false
	} else {
		status = true
	}
	return
}
