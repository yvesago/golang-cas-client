package client

import (
	"github.com/yvesago/golang-cas-client/util"
	"net/http"
)

type CasClientConfig struct {
	Server, Username, Password string
}

func New(server, username, password string) CasClientConfig {
	return CasClientConfig{server, username, password}
}

func (self CasClientConfig) RequestLoginPage(service string, ip string) (*http.Client, string, error) {
	paramsauth := map[string]string{"username": self.Username, "password": self.Password}
	params := map[string]string{"service": service}
	return util.GetResponseForm(self.Server, params, paramsauth, ip)
}

func (self CasClientConfig) RequestSt(client *http.Client, service string) (string, error) {
	params := map[string]string{"service": service}
	return util.GetSt(self.Server+"/login", client, params)
}

func (self CasClientConfig) RequestServiceTicket(service string) (string, error) {
	tgt, err := self.requestTgtLocation()
	if err != nil {
		return "", err
	}

	return self.getServiceTicket(tgt, service)
}

func (self CasClientConfig) requestTgtLocation() (string, error) {
	params := map[string]string{"username": self.Username, "password": self.Password}
	return util.GetResponseHeader(self.Server+"/v1/tickets", "Location", params)
}

func (self CasClientConfig) getServiceTicket(tgt, service string) (string, error) {
	params := map[string]string{"service": service}
	body, err := util.GetResponseBody(tgt, params)

	if err != nil {
		return "", err
	}

	return body, nil
}
