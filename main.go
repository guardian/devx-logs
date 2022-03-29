package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed fluentbit/fluentbit.conf
var fluentbitConfig string

const jsonConfigPath = "/etc/config/tags.json"

func main() {
	rootCmd := RootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RootCmd() *cobra.Command {
	var kinesisStreamNameArg string
	var tagsArg string

	var rootCmd = &cobra.Command{
		Use:   "devx-logs",
		Short: "devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.",
		Long:  "devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.\n\nConfiguration is typically provided by tags on the instance, but flags are also supported to customise behaviour.",
		Run: func(cmd *cobra.Command, args []string) {
			tags, err := getTags(jsonConfigPath, tagsArg)
			check(err, "tags not found")

			kinesisStreamName, err := getKinesisStreamName(tags, kinesisStreamNameArg)
			check(err, "kinesis stream name not found")

			placeholders := map[string]string{"KINESIS_STREAM": kinesisStreamName}
			fluentbitConfig := replaceReplaceholders(fluentbitConfig, placeholders, tags)
			cmd.Print(fluentbitConfig)
		},
	}

	rootCmd.Flags().StringVar(&kinesisStreamNameArg, "kinesisStreamName", "", "Typically configured via a 'LogKinesisStreamName' tag on the instance, but you can override using this flag. To write to Kinesis, your instance will need the following permissions for this stream: kinesis:DescribeStream, kinesis:PutRecord.")
	rootCmd.Flags().StringVar(&tagsArg, "tags", "", "Typically read from /etc/config/tags.json (see Amigo's cdk-base role here for more info), but you can override using this flag. Pass a comma-separated list of Key=Value pairs, to be included on log records.")

	return rootCmd
}

func getKinesisStreamName(tags map[string]string, kinesisStreamNameArg string) (string, error) {
	if kinesisStreamNameArg != "" {
		return kinesisStreamNameArg, nil
	}

	name, ok := tags["LogKinesisStreamName"]
	if ok {
		return name, nil
	}

	return "", fmt.Errorf("Kinesis Stream name was not found: no LogKinesisStreamName tag, and the --kinesisStreamName arg was empty.")
}

func getTags(jsonConfigPath string, tagsArg string) (map[string]string, error) {
	tags := map[string]string{}

	if tagsArg == "" { // lookup from config file
		raw, err := os.ReadFile(jsonConfigPath)
		if err != nil {
			return tags, fmt.Errorf("unable to read tags from %s: %w", jsonConfigPath, err)
		}

		err = json.Unmarshal(raw, &tags)
		if err != nil {
			return tags, fmt.Errorf("unable to unmarshal JSON config from %s: %w", jsonConfigPath, err)
		}

		return tags, nil
	}

	for _, tag := range strings.Split(tagsArg, ",") {
		kv := strings.Split(tag, "=")
		tags[kv[0]] = kv[1]
	}

	return tags, nil
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

	updated = strings.ReplaceAll(updated, "{{TAGS}}", strings.Join(addTags, "\n    "))

	return updated
}

func check(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}
