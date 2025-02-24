package model

const (
	DefaultLang = "markdown"
)

type Lang struct {
	Name        string // Название языка
	Highlighter string // Соответствующий хайлайтер
}

var Langs = map[string]Lang{
	"go": {
		Name:        "GoLang",
		Highlighter: "go",
	},
	"java": {
		Name:        "Java",
		Highlighter: "text/x-java",
	},
	"json": {
		Name:        "JSON",
		Highlighter: "application/json",
	},
	"markdown": {
		Name:        "Markdown",
		Highlighter: "text/x-markdown",
	},
	"markdown_view": {
		Name:        "Markdown View",
		Highlighter: "text/x-markdown",
	},
	"mysql8": {
		Name:        "MySQL 8",
		Highlighter: "sql",
	},
	"php82": {
		Name:        "PHP 8.2",
		Highlighter: "php",
	},
	"postgres13": {
		Name:        "PostgreSQL 13",
		Highlighter: "sql",
	},
}
