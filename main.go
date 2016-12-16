package main

import (
	"errors"
	"fmt"
	"os"

	"regexp"

	"github.com/denverdino/aliyungo/dns"
	"github.com/urfave/cli"
)

// AccessKey from https://ak-console.aliyun.com/#/accesskey
type AccessKey struct {
	ID     string
	Secret string
}

func (ak *AccessKey) isFiled() bool {
	return len(ak.ID) > 0 && len(ak.Secret) > 0
}

func (ak AccessKey) String() string {
	return fmt.Sprintf("Access Key: [ ID: %s ;\t Secret: %s ]", ak.ID, ak.Secret)
}

func (ak *AccessKey) list(domain string) (dnsRecords []dns.RecordType, err error) {
	client := dns.NewClient(ak.ID, ak.Secret)
	res, err := client.DescribeDomainRecords(
		&dns.DescribeDomainRecordsArgs{
			DomainName: domain,
		})
	if err != nil {
		return
	}
	dnsRecords = res.DomainRecords.Record
	return
}

func (ak *AccessKey) update(recordID, rr, value string) (err error) {
	client := dns.NewClient(ak.ID, ak.Secret)
	_, err = client.UpdateDomainRecord(
		&dns.UpdateDomainRecordArgs{
			RecordId: recordID,
			RR:       rr,
			Value:    value,
			Type:     dns.ARecord,
		})
	return
}

var (
	accessKey AccessKey
	version   = "MISSING build version [git hash]"
)

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
					Usage: "specify `DomainName`. like aliyun.com",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"))
				dnsRecords, err := accessKey.list(c.String("domain"))
				if err != nil {
					fmt.Printf("%+v", err)
				} else {
					for _, v := range dnsRecords {
						fmt.Println(v.RR+`.`+v.DomainName, v.Value)
					}
				}
				return nil
			},
		},
		{
			Name:     "update",
			Category: "DDNS",
			Usage:    "Update AliYun's DNS DomainRecords Record",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "specify `DomainName`. like ddns.aliyun.com",
				},
				cli.StringFlag{
					Name:  "ipaddr, i",
					Usage: "specify `IP`. like 1.2.3.4",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"), c.String("ipaddr"))
				rr := regexp.MustCompile(`\.[^\.]*`).ReplaceAllString(c.String("domain"), "")
				domain := regexp.MustCompile(`^[^\.]*\.`).ReplaceAllString(c.String("domain"), "")
				// fmt.Println(rr, domain)
				dnsRecords, err := accessKey.list(domain)
				var target *dns.RecordType
				if err != nil {
					fmt.Printf("%+v", err)
				} else {
					for i := range dnsRecords {
						if dnsRecords[i].RR == rr {
							target = &dnsRecords[i]
							break
						}
					}
				}
				if target != nil {
					ipaddr := getIPipcn()
					if ipaddr == target.Value {
						fmt.Println(target.RR+`.`+target.DomainName, ipaddr)
						return nil
					}
					err = accessKey.update(target.RecordId, target.RR, c.String("ipaddr"))
					if err != nil {
						fmt.Printf("%+v", err)
					} else {
						fmt.Println(target.RR+`.`+target.DomainName, c.String("ipaddr"))
					}
				} else {
					fmt.Println("Can't Find target")
				}
				return nil
			},
		},
		{
			Name:     "auto-update",
			Category: "DDNS",
			Usage:    "Auto-Update AliYun's DNS DomainRecords Record, Get IP using http://ip.cn",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "specify `DomainName`. like ddns.aliyun.com",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"))
				rr := regexp.MustCompile(`\.[^\.]*`).ReplaceAllString(c.String("domain"), "")
				domain := regexp.MustCompile(`^[^\.]*\.`).ReplaceAllString(c.String("domain"), "")
				// fmt.Println(rr, domain)
				dnsRecords, err := accessKey.list(domain)
				var target *dns.RecordType
				if err != nil {
					fmt.Printf("%+v", err)
				} else {
					for i := range dnsRecords {
						if dnsRecords[i].RR == rr {
							target = &dnsRecords[i]
							break
						}
					}
				}
				if target != nil {
					ipaddr := getIPipcn()
					if ipaddr == target.Value {
						fmt.Println(target.RR+`.`+target.DomainName, ipaddr)
						return nil
					}
					err = accessKey.update(target.RecordId, target.RR, ipaddr)
					if err != nil {
						fmt.Printf("%+v", err)
					} else {
						fmt.Println(target.RR+`.`+target.DomainName, ipaddr)
					}
				} else {
					fmt.Println("Can't Find target")
				}
				return nil
			},
		},
		{
			Name:     "getip",
			Category: "GET-IP",
			Usage:    "Get IP using http://ip.cn",
			Action: func(c *cli.Context) error {
				// fmt.Println(c.Command.Name, "task: ", c.Command.Usage)
				fmt.Println(getIPipcn())
				return nil
			},
		},
		{
			Name:     "getip-intl",
			Category: "GET-IP",
			Usage:    "Get IP using http://ipinfo.io",
			Action: func(c *cli.Context) error {
				// fmt.Println(c.Command.Name, "task: ", c.Command.Usage)
				fmt.Println(getIPipio())
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
	}
	app.Action = func(c *cli.Context) error {
		return appInit(c)
	}
	app.Run(os.Args)
}

func appInit(c *cli.Context) error {
	accessKey.ID = c.GlobalString("access-key-id")
	accessKey.Secret = c.GlobalString("access-key-secret")
	if !accessKey.isFiled() {
		cli.ShowAppHelp(c)
		return errors.New("access-key is empty")
	}
	return nil
}
