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
		m.Row(12, func() {
			m.Col(0, func() {
				m.Text("Registration system created by Patrick Lin", props.Text{
					Size:  12,
					Align: consts.Center,
					Top:   4,
				})
			})
		})
	})

	m.Row(28, func() {
		m.Col(6, func() {
			m.Text("Student Name: ", props.Text{
				Top:   10,
				Size:  15,
				Align: consts.Right,
				Style: consts.Bold,
			})
		})
		m.Col(6, func() {
			m.Text(fmt.Sprintf("%s %s", user.firstName, user.lastName), props.Text{
				Top:   10,
				Size:  15,
				Align: consts.Left,
			})
		})
	})
	m.Row(18, func() {
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
	m.Row(18, func() {
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

	teachers := []Teacher{}

	if user.firstTeacher != "" {
		teachers = append(teachers, Teacher{
			period: 1,
			name:   user.firstTeacher,
		})
	}
	if user.secondTeacher != "" {
		teachers = append(teachers, Teacher{
			period: 2,
			name:   user.secondTeacher,
		})
	}
	if user.thirdTeacher != "" {
		teachers = append(teachers, Teacher{
			period: 3,
			name:   user.thirdTeacher,
		})
	}
	if user.fourthTeacher != "" {
		teachers = append(teachers, Teacher{
			period: 4,
			name:   user.fourthTeacher,
		})
	}

	for _, teacher := range teachers {
		m.Row(18, func() {
			m.Col(6, func() {
				m.Text(fmt.Sprintf("Period %d Teacher: ", teacher.period), props.Text{
					Size:  15,
					Align: consts.Right,
					Style: consts.Bold,
				})
			})
			m.Col(6, func() {
				m.Text(teacher.name, props.Text{
					Size:  15,
					Align: consts.Left,
				})
			})
		})
	}

	m.Row(20, func() {
		if len(teachers) == 1 {
			m.Col(3, func() {})
			m.Col(6, func() {
				m.Signature(fmt.Sprintf("Period %d Teacher Signature", teachers[0].period), props.Font{
					Style: consts.Normal,
					Size:  15,
				})
			})
			m.Col(3, func() {})
		} else {
			m.Col(6, func() {
				m.Signature(fmt.Sprintf("Period %d Teacher Signature", teachers[0].period), props.Font{
					Style: consts.Normal,
					Size:  15,
				})
			})
			m.Col(6, func() {
				m.Signature(fmt.Sprintf("Period %d Teacher Signature", teachers[1].period), props.Font{
					Style: consts.Normal,
					Size:  15,
				})
			})
		}
	})
	if len(teachers) > 2 {
		m.Row(10, func() {})
		m.Row(20, func() {
			if len(teachers) == 3 {
				m.Col(3, func() {})
				m.Col(6, func() {
					m.Signature(fmt.Sprintf("Period %d Teacher Signature", teachers[2].period), props.Font{
						Style: consts.Normal,
						Size:  15,
					})
				})
				m.Col(3, func() {})
			} else {
				m.Col(6, func() {
					m.Signature(fmt.Sprintf("Period %d Teacher Signature", teachers[2].period), props.Font{
						Style: consts.Normal,
						Size:  15,
					})
				})
				m.Col(6, func() {
					m.Signature(fmt.Sprintf("Period %d Teacher Signature", teachers[3].period), props.Font{
						Style: consts.Normal,
						Size:  15,
					})
				})
			}
		})
	}

	m.Row(20, func() {
		m.Col(0, func() {
			m.Text("If you do not show up to the contest you may be penalized!", props.Text{
				Top:   20,
				Size:  13,
				Align: consts.Center,
			})
		})
	})

	err := m.OutputFileAndClose(fmt.Sprintf("cache/%s_%s%s.pdf", contest.Name, user.firstName, user.lastName))
	return err == nil
}

type Teacher struct { //temporary struct used for sorting/formatting
	period int
	name   string
}
