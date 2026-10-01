package main

import (
	"Shop-Data-Engine/notionapi" // Matches `go.mod`
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func writeLastUpdated() {
	// Get current timestamp
	timestamp := time.Now().Format("01/02/2006 - 15:04:05")

	// Write timestamp to a second CSV file
	lastUpdatedFile := "last_updated.csv"
	csvFile, err := os.Create(lastUpdatedFile)
	if err != nil {
		fmt.Println("❌ ERROR: Creating Last Updated CSV file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"Last Updated", timestamp})

	fmt.Printf("✅ Timestamp successfully written to %s!\n", lastUpdatedFile)
}

func main() {
	// Load configuration first
	err := notionapi.LoadConfig()
	if err != nil {
		fmt.Println("❌ Error loading config:", err)
	}

	fmt.Println("🚀 Fetching Notion Data...")
	data, err := notionapi.FetchNotionData()
	if err != nil {
		fmt.Println("❌ Error fetching data from Notion:", err)
		return
	}

	totalEntries := len(data)
	if totalEntries == 0 {
		fmt.Println("❌ No data found in Notion.")
		return
	}

	fmt.Printf("✅ Successfully retrieved %d entries from Notion.\n", totalEntries)

	// Save data to CSV
	csvFileName := "notion_data.csv"
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		fmt.Println("❌ ERROR: Creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	headers := []string{"ID", "Name", "Category", "Score", "Employee", "Property", "Date", "Customer Experience Score"}
	writer.Write(headers)

	// Process each entry
	for i, entry := range data {
		props := entry["properties"].(map[string]interface{})

		idStr := notionapi.GetUniqueID(props)
		nameStr := notionapi.GetName(props, "Name")
		categoryStr := notionapi.GetSelectValue(props, "Category")
		scoreStr := notionapi.GetFloatValue(props, "Score")
		employeeStr := notionapi.GetRollupFormulaString(props, "Employee (As Text)")
		propertyStr := notionapi.GetRollupFormulaString(props, "Property (As Text)")
		dateStr := notionapi.GetDateValue(props, "Date")
		customerExpScoreStr := notionapi.GetFloatValue(props, "Customer Experience Score")

		row := []string{idStr, nameStr, categoryStr, scoreStr, employeeStr, propertyStr, dateStr, customerExpScoreStr}
		writer.Write(row)

		fmt.Printf("📊 Progress: %d/%d (%.2f%%)\n", i+1, totalEntries, float64(i+1)/float64(totalEntries)*100)
	}

	fmt.Println("✅ Data successfully written to notion_data.csv!")

	// Write last updated timestamp
	writeLastUpdated()
}
