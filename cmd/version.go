/*
Copyright © 2023 Martin Nirtl <martin.nirtl@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	version = ""
	commit  = ""
	date    = ""
	github  = "github.com/martinnirtl/hosts-cli"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version information",
	Long:  `Print CLI version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			cmd.Printf("Hosts CLI installed via 'go install %s'\n\nView on GitHub > https://%s\n", github, github)

			return
		}

		date := time.Now().UTC().Format(time.ANSIC)
		cmd.Printf("Hosts CLI Version %s built %s (commit: %s)\n\nView on GitHub > %s\n", version, date, commit, github)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
