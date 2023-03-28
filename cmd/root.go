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
	"strings"

	"github.com/martinnirtl/addh/internal/helpers"
	"github.com/martinnirtl/addh/pkg/files"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "addh [hosts...] [alias or IP]",
	Short: "Manage host/ip mappings of '/etc/hosts' and '~/.ssh/config'",
	Long:  `Manage host/ip mappings of '/etc/hosts' and '~/.ssh/config' with one command.`,
	Run: func(cmd *cobra.Command, args []string) {
		hostsFilePath := "/etc/hosts"
		hosts, err := files.GetHosts(hostsFilePath)
		if err != nil {
			cmd.Printf("Error reading file: %v", err)

			os.Exit(1)
		}

		remove, _ := cmd.Flags().GetBool("remove")
		if remove && len(args) > 0 {
			hosts.RemoveHosts(args)
		} else if len(args) > 1 {
			hosts.AddHost(args[0:len(args)-1], args[len(args)-1])
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		if len(args) > 0 && !dryRun {
			if err := hosts.Write(); err != nil {
				cmd.Printf("Error writing file %s: %v", hostsFilePath, err)

				os.Exit(1)
			}
		}

		listHosts, _ := cmd.Flags().GetBool("list-hosts")
		if listHosts {
			hostList := hosts.ListHosts()
			for i, hosts := range hostList {
				cmd.Printf("%d: %s\n", i, strings.Join(hosts, " "))
			}

			os.Exit(0)
		}

		if dryRun || len(args) == 0 {
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

		if remove && len(args) > 0 {
			sshConfig.RemoveHosts(args)
		} else if len(args) > 1 {
			user, _ := cmd.Flags().GetString("user")

			sshConfig.AddHost(args[0:len(args)-1], args[len(args)-1], user)
		}

		if len(args) > 0 && !dryRun {
			sshConfig.Write()
		}

		if listHosts {
			hostList := sshConfig.ListHosts()
			for i, hosts := range hostList {
				cmd.Printf("%d: %s\n", i, strings.Join(hosts, " "))
			}

			os.Exit(0)
		}

		if dryRun || len(args) == 0 {
			cmd.Print(helpers.Header(fmt.Sprintf("%s:", sshConfigPath), "\n--\n"))
			cmd.Print(sshConfig)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Bool("dry-run", false, "Only print updated `/etc/hosts` and `~/.ssh/config` files")
	rootCmd.Flags().BoolP("list-hosts", "l", false, "List hosts")
	rootCmd.Flags().BoolP("remove", "r", false, "Remove host from files")
	rootCmd.Flags().StringP("user", "u", "", "Set 'User' for SSH config file")
	rootCmd.Flags().String("ssh-config", "", "Set SSH Config file (e.g. /etc/ssh/config). Default: ~/.ssh/config")
}
