package pideaapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Privacyidea API: response of '/user?realm=***'
type UsersResponse struct {
	Result struct {
		Value []User `json:"value"`
	} `json:"result"`
}

// Privacyidea API: response of '/user?realm=', element of Value list
// excluded fields are "editable" & "userid" as no valuable info
type User struct {
	// Editable  bool     `json:"editable"`
	Email     string   `json:"email"`
	Givenname string   `json:"givenname"`
	MemberOf  []string `json:"memberOf"`
	Mobile    string   `json:"mobile"`
	Phone     string   `json:"phone"`
	Resolver  string   `json:"resolver"`
	Surname   string   `json:"surname"`
	// Userid    string   `json:"userid"`
	Username string `json:"username"`
}

// Privacyidea API: response of '/auth?username=***&password=***'
type APITokenResponse struct {
	Result struct {
		Value struct {
			Token string `json:"token"`
		} `json:"value"`
	} `json:"result"`
}

// Privacyidea API: response of '/token?realm=***&user=***'
type TokenResponse struct {
	Result struct {
		Value struct {
			Tokens []struct {
				Serial string `json:"serial"`
			} `json:"tokens"`
		} `json:"value"`
	} `json:"result"`
}

// Privacyidea API: response of '/validate/check?user=***&realm=***&otponly=1&pass=***&serial=***'
type ValidateResponse struct {
	Result struct {
		Authentication string `json:"authentication"`
	} `json:"result"`
}

// PrivacyIdea API request: get Pidea API Token for given user(POST)
func GetApiToken(httpClient *http.Client, pideaUrl, apiUser, apiUserPassword string) (string, error) {
	var tokenData APITokenResponse

	// form URL query
	query := fmt.Sprintf("%s/auth?username=%s&password=%s", pideaUrl, apiUser, apiUserPassword)

	// form request
	request, err := http.NewRequest(http.MethodPost, query, nil)
	if err != nil {
		return "", fmt.Errorf("failed to Form getToken POST request,\n\t%v", err)
	}

	// do request
	response, err := httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to Make getToken POST request,\n\t%v", err)
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return "", fmt.Errorf("wrong credentials for getToken POST request StatusCode,\n\t%s", response.Status)
		}
		return "", fmt.Errorf("check getToken POST request StatusCode,\n\t%s", response.Status)
	}

	// decode json response
	json.NewDecoder(response.Body).Decode(&tokenData)

	// check if empty result
	if token := tokenData.Result.Value.Token; len(token) != 0 {
		return token, nil
	}

	return "", fmt.Errorf("token result of getToken POST request is Empty\n\t%+v", tokenData)
}

// PrivacyIdea API request: validate check using realm, user(POST)
func ValidateCheck(httpClient *http.Client, authToken, pideaUrl, realm, userName, serial, otp string) (bool, error) {
	var validateResponse ValidateResponse

	// form URL query
	triggerUrl := fmt.Sprintf(
		"%s/validate/check?user=%s&realm=%s&otponly=1&pass=%s&serial=%s",
		pideaUrl,
		userName,
		realm,
		otp,
		serial,
	)

	// form request
	req, err := http.NewRequest(http.MethodPost, triggerUrl, nil)
	if err != nil {
		return false, fmt.Errorf("failed to form request: %v", err)
	}

	// set auth header(api key)
	req.Header.Set("Authorization", authToken)

	// do request
	response, err := httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != 200 {
		return false, fmt.Errorf("wrong status code: %v", response.Status)
	}

	// unmarshall json response
	json.NewDecoder(response.Body).Decode(&validateResponse)

	// check validation result is ok(ACCEPT)
	if validateResponse.Result.Authentication == "ACCEPT" {
		return true, nil
	}

	return false, nil
}

// PrivacyIdea API request: get users't token serial using using user(GET)
func GetUserTokenSerial(httpClient *http.Client, authToken, pideaUrl, realm, userName string) (string, error) {
	var tokenResponse TokenResponse

	// form URL query
	reqUrl := fmt.Sprintf("%s/token?realm=%s&user=%s", pideaUrl, realm, userName)
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to form request,\n\t%v", err)
	}

	// set auth header(api key)
	req.Header.Set("Authorization", authToken)

	// do request
	response, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request,\n\t%v", err)
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != 200 {
		return "", fmt.Errorf("bad status code to make request,\n\t%s", response.Status)
	}

	// unmarshall json response
	json.NewDecoder(response.Body).Decode(&tokenResponse)

	// check if empty result
	if len(tokenResponse.Result.Value.Tokens[0].Serial) == 0 {
		return "", fmt.Errorf("empty result for token")
	}

	return tokenResponse.Result.Value.Tokens[0].Serial, nil
}

// PrivacyIdea API request: get list of Pidea users of given realm(GET)
func GetPideaUsersByRealm(httpClient *http.Client, authToken, pideaUrl, realm string) ([]User, error) {
	var usersData UsersResponse

	// form URL query
	query := fmt.Sprintf("%s/user?realm=%s", pideaUrl, realm)

	// form request
	request, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to Form getUsersByRealm GET request,\n\t%v", err)
	}

	// set auth header(api key)
	request.Header.Set("Authorization", authToken)

	// do request
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to Make getUsersByRealm GET request,\n\t%v", err)
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return nil, fmt.Errorf("auth failure for getUsersByRealm GET request StatusCode,\n\t%s", response.Status)
		}
		return nil, fmt.Errorf("check getUsersByRealm GET request StatusCode,\n\t%s", response.Status)
	}

	// decode json response
	json.NewDecoder(response.Body).Decode(&usersData)

	// check if empty result
	if users := usersData.Result.Value; len(users) != 0 {
		return users, nil
	}

	return nil, fmt.Errorf("token result of usersByRealm GET request is Empty\n\t%+v", usersData)
}
