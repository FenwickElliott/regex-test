package regex_test

import (
	"errors"
	"reflect"
	"regexp"
	"regexp/syntax"
	"testing"
)

func TestRegex(t *testing.T) {
	const (
		catchAll      = `(.|\n)*?`
		semVarPattern = `^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	)

	_ = errors.New("") // ensures errors is imported in case the deliberate error case is used
	for _, tt := range []struct {
		pattern     string
		testString  string
		shouldMatch bool
		err         error
	}{
		{"", "", true, nil},

		{catchAll, "", true, nil},
		{catchAll, "any", true, nil},
		{catchAll, "thing", true, nil},
		{catchAll, "any-thing", true, nil},
		{catchAll, "any\nthing", true, nil},

		{semVarPattern, "v0.1.2", true, nil},
		{semVarPattern, "broken-v0.1.2", false, nil},

		{semVarPattern, "0.0.0", true, nil},
		{semVarPattern, "0.1.2", true, nil},
		{semVarPattern, "0.11.2", true, nil},
		{semVarPattern, "0.11.2-SNAPSHOT", true, nil},
		{semVarPattern, "0.11.2-OTHER-NOTE", true, nil},

		{semVarPattern, "v0.0.0", true, nil},
		{semVarPattern, "v0.1.2", true, nil},
		{semVarPattern, "v0.11.2", true, nil},
		{semVarPattern, "v0.11.2-SNAPSHOT", true, nil},
		{semVarPattern, "v0.11.2-OTHER-NOTE", true, nil},

		{semVarPattern, "foo", false, nil},
		{semVarPattern, "bar", false, nil},
		{semVarPattern, `no-v0.1.2`, false, nil},
		{semVarPattern, `no-0.1.2`, false, nil},

		{`(`, "", false, &syntax.Error{}},
		// {`)`, "", false, errors.New("SOME_OTHER_ERROR")}, // coment in to validate error checking failure case
	} {

		t.Logf("testing: pattern: '%s', testString: '%s', match: '%t', error: '%v'", tt.pattern, tt.testString, tt.shouldMatch, tt.err)

		matched, err := regexp.MatchString(tt.pattern, tt.testString)
		if err != nil && reflect.TypeOf(err) != reflect.TypeOf(tt.err) {
			t.Errorf("failed with unexpected error:\npattern: '%s'\nexpected '%v', got: '%v'", tt.pattern, tt.err, err)
		} else if matched != tt.shouldMatch {
			t.Errorf("failed with unexpected matching:\npattern: '%s', testString: '%s'\nexpected '%t', got: '%t'", tt.pattern, tt.testString, tt.shouldMatch, matched)
		}
	}
}
