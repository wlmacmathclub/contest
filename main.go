package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"
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
	_, err := os.Stat("cache")
	if err != nil && os.IsNotExist(err) {
		os.Mkdir("cache", 0700)
	}
	_, err2 := os.Stat("contest.log")
	if err2 == nil {
		os.Remove("contest.log")
	}

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

	clearCacheButton := widget.NewButton("Clear Cache", func() {
		dialog.NewConfirm("Confirm Cache Clear", "Are you sure you want to clear the cach? All previously generated PDFs will be deleted", func(del bool) {
			if del {
				os.RemoveAll("cache")
				os.Mkdir("cache", 0700)
			}
		}, appWindow).Show()
	})
	form := makeForm(appWindow, &users, contestApp)
	formTitle := canvas.NewText("Settings", color.White)
	formTitle.TextSize = 20
	formTitle.TextStyle = fyne.TextStyle{Bold: true}
	formContainer := container.NewBorder(container.NewBorder(nil, nil, nil, clearCacheButton, container.NewCenter(formTitle)), nil, nil, nil, form)
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
				logw := app.NewWindow("Send Log")
				logw.Resize(fyne.NewSize(720, 480))
				logw.Show()
				logw.SetOnClosed(func() {
					isSubbed = false
				})
				logtext := widget.NewTextGrid()
				logtext.ShowLineNumbers = true
				logBox := container.NewScroll(logtext)
				logw.SetContent(logBox)
				writeLog(logtext, fmt.Sprintf("Starting generator @ %s ...\n\n[Contest Name]: %s\n[Contest Date]: %s\n[Sending Email]: %s\n[Email Name]: %s\n[Public Key]: %s\n[Private Key]: %s\n[Email Subject]: %s\n[Email Body]: %s\n\n\n", time.Now().Format(time.RFC1123), contest.name, contest.date, config.email, config.name, strings.Repeat("*", len(config.publickey)), strings.Repeat("*", len(config.privatekey)), config.subject, config.body), logBox)
				for i, user := range *users {
					if !isSubbed {
						break
					}
					writeLog(logtext, fmt.Sprintf("=============== USER #%d %s %s ===============\nGenerating PDF...\n", i+1, user.firstName, user.lastName), logBox)
					startpdf := time.Now()
					pdfSuccess := generatePDF(user, contest)
					if !pdfSuccess {
						writeLog(logtext, "There was an error generating the PDF. Skipping mailing. \n\n", logBox)
					} else {
						writeLog(logtext, fmt.Sprintf("Generated PDF in %s\nSending mail to %s...\n", time.Since(startpdf), user.email), logBox)
						mailSuccess := mailUser(user, contest, config)
						if mailSuccess {
							writeLog(logtext, "Email sent successfully, sleeping for 5 seconds\n\n", logBox)
						} else {
							writeLog(logtext, "There was an error sending the email\n\n", logBox)
						}
						time.Sleep(5 * time.Second)
					}
				}
				if !isSubbed {
					writeLog(logtext, "Sender terminated due to log window being closed. \n", logBox)
				} else {
					writeLog(logtext, "Finished sending to all emails! \n\n", logBox)
					time.Sleep(5 * time.Second)
					isSubbed = false
					logw.Close()
				}
			}
		},
		SubmitText: "SEND EMAILS",
	}
	return form
}

func writeLog(logtext *widget.TextGrid, text string, logbox *container.Scroll) {
	logtext.SetText(logtext.Text() + text)
	logbox.ScrollToBottom()
	f, _ := os.OpenFile("contest.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte(text))
	f.Close()
}
