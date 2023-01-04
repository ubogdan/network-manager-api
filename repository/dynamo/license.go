package dynamo

import (
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	expr "github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/ubogdan/network-manager-api/model"
)

type license struct {
	db *dynamodb.DynamoDB
}

const licenseTable = "nm-licenses"

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

// FindBySerial returns a license by serial.
func (s *license) FindBySerial(serial string) (*model.License, error) {
	e, err := expr.NewBuilder().
		WithFilter(expr.Name("Serial").Equal(expr.Value(serial))).
		WithProjection(expr.NamesList(expr.Name("HardwareID"), expr.Name("Serial"), expr.Name("Features"))).
		Build()
	if err != nil {
		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  e.Names(),
		ExpressionAttributeValues: e.Values(),
		FilterExpression:          e.Filter(),
		ProjectionExpression:      e.Projection(),
		TableName:                 aws.String(licenseTable),
	}

	// Make the DynamoDB Query API call
	result, err := s.db.Scan(params)
	if err != nil {
		return nil, err
	}

	if len(result.Items) != 1 {
		return nil, errors.New("no license found")
	}

	var item model.License

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
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
	update := dynamodb.UpdateItemInput{
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
			":features": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: aws.String(licenseTable),
		Key: map[string]*dynamodb.AttributeValue{
			"HardwareID": {
				S: aws.String(license.HardwareID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Serial=:serial, Features=:features,LastIssued=:issued, Expire=:expire, Customer=:customer"),
	}

	for _, feature := range license.Features {
		featureObject := dynamodb.AttributeValue{
			M: map[string]*dynamodb.AttributeValue{
				"Name": {
					S: aws.String(string(feature.Name)),
				},
			},
		}

		if feature.Limit > 0 {
			featureObject.M["Limit"] = &dynamodb.AttributeValue{
				N: aws.String(strconv.FormatInt(feature.Limit, 10)),
			}
		}

		if feature.Expire > 0 {
			featureObject.M["Expire"] = &dynamodb.AttributeValue{
				N: aws.String(strconv.FormatInt(feature.Expire, 10)),
			}
		}

		update.ExpressionAttributeValues[":features"].L = append(update.ExpressionAttributeValues[":features"].L, &featureObject)
	}

	_, err := s.db.UpdateItem(&update)
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
