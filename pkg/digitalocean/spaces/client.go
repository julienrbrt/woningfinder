package spaces

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client interface {
	UploadPicture(prefix, orignalName string, pictureURL *url.URL) (string, error)
}

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	client           *minio.Client
	bucketName       string
}

func NewClient(logger *logging.Logger, c networking.Client, endpoint, bucketName, accessKeyID, secretAccessKey string) (Client, error) {
	// initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating digitalocean spaces client: %w", err)
	}

	if ok, err := minioClient.BucketExists(context.Background(), bucketName); !ok || err != nil {
		return nil, fmt.Errorf("failed reaching digitalocean spaces client: %w", err)
	}

	return &client{
		logger:           logger,
		networkingClient: c,
		client:           minioClient,
		bucketName:       bucketName,
	}, nil
}

// UploadPicture uploads a given picture to DigitalOcean Spaces
// prefix is the folder where the picture must be stored
// orignalName is the orignal file name
func (c *client) UploadPicture(prefix, orignalName string, pictureURL *url.URL) (string, error) {
	fileName := c.buildFileName(prefix, orignalName)
	if c.doesPictureExists(fileName) {
		return fileName, nil
	}

	// download offer
	pictureRaw, pictureSize, err := c.downloadPicture(pictureURL)
	if err != nil {
		return "", err
	}

	// set picture public
	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	// upload the picture to our bucket
	_, err = c.client.PutObject(context.Background(), c.bucketName, fileName, pictureRaw, pictureSize, minio.PutObjectOptions{UserMetadata: userMetaData, ContentType: "image/jpeg"})
	if err != nil {
		return "", fmt.Errorf("failed uploading %s image to digitalocean spaces: %w", orignalName, err)
	}

	return fileName, nil
}

func (c *client) downloadPicture(url *url.URL) (io.Reader, int64, error) {
	response, err := c.networkingClient.Send(&networking.Request{
		Host: url,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("error while downloading %s: %w", url.String(), err)
	}

	if response.StatusCode != 200 {
		return nil, 0, fmt.Errorf("error while downloading %s: response status code %d", url.String(), response.StatusCode)
	}

	return response.RawResponse.Body, response.RawResponse.ContentLength, nil
}

func (c *client) buildFileName(prefix, orignalName string) string {
	fileName := strings.ToLower(base64.StdEncoding.EncodeToString([]byte(orignalName)))

	if prefix == "" {
		return fileName
	}

	return fmt.Sprintf("%s/%s", prefix, fileName)
}

func (c *client) doesPictureExists(fileName string) bool {
	statObject, err := c.client.StatObject(context.Background(), c.bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return false
	}

	return statObject.ContentType != ""
}
