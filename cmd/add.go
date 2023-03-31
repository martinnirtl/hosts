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

var (
	user         string
	identityFile string
	// importIdentityFilesGlob string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add ADDRESS ALIASES...",
	Short: "Add address mappings to ssh-config and hosts file",
	Long: `Add address mappings to ssh-config and hosts file. Address can be an IP or a domain. 
  Makes your life easier!
    Don't forget the sudo!`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string
		if len(args) == 0 {
			comps = cobra.AppendActiveHelp(comps, "Expecting address/IP")
		} else if len(args) == 1 {
			comps = cobra.AppendActiveHelp(comps, "Expecting one or more host names")
		} else {
			comps = cobra.AppendActiveHelp(comps, "Expecting host names or enter key")
		}
		return comps, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
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

			hosts.AddHost(args[0], args[1:])

			if !dryRun {
				if err := hosts.Write(); err != nil {
					cmd.Printf("Error writing file %s: %v", hostsFilePath, err)

					os.Exit(1)
				}
			}

			if dryRun {
				cmd.Print(helpers.PrintFileWithSpacer(hostsFilePath, hosts))
			}
		}

		sshConfig, err := files.GetSSHConfig(sshConfigFilePath)
		if err != nil {
			cmd.Printf("Error reading file: %v", err)

			os.Exit(1)
		}

		sshConfig.AddHost(args[1:], args[0], user, identityFile)

		if !dryRun {
			sshConfig.Write()
		}

		if dryRun {
			cmd.Print(helpers.PrintFile(sshConfigFilePath, sshConfig))
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	flags := addCmd.Flags()

	flags.StringVarP(&user, "user", "u", "", "Set User property in SSH config Host block")
	flags.StringVarP(&identityFile, "identity-file", "i", "", "Use identity file; e.g. ~/.ssh/custom")
	// flags.StringVarP(&importIdentityFilesGlob, "import-idenity-files-glob", "j", "", "Import and use identity file; moves file to ~/.ssh/")
}
