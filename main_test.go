package main

import (
	"fmt"
	"testing"

	e "github.com/aws/aws-lambda-go/events"
)

const metaData = `{"email":{"from":"edakghar@gmail.com","to":"pras.p.in@gmail.com"},
"data":{"name":"Usha Patel","addressLine1":"ketaki","addressLine2":"maneklal","addressLine3":"Ghatkopar",
"city":"mumbai","pinCode":400086,"mobileNumber":9821284567,"policyNumber":1234567890},
"status":{"mailSent":true,"pdfCreated":true}}`

var pmetadata polmetadata

func TestMarshallPolData(t *testing.T) {
	pm, err := marshallReq(metaData)
	if err != nil {
		t.Errorf("marshall err %v", err)
	}
	fmt.Println(pm)
}

func TestGenerateHTML(t *testing.T) {
	pm, err := marshallReq(metaData)
	if err != nil {
		t.Errorf("marshall err %v", err)
	}

	html, err := generateHTML(pm)

	if err != nil {
		t.Errorf("Err while html generation %v", err)
	}

	if len(html) == 0 {
		t.Errorf("html not generated  %v", err)
	}
}

func TestCurrentDate(t *testing.T) {
	date, err := currentdate()

	if err != nil {
		t.Errorf("time formatting error")
	}

	if date != "2019-06-30" {
		t.Errorf("Incorrect date format")
	}
}

func TestProcessEvent(t *testing.T) {
	bucket := e.S3Bucket{Name: "io.kubesure-esyhealth-policy-issued-dev"}
	object := e.S3Object{Key: "unprocessed/1234567890.json"}
	r := e.S3EventRecord{}
	r.S3.Bucket = bucket
	r.S3.Object = object
	err := processEvent(r)

	if err != nil {
		t.Errorf("S3 Event Processed %v", err)
	}
}
