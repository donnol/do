package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/donnol/do"
	"github.com/donnol/do/cmd/letgo/sqlparser"
	"github.com/donnol/do/parser"
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
				pair := c.String("pair")
				if pair != "" {
					wg := new(sync.WaitGroup)
					for _, pai := range strings.Split(pair, ",") {
						wg.Add(1)
						go func(pai string) {
							defer wg.Done()

							parts := strings.Split(pai, "->")

							if err := do.TCPProxy(parts[0], parts[1]); err != nil {
								log.Printf("tcp proxy from %s to %s failed: %v", parts[0], parts[1], err)
								return
							}
						}(pai)
					}
					wg.Wait()

					return nil
				}
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
				&cli.StringFlag{
					Name:  "pair",
					Usage: "specify proxy pairs, not work with localAddr and remoteAddr like ':54388->127.0.0.1:54399,:54389->127.0.0.1:54398'",
				},
			},
		},
		{
			Name:  "sql2struct",
			Usage: `letgo sql2struct 'create table user(id int not null)'`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "ignore",
					DefaultText: "",
					Usage:       "ignore field like order_number",
				},
				&cli.StringFlag{
					Name:        "file",
					Aliases:     []string{"f"},
					DefaultText: "",
					Value:       "",
					Usage:       "specify sql file to input",
				},
				&cli.StringFlag{
					Name:        "output",
					Aliases:     []string{"o"},
					DefaultText: "",
					Value:       "",
					Usage:       "specify output file",
				},
				&cli.StringFlag{
					Name:        "pkg",
					DefaultText: "",
					Value:       "",
					Usage:       "specify package name",
				},
			},
			Action: func(c *cli.Context) error {
				// 标志
				ignoreField := c.String("ignore")
				file := c.String("file")
				output := c.String("output")
				pkg := c.String("pkg")

				sql := ""
				if len(c.Args().Slice()) > 0 {
					sql = c.Args().Slice()[0]
				} else if file != "" {
					data, err := os.ReadFile(file)
					if err != nil {
						fmt.Printf("read file failed: %v\n", err)
						os.Exit(1)
					}
					sql = string(data)
				}

				if sql == "" {
					fmt.Printf("please specify sql like 'create table user(id int not null)' or input file by --file=xxx.sql\n")
					os.Exit(1)
				}

				opt := sqlparser.Option{}
				if ignoreField != "" {
					opt.IgnoreField = append(opt.IgnoreField, ignoreField)
				}

				ss := sqlparser.ParseCreateSQLBatch(sql)
				if ss == nil {
					fmt.Printf("parse sql failed\n")
					os.Exit(1)
				}

				w := os.Stdout
				if output != "" {
					f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
					if err != nil {
						fmt.Printf("open file %s failed: %v\n", output, err)
						os.Exit(1)
					}
					defer f.Close()

					w = f
				}
				buf := new(bytes.Buffer)
				if pkg != "" {
					_, err := buf.WriteString("package " + pkg)
					if err != nil {
						fmt.Printf("write package name failed: %v\n", err)
						os.Exit(1)
					}
				}
				for _, s := range ss {
					if err := s.Gen(buf, opt); err != nil {
						fmt.Printf("gen struct failed: %v\n", err)
						os.Exit(1)
					}
				}
				content, err := parser.Format(output, buf.String(), false)
				if err != nil {
					fmt.Printf("format failed: %v\ncontent: %s\n", err, buf.String())
					os.Exit(1)
				}
				_, err = w.WriteString(content)
				if err != nil {
					fmt.Printf("write to w failed: %v\n", err)
					os.Exit(1)
				}
				fmt.Println()

				return nil
			},
		},
	}
)
