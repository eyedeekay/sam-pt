package sampts

import (
	"github.com/eyedeekay/goSam"
)

func NewSAMServerPlug() (*SAMServerPlug, error) {
	var s SAMServerPlug
	var err error
	s.Client, err = goSam.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// NewSAMClientPlugFromOptions creates a new client, connecting to a specified port
func NewSAMServerPlugFromOptions(opts ...func(*goSam.Client) error) (*SAMServerPlug, error) {
	var c SAMServerPlug
	for _, o := range opts {
		if err := o(c.Client); err != nil {
			return nil, err
		}
	}
	return &c, nil
}
