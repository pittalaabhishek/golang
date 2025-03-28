package school

import (
	"sort"
)

// Grade represents a school grade with its students
type Grade struct {
	Level    int
	Students []string
}

// School represents the entire school roster
type School struct {
	grades map[int]*Grade
}

// New creates a new School instance
func New() *School {
	return &School{
		grades: make(map[int]*Grade),
	}
}

// Add a student to a grade
func (s *School) Add(student string, grade int) {
	if _, exists := s.grades[grade]; !exists {
		s.grades[grade] = &Grade{Level: grade}
	}

	for _, name := range s.grades[grade].Students {
		if name == student {
			return // Student already exists in this grade
		}
	}

	s.grades[grade].Students = append(s.grades[grade].Students, student)
	sort.Strings(s.grades[grade].Students)
}

// Grade returns students in a specific grade
func (s *School) Grade(level int) []string {
	if grade, exists := s.grades[level]; exists {
		return grade.Students
	}
	return []string{}
}

// Enrollment returns all grades with their students, sorted
func (s *School) Enrollment() []Grade {
	var result []Grade

	// Get all grade levels and sort them
	var levels []int
	for level := range s.grades {
		levels = append(levels, level)
	}
	sort.Ints(levels)

	// Build the result in order
	for _, level := range levels {
		result = append(result, *s.grades[level])
	}

	return result
}