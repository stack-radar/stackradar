package models

import (
	"fmt"
	"strings"
)

// Language represents programming language information
type Language struct {
	Name       string `json:"name" yaml:"name"`
	Version    string `json:"version" yaml:"version"`
	BuildTool  string `json:"build_tool" yaml:"build_tool"`
	CIImageTag string `json:"ci_image_tag" yaml:"ci_image_tag"`
}

// TechStack represents the complete tech stack information
type TechStack struct {
	Language Language `json:"language" yaml:"language"`
}

// ToEnv converts TechStack to environment variable format
func (ts *TechStack) ToEnv() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("LANGUAGE_NAME=%s\n", ts.Language.Name))
	sb.WriteString(fmt.Sprintf("LANGUAGE_VERSION=%s\n", ts.Language.Version))
	sb.WriteString(fmt.Sprintf("BUILD_TOOL=%s\n", ts.Language.BuildTool))
	sb.WriteString(fmt.Sprintf("CI_IMAGE_TAG=%s\n", ts.Language.CIImageTag))
	return sb.String()
}
