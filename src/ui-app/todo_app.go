package main

import (
	"fmt"
	
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type TodoItem struct {
	Text      string
	Completed bool
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Todo List")

	// Store todo items
	todos := []TodoItem{}

	// Create list widget
	list := widget.NewList(
		func() int {
			return len(todos)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", func(bool) {}),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			hbox := obj.(*fyne.Container)
			check := hbox.Objects[0].(*widget.Check)
			label := hbox.Objects[1].(*widget.Label)
			
			check.SetChecked(todos[id].Completed)
			label.SetText(todos[id].Text)
			
			check.OnChanged = func(checked bool) {
				todos[id].Completed = checked
				list.Refresh()
			}
		},
	)

	// Create input for new todo
	newTodoEntry := widget.NewEntry()
	newTodoEntry.SetPlaceHolder("Enter new todo...")

	// Create counter label
	counterBinding := binding.NewString()
	updateCounter := func() {
		total := len(todos)
		completed := 0
		for _, todo := range todos {
			if todo.Completed {
				completed++
			}
		}
		counterBinding.Set(fmt.Sprintf("Total: %d | Completed: %d | Pending: %d", 
			total, completed, total-completed))
	}
	counterLabel := widget.NewLabelWithData(counterBinding)

	// Add button
	addButton := widget.NewButton("Add Todo", func() {
		if newTodoEntry.Text != "" {
			todos = append(todos, TodoItem{
				Text:      newTodoEntry.Text,
				Completed: false,
			})
			newTodoEntry.SetText("")
			list.Refresh()
			updateCounter()
		}
	})

	// Delete completed button
	deleteCompletedButton := widget.NewButton("Delete Completed", func() {
		newTodos := []TodoItem{}
		for _, todo := range todos {
			if !todo.Completed {
				newTodos = append(newTodos, todo)
			}
		}
		todos = newTodos
		list.Refresh()
		updateCounter()
	})

	// Clear all button
	clearAllButton := widget.NewButton("Clear All", func() {
		todos = []TodoItem{}
		list.Refresh()
		updateCounter()
	})

	// Add sample todos
	addSampleButton := widget.NewButton("Add Samples", func() {
		samples := []string{
			"Learn Go programming",
			"Build a web application",
			"Write unit tests",
			"Deploy to production",
			"Write documentation",
		}
		for _, sample := range samples {
			todos = append(todos, TodoItem{
				Text:      sample,
				Completed: false,
			})
		}
		list.Refresh()
		updateCounter()
	})

	// Create input row
	inputRow := container.NewBorder(nil, nil, nil, addButton, newTodoEntry)

	// Create button row
	buttonRow := container.NewHBox(
		deleteCompletedButton,
		clearAllButton,
		addSampleButton,
	)

	// Create main content
	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("My Todo List", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			inputRow,
			counterLabel,
		),
		buttonRow,
		nil,
		nil,
		list,
	)

	// Initialize counter
	updateCounter()

	// Set content and show
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(500, 600))
	myWindow.ShowAndRun()
}
