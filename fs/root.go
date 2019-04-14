package fs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/justasbieliauskas/drivemv/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Root represents files and folders in root folder of google drive.
type Root struct {
	service *drive.FilesService
}

// File represents file in google drive.
type File struct {
	*drive.File
	service *drive.FilesService
}

// NewRoot connects to google drive using given environment variables.
func NewRoot(env []string) (*Root, error) {
	creds := credsFromEnv(env)
	json, err := creds.json()
	if err != nil {
		return nil, errors.Nest("Unable to retrieve config json from creds", err)
	}
	config, err := google.ConfigFromJSON(json, drive.DriveScope)
	if err != nil {
		return nil, errors.Nest("Unable to create config from json credentials", err)
	}
	token, err := creds.token()
	if err != nil {
		return nil, errors.Nest("Unable to create oauth token from credentials", err)
	}
	client := config.Client(context.Background(), token)
	service, err := drive.New(client)
	if err != nil {
		return nil, errors.Nest("Unable to new-up a Drive service", err)
	}
	return &Root{service: service.Files}, nil
}

// List returns sample files to make sure the service is working.
func (root *Root) List() (*drive.FileList, error) {
	return root.service.List().Q("name = 'verslumas'").PageSize(1).Fields("files(name)").Do()
}

// UploadFile uploads given file to root folder of drive using the given name.
// Returns ID of uploaded file.
func (root *Root) UploadFile(file *os.File, name string) (*File, error) {
	metadata := &drive.File{Name: name}
	uploadedFile, err := root.service.Create(metadata).Fields("id").Media(file).Do()
	if err != nil {
		return nil, err
	}
	return &File{uploadedFile, root.service}, nil
}

// GetFileByName gets a File by its name.
func (root *Root) GetFileByName(name string) (*File, error) {
	list, err := root.service.List().Q(fmt.Sprintf("name = '%s'", name)).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable fetch files list under name '%s': %v", name, err)
	}
	if len(list.Files) == 0 {
		return nil, fmt.Errorf("Unable to find files with name '%s'", name)
	}
	return &File{list.Files[0], root.service}, nil
}

// Content downloads content of a file.
func (file *File) Content() (string, error) {
	response, err := file.service.Get(file.Id).Download()
	if err != nil {
		return "", errors.Nest("Unable to download file", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Nest("Error while reading http response body", err)
	}
	return string(body), nil
}

// Delete removes file from drive.
func (file *File) Delete() error {
	return file.service.Delete(file.Id).Do()
}
