package parser

import (
	"flag"
)

type Parameters struct {
	EnvPath *string
}

func ParseArgs() (*Parameters, error) {
	params := Parameters{}
	params.EnvPath = flag.String("env-path", "", "Path to .env")
	flag.Parse()

	if *params.EnvPath == "" {
		return nil, flag.ErrHelp
	}
	return &params, nil
}
