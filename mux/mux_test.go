package mux_test

import (
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	re := regexp.MustCompile(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`)
	t.Logf("to string: %s", re.String())

	captures := make(map[string]string)

	match := re.FindStringSubmatch("1234.5678.9")
	if match == nil {
		t.Errorf("regexp error")
	}

	for i, name := range re.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}

	t.Logf("%v", captures)

}
