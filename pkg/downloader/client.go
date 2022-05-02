package downloader

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/julienrbrt/woningfinder/pkg/logging"
	"github.com/julienrbrt/woningfinder/pkg/networking"
	"go.uber.org/zap"
)

type Client interface {
	Download(orignalName string, fileURL *url.URL) (string, error)
	Get(orignalName string) ([]byte, error)
}

type client struct {
	logger           *logging.Logger
	networkingClient networking.Client
	downloadPath     string
}

func NewClient(logger *logging.Logger, networkingClient networking.Client, path string) (Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	return &client{
		logger:           logger,
		networkingClient: networkingClient,
		downloadPath:     fmt.Sprintf("%s/woningfinder/%s", homeDir, path),
	}, nil
}

// Download a given file to disk
// orignalName is the orignal file name
// fileURL is the path to the file to download
func (c *client) Download(originalName string, fileURL *url.URL) (string, error) {
	fileName := c.buildPath(originalName)
	if c.exists(fileName) {
		return fileName, nil
	}

	// download file
	response, err := c.networkingClient.Send(&networking.Request{Host: fileURL})
	if err != nil {
		return "", fmt.Errorf("error while downloading %s: %w", fileURL.String(), err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("error while downloading %s: response status code %d", fileURL.String(), response.StatusCode)
	}

	// create an empty file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// write the bytes to the file
	if _, err = io.Copy(file, response.RawResponse.Body); err != nil {
		return "", err
	}

	return fileName, nil
}

// Get a file from disk
func (c *client) Get(originalName string) ([]byte, error) {
	fileName := c.buildPath(originalName)
	if !c.exists(fileName) {
		return nil, errors.New("file not found")
	}

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *client) buildPath(originalName string) string {
	return fmt.Sprintf("%s/%s", c.downloadPath, strings.ToLower(hex.EncodeToString([]byte(originalName))))
}

func (c *client) exists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		c.logger.Error("error while checking if file exists", zap.String("filename", fileName), zap.Error(err))
		return false
	}
}
