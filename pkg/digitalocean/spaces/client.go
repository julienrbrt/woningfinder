package spaces

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"github.com/woningfinder/woningfinder/pkg/networking"
)

type Client interface {
	UploadPicture(name string, pictureURL *url.URL) (string, error)
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

func (c *client) UploadPicture(name string, pictureURL *url.URL) (string, error) {
	fileName := c.buildFileName(name)
	if c.doesPictureExists(fileName) {
		return fileName, nil
	}

	// download offer
	pictureRaw, pictureSize, err := c.downloadPicture(pictureURL)
	if err != nil {
		return "", err
	}

	// upload the image to our bucket
	_, err = c.client.PutObject(context.Background(), c.bucketName, fileName, pictureRaw, pictureSize, minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		return "", fmt.Errorf("failed uploading %s image to digitalocean spaces: %w", name, err)
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

func (c *client) buildFileName(name string) string {
	return "offers/" + strings.ToLower(base64.StdEncoding.EncodeToString([]byte(name)))
}

func (c *client) doesPictureExists(fileName string) bool {
	statObject, err := c.client.StatObject(context.Background(), c.bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return false
	}

	return statObject.ContentType != ""
}