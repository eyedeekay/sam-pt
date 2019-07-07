package samptc

import (
	"github.com/eyedeekay/goSam"
)

func NewSAMClientPlug() (*SAMClientPlug, error) {
	var s SAMClientPlug
	var err error
	s.Client, err = goSam.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// NewSAMClientPlugFromOptions creates a new client, connecting to a specified port
func NewSAMClientPlugFromOptions(opts ...func(*goSam.Client) error) (*SAMClientPlug, error) {
	var c SAMClientPlug
	for _, o := range opts {
		if err := o(c.Client); err != nil {
			return nil, err
		}
	}
	return &c, nil
}
