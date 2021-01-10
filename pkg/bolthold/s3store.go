package bolthold

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	bolt "go.etcd.io/bbolt"

	"github.com/ubogdan/network-manager-api/service"
)

// S3Client an interface to implement PutObject and GetObject from s3
type S3Client interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// Db client
type Db struct {
	S3Client S3Client
	db       *bolt.DB
	dbPath   string
	dbName   string
	S3bucket string
	S3prefix string
	Log      service.Logger
}

// WithS3 creates an empty s3bolt wrapper
func WithS3(S3Client S3Client, bucket, prefix string) *Db {
	return &Db{
		S3bucket: bucket,
		S3prefix: prefix,
		S3Client: S3Client,
	}
}

// Open wrapper
func (s *Db) Open(path string, mode os.FileMode, options *bolt.Options) (*Store, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	s.dbName = filepath.Base(absPath)
	s.dbPath = filepath.Dir(absPath)

	err = s.load()
	if err != nil {
		log.Fatalf("Can't load %s database from s3: %s", s.dbName, err)
	}

	db, err := bolt.Open(path, mode, options)
	if err != nil {
		return nil, err
	}
	s.db = db

	return &Store{
		db:     s,
		encode: DefaultEncode,
		decode: DefaultDecode,
	}, nil
}

// Close wrapper
func (s *Db) Close() error {
	return s.db.Close()
}

// Update wrapper
func (s *Db) Update(fn func(*bolt.Tx) error) error {
	err := s.db.Update(fn)
	if err != nil {
		return err
	}

	err = s.store()
	if err != nil && s.Log != nil {
		s.Log.Warnf("Unable to store s3bolt db into s3: %s", err)
	}
	return nil
}

// Batch wrapper
func (s *Db) Batch(fn func(*bolt.Tx) error) error {
	err := s.db.Batch(fn)
	if err != nil {
		return err
	}
	err = s.store()
	if err != nil && s.Log != nil {
		s.Log.Warnf("Unable to store s3bolt db into s3: %s", err)
	}
	return nil
}

// View wrapper
func (s *Db) View(fn func(*bolt.Tx) error) error {
	return s.db.Batch(fn)
}

func (s *Db) store() error {
	backup := &bytes.Buffer{}
	err := s.db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(backup)
		return err
	})

	body, _ := ioutil.ReadAll(backup)
	if err != nil {
		return err
	}

	_, err = s.S3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.S3bucket),
		Key:    aws.String(s.s3FilePath()),
		Body:   bytes.NewReader(body),
	})
	return err
}

func (s *Db) load() error {
	getObjectOutput, err := s.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.S3bucket),
		Key:    aws.String(s.s3FilePath()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return err
			case s3.ErrCodeNoSuchKey:
				return nil
			}
		}
		return err
	}
	defer getObjectOutput.Body.Close()

	content, err := ioutil.ReadAll(getObjectOutput.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.dbFilePath(), content, 0644)
}

func (s *Db) dbFilePath() string {
	return fmt.Sprintf("%s/%s", s.dbPath, s.dbName)
}

func (s *Db) s3FilePath() string {
	return fmt.Sprintf("%s/%s", s.S3prefix, s.dbName)
}
