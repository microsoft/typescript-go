package compiler

import (
	"encoding/json"
	"errors"
)

type BuildInfo struct {
	Version string `json:"version"`
}

func GetBuildInfo(buildInfoFile string, buildInfoText string) (*BuildInfo, error) {
	if buildInfoText == "" {
		return nil, errors.New("empty buildInfoText")
	}
	var buildInfo BuildInfo
	if err := json.Unmarshal([]byte(buildInfoText), &buildInfo); err != nil {
		return nil, err
	}
	return &buildInfo, nil
}

func GetBuildInfoText(buildInfo BuildInfo) (string, error) {
	data, err := json.MarshalIndent(buildInfo, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
