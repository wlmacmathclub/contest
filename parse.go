package main

import (
	"encoding/csv"
	"os"
)

func parseCSV(path string) ([]User, string) {
	file, _ := os.Open(path)
	defer file.Close()

	csvRdr := csv.NewReader(file)
	recs, err := csvRdr.ReadAll()

	if err != nil {
		return make([]User, 0), "Error reading CSV file"
	}

	users := make([]User, 0)
	for _, person := range recs {
		if len(person) >= 4 {
			users = append(users, User{
				email:        person[0],
				firstName:    person[1],
				lastName:     person[2],
				firstTeacher: person[3],
				secondTeacher: func() string {
					if len(person) > 4 {
						return person[4]
					} else {
						return ""
					}
				}(),
			})
		}
	}

	return users, ""
}
