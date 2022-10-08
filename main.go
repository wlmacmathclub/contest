package main

import (
	"errors"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	icon, _ := fyne.LoadResourceFromPath("assets/images/private.png")
	contestApp := app.New()
	appWindow := contestApp.NewWindow("Math Contest Permission Form Generator")
	appWindow.SetIcon(icon)
	appWindow.Resize(fyne.NewSize(1280, 720))

	var users []User

	userlisttext := widget.NewLabel("")
	userlisttext.SetText("No file selected. Click the \"Open File\" button")

	diag := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err == nil && uc != nil {
			userlist, err2 := parseCSV(uc.URI().Path())
			if err2 == "" {
				users = userlist
				userlisttext.SetText(func() string {
					if len(users) == 0 {
						return "No user records found. Is the file formatted correctly?"
					} else {
						var formatstr strings.Builder
						for _, usr := range users {
							formatstr.WriteString(usr.email + "\t" + usr.firstName + "\t" + usr.lastName + "\t" + usr.firstTeacher + "\t" + usr.secondTeacher + "\n")
						}
						return formatstr.String()
					}
				}())
			} else {
				userlisttext.SetText(err2)
			}
		} else {
			userlisttext.SetText("Error reading selected file")
		}
	}, appWindow)

	openFileButton := widget.NewButton("Open File", func() {
		diag.Show()
	})

	userBox := container.NewScroll(userlisttext)
	userBoxTitle := canvas.NewText("Contest Participant List", color.White)
	userBoxTitle.TextSize = 20
	userBoxTitle.TextStyle = fyne.TextStyle{Bold: true}
	userBoxContainer := container.NewBorder(container.NewBorder(nil, nil, nil, openFileButton, container.NewCenter(userBoxTitle)), nil, nil, nil, userBox)

	contestName := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "Contest Name",
				Widget: contestName,
			},
		},
		OnSubmit: func() {
			fmt.Println(contestName.Text)
			dialog.NewError(errors.New(contestName.Text), appWindow).Show()
		},
	}
	formTitle := canvas.NewText("Settings", color.White)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formContainer := container.NewBorder(container.NewCenter(formTitle), nil, nil, nil, form)
	grid := container.NewHSplit(userBoxContainer, formContainer)
	appWindow.SetContent(grid)

	appWindow.ShowAndRun()
}
