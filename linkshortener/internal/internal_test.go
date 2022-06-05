package internal

import "testing"

func TestCourseProject(t *testing.T) {
	var s string
	s = CourseProject()
	if s != "I started working on a course project\n" {
		t.Error("Expected 'I started working on a course project', got ", s)
	}
}
