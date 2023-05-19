package main

import (
	"github.com/donnol/do"
	"github.com/urfave/cli/v2"
)

var (
	cmds = []*cli.Command{
		{
			Name:        "proxy",
			Aliases:     []string{},
			Usage:       "letgo proxy --localAddr=':54388' --remoteAddr='127.0.0.1:54399'",
			Description: "tcp proxy",
			Action: func(c *cli.Context) error {
				return do.TCPProxy(c.String("localAddr"), c.String("remoteAddr"))
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "localAddr",
					DefaultText: ":54388",
					Value:       ":54388",
				},
				&cli.StringFlag{
					Name:        "remoteAddr",
					DefaultText: "127.0.0.1:54399",
					Value:       "127.0.0.1:54399",
				},
			},
		},
	}
)
