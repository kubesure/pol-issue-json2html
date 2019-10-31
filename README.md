# pol-issue-json2html

### setup Dev

1. Build function zip
    
```
    go build main.go
    zip function.zip main esyhealth-pdf.html
    unzip -l function.zip
```    

4. Create S3 bucket, permission & folder

```
    aws s3api create-bucket --bucket=io.kubesure-esyhealth-policy-issued-dev --region us-east-1

    aws s3api put-public-access-block --bucket io.kubesure-esyhealth-policy-issued-dev \
    --public-access-block-configuration BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true

    aws s3api put-object --bucket io.kubesure-esyhealth-policy-issued-dev --key unprocessed/blank.txt
    aws s3api put-object --bucket io.kubesure-esyhealth-policy-issued-dev --key processed/blank.txt    
```

2. Create lambda exection role 'lambda_execution_role' add policies s3 full & Lambda basic acces

3. Create Lambda Function

```
    aws lambda create-function --function-name esyhealth-pol-issue-json2html \
    --zip-file fileb://function.zip --handler main --runtime go1.x \
    --role arn:aws:iam::708908412990:role/lambda_execution_role --description "Create a html template file from a esyhealth policy issued data"
```

4. Create S3 permission for lambda function

```
    aws lambda add-permission --function-name esyhealth-pol-issue-json2html \
    --action lambda:InvokeFunction \
    --statement-id s3-account \
    --principal s3.amazonaws.com \
    --source-arn arn:aws:s3:::io.kubesure-esyhealth-policy-issued-dev \
    --source-account 708908412990
```

5. Create S3 trigger/notification for lambda function

```
    aws s3api put-bucket-notification-configuration \
    --bucket io.kubesure-esyhealth-policy-issued-dev \
    --notification-configuration file://s3-notification.json
```