package e2e

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/niemeyer/pretty"
	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/db"
	"github.com/rickschubert/usersgrpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client users.UsersClient

func TestMain(m *testing.M) {
	// In future, we could add credentials
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", config.Port()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client = users.NewUsersClient(conn)

	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func setup() {
	dynamoClient, err := db.New()
	if err != nil {
		log.Fatal("unable to connect to dymamoDB", err)
	}

	createTable(dynamoClient)

	pretty.Println("Seeding Luke Skywalker...")
	_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(config.TableName()),
		Item: map[string]*dynamodb.AttributeValue{
			"Id":        {S: aws.String("34541724-d210-47e0-9cd3-dc950344e421")},
			"FirstName": {S: aws.String("Luke")},
			"LastName":  {S: aws.String("Skywalker")},
			"Nickname":  {S: aws.String("Starkiller")},
			"Password":  {S: aws.String("tw0M00ns")},
			"Email":     {S: aws.String("luke.skywalker@gmail.com")},
			"Country":   {S: aws.String("Tattooine")},
			"CreatedAt": {S: aws.String(time.Now().Format(time.RFC3339))},
		},
	})
	if err != nil {
		log.Fatalf("failed to insert Luke Skywalker into table: %v", err)
	}

	pretty.Println("Seeding Darth Vader...")
	_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(config.TableName()),
		Item: map[string]*dynamodb.AttributeValue{
			"Id":        {S: aws.String("34541724-d210-47e0-9cd3-dc950344e422")},
			"FirstName": {S: aws.String("Darth")},
			"LastName":  {S: aws.String("Vader")},
			"Nickname":  {S: aws.String("Anakin")},
			"Password":  {S: aws.String("iHateSand")},
			"Email":     {S: aws.String("darth.vader@deathstar.com")},
			"Country":   {S: aws.String("Deathstar")},
			"CreatedAt": {S: aws.String(time.Now().Format(time.RFC3339))},
		},
	})
	if err != nil {
		log.Fatalf("failed to insert Darth Vader into table: %v", err)
	}

	pretty.Println("Seeding Storm Trooper 1...")
	_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(config.TableName()),
		Item: map[string]*dynamodb.AttributeValue{
			"Id":        {S: aws.String("34541724-d210-47e0-9cd3-dc950344e473")},
			"FirstName": {S: aws.String("Storm")},
			"LastName":  {S: aws.String("Trooper1")},
			"Nickname":  {S: aws.String("gunfodder1")},
			"Password":  {S: aws.String("iCantShoot")},
			"Email":     {S: aws.String("storm.trooper1@deathstar.com")},
			"Country":   {S: aws.String("Deathstar")},
			"CreatedAt": {S: aws.String(time.Now().Format(time.RFC3339))},
		},
	})
	if err != nil {
		log.Fatalf("failed to insert Storm Trooper 1 into table: %v", err)
	}

	pretty.Println("Seeding Storm Trooper 2...")
	_, err = dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(config.TableName()),
		Item: map[string]*dynamodb.AttributeValue{
			"Id":        {S: aws.String("34541724-d210-47e0-9cd3-dc950344e444")},
			"FirstName": {S: aws.String("Storm")},
			"LastName":  {S: aws.String("Trooper2")},
			"Nickname":  {S: aws.String("gunfodder2")},
			"Password":  {S: aws.String("iCantShoot")},
			"Email":     {S: aws.String("storm.trooper2@deathstar.com")},
			"Country":   {S: aws.String("Deathstar")},
			"CreatedAt": {S: aws.String(time.Now().Format(time.RFC3339))},
		},
	})
	if err != nil {
		log.Fatalf("failed to insert Storm Trooper 2 into table: %v", err)
	}
}

func teardown() {
	dynamoClient, err := db.New()
	if err != nil {
		log.Fatal("unable to connect to dymamoDB", err)
	}

	_, err = dynamoClient.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(config.TableName()),
	})
	if err != nil {
		log.Fatalf("failed to delete '%s' table", config.TableName())
	}

	pretty.Println(fmt.Sprintf("Table '%s' successfully deleted", config.TableName()))
}

func createTable(dynamoClient *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(config.TableName()),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := dynamoClient.CreateTable(input)
	if err != nil && !strings.Contains(err.Error(), "Table already exists") {
		log.Fatalf("failed to create table: %v", err)
	}

	pretty.Println(fmt.Sprintf("Table '%s' created successfully", config.TableName()))
}
