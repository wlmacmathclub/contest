package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mailjet/mailjet-apiv3-go"
)

func mailUser(user User, contest Contest, config MailConfig) bool {
	filename := fmt.Sprintf("%s_%s%s.pdf", contest.Name, user.firstName, user.lastName)
	filebuf, _ := ioutil.ReadFile("cache/" + filename)
	content := base64.StdEncoding.EncodeToString(filebuf)
	mailjetClient := mailjet.NewMailjetClient(config.Publickey, config.Privatekey)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: config.Email,
				Name:  config.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: user.email,
					Name:  fmt.Sprintf("%s %s", user.firstName, user.lastName),
				},
			},
			Subject:  textReplace(user, contest, config.Subject),
			HTMLPart: textReplace(user, contest, config.Body),
			Attachments: &mailjet.AttachmentsV31{
				mailjet.AttachmentV31{
					ContentType:   "application/pdf",
					Filename:      filename,
					Base64Content: content,
				},
			},
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)

	return err == nil
}

func textReplace(user User, contest Contest, text string) string {
	newtext := text
	list := [][]string{
		{"{USER_FIRST_NAME}", user.firstName},
		{"{USER_LAST_NAME}", user.lastName},
		{"{CONTEST_NAME}", contest.Name},
		{"{CONTEST_DATE}", contest.Date},
		{"{USER_EMAIL}", user.email},
		{"{USER_P1_TEACHER}", user.firstTeacher},
		{"{USER_P2_TEACHER}", user.secondTeacher},
		{"{USER_P3_TEACHER}", user.thirdTeacher},
		{"{USER_P4_TEACHER}", user.fourthTeacher},
	}
	for _, item := range list {
		newtext = strings.Replace(newtext, item[0], item[1], -1)
	}
	return newtext
}
