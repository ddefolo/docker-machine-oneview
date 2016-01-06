package client

import (
	"fmt"

	Cli "github.com/docker/docker/cli"
	flag "github.com/docker/docker/pkg/mflag"
)

// CmdUnpause unpauses all processes within a container, for one or more containers.
//
// Usage: docker unpause CONTAINER [CONTAINER...]
func (cli *DockerCli) CmdUnpause(args ...string) error {
	cmd := Cli.Subcmd("unpause", []string{"CONTAINER [CONTAINER...]"}, Cli.DockerCommands["unpause"].Description, true)
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args, true)

	var errNames []string
	for _, name := range cmd.Args() {
		if err := cli.client.ContainerUnpause(name); err != nil {
			fmt.Fprintf(cli.err, "%s\n", err)
			errNames = append(errNames, name)
		} else {
			fmt.Fprintf(cli.out, "%s\n", name)
		}
	}
	if len(errNames) > 0 {
		return fmt.Errorf("Error: failed to unpause containers: %v", errNames)
	}
	return nil
}
