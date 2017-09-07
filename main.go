package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"

	dns "github.com/chenhw2/aliyun-ddns-cli/alidns"
	"github.com/urfave/cli"
)

// AccessKey from https://ak-console.aliyun.com/#/accesskey
type AccessKey struct {
	ID     string
	Secret string
	client *dns.Client
}

func (ak *AccessKey) getClient() *dns.Client {
	if len(ak.ID) <= 0 && len(ak.Secret) <= 0 {
		return nil
	}
	if ak.client == nil {
		ak.client = dns.NewClient(ak.ID, ak.Secret)
		ak.client.SetEndpoint(dns.DNSDefaultEndpointNew)
	}
	return ak.client
}

func (ak AccessKey) String() string {
	return fmt.Sprintf("Access Key: [ ID: %s ;\t Secret: %s ]", ak.ID, ak.Secret)
}

func (ak *AccessKey) listRecord(domain string) (dnsRecords []dns.RecordTypeNew, err error) {
	res, err := ak.getClient().DescribeDomainRecordsNew(
		&dns.DescribeDomainRecordsNewArgs{
			DomainName: domain,
		})
	if err != nil {
		return
	}
	dnsRecords = res.DomainRecords.Record
	return
}

func (ak *AccessKey) delRecord(fulldomain string) (err error) {
	rr := regexp.MustCompile(`\.[^\.]*`).ReplaceAllString(fulldomain, "")
	domain := regexp.MustCompile(`^[^\.]*\.`).ReplaceAllString(fulldomain, "")
	// fmt.Println(rr, domain)
	var target *dns.RecordTypeNew
	if dnsRecords, err := ak.listRecord(domain); err == nil {
		for i := range dnsRecords {
			if dnsRecords[i].RR == rr {
				target = &dnsRecords[i]
				break
			}
		}
	} else {
		return err
	}
	_, err = ak.getClient().DeleteDomainRecord(
		&dns.DeleteDomainRecordArgs{
			RecordId: target.RecordId,
		},
	)
	return
}

func (ak *AccessKey) updateRecord(recordID, rr, value string) (err error) {
	_, err = ak.getClient().UpdateDomainRecord(
		&dns.UpdateDomainRecordArgs{
			RecordId: recordID,
			RR:       rr,
			Value:    value,
			Type:     dns.ARecord,
		})
	return
}

func (ak *AccessKey) addRecord(domain, rr, dmType, value string) (err error) {
	_, err = ak.getClient().AddDomainRecord(
		&dns.AddDomainRecordArgs{
			DomainName: domain,
			RR:         rr,
			Type:       dmType,
			Value:      value,
		})
	return err
}

func (ak *AccessKey) doDDNSUpdate(fulldomain, ipaddr string) (err error) {
	if getDNS(fulldomain) == ipaddr {
		return // Skip
	}
	rr := regexp.MustCompile(`\.[^\.]*`).ReplaceAllString(fulldomain, "")
	domain := regexp.MustCompile(`^[^\.]*\.`).ReplaceAllString(fulldomain, "")
	// fmt.Println(rr, domain)
	var target *dns.RecordTypeNew
	if dnsRecords, err := ak.listRecord(domain); err == nil {
		for i := range dnsRecords {
			if dnsRecords[i].RR == rr {
				target = &dnsRecords[i]
				break
			}
		}
	} else {
		return err
	}
	if target == nil {
		err = ak.addRecord(domain, rr, "A", ipaddr)
	} else if target.Value != ipaddr {
		err = ak.updateRecord(target.RecordId, target.RR, ipaddr)
	}
	return err
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
					Usage: "Specific `DomainName`. like aliyun.com",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"))
				if dnsRecords, err := accessKey.listRecord(c.String("domain")); err != nil {
					fmt.Printf("%+v", err)
				} else {
					for _, v := range dnsRecords {
						fmt.Println(v.RR+`.`+v.DomainName, fmt.Sprintf(" %5s ", v.Type), v.Value)
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
					Name:  "d",
					Usage: "Specific `FullDomainName`. like ddns.aliyun.com",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"))
				if err := accessKey.delRecord(c.String("domain")); err != nil {
					fmt.Printf("%+v", err)
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
					log.Printf("%+v", err)
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
						log.Printf("%+v", err)
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
	rand.Seed(time.Now().UnixNano())
	return nil
}
