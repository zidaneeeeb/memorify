package config

import "strings"

// ReadOption controls the behavior of config reader.
type ReadOption func(r *reader) error

// WithSecrets returns ReadOption to replaces all secret
// placeholders in the config file with the given values
// before parsing it as the desired config struct.
//
// Secret placeholder is a string of the following format:
//
//	${secret_name}
//
// Secret placeholder can be put anywhere in the config file,
// but it is most common to put it as/part of config value.
//
// Parameter:
//  - kv is a map containing all the secrets in key-value
//  format. It must contains at least one key-value pair,
//  otherwise it will be ignored.
func WithSecrets(kv map[string]string) ReadOption {
	return func(r *reader) error {
		if len(kv) == 0 {
			return nil
		}

		alternatingKV := make([]string, 0, 2*len(kv))
		for k, v := range kv {
			alternatingKV = append(alternatingKV, "${"+k+"}", v)
		}

		r.secretsReplacer = strings.NewReplacer(alternatingKV...)
		return nil
	}
}

// WithStrictParsing returns ReadOption to enable strict
// parsing mode.
//
// When strict parsing enabled, read process returns error
// when there are any keys/fields in the config file that
// cannot be mapped to the underlying type.
func WithStrictParsing() ReadOption {
	return func(r *reader) error {
		r.strict = true
		return nil
	}
}
