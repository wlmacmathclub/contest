package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"
)

func mailUser(user User, contest Contest, config MailConfig) bool {
	msgbuf := bytes.NewBuffer(nil)
	message := "From: " + config.from +
		"\nTo: " + user.email +
		"\nSuject: " + textReplace(user, contest, config.subject) +
		"\nMIME-Version: 1.0\n"
	msgbuf.WriteString(message)
	writer := multipart.NewWriter(msgbuf)
	msgbuf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", writer.Boundary()))
	msgbuf.WriteString(fmt.Sprintf("--%s\n", writer.Boundary()))
	msgbuf.WriteString("\n" + textReplace(user, contest, config.body) + "\n")
	filename := fmt.Sprintf("%s_%s%s.pdf", contest.name, user.firstName, user.lastName)
	filebuf, _ := ioutil.ReadFile("cache/" + filename)

	msgbuf.WriteString(fmt.Sprintf("\n\n--%s\n", writer.Boundary()))
	msgbuf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(filebuf)))
	msgbuf.WriteString("Content-Transfer-Encoding: base64\n")
	msgbuf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", filename))
	b := make([]byte, base64.StdEncoding.EncodedLen(len(filebuf)))
	base64.StdEncoding.Encode(b, filebuf)
	msgbuf.Write(b)
	msgbuf.WriteString("--")

	err := smtp.SendMail(config.server+config.port,
		smtp.PlainAuth("", config.from, config.password, config.server),
		config.from, []string{user.email}, []byte(msgbuf.Bytes()))

	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func textReplace(user User, contest Contest, text string) string {
	newtext := text
	list := [][]string{
		[]string{"{USER_FIRST_NAME}", user.firstName},
		[]string{"{USER_LAST_NAME}", user.lastName},
		[]string{"{CONTEST_NAME}", contest.name},
		[]string{"{CONTEST_DATE}", contest.date},
		[]string{"{USER_EMAIL}", user.email},
		[]string{"{USER_P1_TEACHER}", user.firstTeacher},
		[]string{"{USER_P2_TEACHER}", user.secondTeacher},
	}
	for _, item := range list {
		newtext = strings.Replace(newtext, item[0], item[1], -1)
	}
	return newtext
}
