// Package jsonfile provides an implementation for the repository interfaces.
package jsonfile

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aisbergg/keycli/pkg/core"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
)

// jsonfileSessionRepository implements `core.SessionRepository`. It loads and
// stores session information in a json file. The file is locked for exclusive
// access.
type jsonfileSessionRepository struct {
	lFile *lockedfile.File
	path  string
}

// NewJSONFileSessionRepository initializes a `jsonfileSessionRepository`.
func NewJSONFileSessionRepository() core.SessionRepository {
	return &jsonfileSessionRepository{}
}

// Exists returns true, if the session file exists.
func (js *jsonfileSessionRepository) Exists(name string) (bool, error) {
	path := PathFromName(name)
	exists, err := checkFile(path)
	return exists, err
}

// Open opens the session file for exclusive access.
func (js *jsonfileSessionRepository) Open(name string) error {
	path := PathFromName(name)

	// create parent dir
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return errors.Errorf("cannot create directory for session file '%s': %v", path, err)
	}

	// create and open locked session file
	lFile, err := lockedfile.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return errors.Errorf("cannot open session file '%s': %v", path, err)
	}
	js.lFile = lFile
	js.path = name

	return nil
}

// Close closes the session file, so that other instances can access it.
func (js *jsonfileSessionRepository) Close() error {
	if js.lFile == nil {
		return nil
	}

	err := js.lFile.Close()
	if err != nil {
		return errors.Errorf("cannot close session file '%s': %v", js.path, err)
	}

	return nil
}

// Read reads the content of the session file. Returns `nil` if file doesn't exist.
func (js *jsonfileSessionRepository) Read() (*core.Session, error) {
	if js.lFile == nil {
		panic("session file must be opened before use")
	}

	// read file contents
	rawData, err := ioutil.ReadAll(js.lFile)
	if err != nil {
		return nil, errors.Errorf("cannot read from session file '%s': %v", js.path, err)
	}

	// decode json content
	session := &core.Session{}
	err = json.Unmarshal(rawData, session)
	if err != nil {
		return nil, errors.Errorf("cannot decode session file '%s': %v", js.path, err)
	}
	return session, nil
}

// Write writes session information to a file.
func (js *jsonfileSessionRepository) Write(s *core.Session) error {
	if js.lFile == nil {
		panic("session file must be opened before use")
	}

	if s == nil {
		return nil
	}

	// json encode data
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// write to file
	js.lFile.Truncate(0)
	_, err = js.lFile.Write(jsonData)
	if err != nil {
		return errors.Errorf("cannot write to session file '%s': %v", js.path, err)
	}

	return nil
}

// Remove removes the stored file. Has no effect, if the file doesn't exist.
func (js *jsonfileSessionRepository) Remove(name string) error {
	path := PathFromName(name)
	exists, err := checkFile(path)
	if !exists {
		return nil
	}
	if err != nil {
		return errors.Errorf("cannot remove session file '%s': %v", path, err)
	}
	err = os.Remove(path)
	if err != nil {
		return errors.Errorf("cannot remove session file '%s': %v", path, err)
	}

	return nil
}

// PathFromName creates the file path for a given session name.
func PathFromName(name string) string {
	cacheDir, err := os.UserCacheDir()
	// if chache dir cannot be determined use a temp dir
	if err != nil {
		cacheDir = os.TempDir()
	}
	path := filepath.Join(cacheDir, "keycli", "tokens", name+".json")
	return path
}

// checkFile checks whether a given path exists and if it is a file.
func checkFile(path string) (exists bool, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	fileMode := fileInfo.Mode()
	if !fileMode.IsRegular() {
		return true, errors.New("not a file")
	}
	return true, nil
}
