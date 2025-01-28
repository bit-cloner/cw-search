package main

import (
	"context"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func promptInput(message string) string {
	var input string
	if err := survey.AskOne(&survey.Input{Message: message}, &input); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	return input
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
