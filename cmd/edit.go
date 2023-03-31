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
	"os/exec"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [EDITOR]",
	Short: "Edit host entries of ssh-config and optionally hosts file", // TODO fix descriptions
	Long:  `Edit host entries of ssh-config and optionally hosts file. Remember :wq to escape vim!`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string
		if len(args) == 0 {
			comps = cobra.AppendActiveHelp(comps, "Provide an editor or hit enter")
		}
		if len(args) == 1 {
			comps = cobra.AppendActiveHelp(comps, "Hit it!")
		}
		if len(args) > 1 {
			comps = cobra.AppendActiveHelp(comps, "Too many arguments specified!")
		}
		return comps, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")
		if len(args) == 1 {

			if _, err := exec.LookPath(args[0]); err != nil {
				cmd.Printf("Executable '%s' not found in $PATH. Try nano or vi!\n", args[0])

				os.Exit(1)
			}

			editor = args[0]
		} else if editor == "" {
			editor = "vi"
		}

		err := getFilePaths()
		if err != nil {
			cmd.Printf("Error retrieving file paths: %v", err)

			os.Exit(1)
		}

		if etcHosts {
			vi := exec.Command(editor, hostsFilePath)
			vi.Stdin = os.Stdin
			vi.Stdout = os.Stdout
			if err := vi.Start(); err != nil {
				cmd.Printf("Error opening file with vi: %v", err)

				os.Exit(1)
			}
			if err := vi.Wait(); err != nil {
				cmd.Printf("Unexpected error occurred: %v", err)

				os.Exit(1)
			}
		}

		vi := exec.Command(editor, sshConfigFilePath)
		vi.Stdin = os.Stdin
		vi.Stdout = os.Stdout
		if err := vi.Start(); err != nil {
			cmd.Printf("Error opening file with vi: %v", err)

			os.Exit(1)
		}
		if err := vi.Wait(); err != nil {
			cmd.Printf("Unexpected error occurred: %v", err)

			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
