package main

import (
	"fmt"
	"strconv"
	
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Calculator")

	// Display
	display := widget.NewEntry()
	display.SetText("0")
	display.Disable() // Make it read-only

	// Calculator state
	var firstNumber float64
	var operator string
	var newNumber bool = true

	// Helper functions
	appendToDisplay := func(digit string) {
		if newNumber {
			display.SetText(digit)
			newNumber = false
		} else {
			if display.Text == "0" {
				display.SetText(digit)
			} else {
				display.SetText(display.Text + digit)
			}
		}
	}

	calculate := func() {
		if operator == "" {
			return
		}

		current, err := strconv.ParseFloat(display.Text, 64)
		if err != nil {
			display.SetText("Error")
			return
		}

		var result float64
		switch operator {
		case "+":
			result = firstNumber + current
		case "-":
			result = firstNumber - current
		case "*":
			result = firstNumber * current
		case "/":
			if current == 0 {
				display.SetText("Error: Div by 0")
				operator = ""
				firstNumber = 0
				newNumber = true
				return
			}
			result = firstNumber / current
		}

		display.SetText(fmt.Sprintf("%g", result))
		operator = ""
		firstNumber = 0
		newNumber = true
	}

	setOperator := func(op string) {
		if !newNumber {
			num, err := strconv.ParseFloat(display.Text, 64)
			if err == nil {
				firstNumber = num
			}
		}
		operator = op
		newNumber = true
	}

	// Create number buttons
	createNumberButton := func(number string) *widget.Button {
		return widget.NewButton(number, func() {
			appendToDisplay(number)
		})
	}

	// Create operator buttons
	createOperatorButton := func(op string) *widget.Button {
		return widget.NewButton(op, func() {
			setOperator(op)
		})
	}

	// Number buttons
	btn0 := createNumberButton("0")
	btn1 := createNumberButton("1")
	btn2 := createNumberButton("2")
	btn3 := createNumberButton("3")
	btn4 := createNumberButton("4")
	btn5 := createNumberButton("5")
	btn6 := createNumberButton("6")
	btn7 := createNumberButton("7")
	btn8 := createNumberButton("8")
	btn9 := createNumberButton("9")
	
	// Operator buttons
	btnPlus := createOperatorButton("+")
	btnMinus := createOperatorButton("-")
	btnMultiply := createOperatorButton("*")
	btnDivide := createOperatorButton("/")

	// Special buttons
	btnDecimal := widget.NewButton(".", func() {
		if newNumber {
			display.SetText("0.")
			newNumber = false
		} else if !contains(display.Text, ".") {
			display.SetText(display.Text + ".")
		}
	})

	btnEquals := widget.NewButton("=", func() {
		calculate()
	})

	btnClear := widget.NewButton("C", func() {
		display.SetText("0")
		firstNumber = 0
		operator = ""
		newNumber = true
	})

	btnClearEntry := widget.NewButton("CE", func() {
		display.SetText("0")
		newNumber = true
	})

	// Create button grid layout
	buttons := container.NewGridWithColumns(4,
		btnClear, btnClearEntry, btnDivide, btnMultiply,
		btn7, btn8, btn9, btnMinus,
		btn4, btn5, btn6, btnPlus,
		btn1, btn2, btn3, btnEquals,
		btn0, btnDecimal, widget.NewLabel(""), widget.NewLabel(""),
	)

	// Create main content
	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Simple Calculator", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			display,
		),
		nil,
		nil,
		nil,
		buttons,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 400))
	myWindow.ShowAndRun()
}

// Helper function to check if string contains a character
func contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if string(s[i]) == substr {
			return true
		}
	}
	return false
}
