package institutions

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	domain "obas/domain/institutions"
	"testing"
)

func TestCreateInstitution(t *testing.T) {
	obj := domain.Institution{"U111", "CPUT", "UNIVERSITY"}
	resp, err := CreateInstitution(obj)
	assert.True(t, resp)
	fmt.Println(" The Results", resp)
	assert.Nil(t, err)
}
func TestGetInstitutions(t *testing.T) {
	resp, err := GetInstitutions()
	assert.Nil(t, err)
	fmt.Println(" The Results", resp)
	assert.NotNil(t, len(resp) > 0)
}
func TestDeleteInstitution(t *testing.T) {
	obj := domain.Institution{"U111", "CPUT", "UNIVERSITY"}
	resp, err := DeleteInstitution(obj)
	assert.Nil(t, err)
	fmt.Println(" The Results", resp)
	assert.NotNil(t, resp)
}
func TestGetInstitution(t *testing.T) {
	resp, err := GetInstitution("U111")
	assert.Nil(t, err)
	fmt.Println(" The Results", resp)
	assert.NotNil(t, resp)
}