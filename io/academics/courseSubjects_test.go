package academics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	domain "obas/domain/academics"
	"testing"
)

func TestCreateCourseSubject(t *testing.T) {
	obj := domain.CourseSubject{"0001", "A0000"}
	result, err := CreateCourseSubject(obj)
	assert.Nil(t, err)
	fmt.Println(" The Results", result)
	assert.Equal(t, obj, result)
}
func TestGetCourseSubject(t *testing.T) {
	result, err := GetCourseSubject("0001", "002")
	assert.Nil(t, err)
	fmt.Println(" The Results", result)
	assert.NotNil(t, result)
}
func TestDeleteCourseSubject(t *testing.T) {
	obj := domain.CourseSubject{"0001", "A0000"}
	result, err := DeleteCourseSubject(obj)
	assert.Nil(t, err)
	fmt.Println(" The Results", result)
	assert.NotNil(t, result)
}
func TestGetCourseSubjects(t *testing.T) {
	result, err := GetCourseSubjects()
	assert.Nil(t, err)
	fmt.Println(" The Results", result)
	assert.NotNil(t, result)
}
func TestUpdateCourseSubject(t *testing.T) {
	//obj:=domain.CourseSubject{"0001","A0000"}
	result, err := UpdateCourseSubject("00003")
	assert.Nil(t, err)
	fmt.Println(" The Results", result)
	assert.NotNil(t, result)
}