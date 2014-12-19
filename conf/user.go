// Copyright (c) 2014, B3log
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// LatestSessionContent represents the latest session content.
type LatestSessionContent struct {
	FileTree    []string // paths of expanding nodes of file tree
	Files       []string // paths of files of opening editor tabs
	CurrentFile string   // path of file of the current focused editor tab
}

// User configuration.
type User struct {
	Name                 string
	Password             string
	Email                string
	Gravatar             string // see http://gravatar.com
	Workspace            string // the GOPATH of this user
	Locale               string
	GoFormat             string
	FontFamily           string
	FontSize             string
	Theme                string
	Editor               *editor
	LatestSessionContent *LatestSessionContent
}

// Editor configuration of a user.
type editor struct {
	FontFamily string
	FontSize   string
	LineHeight string
	Theme      string
	TabSize    string
}

// NewUser creates a user with the specified username, password, email and workspace.
func NewUser(username, password, email, workspace string) *User {
	hash := md5.New()
	hash.Write([]byte(email))
	gravatar := hex.EncodeToString(hash.Sum(nil))

	return &User{Name: username, Password: password, Email: email, Gravatar: gravatar, Workspace: workspace,
		Locale: Wide.Locale, GoFormat: "gofmt", FontFamily: "Helvetica", FontSize: "13px", Theme: "default",
		Editor: &editor{FontFamily: "Consolas, 'Courier New', monospace", FontSize: "inherit", LineHeight: "17px",
			Theme: "wide", TabSize: "4"}}
}

// Save saves the user's configurations in conf/users/{username}.json.
func (u *User) Save() bool {
	bytes, err := json.MarshalIndent(u, "", "    ")

	if nil != err {
		logger.Error(err)

		return false
	}

	if err = ioutil.WriteFile("conf/users/"+u.Name+".json", bytes, 0644); nil != err {
		logger.Error(err)

		return false
	}

	return true
}

// GetWorkspace gets workspace path of the user.
//
// Compared to the use of Wide.Workspace, this function will be processed as follows:
//  1. Replace {WD} variable with the actual directory path
//  2. Replace ${GOPATH} with enviorment variable GOPATH
//  3. Replace "/" with "\\" (Windows)
func (u *User) GetWorkspace() string {
	w := strings.Replace(u.Workspace, "{WD}", Wide.WD, 1)
	w = strings.Replace(w, "${GOPATH}", os.Getenv("GOPATH"), 1)

	return filepath.FromSlash(w)
}

// GetOwner gets the user the specified path belongs to. Returns "" if not found.
func GetOwner(path string) string {
	for _, user := range Users {
		if strings.HasPrefix(path, user.GetWorkspace()) {
			return user.Name
		}
	}

	return ""
}
