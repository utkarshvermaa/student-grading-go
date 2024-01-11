package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cast"
)

const (
	filename string = "grades.csv"
)

type Header string

const (
	FirstName  Header = "FirstName"
	LastName   Header = "LastName"
	University Header = "University"
	Test1      Header = "Test1"
	Test2      Header = "Test2"
	Test3      Header = "Test3"
	Test4      Header = "Test4"
)

func (h Header) ToString() string {
	return string(h)
}

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

var (
	firstNameIndex,
	lastNameIndex,
	uniIndex,
	test1ScoreIndex,
	test2ScoreIndex,
	test3ScoreIndex,
	test4ScoreIndex int
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

func (s student) finalScore() float32 {
	final := s.test1Score + s.test2Score + s.test3Score + s.test4Score
	return float32(final) / 4
}

func (s student) grade() Grade {
	finalScore := s.finalScore()
	switch {
	case finalScore < 35:
		return F
	case finalScore >= 35 && finalScore < 50:
		return C
	case finalScore >= 50 && finalScore < 70:
		return B
	case finalScore >= 70:
		return A
	default:
		return F
	}
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {
	resp := make([]student, 0)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error: %v while trying to open file: %v", err, filename)
		return resp
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("error: %v while trying to read records", err)
		return resp
	}

	indMap := map[Header]int{
		FirstName:  -1,
		LastName:   -1,
		University: -1,
		Test1:      -1,
		Test2:      -1,
		Test3:      -1,
		Test4:      -1,
	}
	for rowNum, row := range records {
		if rowNum == 0 {
			for colNum, cell := range row {
				if strings.Compare(cell, FirstName.ToString()) == 0 {
					firstNameIndex = colNum
					delete(indMap, FirstName)
				}
				if strings.Compare(cell, LastName.ToString()) == 0 {
					lastNameIndex = colNum
					delete(indMap, LastName)
				}
				if strings.Compare(cell, University.ToString()) == 0 {
					uniIndex = colNum
					delete(indMap, University)
				}
				if strings.Compare(cell, Test1.ToString()) == 0 {
					test1ScoreIndex = colNum
					delete(indMap, Test1)
				}
				if strings.Compare(cell, Test2.ToString()) == 0 {
					test2ScoreIndex = colNum
					delete(indMap, Test2)
				}
				if strings.Compare(cell, Test3.ToString()) == 0 {
					test3ScoreIndex = colNum
					delete(indMap, Test3)
				}
				if strings.Compare(cell, Test4.ToString()) == 0 {
					test4ScoreIndex = colNum
					delete(indMap, Test4)
				}
			}
			continue
		}
		if len(indMap) > 0 {
			missingHeaders := ""
			for h := range indMap {
				missingHeaders += h.ToString() + ", "
			}
			fmt.Printf("missing headers: %v", missingHeaders)
			return resp
		}
		s := student{
			firstName:  row[firstNameIndex],
			lastName:   row[lastNameIndex],
			university: row[uniIndex],
			test1Score: cast.ToInt(row[test1ScoreIndex]),
			test2Score: cast.ToInt(row[test2ScoreIndex]),
			test3Score: cast.ToInt(row[test3ScoreIndex]),
			test4Score: cast.ToInt(row[test4ScoreIndex]),
		}
		resp = append(resp, s)
	}
	return resp
}

func calculateGrade(students []student) []studentStat {
	resp := make([]studentStat, 0)
	for _, s := range students {
		stat := studentStat{
			student:    s,
			finalScore: s.finalScore(),
			grade:      s.grade(),
		}
		resp = append(resp, stat)
	}
	return resp
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	resp := studentStat{}
	for _, gs := range gradedStudents {
		if resp.finalScore < gs.finalScore {
			resp = gs
		}
	}
	return resp
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	resp := make(map[string]studentStat, 0)
	for _, s := range gs {
		fs, found := resp[s.university]
		if !found {
			resp[s.university] = s
			continue
		}
		if fs.finalScore < s.finalScore {
			resp[s.university] = s
		}
	}
	return resp
}
