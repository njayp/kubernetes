package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/njayp/ophis/bridge"
	"github.com/njayp/ophis/tools"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/kubectl/pkg/cmd/plugin"
)

type CmdFactory struct {
	rootCmd *cobra.Command
}

// RegistrationCommand creates a Helm command tree for MCP tool registration.
func (f *CmdFactory) Tools() []tools.Tool {
	return tools.FromRootCmd(f.rootCmd)
}

// New creates a fresh Helm command instance and its execution function.
func (f *CmdFactory) New() (*cobra.Command, bridge.CommandExecFunc) {
	var output strings.Builder
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: &output, ErrOut: &output}
	rootCmd := NewDefaultKubectlCommandWithArgs(KubectlOptions{
		PluginHandler: NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     os.Args,
		ConfigFlags:   defaultConfigFlags().WithWarningPrinter(ioStreams),
		IOStreams:     ioStreams,
	})

	exec := func(ctx context.Context, cmd *cobra.Command) *mcp.CallToolResult {
		err := cmd.ExecuteContext(ctx)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("execution error:", err)
		}
		return mcp.NewToolResultText(output.String())
	}

	return rootCmd, exec
}
