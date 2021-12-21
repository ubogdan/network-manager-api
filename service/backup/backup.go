package backup

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type backup struct {
	BucketName string
}

// New returns license service implementation.
func New(bucketName string) *backup {
	return &backup{
		BucketName: bucketName,
	}
}

func (svc *backup) GetUploadLink(hardwareID, fileName string, validFor time.Duration) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	h := md5.New() //nolint:gosec
	h.Write([]byte("network-manager:" + svc.BucketName + ":backup-service:" + hardwareID))

	req, _ := s3.New(sess).PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(svc.BucketName),
		Key:    aws.String(fmt.Sprintf("%x/%s", h.Sum(nil), fileName)),
	})

	return req.Presign(validFor)
}
