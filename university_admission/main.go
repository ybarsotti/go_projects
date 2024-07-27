package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const NumDepartmentsPerApplicant = 3

type Person struct {
	Firstname string
	Lastname  string
}

type Department string
type Subject int

const (
	physics Subject = iota
	chemistry
	math
	computerScience
)

type Applicant struct {
	Person
	ExamScores map[Subject]float64
	Deps       []Department
	GPA        float64
	SpecialScore float64
}

func (p *Person) Fullname() string {
	return p.Firstname + " " + p.Lastname
}

func (a *Applicant) GetScore(dep Department) float64 {
	switch dep {
	case "Physics":
		return (a.ExamScores[physics] + a.ExamScores[math]) / 2
	case "Chemistry":
		return a.ExamScores[chemistry]
	case "Mathematics":
		return a.ExamScores[math]
	case "Engineering":
		return (a.ExamScores[computerScience] + a.ExamScores[math]) / 2
	case "Biotech":
		return (a.ExamScores[chemistry] + a.ExamScores[physics]) / 2
	default:
		panic("Unknown department")
	}
}

func main() {
	applicants := loadInfo()

	var maxApplicants int
	fmt.Scan(&maxApplicants)
	
	departments := getDepartments(applicants)
	applicantsByDepartment := distributeApplicants(departments, applicants, maxApplicants)

	for _, department := range departments {
		file, err := os.Create(strings.ToLower(string(department)) + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		for _, applicant := range applicantsByDepartment[department] {
			score := applicant.GetScore(department)
			if applicant.SpecialScore > score {
				score = applicant.SpecialScore
			}
			fmt.Fprintf(file, "%s %.2f\n", applicant.Fullname(), score)
		}
		file.Close()
	}
}

func loadInfo() (applicants []Applicant) {
	file, err := os.Open("applicants.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	applicants = make([]Applicant, 0, 30)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if applicant, err := getApplicant(strings.Split(scanner.Text(), " ")); err == nil {
			applicants = append(applicants, applicant)
		}
	}

	return
}

func getApplicant(s []string) (applicant Applicant, err error) {
	if len(s) != 10 {
		err = errors.New("Invalid line")
		return
	}

	applicant.Firstname = s[0]
	applicant.Lastname = s[1]
	applicant.ExamScores = make(map[Subject]float64, 4)
	applicant.ExamScores[physics], _ = strconv.ParseFloat(s[2], 64)
	applicant.ExamScores[chemistry], _ = strconv.ParseFloat(s[3], 64)
	applicant.ExamScores[math], _ = strconv.ParseFloat(s[4], 64)
	applicant.ExamScores[computerScience], _ = strconv.ParseFloat(s[5], 64)
	applicant.SpecialScore, _ = strconv.ParseFloat(s[6], 64)
	applicant.Deps = make([]Department, 0, 3)
	applicant.Deps = append(applicant.Deps, Department(s[7]))
	applicant.Deps = append(applicant.Deps, Department(s[8]))
	applicant.Deps = append(applicant.Deps, Department(s[9]))
	

	gpa := applicant.ExamScores[physics]
	gpa += applicant.ExamScores[chemistry]
	gpa += applicant.ExamScores[math]
	gpa += applicant.ExamScores[computerScience]

	applicant.GPA = gpa / 4

	return
}

func distributeApplicants(departments []Department, applicants []Applicant, maxApplicants int) map[Department][]Applicant {
	deps := make(map[Department][]Applicant, len(departments))

	for i := 0; i < NumDepartmentsPerApplicant; i++ {
		for _, dep := range departments {
			deps[dep] = append(deps[dep], selectApplicantsByDepartment(dep, &applicants, i, maxApplicants-len(deps[dep]))...)
		}
	}

	for dep, selectedApplicants := range deps {
		sortApplicantsByDepartment(dep, &selectedApplicants)
		deps[dep] = selectedApplicants
	}

	return deps
}

func selectApplicantsByDepartment(department Department, applicants *[]Applicant, priority, maxApplicants int) []Applicant {
	if maxApplicants <= 0 {
		return nil
	}

	selectedApplicants := make([]Applicant, 0, maxApplicants)

	sortApplicantsByDepartment(department, applicants)

	for k, j := 0, len(*applicants); j > 0; j-- {
		applicant := (*applicants)[k]
		dep := applicant.Deps[priority]
		if dep != department {
			k++
			continue
		}

		selectedApplicants = append(selectedApplicants, applicant)
		*applicants = append((*applicants)[:k], (*applicants)[k+1:]...)

		if len(selectedApplicants) == maxApplicants {
			break
		}
	}

	return selectedApplicants
}

func sortApplicantsByDepartment(department Department, applicants *[]Applicant) {
	sort.Slice(*applicants, func(i, j int) bool {
		a, b := (*applicants)[i], (*applicants)[j]
		score1, score2 := a.GetScore(department), b.GetScore(department)
		bestScore1, bestScore2 := a.SpecialScore, b.SpecialScore

		if bestScore1 > score1 {
			score1 = bestScore1
		}
		if bestScore2 > score2 {
			score2 = bestScore2
		}
		if score1 != score2 {
			return score1 > score2
		}
		return a.Fullname() < b.Fullname()
	})
}

func getDepartments(applicants []Applicant) []Department {
	departments := make([]Department, 0, 4)

	for _, applicant := range applicants {
		for _, dep := range applicant.Deps {
			if !slices.Contains(departments, dep) {
				departments = append(departments, dep)
			}
		}
	}

	sort.Slice(departments, func(i, j int) bool {
		return departments[i] < departments[j]
	})

	return departments
}
