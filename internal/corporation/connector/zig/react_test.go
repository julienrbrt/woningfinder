package zig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReactToOfferRequest(t *testing.T) {
	a := assert.New(t)
	expected := `__hash__=hash&__id__=Portal_Form_SubmitOnly&add=assignmentID&dwellingID=dwellingID`
	req := reactRequest("assignmentID", "dwellingID", "hash")
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(expected, string(body))
}
