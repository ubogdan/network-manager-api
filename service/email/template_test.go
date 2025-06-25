package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubogdan/network-manager-api/model"
)

func TestTemplate_GenerateHTML(t *testing.T) {
	template := Template{
		// Optional Theme
		// Theme: new(Default)
		Product: Product{
			Name: "Network Manager Pro",
			Link: "https://lfpanels.com/",
			Logo: "https://lfpanels.com/logo.png",
		},
	}

	email := model.Email{
		Name: "Jon Deer",
		Intros: []string{
			"Welcome! We're very excited to have you on board.",
		},
		Actions: []model.Action{
			{
				Instructions: "To get started with Network Manager, please click here:",
				Button: model.Button{
					Color: "#22BC66", // Optional action button color
					Text:  "Confirm your account",
					Link:  "https://lfpanels.com/confirm?token=d9729feb74992cc3482b350163a1a010",
				},
			},
		},
		Outros: []string{
			"Need help, or have questions? Just reply to this email, we'd love to help.",
		},
	}

	out, err := template.GenerateHTML(email)
	assert.NoError(t, err)

	assert.Contains(t, out, "Network Manager Pro")
	assert.Contains(t, out, "Jon Deer")
	assert.Contains(t, out, "To get started with Network Manager, please click here:")
	assert.Contains(t, out, "Confirm your account")
	assert.Contains(t, out, "https://lfpanels.com/confirm?token=d9729feb74992cc3482b350163a1a010")
}

func TestTemplate_GeneratePlainText(t *testing.T) {
	template := Template{
		// Optional Theme
		// Theme: new(Default)
		Product: Product{
			Name: "Network Manager Pro",
			Link: "https://lfpanels.com/",
			Logo: "https://lfpanels.com/logo.png",
		},
	}

	email := model.Email{
		Name: "Jon Deer",
		Intros: []string{
			"Welcome! We're very excited to have you on board.",
		},
		Actions: []model.Action{
			{
				Instructions: "To get started with Network Manager, please click here:",
				Button: model.Button{
					Color: "#22BC66", // Optional action button color
					Text:  "Confirm your account",
					Link:  "https://lfpanels.com/confirm?token=d9729feb74992cc3482b350163a1a010",
				},
			},
		},
		Outros: []string{
			"Need help, or have questions? Just reply to this email, we'd love to help.",
		},
	}

	out, err := template.GeneratePlainText(email)
	assert.NoError(t, err)

	assert.Contains(t, out, "Network Manager Pro")
	assert.Contains(t, out, "Jon Deer")
	assert.Contains(t, out, "To get started with Network Manager, please click here:")
	assert.Contains(t, out, "Confirm your account")
	assert.Contains(t, out, "https://lfpanels.com/confirm?token=d9729feb74992cc3482b350163a1a010")
}
