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
	var kinesisStreamName string
	var systemdUnit string
	var tagsArg string

	var rootCmd = &cobra.Command{
		Use:   "devx-logs",
		Short: "devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.",
		Long:  `devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.`,
		Run: func(cmd *cobra.Command, args []string) {
			tags, err := getTags(jsonConfigPath, tagsArg)
			check(err, "tags not found")

			placeholders := map[string]string{
				"KINESIS_STREAM": kinesisStreamName,
				"SYSTEMD_UNIT":   systemdUnit,
			}

			fluentbitConfig := replaceReplaceholders(fluentbitConfig, placeholders, tags)
			cmd.Print(fluentbitConfig)
		},
	}

	rootCmd.Flags().StringVar(&kinesisStreamName, "kinesisStreamName", "", "Set to a Kinesis log stream name. Your instance will need the following permissions for this stream: kinesis:DescribeStream, kinesis:PutRecord.")
	rootCmd.MarkFlagRequired("kinesisStreamName")

	rootCmd.Flags().StringVar(&systemdUnit, "systemdUnit", "", "Set to the name of your app's systemd service. I.e. 'name' from [name].service")
	rootCmd.MarkFlagRequired("systemdUnits")

	rootCmd.Flags().StringVar(&tagsArg, "tags", "", "Set a comma-separated list of Key=Value pairs, to be included on log records. At the least, this should include App, Stack, and Stage. Eg. 'App=foo,Stage=PROD,Stack=bar'. If empty, tags will be sourced from /etc/config/tags.json (see the cdk-base Amigo role).")

	return rootCmd
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
