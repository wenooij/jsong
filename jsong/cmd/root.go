package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/cobra"
)

var cpuProfile *os.File

var rootFlags = struct {
	CPUProfile string
	MemProfile string
}{}

var rootCmd = &cobra.Command{
	Use:   "jsong",
	Short: "jsong executes JSON operations and aggregations",
	Long:  `jsong executes JSON operations and aggregations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if rootFlags.CPUProfile != "" {
			f, err := os.Create(rootFlags.CPUProfile)
			if err != nil {
				return fmt.Errorf("failed to create file for CPU profile: %v", err)
			}
			cpuProfile = f
			pprof.StartCPUProfile(cpuProfile)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if cpuProfile != nil {
			pprof.StopCPUProfile()
			cpuProfile.Close()
		}
		if rootFlags.MemProfile != "" {
			f, err := os.Create(rootFlags.MemProfile)
			if err != nil {
				return fmt.Errorf("failed to create file for Mem profile: %v", err)
			}
			defer f.Close()
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}
		return nil
	},
}

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVar(&rootFlags.CPUProfile, "cpuprofile", "", "CPU profile")
	fs.StringVar(&rootFlags.MemProfile, "memprofile", "", "Mem profile")
	rootCmd.AddCommand(
		extractCmd,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
