package config

import "time"

// Duration is an alias for time.Duration that implements
// custom YAML Unmarshaler and Marshaler to allows marshal
// and unmarshal a time.Duration value to and from its
// string form.
type Duration time.Duration

// UnmarshalYAML implements custom YAML unmarshaler.
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

// MarshalYAML implements custom YAML marshaler.
func (d Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(d).String(), nil
}
