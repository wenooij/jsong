package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wenooij/jsong"
)

var extractFlags struct {
	Input  string
	Format string
	Path   string
}

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Print extracted JSON values",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(extractFlags.Input)
		if err != nil {
			return fmt.Errorf("failed to read from file: %v", err)
		}
		defer f.Close()
		var v any
		switch strings.ToLower(extractFlags.Format) {
		case "", "json":
			v, err = jsong.NewDecoder(f).Decode()
			if err != nil {
				return fmt.Errorf("failed to decode file: %v", err)
			}
		default:
			return fmt.Errorf("unexpected format: %q", extractFlags.Format)
		}

		res := jsong.Extract(v, extractFlags.Path)
		data, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed to marshal output: %v", err)
		}

		fmt.Println(string(data))
		return nil
	},
}

func init() {
	fs := extractCmd.Flags()
	fs.StringVarP(&extractFlags.Input, "input", "i", "", "Input file name")
	fs.StringVarP(&extractFlags.Format, "format", "f", "json", "Input file format")
	fs.StringVarP(&extractFlags.Path, "path", "p", "", "Path to extract")
	extractCmd.MarkFlagRequired("input")
}
