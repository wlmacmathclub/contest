package main

import (
	"image/color"
	"os"
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
	userlisttext.SetText("Waiting for CSV file selection...")
	userBox := container.NewScroll(userlisttext)
	userBoxTitle := canvas.NewText("Contest Participant List", color.White)
	userBoxTitle.TextSize = 20
	userBoxTitle.TextStyle = fyne.TextStyle{Bold: true}
	userBoxContainer := container.NewBorder(container.NewCenter(userBoxTitle), nil, nil, nil, userBox)
	text2 := canvas.NewText("2", color.White)
	grid := container.NewHSplit(userBoxContainer, text2)
	appWindow.SetContent(grid)

	diag := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err == nil && uc != nil {
			//fmt.Println(uc.URI().Path())
			userlist, err2 := parseCSV(uc.URI().Path())
			if err2 == "" {
				//fmt.Print(userlist)
				users = userlist
				/*userTable := widget.NewTable(
					func() (int, int) {
						return len(userlist), 5
					},
					func() fyne.CanvasObject {
						return widget.NewLabel("N/A")
					},
					func(i widget.TableCellID, o fyne.CanvasObject) {
						o.(*widget.Label).SetText(func() string {
							//man fyne sucks balls lol
							if i.Col == 0 {
								return userlist[i.Row].email
							} else if i.Col == 1 {
								return userlist[i.Row].firstName
							} else if i.Col == 2 {
								return userlist[i.Row].lastName
							} else if i.Col == 3 {
								return userlist[i.Row].firstTeacher
							} else {
								if userlist[i.Row].secondTeacher == "" {
									return "N/A"
								} else {
									return userlist[i.Row].secondTeacher
								}
							}
						}())
					})
				userBox = container.NewScroll(userTable)
				userBox.Refresh()*/
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
			userlisttext.SetText("File dialog closed. Terminating process...")
			time.Sleep(3 * time.Second)
			os.Exit(0)
		}
	}, appWindow)
	diag.Show()

	appWindow.ShowAndRun()
}
