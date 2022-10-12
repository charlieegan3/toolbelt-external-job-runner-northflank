package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs/v2"

	"github.com/charlieegan3/toolbelt/pkg/apis"
)

// Northflank is an external job runner for northflank
type Northflank struct {
	// APIToken is the token used to authenticate with the Northflank API
	APIToken string
}

func (n *Northflank) Configure(config map[string]any) error {
	tokenPath := "token"

	var ok bool
	n.APIToken, ok = config[tokenPath].(string)
	if !ok {
		return fmt.Errorf("missing required config path: %s", tokenPath)
	}

	return nil
}

func (n *Northflank) RunJob(job apis.ExternalJob) error {
	cfg := gabs.Wrap(job.Config())

	var path string
	path = "job_id"
	jobID, ok := cfg.Path(path).Data().(string)
	if !ok {
		return fmt.Errorf("missing required config path: %s", path)
	}

	path = "project_id"
	projectID, ok := cfg.Path(path).Data().(string)
	if !ok {
		return fmt.Errorf("missing required config path: %s", path)
	}

	p := createJobPayload{}

	path = "command"
	p.Deployment.CMDOverride, ok = cfg.Path(path).Data().(string)
	if !ok {
		return fmt.Errorf("missing required config path: %s", path)
	}

	path = "env"
	env, ok := cfg.Path(path).Data().(map[string]any)
	if !ok {
		return fmt.Errorf("missing required config path: %s", path)
	}
	p.RuntimeEnvironment = make(map[string]string)
	for k, v := range env {
		p.RuntimeEnvironment[strings.ToUpper(k)] = v.(string)
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.northflank.com/v1/projects/%s/jobs/%s/runs", projectID, jobID),
		bytes.NewBuffer(body),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.APIToken))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create job instance: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to create job instance: %s", resp.Status)
	}

	return nil
}

type createJobPayload struct {
	Deployment struct {
		CMDOverride        string `json:"cmdOverride,omitempty"`
		EntrypointOverride string `json:"entrypointOverride,omitempty"`
	} `json:"deployment,omitempty"`
	RuntimeEnvironment map[string]string `json:"runtimeEnvironment,omitempty"`
}
