package files

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/martinnirtl/hosts-cli/internal/helpers"
)

type SSHConfig struct {
	filepath string
	read     bool
	blocks   []*HostBlock
}

type HostBlock struct {
	Kind  string // Host or Match
	Hosts []string
	Props []*HostBlockProp
}

type HostBlockProp struct {
	Kind  string
	Value string
}

type EmptyLineBlock struct {
	Kind  string
	count int
}

type CommentBlock struct {
	Kind    string
	comment string
}

func (prop *HostBlockProp) String() string {
	return fmt.Sprintf("  %s %s", prop.Kind, prop.Value)
}

func (block *HostBlock) String() string {
	output := fmt.Sprintf("%s %s\n", block.Kind, strings.Join(block.Hosts, " "))

	for _, prop := range block.Props {
		output = output + prop.String() + "\n"
	}

	return output
}

func (sshConfig *SSHConfig) String() string {
	stringifiedBlocks := make([]string, len(sshConfig.blocks))

	for i, block := range sshConfig.blocks {
		stringifiedBlocks[i] = block.String()
	}

	return strings.Join(stringifiedBlocks, "\n")
}

func (sshConfig *SSHConfig) Read() error {
	file, err := os.Open(sshConfig.filepath)
	if err != nil {
		return fmt.Errorf("Failed to open '%s': %v", sshConfig.filepath, err)
	}
	defer file.Close()

	rd := bufio.NewReader(file)

	var currentBlock *HostBlock
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("Failed reading file: %v", err)
		} else if err == io.EOF && len(line) == 0 {
			break
		}

		fields := strings.Fields(line)
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue // skip empty lines and comments
		} // TODO think about keeping empty lines and comments. use an interface with field Kind to identify blockType.. make SSHConfig hold array of interface

		key := fields[0]

		switch strings.ToUpper(key) {
		case "HOST", "MATCH":
			currentBlock = &HostBlock{Kind: key, Hosts: fields[1:]}
			// if len(dropList) > 0 {
			// 	for _, host := range fields[1:] {
			// 		if helpers.SliceContains(dropList, host) {
			// 			continue
			// 		}
			// 	}
			// }

			sshConfig.blocks = append(sshConfig.blocks, currentBlock)

		default:
			value := strings.Join(fields[1:], " ")
			currentBlock.Props = append(currentBlock.Props, &HostBlockProp{Kind: key, Value: value})
		}
	}

	return nil
}

func (sshConfig *SSHConfig) ListHosts() [][]string {
	list := make([][]string, len(sshConfig.blocks))
	for i, entry := range sshConfig.blocks {
		list[i] = entry.Hosts
	}

	return list
}

func (sshConfig *SSHConfig) AddHost(hosts []string, hostname string, user string) error {
	configBlockProps := make([]*HostBlockProp, 0)
	configBlockProps = append(configBlockProps, &HostBlockProp{Kind: "HostName", Value: hostname})
	if user != "" {
		configBlockProps = append(configBlockProps, &HostBlockProp{Kind: "User", Value: user})
	}

	configBlock := &HostBlock{
		Kind:  "Host",
		Hosts: hosts,
		Props: configBlockProps,
	}

	// TODO check if there is already such an entry

	sshConfig.blocks = append(sshConfig.blocks, configBlock)

	return nil
}

func (sshConfig *SSHConfig) RemoveHosts(hosts []string) []*HostBlock {
	new := make([]*HostBlock, 0, 10)
	removed := make([]*HostBlock, 0, 10)

	for _, block := range sshConfig.blocks {
		keep := true
		for _, host := range hosts {
			if helpers.SliceContains(block.Hosts, host) {
				removed = append(removed, block)
				keep = false
				break
			}
		}

		if keep {
			new = append(new, block)
		}
	}

	sshConfig.blocks = new

	return removed
}

func (sshConfig *SSHConfig) Write() error {
	file, err := os.OpenFile(sshConfig.filepath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open '%s': %v", sshConfig.filepath, err)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("Failed writing file '%s': %v", file.Name(), err)
	}

	_, err = file.WriteAt([]byte(sshConfig.String()+"\n"), 0)
	if err != nil {
		return fmt.Errorf("Failed writing file '%s': %v", file.Name(), err)
	}

	return nil
}

func GetSSHConfig(filepath string) (*SSHConfig, error) {
	sshConfig := &SSHConfig{
		filepath: filepath,
		blocks:   make([]*HostBlock, 0, 10), // TODO test with capacity 1
	}

	err := sshConfig.Read()

	return sshConfig, err
}
