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
	"os"
	"os/signal"
	"syscall"
)

// DeviceState This is the main function of the devcie-state-service.
// The following tasks are executed:
// * Setup signal handler
// * Startup processes
// * Wait for processes to terminate
func DeviceState(gpioChip string, lineNr int, checkIntervalMs int) {

	// Create channel for transmiision of connection state to led service
	connectionStateChannel := make(chan bool)

	// Create LED service instance
	ledService, err := NewLedService(connectionStateChannel, gpioChip, lineNr)
	if err != nil {
		// Nothing needs to be closed before termination
		// print error and terminate
		fmt.Println(err)
		fmt.Println("device-state-service terminated.")
		os.Exit(1)
	}

	// Create connection state instance
	connectionState := NewConnectionState(connectionStateChannel, checkIntervalMs)

	// Change led status depending on connection state in goroutine
	go ledService.Run()

	// Observe connection status in goroutine
	go connectionState.Run()

	// wait for termination signal
	signalHandler()

	// terminate goroutines and wait for their termination
	ledService.Close()
	connectionState.Close()

	fmt.Println("device-state-service terminated.")
}

// signalHandler Wait for termination signal (SIGINT and SIGTERM) to come up.
func signalHandler() {
	// setup signal catching
	signChan := make(chan os.Signal, 1)

	// catch listed signals SIGINT and SIGTERM
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	<-signChan
	fmt.Println("Termination signal received.")
}
