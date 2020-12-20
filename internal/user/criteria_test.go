package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MatchCriteria_Age(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	testOffer.MinAge = 55
	a.False(testUser.MatchCriteria(testOffer))
	testOffer.MinAge = 18
	testOffer.MaxAge = 99
	testUser.FamilySize = 2
	testOffer.MaxIncome = 0
	a.True(testUser.MatchCriteria(testOffer))
}

func Test_MatchCriteria_FamilySize(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	testOffer.MinAge = 0
	testOffer.MaxAge = 0
	a.False(testUser.MatchCriteria(testOffer))
	testUser.FamilySize = 2
	testOffer.MaxIncome = 0
	a.True(testUser.MatchCriteria(testOffer))
}

func Test_MatchCriteria_Income(t *testing.T) {
	a := assert.New(t)
	testUser := getUser()
	testOffer := getOffer()

	testOffer.MinAge = 0
	testOffer.MaxAge = 0
	testUser.FamilySize = 2
	a.False(testUser.MatchCriteria(testOffer))
	testOffer.MaxIncome = 40000
	a.True(testUser.MatchCriteria(testOffer))
}
