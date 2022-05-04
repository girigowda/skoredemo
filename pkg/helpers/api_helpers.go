package helpers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Generatevidatoken() string {

	apiUrl := "https://qa-sso.vida.id"
	resource := "/auth/realms/vida/protocol/openid-connect/token"
	data := url.Values{}

	data.Set("grant_type", GetEnvVariable("GRANT_TYPE"))
	data.Set("scope", GetEnvVariable("SCOPE"))
	data.Set("client_id", GetEnvVariable("CLIENT_ID"))
	data.Set("client_secret", GetEnvVariable("CLIENT_SECRET"))

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(r)

	body, _ := ioutil.ReadAll(resp.Body)

	var resp_data map[string]string
	json.Unmarshal([]byte(body), &resp_data)

	return resp_data["access_token"]
}
