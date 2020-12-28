package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/TheYeung1/yata-server/config"
	"github.com/TheYeung1/yata-server/database"
	"github.com/TheYeung1/yata-server/server"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	log.SetReportCaller(true)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "yata"),
	})

	if err != nil {
		log.Fatal(err)
	}

	yataDynamo := &database.DynamoDbYataDatabase{
		Dynamo: dynamodb.New(sess),
	}

	cognitoCfgFile, err := ioutil.ReadFile("env/CognitoConfig.json")
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
