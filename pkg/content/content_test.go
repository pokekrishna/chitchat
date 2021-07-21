package content

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateType(t *testing.T) {
	testcases := []struct{
		contentType string
		expect string
		description string
	}{
		{
			contentType: "xml",
			expect:      TypeNotSupported,
			description: "Type xml should not be supported",
		},
		{
			contentType: "foo",
			expect:      TypeNotSupported,
			description: "Type foo should not be supported",
		},
		{
			contentType: "Application/JSoN",
			expect:      TypeJSON,
			description: "Type checks should be case insensitive",
		},
		{
			contentType: "",
			expect:      TypeNotSupported,
			description: "blank type should not be supported",
		},
	}

	for _, tc := range testcases{
		t.Run(tc.description, func(t *testing.T) {
			assert.Equal(t, tc.expect, ValidateType(tc.contentType))
		})
	}
}
