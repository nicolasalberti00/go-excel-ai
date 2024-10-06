package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/nicolasalberti00/go-excel-ai/internal/providers"
	"github.com/xuri/excelize/v2"
)

const (
	OpenAIAPIKey = "OPENAI_API_KEY"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	modelName := flag.String("model", "", "Preferred AI model name (e.g., gpt-4o, llama3.1)")
	filePath := flag.String("file", "", "Path to the Excel sheet file")
	sheetName := flag.String("sheetName", "", "Name of the sheet chosen")
	cell := flag.String("cell", "", "Cell to modify (e.g. A1)")
	prompt := flag.String("prompt", "", "AI prompt for generating text")
	flag.Parse()

	apiKey := os.Getenv(OpenAIAPIKey)

	//Input validation
	if *modelName == "" || *filePath == "" || *sheetName == "" || *cell == "" || *prompt == "" {
		fmt.Printf("Error: all flags (model, file, sheetName, cell, prompt) are required")
		return
	}

	ctx := context.Background()

	// Load the Excel file
	f, err := excelize.OpenFile(*filePath)
	if err != nil {
		fmt.Printf("Error opening Excel file: %v", err)
		return
	}

	// Fetch the current value in the specified cell
	currentValue, err := f.GetCellValue(*sheetName, *cell)
	if err != nil {
		fmt.Printf("Error reading cell %s from sheet %s: %v", *sheetName, *cell, err)
		return
	}
	slog.InfoContext(ctx, "Value for chosen cell", "cell", *cell, "currentValue", currentValue)

	// Select AI Model
	aiProvider, err := providers.AIProviderFactory(apiKey, *modelName)
	if err != nil {
		fmt.Printf("Error selecting AI provider: %v", err)
		return
	}

	newValue, err := aiProvider.GenerateText(ctx, *prompt)
	if err != nil {
		fmt.Printf("Error generating text: %v", err)
		return
	}

	err = f.SetCellValue(*sheetName, *cell, newValue)
	if err != nil {
		fmt.Printf("Error when setting cell %s in sheet %s: %v", *cell, *sheetName, err)
		return
	}

	err = f.Save()
	if err != nil {
		fmt.Printf("Error saving Excel file: %v", err)
		return
	}

	slog.InfoContext(ctx, "Value updated successfully", "newValue", newValue)

}
