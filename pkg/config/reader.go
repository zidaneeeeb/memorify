package config

import (
	"io"
	"strings"

	"gopkg.in/yaml.v2"
)

// reader is a config reader.
//
// Currently only supports config file in yaml format.
type reader struct {
	secretsReplacer *strings.Replacer
	strict          bool
}

// newReader returns a new reader with the given options.
func newReader(options ...ReadOption) (*reader, error) {
	r := &reader{}

	// apply options
	for _, opt := range options {
		if err := opt(r); err != nil {
			return nil, nil
		}
	}

	return r, nil
}

// read reads config from the given io.Reader and parses it
// in to the desired config struct.
func (r *reader) read(re io.Reader, cfg interface{}) error {
	// replace secret placeholder
	if r.secretsReplacer != nil {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, re)
		if err != nil {
			return err
		}
		re = strings.NewReader(
			r.secretsReplacer.Replace(buf.String()),
		)
	}

	// create yaml decoder
	dec := yaml.NewDecoder(re)
	if r.strict {
		dec.SetStrict(true)
	}

	// decode in to the desired config struct.
	return dec.Decode(cfg)
}
