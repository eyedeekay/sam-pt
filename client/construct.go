package samptc

import (
	"github.com/eyedeekay/goSam"
)

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
