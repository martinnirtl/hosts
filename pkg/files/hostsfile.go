package files

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/martinnirtl/addh/internal/helpers"
)

var (
	requiredHostNames = []string{
		"localhost",
		"broadcasthost",
	}
	requiredHostEntries = []string{
		"127.0.0.1 localhost",
		"255.255.255.255 broadcasthost",
		"::1 localhost",
	}
)

type Hosts struct {
	filepath    string
	hostEntries []*Host
}

type Host struct {
	address string
	aliases []string
	comment string
}

func (host *Host) String() string {
	output := host.address + " " + strings.Join(host.aliases, " ")
	if host.comment != "" {
		output = output + " " + host.comment
	}

	return output
}

func (hosts *Hosts) String() string {
	sort.Slice(hosts.hostEntries, func(i, j int) bool {
		if hosts.hostEntries[i].address == "127.0.0.1" || hosts.hostEntries[i].address == "255.255.255.255" || hosts.hostEntries[i].address == "::1" {
			return true // always print on top
		}

		return hosts.hostEntries[j].address > hosts.hostEntries[i].address
	})

	output := make([]string, len(hosts.hostEntries))
	for i, entry := range hosts.hostEntries {
		output[i] = entry.String()
	}

	return strings.Join(output, "\n") + "\n"
}

func (hosts *Hosts) Read() error {
	file, err := os.Open(hosts.filepath)
	if err != nil {
		return fmt.Errorf("Failed to open '%s': %v", hosts.filepath, err)
	}
	defer file.Close()

	rd := bufio.NewReader(file)

	var currentHost *Host
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("Failed reading file: %v", err)
		} else if err == io.EOF && len(line) == 0 {
			break
		}

		entry, comment, _ := strings.Cut(line, "#")

		fields := strings.Fields(entry)
		if len(fields) == 0 {
			continue // skip empty lines
		}

		currentHost = &Host{address: fields[0], aliases: fields[1:], comment: comment}
		hosts.hostEntries = append(hosts.hostEntries, currentHost)
	}

	return nil
}

func (hosts *Hosts) ListHosts() [][]string {
	list := make([][]string, len(hosts.hostEntries))
	for _, entry := range hosts.hostEntries {
		list = append(list, entry.aliases)
	}

	return list
}

func (h *Hosts) AddHost(hosts []string, ipOrAlias string) error {
	host := &Host{
		address: ipOrAlias,
		aliases: hosts,
		comment: "",
	}

	// TODO check if there is already such an entry by ipOrAlias and extend it

	h.hostEntries = append(h.hostEntries, host)

	return nil
}

func (h *Hosts) RemoveHosts(hosts []string) []*Host {
	new := make([]*Host, 0, 10)
	removed := make([]*Host, 0, 10)

	// filter out required hostnames like localhost and broadcasthost
	removeHosts := make([]string, 10)
	for _, host := range hosts {
		if !helpers.SliceContains(requiredHostNames, host) {
			removeHosts = append(removeHosts, host)
		} else {
			// TODO print warning that e.g. localhost will not be removed
		}
	}

	for _, entry := range h.hostEntries {
		keep := true
		for _, host := range removeHosts {
			if helpers.SliceContains(entry.aliases, host) {
				removed = append(removed, entry)
				keep = false
				break
			}
		}

		if keep {
			new = append(new, entry)
		}
	}

	h.hostEntries = new

	return removed
}

func (hosts *Hosts) Write() error {
	file, err := os.OpenFile(hosts.filepath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open '%s': %v", hosts.filepath, err)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("Failed writing file '%s': %v", file.Name(), err)
	}

	_, err = file.WriteAt([]byte(hosts.String()+"\n"), 0)
	if err != nil {
		return fmt.Errorf("Failed writing file '%s': %v", file.Name(), err)
	}

	return nil
}

func GetHosts(filepath string) (*Hosts, error) {
	hosts := &Hosts{
		filepath:    filepath,
		hostEntries: make([]*Host, 0, 10), // TODO test with capacity 1
	}

	err := hosts.Read()

	return hosts, err
}
