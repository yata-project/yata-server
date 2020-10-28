package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/gorilla/mux"
)

type server struct {
	dynamo *dynamodb.DynamoDB
}

type YataList struct {
	UserID string
	ListID string
	Title  string
}

type InsertListInput struct {
	ListID string
	Title  string
}

type InsertListOutput struct {
	ListID string
}

func writeInternalErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Sorry! Something went wrong"))
}

func (s *server) GetLists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetLists")

	userID := r.Header.Get("User")
	// TODO: use QueryPages to go over paginated results
	queryResults, err := s.dynamo.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ListTable"),
		KeyConditionExpression: aws.String("UserID = :user"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": &dynamodb.AttributeValue{
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
	}

	yl := []YataList{}
	err = dynamodbattribute.UnmarshalListOfMaps(queryResults.Items, &yl)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
	}

	res, err := json.Marshal(yl)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
	}

	_, err = w.Write(res)
	if err != nil {
		log.Println(err)
	}
}

func (s *server) InsertList(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
		return
	}

	var in InsertListInput
	err = json.Unmarshal(b, &in)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
		return
	}

	uid, ok := r.Header["User"]
	if !ok {
		log.Println("UserId not found")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("UserId is missing"))
		return
	}

	// TODO: assert input lengths
	yl := YataList{
		UserID: uid[0],
		ListID: in.ListID,
		Title:  in.Title,
	}

	av, err := dynamodbattribute.MarshalMap(yl)
	if err != nil {
		log.Println(err)
		writeInternalErrorResponse(w)
		return
	}
	_, err = s.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName:           aws.String("ListTable"),
		ConditionExpression: aws.String("attribute_not_exists(ListID)"),
		Item:                av,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println("Error: %+v", err)
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				// TODO: get the item and check that the title matches before throwing a 409
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("List already exists"))
				return
			default:
				writeInternalErrorResponse(w)
				return
			}
		}
		log.Println(err)
		writeInternalErrorResponse(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte{})
	if err != nil {
		log.Println(err)
	}
}

func (s *server) start() {
	r := mux.NewRouter()
	r.HandleFunc("/lists", s.GetLists).Methods(http.MethodGet)
	r.HandleFunc("/lists", s.InsertList).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8888", r))
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "yata"),
	})

	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s := server{
		dynamo: dynamodb.New(sess),
	}
	s.start()
}
