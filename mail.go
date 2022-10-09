package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mailjet/mailjet-apiv3-go"
)

func mailUser(user User, contest Contest, config MailConfig) bool {
	filename := fmt.Sprintf("%s_%s%s.pdf", contest.name, user.firstName, user.lastName)
	filebuf, _ := ioutil.ReadFile("cache/" + filename)
	content := base64.StdEncoding.EncodeToString(filebuf)
	mailjetClient := mailjet.NewMailjetClient(config.publickey, config.privatekey)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: config.email,
				Name:  config.name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: user.email,
					Name:  fmt.Sprintf("%s %s", user.firstName, user.lastName),
				},
			},
			Subject:  textReplace(user, contest, config.subject),
			HTMLPart: textReplace(user, contest, config.body),
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
		{"{CONTEST_NAME}", contest.name},
		{"{CONTEST_DATE}", contest.date},
		{"{USER_EMAIL}", user.email},
		{"{USER_P1_TEACHER}", user.firstTeacher},
		{"{USER_P2_TEACHER}", user.secondTeacher},
	}
	for _, item := range list {
		newtext = strings.Replace(newtext, item[0], item[1], -1)
	}
	return newtext
}
