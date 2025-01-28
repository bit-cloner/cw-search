package main

import (
	"context"
	"log"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func promptInput(message string) string {
	var input string
	if err := survey.AskOne(&survey.Input{Message: message}, &input); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	return input
}

func promptSelectRegions(regions []string) []string {
	var selectedRegions []string
	prompt := &survey.MultiSelect{
		Message:  "Select Regions:Default is eu-west-1",
		Options:  regions,
		PageSize: 15,
		Default:  "eu-west-1",
	}
	if err := survey.AskOne(prompt, &selectedRegions); err != nil {
		log.Fatalf("Error selecting regions: %v", err)
	}
	return selectedRegions
}

// getAllRegions retrieves all AWS regions
func getAllRegions(ctx context.Context, cfg aws.Config) ([]string, error) {
	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	var regions []string
	for _, region := range output.Regions {
		regions = append(regions, *region.RegionName)
	}
	sort.Strings(regions)
	return regions, nil
}

// listLogGroups lists all log groups in a specific region
func listLogGroups(ctx context.Context, client *cloudwatchlogs.Client) ([]string, error) {
	req := &cloudwatchlogs.DescribeLogGroupsInput{}
	var logGroups []string

	for {
		resp, err := client.DescribeLogGroups(ctx, req)
		if err != nil {
			return nil, err
		}

		for _, group := range resp.LogGroups {
			logGroups = append(logGroups, *group.LogGroupName)
		}

		if resp.NextToken == nil {
			break
		}
		req.NextToken = resp.NextToken
	}

	return logGroups, nil
}

// filterLogEvents retrieves events matching a filter pattern in a log group
func filterLogEvents(ctx context.Context, client *cloudwatchlogs.Client, region, logGroupName, filterPattern string) ([]string, error) {
	req := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		FilterPattern: aws.String(filterPattern),
	}
	var events []string

	for {
		resp, err := client.FilterLogEvents(ctx, req)
		if err != nil {
			return nil, err
		}

		for _, event := range resp.Events {
			events = append(events, *event.Message)
		}

		if resp.NextToken == nil {
			break
		}
		req.NextToken = resp.NextToken
	}

	return events, nil
}
