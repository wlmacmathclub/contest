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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	nativedialog "github.com/sqweek/dialog"
)

func main() {
	icon, _ := fyne.LoadResourceFromPath("assets/images/private.png")
	contestApp := app.New()
	appWindow := contestApp.NewWindow("Math Contest Permission Form Generator " + version)
	appWindow.SetIcon(icon)
	appWindow.Resize(fyne.NewSize(1280, 720))

	darkMode := contestApp.Settings().ThemeVariant() == theme.VariantDark

	var users []User

	usertable := widget.NewTable(
		func() (int, int) {
			entries := 4
			if len(users) >= 1 {
				if users[0].secondTeacher != "" {
					entries = 5
				}
				if users[0].thirdTeacher != "" {
					entries = 6
				}
				if users[0].fourthTeacher != "" {
					entries = 7
				}
			}
			return len(users), entries
		},
		func() fyne.CanvasObject {
			//this determines starting size (todo: change this)
			return widget.NewLabel("template content here")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(users[i.Row].email)
			case 1:
				o.(*widget.Label).SetText(users[i.Row].firstName)
			case 2:
				o.(*widget.Label).SetText(users[i.Row].lastName)
			case 3:
				o.(*widget.Label).SetText(users[i.Row].firstTeacher)
			case 4:
				o.(*widget.Label).SetText(users[i.Row].secondTeacher)
			case 5:
				o.(*widget.Label).SetText(users[i.Row].thirdTeacher)
			case 6:
				o.(*widget.Label).SetText(users[i.Row].fourthTeacher)
			}
		})
	usertable.Hide()
	usertable.SetColumnWidth(0, 500)

	userlisttext := widget.NewLabel("")
	userlisttext.SetText("No file selected. Click the \"Open File\" button")

	openFileButton := widget.NewButton("Open File", func() {
		filename, err := nativedialog.File().Filter("Comma Separated Values (.csv)", "csv").Title("Open Contest Participants List").Load()
		if err == nil && filename != "" {
			userlist, err2 := parseCSV(filename)
			if err2 == "" {
				users = userlist
				userlisttext.SetText(func() string {
					if len(users) == 0 {
						return "No user records found. Is the file formatted correctly?"
					} else {
						//userlisttext.Hide()
						usertable.Show()
						return fmt.Sprintf("%d user records loaded\nTable format: email, first name, last name, p1 teacher, p2 teacher, p3 teacher, p4 teacher", len(users))
					}
				}())
			} else {
				userlisttext.SetText(err2)
			}
		}
	})

	userBox := container.NewScroll(container.NewBorder(userlisttext, nil, nil, nil, usertable))
	userBoxTitle := canvas.NewText("Contest Participant List", func() color.Color {
		if darkMode {
			return color.White
		} else {
			return color.Black
		}
	}())
	userBoxTitle.TextSize = 20
	userBoxTitle.TextStyle = fyne.TextStyle{Bold: true}
	userBoxContainer := container.NewBorder(container.NewBorder(nil, nil, nil, openFileButton, container.NewCenter(userBoxTitle)), nil, nil, nil, userBox)

	clearCacheButton := widget.NewButton("Clear Cache", func() {
		dialog.NewConfirm("Confirm Cache Clear", "Are you sure you want to clear the cache? All previously generated PDFs will be deleted", func(del bool) {
			if del {
				os.RemoveAll("cache")
				os.Mkdir("cache", 0700)
			}
		}, appWindow).Show()
	})
	form := makeForm(appWindow, &users, contestApp)
	formTitle := canvas.NewText("Settings", func() color.Color {
		if darkMode {
			return color.White
		} else {
			return color.Black
		}
	}())
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
	mailRegExp, _ := regexp.Compile(`^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	emailName := widget.NewEntry()
	emailPubKey := widget.NewEntry()
	emailPrivKey := widget.NewEntry()
	emailSubject := widget.NewEntry()
	emailBody := widget.NewMultiLineEntry()

	iniContest, iniConfig, iniErr := initialize()
	if iniErr == nil {
		contestName.SetText(iniContest.Name)
		contestDate.SetText(iniContest.Date)
		email.SetText(iniConfig.Email)
		emailName.SetText(iniConfig.Name)
		emailPubKey.SetText(iniConfig.Publickey)
		emailPrivKey.SetText(iniConfig.Privatekey)
		emailSubject.SetText(iniConfig.Subject)
		emailBody.SetText(iniConfig.Body)
	}

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
					Name: contestName.Text,
					Date: contestDate.Text,
				}
				config := MailConfig{
					Email:      email.Text,
					Name:       emailName.Text,
					Publickey:  emailPubKey.Text,
					Privatekey: emailPrivKey.Text,
					Subject:    emailSubject.Text,
					Body:       emailBody.Text,
				}
				saveJson(contest, config)
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
				writeLog(logtext, fmt.Sprintf("Starting generator @ %s ...\n\n[Contest Name]: %s\n[Contest Date]: %s\n[Sending Email]: %s\n[Email Name]: %s\n[Public Key]: %s\n[Private Key]: %s\n[Email Subject]: %s\n[Email Body]: %s\n\n\n", time.Now().Format(time.RFC1123), contest.Name, contest.Date, config.Email, config.Name, strings.Repeat("*", len(config.Publickey)), strings.Repeat("*", len(config.Privatekey)), config.Subject, config.Body), logBox)
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
