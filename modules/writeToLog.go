package modules

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func printWithTimestamp(file *os.File, IFFID, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, IFFID, message)

	// Print to console
	fmt.Println(formattedMessage)

	// Append to the CSV file
	csvWriter := csv.NewWriter(file)
	err := csvWriter.Write([]string{timestamp, IFFID, message})
	if err != nil {
		fmt.Println("Error writing to CSV file:", err)
		return
	}
	csvWriter.Flush()
}
