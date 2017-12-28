package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// https://github.com/aliyun/alibaba-cloud-sdk-go/blob/079e0a218a45e92dfdbde0bc35d26973b59091d5/services/alidns/describe_domain_records.go#L80
type RecordType struct {
	DomainName string `json:"DomainName" xml:"DomainName"`
	RecordId   string `json:"RecordId" xml:"RecordId"`
	RR         string `json:"RR" xml:"RR"`
	Type       string `json:"Type" xml:"Type"`
	Value      string `json:"Value" xml:"Value"`
}

type DomainRecordsType struct {
	Record []RecordType
}

type DomainRecordsResponseType struct {
	DomainRecords DomainRecordsType
}

// AccessKey from https://ak-console.aliyun.com/#/accesskey
type AccessKey struct {
	ID     string
	Secret string
	client *alidns.Client
}

func (ak *AccessKey) getClient() *alidns.Client {
	if len(ak.ID) <= 0 && len(ak.Secret) <= 0 {
		return nil
	}
	if ak.client == nil {
		ak.client, _ = alidns.NewClientWithAccessKey("cn-shenzhen", ak.ID, ak.Secret)
	}
	// fmt.Printf("%+v\n", *ak.client)
	return ak.client
}

func (ak AccessKey) String() string {
	return fmt.Sprintf("Access Key: [ ID: %s ;\t Secret: %s ]", ak.ID, ak.Secret)
}

func (ak *AccessKey) listRecord(domain string) (dnsRecords []RecordType, err error) {
	req := alidns.CreateDescribeDomainRecordsRequest()
	req.DomainName = domain
	resp, _ := ak.client.DescribeDomainRecords(req)

	respType := &DomainRecordsResponseType{}
	err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &respType)
	dnsRecords = respType.DomainRecords.Record
	return
}

func (ak *AccessKey) delRecord(fulldomain string) (err error) {
	rr := regexp.MustCompile(`\.[^\.]*`).ReplaceAllString(fulldomain, "")
	domain := regexp.MustCompile(`^[^\.]*\.`).ReplaceAllString(fulldomain, "")
	// fmt.Println(rr, domain)
	var target *RecordType
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
	if nil == target {
		return errors.New("Error: Cant Find Target Domain")
	}
	req := alidns.CreateDeleteDomainRecordRequest()
	req.RecordId = target.RecordId
	_, err = ak.client.DeleteDomainRecord(req)
	return
}

func (ak *AccessKey) updateRecord(dnsRecord *RecordType, value string) (err error) {
	req := alidns.CreateUpdateDomainRecordRequest()
	req.Value = value
	req.RR = dnsRecord.RR
	req.Type = dnsRecord.Type
	req.RecordId = dnsRecord.RecordId
	_, err = ak.client.UpdateDomainRecord(req)
	return
}

func (ak *AccessKey) addRecord(domain, rr, dmType, value string) (err error) {
	req := alidns.CreateAddDomainRecordRequest()
	req.RR = rr
	req.Type = dmType
	req.Value = value
	req.DomainName = domain
	_, err = ak.client.AddDomainRecord(req)
	return err
}

func (ak *AccessKey) doDDNSUpdate(fulldomain, ipaddr string) (err error) {
	if getDNS(fulldomain) == ipaddr {
		return // Skip
	}
	rr := regexp.MustCompile(`\.[^\.]*`).ReplaceAllString(fulldomain, "")
	domain := regexp.MustCompile(`^[^\.]*\.`).ReplaceAllString(fulldomain, "")
	// fmt.Println(rr, domain)
	var target *RecordType
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
		err = ak.updateRecord(target, ipaddr)
	}
	return err
}
