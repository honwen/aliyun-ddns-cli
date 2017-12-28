package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/urfave/cli"
)

var (
	accessKey AccessKey
	version   = "MISSING build version [git hash]"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app := cli.NewApp()
	app.Name = "aliddns"
	app.Usage = "aliyun-ddns-cli"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:     "list",
			Category: "DDNS",
			Usage:    "List AliYun's DNS DomainRecords Record",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "Specific `DomainName`. like aliyun.com",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"))
				if dnsRecords, err := accessKey.listRecord(c.String("domain")); err != nil {
					log.Printf("%+v\n", err)
				} else {
					fmt.Println(fmt.Sprintf("%16s %8s   %s", "DomainName", "Type", "Value"))
					fmt.Println("==========================================================")
					for _, v := range dnsRecords {
						fmt.Println(fmt.Sprintf("%16s %8s   %s", v.RR+`.`+v.DomainName, v.Type, v.Value))
					}
				}
				return nil
			},
		},
		{
			Name:     "delete",
			Category: "DDNS",
			Usage:    "Delete AliYun's DNS DomainRecords Record",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "Specific `FullDomainName`. like ddns.aliyun.com",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"))
				if err := accessKey.delRecord(c.String("domain")); err != nil {
					log.Printf("%+v\n", err)
				} else {
					fmt.Println(c.String("domain"), "Deleted")
				}
				return nil
			},
		},
		{
			Name:     "update",
			Category: "DDNS",
			Usage:    "Update AliYun's DNS DomainRecords Record, Create Record if not exist",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "Specific `DomainName`. like ddns.aliyun.com",
				},
				cli.StringFlag{
					Name:  "ipaddr, i",
					Usage: "Specific `IP`. like 1.2.3.4",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"), c.String("ipaddr"))
				if err := accessKey.doDDNSUpdate(c.String("domain"), c.String("ipaddr")); err != nil {
					log.Printf("%+v\n", err)
				} else {
					log.Println(c.String("domain"), c.String("ipaddr"))
				}
				return nil
			},
		},
		{
			Name:     "auto-update",
			Category: "DDNS",
			Usage:    "Auto-Update AliYun's DNS DomainRecords Record, Get IP using its getip",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "Specific `DomainName`. like ddns.aliyun.com",
				},
				cli.Int64Flag{
					Name:  "redo, r",
					Value: 0,
					Usage: "redo Auto-Update, every N `Seconds`; Disable if N less than 10",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"), c.Int64("redo"))
				redoDurtion := c.Int64("redo")
				for {
					autoip := getIP()
					if err := accessKey.doDDNSUpdate(c.String("domain"), autoip); err != nil {
						log.Printf("%+v\n", err)
					} else {
						log.Println(c.String("domain"), autoip)
					}
					if redoDurtion < 10 {
						break // Disable if N less than 10
					}
					time.Sleep(time.Duration(redoDurtion) * time.Second)
				}
				return nil
			},
		},
		{
			Name:     "getip",
			Category: "GET-IP",
			Usage:    fmt.Sprintf("      Get IP Combine %d different Web-API", len(ipAPI)),
			Action: func(c *cli.Context) error {
				// fmt.Println(c.Command.Name, "task: ", c.Command.Usage)
				fmt.Println(getIP())
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "access-key-id, id",
			Usage: "AliYun's Access Key ID",
		},
		cli.StringFlag{
			Name:  "access-key-secret, secret",
			Usage: "AliYun's Access Key Secret",
		},
		cli.StringSliceFlag{
			Name:  "ipapi, api",
			Usage: "Web-API to Get IP, like: http://myip.ipip.net",
		},
	}
	app.Action = func(c *cli.Context) error {
		return appInit(c)
	}
	app.Run(os.Args)
}

func appInit(c *cli.Context) error {
	accessKey.ID = c.GlobalString("access-key-id")
	accessKey.Secret = c.GlobalString("access-key-secret")
	if accessKey.getClient() == nil {
		cli.ShowAppHelp(c)
		return errors.New("access-key is empty")
	}

	newIPAPI := make([]string, 0)
	for _, api := range c.GlobalStringSlice("ipapi") {
		if !regexp.MustCompile(`^https?://.*`).MatchString(api) {
			api = "http://" + api
		}
		if regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`).MatchString(api) {
			newIPAPI = append(newIPAPI, api)
		}
	}
	if len(newIPAPI) > 0 {
		ipAPI = newIPAPI
	}

	return nil
}
