package oauth

import (
	"reflect"
	"testing"
)

func TestParseIDToken(t *testing.T) {
	testCases := []struct {
		token     string
		expected  *IDToken
		expectErr bool
	}{
		// Invalid input token, returns error
		{
			token:     "foo.bar.baz",
			expectErr: true,
		},
		// Token with the wrong fields returns empty
		{
			token:    "foo.eyJ1c2VySWQiOiJiMDhmODZhZi0zNWRhLTQ4ZjItOGZhYi1jZWYzOTA0NjYwYmQifQ.bar",
			expected: &IDToken{},
		},
		// Token with the correct fields marshals into the IDToken
		{
			token: "foo.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXVkIjoiZm9vIiwic3ViIjoiMTAiLCJlbWFpbCI6ImZvb0BiYXIuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWV9.bar",
			expected: &IDToken{
				Iss:      "accounts.google.com",
				Aud:      "foo",
				Sub:      "10",
				Email:    "foo@bar.com",
				Verified: true,
			},
		},
	}

	for i, testCase := range testCases {
		outToken, err := ParseIDToken(testCase.token)
		if err != nil {
			if !testCase.expectErr {
				t.Errorf("[%d] - Didn't expect an error, but got one: %s\n", i, err.Error())
			}
		} else {
			if testCase.expectErr {
				t.Errorf("[%d] - Expected an error but got none\n", i)
			}
		}

		if !reflect.DeepEqual(outToken, testCase.expected) {
			t.Errorf("[%d] - Expected %+v\n, got: %+v\n", i, testCase.expected, outToken)
		}
	}
}
