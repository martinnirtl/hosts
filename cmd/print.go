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
	"os"

	"github.com/martinnirtl/hosts-cli/internal/helpers"
	"github.com/martinnirtl/hosts-cli/pkg/files"
	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print contents of ssh-config and hosts file",
	Long:  `Print contents of ssh-config and hosts file. Gonna stay up to date on their content!`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string
		if len(args) >= 0 {
			comps = cobra.AppendActiveHelp(comps, "No args expected - hit it!")
		}
		return comps, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.ExactArgs(0),
	Run:  Print,
}

func init() {
	rootCmd.AddCommand(printCmd)
}

func Print(cmd *cobra.Command, args []string) {
	err := getFilePaths()
	if err != nil {
		cmd.Printf("Error retrieving file paths: %v", err)

		os.Exit(1)
	}

	if etcHosts {
		hosts, err := files.GetHosts(hostsFilePath)
		if err != nil {
			cmd.Printf("Error reading file: %v", err)

			os.Exit(1)
		}

		cmd.Print(helpers.PrintFileWithSpacer(hostsFilePath, hosts))
	}

	sshConfig, err := files.GetSSHConfig(sshConfigFilePath)
	if err != nil {
		cmd.Printf("Error reading file: %v", err)

		os.Exit(1)
	}

	cmd.Print(helpers.PrintFile(sshConfigFilePath, sshConfig))
}
