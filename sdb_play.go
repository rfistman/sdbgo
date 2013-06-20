package main

import (
	"./s3go" // import straight from github? commit?
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "strconv"
	"time"
)

func dothing(r *s3.SDBRequest, cred *s3.SecurityCredentials) {
	r.AddCredentials(cred)

	req, err := r.HttpRequest()

	if err != nil {
		log.Fatal(err)
	}
	log.Println(req)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
		// TODO: interpret errors here
	}
	log.Println(string(contents))
}

func main() {
	cred := &s3.SecurityCredentials{AWSAccessKeyId: os.Getenv("AWS_ACCESS_KEY"), AWSSecretAccessKey: os.Getenv("AWS_SECRET_KEY")}

	//2. get xml response
	m := s3.Strmap{
		"Action": "ListDomains",
	}

	r := s3.NewSDBRequest(m)
	dothing(r, cred)

	nplays := 1234

	m = s3.Strmap{
		"Action":           "PutAttributes",
		"DomainName":       "bn_bs",
		"ItemName":         "1777xy",
		"Attribute.1.Name": "upload_date",
		//"Attribute.1.Value": strconv.FormatInt(time.Now().Unix(), 10),
		// lexographically sortable date. not exactly RFC8601 but looks ok to me
		// http://docs.aws.amazon.com/AmazonSimpleDB/latest/DeveloperGuide/Dates.html
		// choose UTC/GMT so dates in different timezones are comparable
		"Attribute.1.Value": time.Now().UTC().Format(time.RFC3339),
		// lexographically pad out for 1billion plays
		"Attribute.2.Name":  "plays",
		"Attribute.2.Value": fmt.Sprintf("%.9d", nplays),
	}
	r = s3.NewSDBRequest(m)
	dothing(r, cred)

	// sort attr must must appear in predicate
	top10 := "select * from `bn_bs` where `plays` is not null order by `plays` desc limit 10"
	m = s3.Strmap{
		"Action":           "Select",
		"SelectExpression": top10,
	}
	r = s3.NewSDBRequest(m)
	dothing(r, cred)
}
