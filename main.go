package main

import (
	"errors"
	"image/color"
	"regexp"
	"strings"
	"time"

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

	form := makeForm(appWindow, &users, contestApp)
	formTitle := canvas.NewText("Settings", color.White)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formContainer := container.NewBorder(container.NewCenter(formTitle), nil, nil, nil, form)
	grid := container.NewHSplit(userBoxContainer, formContainer)
	appWindow.SetContent(grid)

	appWindow.ShowAndRun()
}

func makeForm(appWindow fyne.Window, users *[]User, app fyne.App) *widget.Form {
	contestName := widget.NewEntry()
	contestDate := widget.NewEntry()
	email := widget.NewEntry()
	mailRegExp, _ := regexp.Compile("^\\w+([-+.']\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$")
	emailName := widget.NewEntry()
	emailPubKey := widget.NewEntry()
	emailPrivKey := widget.NewEntry()
	emailSubject := widget.NewEntry()
	emailBody := widget.NewMultiLineEntry()
	isSubbed := false
	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "Contest Name",
				Widget: contestName,
			},
			{
				Text:   "Contest Date",
				Widget: contestDate,
			},
			{
				Text:   "Email From",
				Widget: email,
			},
			{
				Text:   "Email Name",
				Widget: emailName,
			},
			{
				Text:   "Email Public Key",
				Widget: emailPubKey,
			},
			{
				Text:   "Email Private Key",
				Widget: emailPrivKey,
			},
			{
				Text:   "Email Subject",
				Widget: emailSubject,
			},
			{
				Text:   "Email Body",
				Widget: emailBody,
			},
		},
		OnSubmit: func() {
			if isSubbed {
				dialog.NewError(errors.New("sending in progress"), appWindow).Show()
			} else if contestName.Text == "" || contestDate.Text == "" || email.Text == "" || emailName.Text == "" || emailPubKey.Text == "" || emailPrivKey.Text == "" || emailSubject.Text == "" || emailBody.Text == "" {
				dialog.NewError(errors.New("cannot have empty field"), appWindow).Show()
			} else if len(*users) == 0 {
				dialog.NewError(errors.New("no user to send to"), appWindow).Show()
			} else if !mailRegExp.MatchString(email.Text) {
				dialog.NewError(errors.New("invalid sender email"), appWindow).Show()
			} else {
				isSubbed = true
				contest := Contest{
					name: contestName.Text,
					date: contestDate.Text,
				}
				config := MailConfig{
					email:      email.Text,
					name:       emailName.Text,
					publickey:  emailPubKey.Text,
					privatekey: emailPrivKey.Text,
					subject:    emailSubject.Text,
					body:       emailBody.Text,
				}
				for _, user := range *users {
					generatePDF(user, contest)
					mailUser(user, contest, config)
					time.Sleep(5 * time.Second)
				}
			}
		},
		SubmitText: "SEND EMAILS",
	}
	return form
}
