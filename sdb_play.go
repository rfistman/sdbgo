package main

import (
	"./s3go" // import straight from github? commit?
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strconv"
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
	//3. PutAttributes (still a GET)
	m := s3.Strmap{
		"Action": "ListDomains",
	}

	r := s3.NewSDBRequest("GET", m)
	dothing(r, cred)

	m = s3.Strmap{
		"Action":     "PutAttributes",
		"DomainName": "bn_bs",
		//"Attribute.1.Name":  "slug",
		//"Attribute.1.Value": "12afc3",
		"ItemName":          "12afcd",
		"Attribute.1.Name":  "upload_date",
		"Attribute.1.Value": strconv.FormatInt(time.Now().Unix(), 10),
		"Attribute.2.Name":  "plays",
		"Attribute.2.Value": "0",
	}
	r = s3.NewSDBRequest("GET", m)
	dothing(r, cred)

}
