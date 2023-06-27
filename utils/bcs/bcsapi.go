package bcs

import (
	"context"
	"crypto/tls"
	"fmt"
	"kconsole/config"
	"net/http"

	"github.com/carlmjohnson/requests"
)

type BCSBaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BcsProject struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectID"`
}

type BCSProjectsResponse struct {
	BCSBaseResponse
	Data struct {
		Total   int          `json:"total"`
		Results []BcsProject `json:"results"`
	}
}

type BcsCluster struct {
	ClusterID   string `json:"clusterID"`
	ClusterName string `json:"clusterName"`
}

type BCSClustersResponse struct {
	BCSBaseResponse
	Data []BcsCluster `json:"data"`
}

var InsecureHttpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func UserBCSProjects(ctx context.Context) (resp BCSProjectsResponse, err error) {
	err = requests.URL(fmt.Sprintf("%s/bcsapi/v4/bcsproject/v1/authorized_projects", bcsHost())).
		Bearer(userToken()).Client(&InsecureHttpClient).ToJSON(&resp).Fetch(ctx)
	if err != nil {
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("Unexpected response:%v", resp)
	}
	return
}

func UserBCSCluster(ctx context.Context, projectid string) (resp BCSClustersResponse, err error) {
	err = requests.URL(fmt.Sprintf("%s/bcsapi/v4/clustermanager/v1/cluster?projectID=%s", bcsHost(), projectid)).
		Bearer(userToken()).Client(&InsecureHttpClient).ToJSON(&resp).Fetch(ctx)
	if err != nil {
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("Unexpected response:%v", resp)
	}
	return
}

func userToken() string {
	c := config.GetKconsoleConfig()
	return c.BCSToken
}

func bcsHost() string {
	c := config.GetKconsoleConfig()
	return c.BCSHost
}
