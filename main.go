package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/TheYeung1/yata-server/config"
	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/server"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

var (
	awsRegion            = flag.String("aws-region", "us-west-2", "aws region")
	awsCredentialProfile = flag.String("aws-profile", "yata", "aws credential profile; create with 'aws configure --profile <name>'")
	cognitoConfigFile    = flag.String("cognito-config", "env/CognitoConfig.json", "cognito config file; see env/SampleConfig.json for reference")
	listsTableName       = flag.String("lists-table", "ListTable", "lists DynamoDB table name")
	itemsTableName       = flag.String("items-table", "ItemsTable", "items DynamoDB table name")
)

func init() {
	flag.Parse()
}

func main() {
	log.SetReportCaller(true)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*awsRegion),
		Credentials: credentials.NewSharedCredentials("", *awsCredentialProfile),
	})

	if err != nil {
		log.Fatal(err)
	}

	yataDynamo := &database.DynamoDbYataDatabase{
		Dynamo:         dynamodb.New(sess),
		ListsTableName: *listsTableName,
		ItemsTableName: *itemsTableName,
	}

	cognitoCfgFile, err := ioutil.ReadFile(*cognitoConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	var cognitoConfig config.AwsCognitoUserPoolConfig
	if err := json.Unmarshal(cognitoCfgFile, &cognitoConfig); err != nil {
		log.Fatal(err)
	}

	s := server.Server{
		CognitoCfg: cognitoConfig,
		Ydb:        yataDynamo,
	}
	s.Start()
}
