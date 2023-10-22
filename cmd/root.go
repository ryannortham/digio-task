package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "digio-task",
	Short: "Parses a log file containing HTTP requests and to reports on its contents",
	Long: `
Parses a log file containing HTTP requests and to reports on its contents

For a given log file we want to know:
  - The number of unique IP addresses
  - The top 3 most visited URLs
  - The top 3 most active IP addresses


         xxxxxx                                    $$$$$$   $$$$$$                          $$$$$$                      
         xxxxxx       :                            $$$$$$   $$$$$$                          $$$$$$                      
        xxxxxx    :::::                            $$$$$$   $$$$$$                          $$$$$$                      
      xxxxxxxx  ::::::::                           $$$$$$                                                               
  xxxxxxxxxxx ::::::::::   ++            $$$$$$$$$ $$$$$$   $$$$$$      $$$$$$$$$$ $$$$$$   $$$$$$       $$$$$$$$$$     
xxxxxxxxxxx  ::::::::    +++++         $$$$$$$$$$$$$$$$$$   $$$$$$    $$$$$$$$$$$$$$$$$$$   $$$$$$    $$$$$$$$$$$$$$$$  
xxxxxxxxx    ::::::     +++++++      $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$ 
xxxxxx      ::::::     ++++++++      $$$$$$       $$$$$$$   $$$$$$  $$$$$$$      $$$$$$$$   $$$$$$  $$$$$$$      $$$$$$$
            ::::::    +++++++       $$$$$$         $$$$$$   $$$$$$  $$$$$$         $$$$$$   $$$$$$  $$$$$$        $$$$$$
   ;;;      ::::::   +++++++        $$$$$$         $$$$$$   $$$$$$  $$$$$$         $$$$$$   $$$$$$  $$$$$          $$$$$
 ;;;;;;     ::::::   ++++++         $$$$$$         $$$$$$   $$$$$$  $$$$$$        $$$$$$$   $$$$$$  $$$$$$        $$$$$$
 ;;;;;;;;   ::::     ++++++          $$$$$$       $$$$$$$   $$$$$$   $$$$$$$$$ $$$$$$$$$$   $$$$$$  $$$$$$$      $$$$$$$
  ;;;;;;;;;          ++++++          $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$ 
    ;;;;;;;;;;;;;;   +++++++           $$$$$$$$$$$$$$$$$$   $$$$$$     $$$$$$$$$$$$$$$$$$   $$$$$$    $$$$$$$$$$$$$$$$  
     ;;;;;;;;;;;;;;   +++++              $$$$$$$$$ $$$$$$   $$$$$$        $$$$$$   $$$$$$   $$$$$$       $$$$$$$$$$$    
        ;;;;;;;;;;;    +                                                $$$       $$$$$$                                
                                                                      $$$$$$$$$$$$$$$$$$                                
                                                                       $$$$$$$$$$$$$$$$                                 
                                                                         $$$$$$$$$$$$                                   


`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s", err)
	}

	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
}
