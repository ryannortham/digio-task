/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// analyseCmd represents the analyse command
var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.GetString("log-dir")
		fileName := viper.GetString("log-file")
		filePath := filepath.Join(dir, fileName)
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening the file: %v\n", err)
			return
		}
		defer file.Close() // Ensure the file is closed when we're done

		// Read the file's contents
		contents, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error reading the file: %v\n", err)
			return
		}

		// Print the file's contents as a string
		fmt.Printf(string(contents))
	},
}

func init() {
	rootCmd.AddCommand(analyseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analyseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analyseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
