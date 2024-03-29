package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func StringRequired(prop, value string) Fn {
	return func() error {
		if strings.Trim(value, " ") == "" {
			return fmt.Errorf("%s is required", prop)
		}
		return nil
	}
}

func StringStartsWithIfNotEmpty(prop, value, prefix string) Fn {
	v := strings.Trim(value, " ")

	return func() error {
		if len(v) == 0 {
			return nil
		}

		if !strings.HasPrefix(v, prefix) {
			return fmt.Errorf("%s (value:%s) must be started with %s", prop, value, prefix)
		}

		return nil
	}
}

func StringStartsWith(prop, value, prefix string) Fn {
	return func() error {
		if err := StringRequired(prop, value)(); err != nil {
			return err
		}
		if err := StringStartsWithIfNotEmpty(prop, value, prefix)(); err != nil {
			return err
		}
		return nil
	}
}

func StringStartsWithOneOf(prop, value string, prefixes []string) Fn {
	return func() error {
		for i := range prefixes {
			if err := StringStartsWithIfNotEmpty(prop, value, prefixes[i])(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("%s (value:%s) prefix must be started with one of %s", prop, value, prefixes)
	}
}

func StringUri(prop, value string) Fn {
	return func() error {
		if err := StringRequired(prop, value)(); err != nil {
			return err
		}

		if _, err := url.ParseRequestURI(value); err != nil {
			return fmt.Errorf("%s (error:%s) is not a valid uri", prop, err.Error())
		}
		return nil
	}
}

func StringLen(prop, value string, min, max int) Fn {
	return func() error {
		if len(value) < min {
			return fmt.Errorf("%s (len:%d) length must be greater than or equal %d", prop, len(value), min)
		}
		if len(value) > max {
			return fmt.Errorf("%s (len:%d) length must be less than or equal %d", prop, len(value), max)
		}
		return nil
	}
}

func StringLenIfNotEmpty(prop, value string, min, max int) Fn {
	return func() error {
		if len(value) == 0 {
			return nil
		}
		return StringLen(prop, value, min, max)()
	}
}

func StringOneOf(prop, value string, oneOf []string) Fn {
	m := map[string]bool{}
	for _, o := range oneOf {
		m[o] = true
	}

	return func() error {
		if err := StringRequired(prop, value)(); err != nil {
			return err
		}

		if _, has := m[value]; !has {
			return fmt.Errorf("%s (value:%s) must be one of %q", prop, value, oneOf)
		}

		return nil
	}
}

func StringAlphaNumericUnderscore(prop, value string) Fn {
	pattern := "^[0-9a-zA-z_]+$"
	r := regexp.MustCompile(pattern)
	return func() error {
		if !r.MatchString(value) {
			return fmt.Errorf("%s (value:%s) is not matched the pattern %q", prop, value, pattern)
		}

		return nil
	}
}

func StringAlphaNumericUnderscoreHyphenDot(prop, value string) Fn {
	pattern := "^[0-9a-zA-z_.-]+$"
	r := regexp.MustCompile(pattern)
	return func() error {
		if !r.MatchString(value) {
			return fmt.Errorf("%s (value:%s) is not matched the pattern %q", prop, value, pattern)
		}

		return nil
	}
}
