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
	"fmt"
	"os"

	"github.com/martinnirtl/hosts-cli/internal/helpers"
	"github.com/martinnirtl/hosts-cli/pkg/files"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm HOST...",
	Short: "Remove one or more host entries from ssh-config and hosts file",
	Long:  `Remove one or more host entries from ssh-config and hosts file. Gonna keep those files clean!`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string
		if len(args) == 0 {
			comps = cobra.AppendActiveHelp(comps, "Expecting one or more host names")
		}
		return comps, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostsFilePath, _ := cmd.PersistentFlags().GetString("hosts-file")
		if hostsFilePath == "" {
			hostsFilePath = "/etc/hosts"
		}
		hosts, err := files.GetHosts(hostsFilePath)
		if err != nil {
			cmd.Printf("Error reading file: %v", err)

			os.Exit(1)
		}

		hosts.RemoveHosts(args)

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		if !dryRun {
			if err := hosts.Write(); err != nil {
				cmd.Printf("Error writing file %s: %v", hostsFilePath, err)

				os.Exit(1)
			}
		}

		if dryRun {
			cmd.Print(helpers.Header(fmt.Sprintf("%s:", hostsFilePath), ""))
			cmd.Print(hosts)
		}

		sshConfigPath, _ := cmd.Flags().GetString("ssh-config")
		if sshConfigPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				cmd.Printf("Error retrieving user's home directory: %v", err)

				os.Exit(1)
			}
			sshConfigPath = fmt.Sprintf("%s/.ssh/config", homeDir)
		}
		sshConfig, err := files.GetSSHConfig(sshConfigPath)
		if err != nil {
			cmd.Printf("Error reading file: %v", err)

			os.Exit(1)
		}

		sshConfig.RemoveHosts(args)

		if !dryRun {
			sshConfig.Write()
		}

		if dryRun {
			cmd.Print(helpers.Header(fmt.Sprintf("%s:", sshConfigPath), "\n--\n"))
			cmd.Print(sshConfig)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
