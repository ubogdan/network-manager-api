package dynamo

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
)

type license struct {
	db *dynamodb.DynamoDB
}

const licenseTable = "nm-licenses"

var _ repository.License = License(nil)

// License return a license repository.
func License(database *dynamodb.DynamoDB) *license {
	return &license{
		db: database,
	}
}

// FindAll returns a list of licenses.
func (s *license) FindAll() ([]model.License, error) {
	out, err := s.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(licenseTable),
	})
	if err != nil {
		return nil, err
	}

	licenses := make([]model.License, len(out.Items))
	for idx, item := range out.Items {
		var license model.License
		err = dynamodbattribute.UnmarshalMap(item, &license)
		if err != nil {
			return nil, err
		}

		licenses[idx] = license
	}

	return licenses, nil
}

// Find returns a license by id (hradwareID).
func (s *license) Find(hardwareID string) (*model.License, error) {
	result, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(licenseTable),
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
		TableName: aws.String(licenseTable),
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
				S: aws.String(license.Serial),
			},
			":issued": {
				N: aws.String(strconv.FormatInt(license.LastIssued, 10)),
			},
			":expire": {
				N: aws.String(strconv.FormatInt(license.Expire, 10)),
			},
			":customer": {
				M: map[string]*dynamodb.AttributeValue{
					"Name": {
						S: aws.String(license.Customer.Name),
					},
					"Country": {
						S: aws.String(license.Customer.Country),
					},
					"City": {
						S: aws.String(license.Customer.City),
					},
					"Organization": {
						S: aws.String(license.Customer.Organization),
					},
					"OrganizationalUnit": {
						S: aws.String(license.Customer.OrganizationalUnit),
					},
				},
			},
		},
		TableName: aws.String(licenseTable),
		Key: map[string]*dynamodb.AttributeValue{
			"HardwareID": {
				S: aws.String(license.HardwareID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Serial=:serial, LastIssued=:issued, Expire=:expire, Customer=:customer"),
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
		TableName: aws.String(licenseTable),
	})
	if err != nil {
		return err
	}

	return nil
}
