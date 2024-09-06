package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/xuri/excelize/v2"
)

const (
	excelFilePathEnvKey = "EXCEL_FILE_PATH"
	sheetName           = "Sheet1"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	modelName := flag.String("model", "", "preferred ollama model name")
	flag.Parse()

	ctx := context.Background()

	var llm *openai.LLM
	var err error

	if *modelName == "gpt" {
		llm, err = openai.New()
		if err != nil {
			slog.ErrorContext(ctx, "Error when starting model", "error", err)
			return
		}
	}
	//if *modelName == "llama3.1" {
	//	localLLM, err := ollama.New(ollama.WithModel(*modelName))
	//}

	excelFilePath := os.Getenv(excelFilePathEnvKey)

	file, err := excelize.OpenFile(excelFilePath)
	if err != nil {
		slog.Error("Error while opening excel file", "error", err, "path", excelFilePath)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			slog.Error("Error while closing excel file", "error", err)
			return
		}
	}()

	// Text will be found from excel file
	rows, err := file.GetRows(sheetName)
	slog.Info("Columns length", "length", len(rows))
	for i, row := range rows {
		if len(row) > 0 {
			textValue := row[0]
			response, err := generateStreamedResponse(ctx, llm, textValue)
			err = addResponseToCell(file, response, sheetName, fmt.Sprintf("B%d", i+1))
			if err != nil {
				return
			}
		}
	}

	err = file.Save()
	if err != nil {
		slog.Error("Error while saving sheet", "Error", err)
		return
	}
	slog.Info("Sheet saved correctly")
	return
}

func generateStreamedResponse(ctx context.Context, llm *openai.LLM, textToAnalyze string) (string, error) {

	defaultPrompt := `Sei un ecologo ed esperto di scienze forestali con una vasta conoscenza della lingua italiana.
	Il tuo compito è correggere e migliorare il seguente testo, concentrandoti su:
	1. Accuratezza scientifica nel campo delle scienze forestali
	2. Correttezza grammaticale e sintattica
	3. Chiarezza e fluidità espositiva
	4. Uso appropriato della terminologia specifica del settore
	5. Coerenza logica e strutturale del testo
	6. Scrittura in linguaggio tecnico, specialmente con riferimento a Natura 2000

	All'interno di questo testo sono presenti alcuni acronimi, qui la definizione:
	- ATO, Ambito Territoriale Omogeneo;

	Non è necessaria nessuna interazione o richiesta di correzione, rispondi con ciò che ritieni corretto.
	Non servono annotazioni, deve essere semplicemente dato come risposta il testo che hai corretto.
	La risposta non deve superare le 200 parole, mantenendo i punti principali del testo che analizzerai.
	Eventuali parole in maisucolo vanno trattate come nomi propri, quindi riportate alla formattazione normale.
	Ecco il testo da analizzare e correggere: `

	query := defaultPrompt + textToAnalyze

	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, query,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Printf("%s", chunk)
			return nil
		}))
	if err != nil {
		slog.ErrorContext(ctx, "Error while returning completion", "err", err)
		return "", err
	}

	return answer, nil
}

func addResponseToCell(file *excelize.File, response, sheetName, cellName string) error {
	err := file.SetCellValue(sheetName, cellName, response)
	if err != nil {
		slog.Error("Could not set value for the specified cell", "error", err, "sheetName", sheetName, "cellName", cellName, "value", response)
		return err
	}
	return nil
}
