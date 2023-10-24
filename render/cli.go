package render

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/viper"

	"github.com/ryannortham/digio-task/log"
)

func PrintAnalysisResults(logAnalysis *log.LogAnalysis) {
	printDigioLogo()

	fmt.Printf("Analysis Results of Log File: %s\n\n", viper.GetString("log-file"))

	fmt.Printf("Unique IP addresses: %d\n\n", logAnalysis.UniqueIPCount)

	fmt.Printf("Top %d most visited URLs:\n", viper.GetInt("top-n"))
	printTable(logAnalysis.TopNMostVisitedURLs)

	fmt.Printf("Top %d most active IPs:\n", viper.GetInt("top-n"))
	printTable(logAnalysis.TopNMostActiveIPs)
}

func printTable(results [][]string) {
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgHiBlue).SprintfFunc()
	tbl := table.New(results[0][0], results[0][1])
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i, row := range results {
		if i == 0 {
			continue
		}

		count, err := log.ParseInt(row[1])
		if err != nil {
			// this should never happen
			fmt.Printf("error parsing count: %v\n", err)
		}

		tbl.AddRow(row[0], count)
	}

	tbl.Print()
	fmt.Println()
}

// printDigioLogo prints the Digio logo with colors.
func printDigioLogo() {
	const digioLogo = `
    
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

`

	// Define the colors
	colors := map[rune]*color.Color{
		'x': color.New(color.FgHiMagenta),
		';': color.New(color.FgHiRed),
		':': color.New(color.FgHiYellow),
		'+': color.New(color.FgHiGreen),
		'$': color.New(color.FgHiBlue),
	}

	// Print the ASCII art with colors
	for _, line := range digioLogo {
		for _, char := range string(line) {
			if color, ok := colors[char]; ok {
				color.Print(string(char))
			} else {
				fmt.Print(string(char))
			}
		}
	}
}
