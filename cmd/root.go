package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "r1wallet",
	Short: "application that handles customers wallet",
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
