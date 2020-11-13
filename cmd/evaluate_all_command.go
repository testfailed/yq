package cmd

import (
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func createEvaluateAllCommand() *cobra.Command {
	var cmdEvalAll = &cobra.Command{
		Use:     "eval-all [expression] [yaml_file1]...",
		Aliases: []string{"ea"},
		Short:   "Loads all yaml documents of all yaml files and runs expression once",
		Example: `
yq es '.a.b | length' file1.yml file2.yml
yq es < sample.yaml
yq es -n '{"a": "b"}'
`,
		Long: "Evaluate All:\nUseful when you need to run an expression across several yaml documents or files. Consumes more memory than eval-seq",
		RunE: evaluateAll,
	}
	return cmdEvalAll
}
func evaluateAll(cmd *cobra.Command, args []string) error {
	// 0 args, read std in
	// 1 arg, null input, process expression
	// 1 arg, read file in sequence
	// 2+ args, [0] = expression, file the rest

	var err error
	stat, _ := os.Stdin.Stat()
	pipingStdIn := (stat.Mode() & os.ModeCharDevice) == 0

	out := cmd.OutOrStdout()

	fileInfo, _ := os.Stdout.Stat()

	if forceColor || (!forceNoColor && (fileInfo.Mode()&os.ModeCharDevice) != 0) {
		colorsEnabled = true
	}
	printer := yqlib.NewPrinter(out, outputToJSON, unwrapScalar, colorsEnabled, indent, !noDocSeparators)

	switch len(args) {
	case 0:
		if pipingStdIn {
			err = yqlib.EvaluateAllFileStreams("", []string{"-"}, printer)
		} else {
			cmd.Println(cmd.UsageString())
			return nil
		}
	case 1:
		if nullInput {
			err = yqlib.EvaluateAllFileStreams(args[0], []string{}, printer)
		} else {
			err = yqlib.EvaluateAllFileStreams("", []string{args[0]}, printer)
		}
	default:
		err = yqlib.EvaluateAllFileStreams(args[0], args[1:len(args)], printer)
	}

	cmd.SilenceUsage = true
	return err
}
