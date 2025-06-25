package model

import (
	"html/template"

	"github.com/russross/blackfriday/v2"
)

// Email godoc.
type Email struct {
	Name         string   // The name of the contacted person
	Intros       []string // Intro sentences, first displayed in the email
	Dictionary   []Entry  // A list of key+value (useful for displaying parameters/settings/personal info)
	Table        Table    // Table is an table where you can put data (pricing grid, a bill, and so on)
	Actions      []Action // Actions are a list of actions that the user will be able to execute via a button click
	Outros       []string // Outro sentences, last displayed in the email
	Greeting     string   // Greeting for the contacted person (default to 'Hi')
	Signature    string   // Signature for the contacted person (default to 'Yours truly')
	Title        string   // Title replaces the greeting+name when set
	FreeMarkdown Markdown // Free markdown content that replaces all content other than header and footer
}

// Entry is a simple entry of a map
// Allows using a slice of entries instead of a map
// Because Golang maps are not ordered
type Entry struct {
	Key   string
	Value string
}

// Table is an table where you can put data (pricing grid, a bill, and so on)
type Table struct {
	Data    [][]Entry // Contains data
	Columns Columns   // Contains meta-data for display purpose (width, alignement)
}

// Columns contains meta-data for the different columns
type Columns struct {
	CustomWidth     map[string]string
	CustomAlignment map[string]string
}

// Action is anything the user can act on (i.e., click on a button, view an invite code)
type Action struct {
	Instructions string
	Button       Button
	InviteCode   string
}

// Button defines an action to launch
type Button struct {
	Color     string
	TextColor string
	Text      string
	Link      string
}

type Markdown template.HTML

// ToHTML converts Markdown to HTML
func (c Markdown) ToHTML() template.HTML {
	return template.HTML(blackfriday.Run([]byte(string(c))))
}
