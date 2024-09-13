package parser

import (
	"flag"
)

type Parameters struct {
	EnvPath *string
	Port    *int64
}

func ParseArgs() (*Parameters, error) {
	params := Parameters{}
	params.EnvPath = flag.String("env-path", "", "Path to .env")
	params.Port = flag.Int64("port", 8000, "Api port")
	flag.Parse()

	if *params.EnvPath == "" {
		return nil, flag.ErrHelp
	}
	return &params, nil
}
