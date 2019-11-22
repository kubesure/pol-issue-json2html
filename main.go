package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	io "io/ioutil"
	"log"
	"strconv"
	"time"

	e "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//composite struct represent json meta file in S3 bucket put by policy issued event handler
type polmetadata struct {
	Email email      `json:"email"`
	Data  policydata `json:"data"`
}

//data of the policy issue
type policydata struct {
	Name         string `json:"name"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	City         string `json:"city"`
	PinCode      int    `json:"pinCode"`
	MobileNumber int    `json:"mobileNumber"`
	PolicyNumber int    `json:"policyNumber"`
}

//email address of the policy holder
type email struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event e.S3Event) (string, error) {

	for _, record := range event.Records {
		log.Println("bucket" + record.S3.Bucket.Name)
		log.Println("object " + record.S3.Object.Key)
		err := processEvent(record)
		if err != nil {
			log.Println(err)
			return "processing error ", err
		}
	}

	return fmt.Sprintf("HTML generated sucessfully..."), nil
}

//Generates html template for the PDF service for PDF generation
func processEvent(record e.S3EventRecord) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(record.S3.Bucket.Name),
		Key:    aws.String(record.S3.Object.Key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		log.Println("error getting object " + err.Error())
		return err
	}
	defer result.Body.Close()
	bodyBytes, err := io.ReadAll(result.Body)
	metaData := string(bodyBytes)
	log.Println(metaData)
	pm, err := marshallReq(metaData)
	if err != nil {
		return err
	}

	htmlBytes, errhtml := generateHTML(pm)
	if errhtml != nil {
		log.Println("err while generating html")
		return errhtml
	}

	newKey := "unprocessed/" + strconv.Itoa(pm.Data.PolicyNumber) + ".html"
	log.Println(newKey)

	pinput := s3.PutObjectInput{
		Bucket: aws.String(record.S3.Bucket.Name),
		Key:    aws.String(newKey),
		Body:   bytes.NewReader(htmlBytes),
	}

	presult, putErr := svc.PutObject(&pinput)

	if putErr != nil {
		return putErr
	}
	log.Println(presult)

	return nil
}

//generates HTML from policy metadata
func generateHTML(metaData *polmetadata) ([]byte, error) {
	fmap := template.FuncMap{
		"currentdate": currentdate,
	}

	t, errp := template.New("esyhealth-pdf.html").Funcs(fmap).ParseFiles("esyhealth-pdf.html")
	if errp != nil {
		return nil, errp
	}
	buff := new(bytes.Buffer)
	err := t.Execute(buff, metaData)
	if err != nil {
		return nil, err
	}
	log.Println(buff.String())
	return buff.Bytes(), nil
}

func marshallReq(data string) (*polmetadata, error) {
	var pd polmetadata
	err := json.Unmarshal([]byte(data), &pd)
	if err != nil {

		return nil, err
	}
	return &pd, nil
}

func currentdate() (string, error) {
	const layoutISO = "2006-01-02"
	const custom = "Mon Jan _2 15:04:05 2006"
	currentDate := time.Now().Format(custom)
	return currentDate, nil
}
