package database

import (
	"github.com/TheYeung1/yata-server/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

type DynamoDbYataDatabase struct {
	ListsTableName string
	ItemsTableName string
	Dynamo         *dynamodb.DynamoDB
}

func (db *DynamoDbYataDatabase) GetList(uid model.UserID, lid model.ListID) (model.YataList, error) {
	queryResults, err := db.Dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.ListsTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(string(uid)),
			},
			"ListID": {
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
		TableName:              aws.String(db.ListsTableName),
		KeyConditionExpression: aws.String("UserID = :user"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
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
		TableName:           aws.String(db.ListsTableName),
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

func (db *DynamoDbYataDatabase) GetAllItems(uid model.UserID) ([]model.YataItem, error) {
	queryResults, err := db.Dynamo.Query(&dynamodb.QueryInput{
		TableName:              aws.String(db.ItemsTableName),
		KeyConditionExpression: aws.String("UserID = :user"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
				S: aws.String(string(uid)),
			},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	items := []model.YataItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(queryResults.Items, &items)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return items, nil
}

func (db *DynamoDbYataDatabase) GetListItems(uid model.UserID, lid model.ListID) ([]model.YataItem, error) {
	queryResults, err := db.Dynamo.Query(&dynamodb.QueryInput{
		TableName:              aws.String(db.ItemsTableName),
		KeyConditionExpression: aws.String("UserID = :user AND begins_with(#listIDuserID, :list)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
				S: aws.String(string(uid)),
			},
			":list": {
				S: aws.String(string(lid)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#listIDuserID": aws.String("ListID-ItemID"),
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	items := []model.YataItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(queryResults.Items, &items)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return items, nil
}

func (db *DynamoDbYataDatabase) InsertItem(item model.YataItem) error {
	// TODO: make sure that the list exists first

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Errorln(err)
		return err
	}
	av["ListID-ItemID"] = &dynamodb.AttributeValue{
		S: aws.String(string(item.ListID) + ":" + string(item.ItemID)),
	}
	_, err = db.Dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(db.ItemsTableName),
		Item:      av,
	})
	return err
}
