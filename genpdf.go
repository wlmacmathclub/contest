package main

import (
	"fmt"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func generatePDF(user User, contest Contest) {
	begin := time.Now()
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.SetAliasNbPages("{nb}")
	m.SetFirstPageNb(1)

	m.RegisterHeader(func() {
		m.Row(22, func() {
			m.Col(0, func() {
				m.Text("Math Contest Permission Form", props.Text{
					Size:  15,
					Style: consts.Bold,
					Align: consts.Center,
					Top:   4,
				})
				m.Text("William Lyon Mackenzie C.I.", props.Text{
					Size:  14,
					Align: consts.Center,
					Top:   12,
				})
			})
		})

		m.Line(1.0, props.Line{
			Style: consts.Solid,
			Width: 0.5,
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Student Name: %s %s", user.firstName, user.lastName), props.Text{
				Top:   10,
				Size:  12,
				Align: consts.Center,
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Contest Name: %s", contest.name), props.Text{
				Size:  12,
				Align: consts.Center,
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Contest Date: %s", contest.date.Format(time.RFC822)), props.Text{
				Size:  12,
				Align: consts.Center,
			})
		})
	})

	m.Row(20, func() {
		if user.secondTeacher == "" {
			m.Col(3, func() {})
			m.Col(6, func() {
				m.Signature(fmt.Sprintf("Signature of %s (Period 1)", user.firstTeacher))
			})
			m.Col(3, func() {})
		} else {
			m.Col(6, func() {
				m.Signature(fmt.Sprintf("Signature of %s (Period 1)", user.firstTeacher))
			})
			m.Col(6, func() {
				m.Signature(fmt.Sprintf("Signature of %s (Period 2)", user.secondTeacher))
			})
		}
	})

	m.OutputFileAndClose("test.pdf")

	end := time.Now()
	fmt.Println(end.Sub(begin))
}
