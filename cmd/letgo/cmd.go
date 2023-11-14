package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
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

				tcpProxyHandler := func(lconn, rconn net.Conn) {
					go func() {
						defer func() {
							rconn.Close()
							log.Printf("close remote conn\n")
						}()

						n, err := copyBuffer(lconn, rconn, nil)
						if err == io.EOF {
							log.Printf("remote conn read EOF")
							return
						}
						if err != nil {
							log.Printf("copy from remote to local failed: %v\n", err)
							return
						}
						log.Printf("copy %d bytes from remote to local\n", n)
					}()
					go func() {
						defer func() {
							lconn.Close()
							log.Printf("close local conn\n")
						}()

						n, err := copyBuffer(rconn, lconn, nil)
						if err == io.EOF {
							log.Printf("local conn read EOF")
							return
						}
						if err != nil {
							log.Printf("copy from local to remote failed: %v\n", err)
							return
						}
						log.Printf("copy %d bytes from local to remote\n", n)
					}()
				}

				pair := c.String("pair")
				if pair != "" {
					wg := new(sync.WaitGroup)
					for _, pai := range strings.Split(pair, ",") {
						wg.Add(1)
						go func(pai string) {
							defer wg.Done()

							parts := strings.Split(pai, "->")

							if err := do.TCPProxy(parts[0], parts[1], tcpProxyHandler); err != nil {
								log.Printf("tcp proxy from %s to %s failed: %v", parts[0], parts[1], err)
								return
							}
						}(pai)
					}
					wg.Wait()

					return nil
				}
				return do.TCPProxy(c.String("localAddr"), c.String("remoteAddr"), tcpProxyHandler)
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
				&cli.StringFlag{
					Name:        "tablePrefix",
					Aliases:     []string{"tp"},
					DefaultText: "",
					Value:       "",
					Usage:       "specify table prefix, trim it when convert to struct name",
				},
			},
			Action: func(c *cli.Context) error {
				// 标志
				ignoreField := c.String("ignore")
				file := c.String("file")
				output := c.String("output")
				pkg := c.String("pkg")
				tp := c.String("tablePrefix")

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
				if tp != "" {
					opt.TrimTablePrefix = tp
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
				if pkg == "" {
					pkg = "main"
				}
				if _, err := buf.WriteString("package " + pkg); err != nil {
					fmt.Printf("write package name failed: %v\n", err)
					os.Exit(1)
				}

				haveEnum := false
				body := new(bytes.Buffer)
				for _, s := range ss {
					if err := s.Gen(body, opt); err != nil {
						fmt.Printf("gen struct failed: %v\n", err)
						os.Exit(1)
					}
					if s.HaveEnum {
						haveEnum = true
					}
				}
				if haveEnum {
					if _, err := buf.Write([]byte("\nimport \"github.com/donnol/do\"\n")); err != nil {
						return err
					}
				}
				_, err := buf.Write(body.Bytes())
				if err != nil {
					fmt.Printf("write body failed: %v\n", err)
					os.Exit(1)
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
		{
			Name:  "mock",
			Usage: `letgo mock -p=github.com/xxx/yyy -r`,
			Description: `gen interface mock struct, like: type I interface { String() string }, 
			gen mock: 
				type Mock struct { StringFunc func() string } 
				var _ I = &Mock{}
				func (mock *Mock) String() string {
					return mock.StringFunc()
				}
			after that, you can use like below:
				var mock = &Mock{
					// init the func like the normal field
					StringFunc: func() string {
						return "do"
					},	
				}
				fmt.Println(mock.String())`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "path",
					Aliases:     []string{"p"},
					DefaultText: "",
					Usage:       "specify package path",
				},
				&cli.BoolFlag{
					Name:        "recursive",
					Aliases:     []string{"r"},
					DefaultText: "",
					Usage:       "if iterate recursive",
				},
				&cli.StringFlag{
					Name:        "mode",
					Aliases:     []string{""},
					DefaultText: "normal",
					Value:       "",
					Usage:       "specify mode like normal or offsite",
				},
			},
			Action: func(c *cli.Context) error {
				path := c.String("path")
				rec := c.Bool("recursive")
				mode := c.String("mode")

				ip := &parser.ImportPath{}
				curdir, paths, err := getPaths(ip, path, rec)
				if err != nil {
					log.Fatal(err)
				}
				if len(paths) == 0 {
					log.Fatalf("找不到有效路径，请使用-p指定或设置-r！")
				}
				// fmt.Printf("paths: %+v\n", paths)

				// 解析
				p := parser.NewParser(parser.Option{
					Op:                parser.OpMock,
					UseSourceImporter: true,
				})
				pkgs, err := p.ParseByGoPackages(paths...)
				if err != nil {
					log.Fatal(err)
				}
				for _, pkg := range pkgs.Pkgs {
					if err = pkg.SaveMock(mode, curdir, ""); err != nil {
						log.Printf("gen mock failed, pkg: %+v, err: %+v\n", pkg, err)
					}
				}

				return nil
			},
		},
	}
)

func getPaths(ip *parser.ImportPath, path string, rec bool) (string, []string, error) {
	var err error

	dir, err := os.Getwd()
	if err != nil {
		return dir, nil, err
	}

	var paths []string

	if path == "" {
		path, err = ip.GetByCurrentDir()
		if err != nil {
			return dir, nil, err
		}

		haveGoFile, err := checkDirHaveGoFile(dir)
		if err != nil {
			return dir, nil, err
		}
		if haveGoFile {
			paths = append(paths, path)
		}
	} else {
		// 手动指定的path，不校验是否存在go文件，由用户自己保证
		paths = append(paths, path)
	}

	modDir, modPath, err := ip.GetModFilePath(dir)
	if err != nil {
		return dir, nil, err
	}
	fmt.Printf("dir: %s, modDir: %s, modPath: %s\n", dir, modDir, modPath)

	if rec {
		dirs, err := collectGoFileDir(dir)
		if err != nil {
			return dir, nil, err
		}
		for _, d := range dirs {
			paths = append(paths, strings.Replace(d, dir, modDir, -1))
		}
	}
	return dir, paths, nil
}

// collectGoFileDir 在指定目录下收集含有go文件的子目录
func collectGoFileDir(dir string) ([]string, error) {
	var dirs []string
	if err := filepath.Walk(dir, filepath.WalkFunc(func(childDir string, info os.FileInfo, ierr error) error {
		if ierr != nil {
			fmt.Printf("walk got err: %+v\n", ierr)
		}

		if childDir == dir {
			return nil
		}
		// 获取所需目录
		if !info.IsDir() {
			return nil
		}
		haveGoFile, err := checkDirHaveGoFile(childDir)
		if err != nil {
			return err
		}
		// 过滤没有go文件的
		if !haveGoFile {
			return nil
		}

		dirs = append(dirs, childDir)

		return nil
	})); err != nil {
		return nil, err
	}

	return dirs, nil
}

func checkDirHaveGoFile(dir string) (bool, error) {

	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	haveGoFile := false
	for _, fi := range fileInfos {
		ext := filepath.Ext(fi.Name())
		if ext == ".go" {
			haveGoFile = true
			break
		}
	}

	return haveGoFile, nil
}
