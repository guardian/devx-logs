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

type FluentbitConfig struct {
	MainConfigFile string
	ApplicationConfigFile string
}

//go:embed fluentbit/fluentbit.conf
var fluentbitConfig string

//go:embed fluentbit/application-logs.conf
var applicationLogsConfig string

const jsonConfigPath = "/etc/config/tags.json"
const applicationLogsConfigPath = "/etc/td-agent-bit/application-logs.conf"
const fluentbitConfigPath = "/etc/td-agent-bit/td-agent-bit.conf"

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
	var dryRun bool

	var rootCmd = &cobra.Command{
		Use:   "devx-logs",
		Short: "devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.",
		Long:  "devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.\n\nConfiguration is typically provided by tags on the instance, but flags are also supported to customise behaviour.",
		Run: func(cmd *cobra.Command, args []string) {

			config := generateConfigs(tagsArg, kinesisStreamNameArg)
			printableConfig := fmt.Sprintf("Main config:\n%s\nApplication config:%s", config.MainConfigFile, config.ApplicationConfigFile)

			if dryRun {
				cmd.Print(printableConfig)
				return
			}

			err := os.WriteFile(fluentbitConfigPath, []byte(config.MainConfigFile), 0644)
			check(err, fmt.Sprintf("unable to write config file to %s: %v", fluentbitConfigPath, err))

			err = os.WriteFile(applicationLogsConfigPath, []byte(config.ApplicationConfigFile), 0644)
			check(err, fmt.Sprintf("unable to write config file to %s: %v", applicationLogsConfigPath, err))

		},
	}

	rootCmd.Flags().StringVar(&kinesisStreamNameArg, "kinesisStreamName", "", "Typically configured via a 'LogKinesisStreamName' tag on the instance, but you can override using this flag. To write to Kinesis, your instance will need the following permissions for this stream: kinesis:DescribeStream, kinesis:PutRecord.")
	rootCmd.Flags().StringVar(&tagsArg, "tags", "", "Typically read from /etc/config/tags.json (see Amigo's cdk-base role here for more info), but you can override using this flag. Pass a comma-separated list of Key=Value pairs, to be included on log records.")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Set to true to print config to stdout rather than write to file.")

	return rootCmd
}

func generateConfigs(tagsArg string, kinesisStreamNameArg string) FluentbitConfig {
	tags, err := getTags(jsonConfigPath, tagsArg)
	check(err, "tags not found")

	kinesisStreamName, err := getKinesisStreamName(tags, kinesisStreamNameArg)
	check(err, "kinesis stream name not found")

	systemdUnit, systemdUnitLookupError := getSystemdUnit(tags)

	if systemdUnitLookupError != nil {
		fmt.Printf("%s: %v\n", "Application log shipping will not be configured", systemdUnitLookupError)
		placeholders := map[string]string{"KINESIS_STREAM": kinesisStreamName, "APPLICATION_LOGS": ""}
		config := replaceReplaceholders(fluentbitConfig, placeholders, tags)
		return FluentbitConfig{
			MainConfigFile: config,
		}
	}

	fluentbitConfigPlaceholders := map[string]string{"KINESIS_STREAM": kinesisStreamName, "APPLICATION_LOGS": "\n@INCLUDE application-logs.conf\n"}
	fluentbitConfig := replaceReplaceholders(fluentbitConfig, fluentbitConfigPlaceholders, tags)
	applicationLogPlaceholders := map[string]string{"SYSTEMD_UNIT": systemdUnit}
	applicationLogsConfig := replaceReplaceholders(applicationLogsConfig, applicationLogPlaceholders, tags)
	return FluentbitConfig{
		MainConfigFile: fluentbitConfig,
		ApplicationConfigFile: applicationLogsConfig,
	}

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

func getSystemdUnit(tags map[string]string) (string, error) {
	name, ok := tags["SystemdUnit"]
	if ok {
		return name, nil
	}
	return "", fmt.Errorf("SystemdUnit tag was not found")
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

		return normaliseTags(tags), nil
	}

	for _, tag := range strings.Split(tagsArg, ",") {
		kv := strings.Split(tag, "=")

		tags[kv[0]] = kv[1]
	}

	return tags, nil
}

func normaliseTags(tags map[string]string) map[string]string {
	for tag, value := range tags {
		if tag == "App" || tag == "Stage" || tag == "Stack" {
			tags[strings.ToLower(tag)] = value
			delete(tags, tag)
		}
	}

	return tags
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
