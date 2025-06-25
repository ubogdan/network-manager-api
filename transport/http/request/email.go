package request

import (
	"github.com/ubogdan/network-manager-api/model"
)

type Email struct {
	Product    string   `json:"product"`
	License    string   `json:"license"`
	Subject    string   `json:"subject"`
	To         string   `json:"to"`
	Name       string   `json:"name"`
	Intros     []string `json:"intros"`
	Dictionary []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"dictionary,omitempty"`
	Table struct {
		Data [][]struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"data"`
		CustomWidth     map[string]string `json:"customWidth,omitempty"`
		CustomAlignment map[string]string `json:"customAlignment,omitempty"`
	} `json:"table,omitempty"`
	Actions []struct {
		Instructions string `json:"instructions"`
		Button       struct {
			Color     string `json:"color"`
			TextColor string `json:"text_color"`
			Text      string `json:"text"`
			Link      string `json:"link"`
		} `json:"button"`
		InviteCode string `json:"invite_code,omitempty"`
	} `json:"actions,omitempty"`
	Outros       []string `json:"outros"`
	Greeting     string   `json:"greeting"`
	Signature    string   `json:"signature"`
	Title        string   `json:"title"`
	FreeMarkdown string   `json:"free_markdown"`
}

func (e Email) ToModel() model.Email {
	email := model.Email{
		Name:      e.Name,
		Intros:    e.Intros,
		Outros:    e.Outros,
		Greeting:  e.Greeting,
		Signature: e.Signature,
		Title:     e.Title,
		//FreeMarkdown: e.FreeMarkdown,
	}

	// Translate dictionary
	for _, d := range e.Dictionary {
		email.Dictionary = append(email.Dictionary, model.Entry{
			Key:   d.Key,
			Value: d.Value,
		})
	}

	// Translate table
	for _, d := range e.Table.Data {
		var data []model.Entry
		for _, e := range d {
			data = append(data, model.Entry{
				Key:   e.Key,
				Value: e.Value,
			})
		}
		email.Table.Data = append(email.Table.Data, data)
	}

	// Translate custom width
	for k, v := range e.Table.CustomWidth {
		email.Table.Columns.CustomWidth[k] = v
	}

	// Translate custom alignment
	for k, v := range e.Table.CustomAlignment {
		email.Table.Columns.CustomAlignment[k] = v
	}

	// Translate actions
	for _, a := range e.Actions {
		email.Actions = append(email.Actions, model.Action{
			Instructions: a.Instructions,
			Button: model.Button{
				Color:     a.Button.Color,
				TextColor: a.Button.TextColor,
				Text:      a.Button.Text,
				Link:      a.Button.Link,
			},
			InviteCode: a.InviteCode,
		})
	}

	return email
}
