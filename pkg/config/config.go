package config

import "os"

// ReadFile reads config file in the given path and parses it
// into the desired config struct as cfg.
func ReadFile(path string, cfg interface{}, options ...ReadOption) error {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// create reader
	r, err := newReader(options...)
	if err != nil {
		return err
	}

	return r.read(f, cfg)
}
