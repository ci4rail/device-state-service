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

package cmd

import (
	"fmt"
	"os"

	"github.com/ci4rail/device-state-service/internal/devicestate"
	"github.com/spf13/cobra"
)

const (
	defaultGpioChip        = "gpiochip4"
	defaultLineNr          = 26
	defaultCheckIntervalMs = 500
)

var (
	gpioChip        string
	lineNr          int
	checkIntervalMs int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "device-state-service",
	Short: "Read and provide edge device status",
	Long: `Read and provide edge device status

Device status is indicated via LED on edge device:
* LED on: Successful connection to cloud
* LED blinking (1 sec): Device tries to connect to cloud
* LED off: device-state-service terminated
`,
	Run: func(cmd *cobra.Command, args []string) {

		// start main function
		devicestate.DeviceState(gpioChip, lineNr, checkIntervalMs)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().StringVarP(&gpioChip, "chip", "c", defaultGpioChip, "use alternative GPIO chip / bank")
	rootCmd.Flags().IntVarP(&lineNr, "line", "l", defaultLineNr, "use alternative GPIO chip line")
	rootCmd.Flags().IntVarP(&checkIntervalMs, "interval", "i", defaultCheckIntervalMs, "interval in milliseconds for checking connection state")
}
