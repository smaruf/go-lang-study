package main

import (
	"fmt"
	
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("User Registration Form")

	// Create form inputs
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Enter your email")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter your password")

	ageEntry := widget.NewEntry()
	ageEntry.SetPlaceHolder("Enter your age")

	// Create dropdown for country
	countrySelect := widget.NewSelect(
		[]string{"USA", "Canada", "UK", "Australia", "Other"},
		func(value string) {
			fmt.Println("Selected country:", value)
		},
	)

	// Create checkbox
	termsCheck := widget.NewCheck("I agree to terms and conditions", func(checked bool) {
		fmt.Println("Terms accepted:", checked)
	})

	// Create result label
	resultLabel := widget.NewLabel("")

	// Create submit button
	submitButton := widget.NewButton("Submit", func() {
		// Validate inputs
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("name is required"), myWindow)
			return
		}
		if emailEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("email is required"), myWindow)
			return
		}
		if passwordEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("password is required"), myWindow)
			return
		}
		if !termsCheck.Checked {
			dialog.ShowError(fmt.Errorf("you must accept terms and conditions"), myWindow)
			return
		}

		// Show success message
		message := fmt.Sprintf("Registration Successful!\nName: %s\nEmail: %s\nAge: %s\nCountry: %s",
			nameEntry.Text, emailEntry.Text, ageEntry.Text, countrySelect.Selected)
		
		resultLabel.SetText(message)
		dialog.ShowInformation("Success", message, myWindow)
	})

	// Create reset button
	resetButton := widget.NewButton("Reset", func() {
		nameEntry.SetText("")
		emailEntry.SetText("")
		passwordEntry.SetText("")
		ageEntry.SetText("")
		countrySelect.SetSelected("")
		termsCheck.SetChecked(false)
		resultLabel.SetText("")
	})

	// Create form using Form widget for better layout
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Age", Widget: ageEntry},
			{Text: "Country", Widget: countrySelect},
			{Text: "Terms", Widget: termsCheck},
		},
	}

	// Create button container
	buttonBox := container.NewHBox(
		submitButton,
		resetButton,
	)

	// Create main container
	content := container.NewVBox(
		widget.NewLabelWithStyle("User Registration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
		buttonBox,
		resultLabel,
	)

	// Set content and show window
	myWindow.SetContent(container.NewVScroll(content))
	myWindow.Resize(fyne.NewSize(500, 600))
	myWindow.ShowAndRun()
}
