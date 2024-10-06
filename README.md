# Go-excel-AI

This project aims to be able to edit cells in Excel sheets after an AI prompt. I'm committing to this project whenever I can,
as I intend to use it also as a way to improve my skills.

## Prerequisites

In order for this program to work you will need to have some softwares installed in your machine.

- [Go](https://golang.org/)
- [Ollama](https://ollama.com/)

You can also use your OpenAI API key and use the `gpt-4o-mini` model, or just tweak the `factory.go` file and use whichever model you want.

## Usage

From the `cmd` folder you need to run:

```go
go run main.go --model MODEL_NAME --file FILEPATH --sheetName SHEET_NAME --cell CELL_NAME --prompt 'PROMPT'
```
