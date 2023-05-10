package wss

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func NewUpdateRequestFromFile(filepath string) UpdateRequestOriginal {
	var updateRequestOrigin UpdateRequestOriginal
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(data, &updateRequestOrigin)
	if err != nil {
		panic(err)
	}
	return updateRequestOrigin
}

func (u UpdateRequestOriginal) GetValues() url.Values {
	diff_data, err := json.Marshal(u.Diff)
	if err != nil {
		panic(err)
	}
	values := url.Values{}
	values.Set("updateType", u.UpdateType)
	values.Set("type", u.Type)
	values.Set("agent", u.Agent)
	values.Set("agentVersion", u.AgentVersion)
	values.Set("token", u.Token)
	values.Set("userKey", u.UserKey)
	values.Set("timeStamp", strconv.Itoa(u.TimeStamp))
	values.Set("product", u.Product)
	values.Set("diff", string(diff_data))

	return values
}

func (u *UpdateRequestOriginal) LoadUpdateRequest(filepath string) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &u)
	if err != nil {
		panic(err)
	}
}

func (u UpdateRequestOriginal) SendUploadRequest(wssurl string) (resp *http.Response, err error) {

	vals := u.GetValues()
	qs, _ := url.ParseQuery(vals.Encode())

	resp, err = http.PostForm(wssurl, qs)

	return resp, err
}

func (u *UpdateRequestOriginal) FromFile(fromfile string) bool {
	var status bool = true

	bytes, err := ioutil.ReadFile(fromfile)
	if err != nil {
		status = false
	}

	err = json.Unmarshal(bytes, &u)
	if err != nil {
		status = false
	}

	return status
}

func (ud UploadResponseData) GetProjectName() string {
	var projectName string = ""

	for k := range ud.ProjectNamesToDetails {
		projectName = k
	}

	return projectName
}

func (us UploadResponseStatus) GetJson() []byte {
	data, err := json.Marshal(us)

	if err != nil {
		panic(err)
	}

	return data
}

func (ud UploadResponseData) GetJson() []byte {
	data, err := json.Marshal(ud)

	if err != nil {
		panic(err)
	}

	return data
}

func (us UploadResponseStatus) ToFile(destination string) bool {
	var status bool = true

	file, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		status = false
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	data := us.GetJson()
	_, err = writer.Write(data)
	if err != nil {
		status = false
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
	}

	return status
}

func (ud UploadResponseData) ToFile(destination string) bool {
	var status bool = true

	file, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, os.FileMode(0644))

	if err != nil {
		status = false
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	data := ud.GetJson()
	_, err = writer.Write(data)
	if err != nil {
		status = false
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
	return status
}

func (ud *UploadResponseData) FromFile(fromfile string) bool {
	var status bool = true

	bytes, err := ioutil.ReadFile(fromfile)
	if err != nil {
		status = false
		panic(err)
	}
	err = json.Unmarshal(bytes, &ud)
	if err != nil {
		status = false

		panic(err)
	}

	return status

}

func (us *UploadResponseStatus) FromFile(fromfile string) bool {
	var status bool = true

	bytes, err := ioutil.ReadFile(fromfile)
	if err != nil {
		status = false
		panic(err)
	}

	err = json.Unmarshal(bytes, &us)
	if err != nil {
		status = false
		panic(err)
	}

	return status
}
