package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
)

type license struct {
	db *dynamodb.DynamoDB
}

const tableName = "nm-licenses"

var _ repository.License = License(nil)

// License return a license repository.
func License(database *dynamodb.DynamoDB) *license {
	return &license{
		db: database,
	}
}

// FindAll returns a list of licenses.
func (s *license) FindAll() ([]model.License, error) {
	var licenses []model.License
	return licenses, nil
}

// Find returns a license by id (hradwareID).
func (s *license) Find(hardwareID string) (*model.License, error) {
	result, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"HardwareID": {
				S: aws.String(hardwareID),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var item model.License

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Create a new license record.
func (s *license) Create(license *model.License) error {

	av, err := dynamodbattribute.MarshalMap(license)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = s.db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// Update a license record.
func (s *license) Update(license *model.License) error {
	_, err := s.db.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":serial": {
				S: aws.String("40DC4-739A8-87BF2-13698-3C60E-BLAH"),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"HardwareID": {
				S: aws.String(license.HardwareID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Serial = :serial"),
	})
	if err != nil {
		return err
	}

	return nil
}

// Delete a license record.
func (s *license) Delete(hardwareID string) error {
	_, err := s.db.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"HardwareID": {
				S: aws.String(hardwareID),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	return nil
}
