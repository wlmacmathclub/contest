package main

import (
	"fmt"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func generatePDF(user User, contest Contest) bool {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	m.SetFontLocation("assets/font/")
	m.AddUTF8Font("TeXFont", consts.Normal, "lmroman10-regular.ttf")
	m.AddUTF8Font("TeXFont", consts.Bold, "lmroman10-bold.ttf")
	m.SetDefaultFontFamily("TeXFont")

	m.SetAliasNbPages("{nb}")
	m.SetFirstPageNb(1)

	m.RegisterHeader(func() {
		m.Row(40, func() {
			m.Col(0, func() {
				m.Text("Math Contest Permission Form", props.Text{
					Size:  20,
					Style: consts.Bold,
					Align: consts.Center,
					Top:   4,
				})
				m.Text("William Lyon Mackenzie C.I.", props.Text{
					Size:  17,
					Align: consts.Center,
					Top:   20,
				})
			})
		})

		m.Line(2, props.Line{
			Style: consts.Solid,
			Width: 0.5,
		})
	})

	m.RegisterFooter(func() {
		m.Row(20, func() {
			m.Col(0, func() {
				m.Text("Registration system created by Patrick Lin", props.Text{
					Size:  12,
					Align: consts.Center,
					Top:   4,
				})
			})
		})
	})

	m.Row(40, func() {
		m.Col(6, func() {
			m.Text("Student Name: ", props.Text{
				Top:   20,
				Size:  15,
				Align: consts.Right,
				Style: consts.Bold,
			})
		})
		m.Col(6, func() {
			m.Text(fmt.Sprintf("%s %s", user.firstName, user.lastName), props.Text{
				Top:   20,
				Size:  15,
				Align: consts.Left,
			})
		})
	})
	m.Row(20, func() {
		m.Col(6, func() {
			m.Text("Contest Name: ", props.Text{
				Size:  15,
				Align: consts.Right,
				Style: consts.Bold,
			})
		})
		m.Col(6, func() {
			m.Text(contest.Name, props.Text{
				Size:  15,
				Align: consts.Left,
			})
		})
	})
	m.Row(20, func() {
		m.Col(6, func() {
			m.Text("Contest Date: ", props.Text{
				Size:  15,
				Align: consts.Right,
				Style: consts.Bold,
			})
		})
		m.Col(6, func() {
			m.Text(contest.Date, props.Text{
				Size:  15,
				Align: consts.Left,
			})
		})
	})
	m.Row(20, func() {
		m.Col(6, func() {
			m.Text("Period 1 Teacher: ", props.Text{
				Size:  15,
				Align: consts.Right,
				Style: consts.Bold,
			})
		})
		m.Col(6, func() {
			m.Text(user.firstTeacher, props.Text{
				Size:  15,
				Align: consts.Left,
			})
		})
	})
	if user.secondTeacher != "" {
		m.Row(20, func() {
			m.Col(6, func() {
				m.Text("Period 2 Teacher: ", props.Text{
					Size:  15,
					Align: consts.Right,
					Style: consts.Bold,
				})
			})
			m.Col(6, func() {
				m.Text(user.secondTeacher, props.Text{
					Size:  15,
					Align: consts.Left,
				})
			})
		})
	}

	m.Row(20, func() {
		if user.secondTeacher == "" {
			m.Col(3, func() {})
			m.Col(6, func() {
				m.Signature("Period 1 Teacher Signature", props.Font{
					Style: consts.Normal,
					Size:  15,
				})
			})
			m.Col(3, func() {})
		} else {
			m.Col(6, func() {
				m.Signature("Period 1 Teacher Signature", props.Font{
					Style: consts.Normal,
					Size:  15,
				})
			})
			m.Col(6, func() {
				m.Signature("Period 2 Teacher Signature", props.Font{
					Style: consts.Normal,
					Size:  15,
				})
			})
		}
	})

	m.Row(40, func() {
		m.Col(0, func() {
			m.Text("If you do not show up to the contest you may penalized!", props.Text{
				Top:   30,
				Size:  13,
				Align: consts.Center,
			})
		})
	})

	err := m.OutputFileAndClose(fmt.Sprintf("cache/%s_%s%s.pdf", contest.Name, user.firstName, user.lastName))
	return err == nil
}
