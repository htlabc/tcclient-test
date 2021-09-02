package options

import (
	"github.com/spf13/pflag"
)

type HttpOptions struct {
}

func (o HttpOptions) AddFlags(set *pflag.FlagSet) {

}

func (o HttpOptions) Validate() []error {
	return nil
}

func NewHttpOptions() *HttpOptions {
	return &HttpOptions{}
}
