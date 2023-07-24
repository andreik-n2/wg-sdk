package codegeneration

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	// Legacy entrypoint, used by default
	configEntryPointFilename = "wundergraph.config.ts"
	// Entrypoint for applications that use `export default defineConfig` inside generatedDirectory
	applicationEntryPointFilename = "wundergraph.application.ts"

	serverConfigFilename = "wundergraph.server.ts"

	// generatedDirectory is the relative path to the directory with generated
	// from $WUNDERGRAPH_DIR
	generatedDirectory = "generated"
)

var (
	wunderGraphFactoryTemplate = template.Must(template.New(applicationEntryPointFilename).Parse(`
// Code generated by wunderctl. DO NOT EDIT.

import { createWunderGraphApplication } from "@wundergraph/sdk";
import config from '../wundergraph.config';
{{ if .HasWunderGraphServerTs }}import server from '../wundergraph.server';{{ end }}

createWunderGraphApplication(config{{ if .HasWunderGraphServerTs }}, server{{ end }});`))
)

type wunderGraphApplicationTemplateData struct {
	HasWunderGraphServerTs bool
}

func generateWunderGraphApplicationTS(wunderGraphDir string) (string, error) {
	st, err := os.Stat(filepath.Join(wunderGraphDir, "wundergraph.server.ts"))
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	data := &wunderGraphApplicationTemplateData{
		HasWunderGraphServerTs: err == nil && !st.IsDir(),
	}
	generated := filepath.Join(wunderGraphDir, generatedDirectory)
	var buf bytes.Buffer
	if err := wunderGraphFactoryTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	if err := os.MkdirAll(generated, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating %s: %s", generated, err)
	}
	entryPointFilename := filepath.Join(generated, applicationEntryPointFilename)
	if err := os.WriteFile(entryPointFilename, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("error creating %s: %s", entryPointFilename, err)
	}
	return entryPointFilename, nil
}

func hasApplicationConfig(configEntry string) (bool, error) {
	b, err := os.ReadFile(configEntry)
	if err != nil {
		return false, err
	}
	s := string(b)
	return strings.Contains(s, "export default defineConfig") || strings.Contains(s, "WunderGraphConfig"), nil
}

func ApplicationEntryPoint(wunderGraphDir string) (string, error) {
	defaultEntryPoint := filepath.Join(wunderGraphDir, configEntryPointFilename)
	hasApplication, err := hasApplicationConfig(defaultEntryPoint)
	if err != nil {
		return "", err
	}
	if hasApplication {
		return generateWunderGraphApplicationTS(wunderGraphDir)
	}
	return defaultEntryPoint, nil

}
