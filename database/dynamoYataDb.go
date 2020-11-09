package database

import (
	log "github.com/sirupsen/logrus"

	"github.com/TheYeung1/yata-server/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDbYataDatabase struct {
	Dynamo *dynamodb.DynamoDB
}

func (db *DynamoDbYataDatabase) GetList(uid model.UserID, lid model.ListID) (model.YataList, error) {
	queryResults, err := db.Dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("ListTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": &dynamodb.AttributeValue{
				S: aws.String(string(uid)),
			},
			"ListID": &dynamodb.AttributeValue{
				S: aws.String(string(lid)),
			},
		},
	})
	if err != nil {
		log.Println(err)
		return model.YataList{}, err
	}

	if queryResults.Item == nil {
		return model.YataList{}, ListNotFoundError{
			uid: uid,
			lid: lid,
		}
	}

	yl := model.YataList{}
	err = dynamodbattribute.UnmarshalMap(queryResults.Item, &yl)
	if err != nil {
		log.Println(err)
		return model.YataList{}, err
	}
	return yl, nil
}

func (db *DynamoDbYataDatabase) GetLists(uid model.UserID) ([]model.YataList, error) {
	queryResults, err := db.Dynamo.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ListTable"),
		KeyConditionExpression: aws.String("UserID = :user"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": &dynamodb.AttributeValue{
				S: aws.String(string(uid)),
			},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	yl := []model.YataList{}
	err = dynamodbattribute.UnmarshalListOfMaps(queryResults.Items, &yl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return yl, nil
}

func (db *DynamoDbYataDatabase) InsertList(uid model.UserID, yl model.YataList) error {
	av, err := dynamodbattribute.MarshalMap(yl)
	if err != nil {
		log.Errorln(err)
		return err
	}
	_, err = db.Dynamo.PutItem(&dynamodb.PutItemInput{
		TableName:           aws.String("ListTable"),
		ConditionExpression: aws.String("attribute_not_exists(ListID)"),
		Item:                av,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				// TODO: get the item and check that the title matches before throwing a 409
				log.Warnln(aerr)
				return ListExistsError{
					uid: uid,
					lid: yl.ListID,
				}
			default:
				log.Errorln(aerr)
				return err
			}
		}
		log.Errorln(err)
		return err
	}
	return nil
}
