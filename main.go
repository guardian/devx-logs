package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

//go:embed fluentbit.conf
var rawConfig string

func main() {
	ec2MetaURL := "http://169.254.169.254"
	rootCmd := RootCmd(ec2MetaURL)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RootCmd(metadataURL string) *cobra.Command {
	var logStreamARN string
	var systemdUnit string

	var rootCmd = &cobra.Command{
		Use:   "devx-logs",
		Short: "devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.",
		Long:  `devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.`,
		Run: func(cmd *cobra.Command, args []string) {
			tags, err := getTags(metadataURL)
			warn(err, "Unable to read tags")

			placeholders := map[string]string{
				"KINESIS_STREAM": logStreamARN,
				"SYSTEMD_UNIT":   systemdUnit,
			}

			config := replaceReplaceholders(rawConfig, placeholders, tags)
			cmd.Print(config)
		},
	}

	rootCmd.Flags().StringVar(&logStreamARN, "logStreamARN", "", "Set to a Kinesis log stream ARN. Your instance will need the following permissions for this stream: kinesis:DescribeStream, kinesis:PutRecord.")
	rootCmd.MarkFlagRequired("logStreamARN")

	rootCmd.Flags().StringVar(&systemdUnit, "systemdUnit", "", "Set to a SystemD Unit. Used to filter JournalD records.")
	rootCmd.MarkFlagRequired("systemdUnit")

	return rootCmd
}

func getTags(metadataURL string) (map[string]string, error) {
	tags := map[string]string{}

	client := http.DefaultClient
	client.Timeout = time.Second * 2

	token, err := getMetadataToken(client, metadataURL)
	if err != nil {
		return tags, err
	}

	tagData, err := getMetadata(client, metadataURL+"/latest/meta-data/tags/instance", token)
	tagNames := strings.Split(strings.TrimSpace(string(tagData)), "\n")

	for _, name := range tagNames {
		value, err := getMetadata(client, fmt.Sprintf("%s/latest/meta-data/tags/instance/%s", metadataURL, name), token)
		if err != nil {
			return tags, err
		}

		tags[name] = strings.TrimSpace(string(value))
	}

	return tags, err
}

func getMetadataToken(client *http.Client, baseURL string) (string, error) {
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/latest/api/token", baseURL), nil)
	req.Header.Add("X-aws-ec2-metadata-token-ttl-seconds", "60")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	token, err := io.ReadAll(resp.Body)
	return string(token), err
}

func getMetadata(client *http.Client, URL string, token string) ([]byte, error) {
	resp, err := client.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func replaceReplaceholders(config string, info map[string]string, tags map[string]string) string {
	updated := config

	for old, new := range info {
		updated = strings.ReplaceAll(updated, "{{"+old+"}}", new)
	}

	addTags := []string{}
	for name, value := range tags {
		addTags = append(addTags, fmt.Sprintf("Add %s %s", name, value))
	}

	// range over map iterates in random order so let's make it deterministic
	// here (partly for tests).
	sort.StringSlice(addTags).Sort()

	updated = strings.ReplaceAll(updated, "{{TAGS}}", strings.Join(addTags, "\n  "))

	return updated
}

func warn(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
	}
}
