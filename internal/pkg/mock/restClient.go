package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const unsupportedMockRequest = "unsupported mock request"

func restClient() *http.Client {
	return &http.Client{Transport: RoundTripperFunc(discordAPIResponse)}
}

func discordAPIResponse(r *http.Request) (*http.Response, error) {
	pathTokens := strings.Split(r.URL.Path, "/")
	switch pathTokens[3] {
	case "users":
		return usersResponse(r), nil
	case "members":
		return membersResponse(r), nil
	case "roles":
		return rolesResponse(r), nil
	case "channels":
		return channelsResponse(r), nil
	case "guilds":
		return guildsResponse(r), nil
	}

	return nil, fmt.Errorf(unsupportedMockRequest)
}

func usersResponse(r *http.Request) *http.Response {
	pathTokens := strings.Split(r.URL.Path, "/")
	userID := pathTokens[len(pathTokens)-1]
	userName := Test_UserName

	respBody, err := json.Marshal(mockUser(userID, userName))
	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func membersResponse(r *http.Request) *http.Response {
	pathTokens := strings.Split(r.URL.Path, "/")
	userID := pathTokens[len(pathTokens)-1]
	userName := Test_UserName
	guildID := pathTokens[len(pathTokens)-2]

	var (
		respBody []byte
		err      error
	)

	if userID == "members" {
		if guildID == Test_GuildID {
			respBody, err = json.Marshal(mockMembers(guildID))
		}
	} else {
		respBody, err = json.Marshal(mockMember(guildID, userID, userName))
	}

	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func rolesResponse(r *http.Request) *http.Response {
	switch r.Method {
	case http.MethodGet:
		respBody, err := json.Marshal(mockRoles())
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		return newResponse(http.StatusOK, respBody)
	case http.MethodPost:
		respBody, err := json.Marshal(mockRole(Test_RoleID, Test_RoleName))
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		return newResponse(http.StatusOK, respBody)
	case http.MethodPatch:
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		err = r.Body.Close()
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		return newResponse(http.StatusOK, reqBody)
	case http.MethodDelete:
		return newResponse(http.StatusOK, nil)
	}

	return newResponse(http.StatusMethodNotAllowed, []byte{})
}

func channelsResponse(r *http.Request) *http.Response {
	var (
		respBody []byte
		err      error
	)

	if strings.Contains(r.URL.Path, "guilds") {
		respBody, err = json.Marshal(mockChannels(Test_GuildID))
	} else {
		respBody, err = json.Marshal(mockChannel(Test_GuildID, Test_ChannelID, Test_ChannelName))
	}

	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func guildsResponse(r *http.Request) *http.Response {
	pathTokens := strings.Split(r.URL.Path, "/")
	guildID := pathTokens[len(pathTokens)-1]
	guildName := Test_GuildName

	respBody, err := json.Marshal(mockGuild(guildID, guildName))
	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func newResponse(status int, respBody []byte) *http.Response {
	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}
}
