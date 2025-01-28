package main

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/fatih/color"
)

func main() {
	ctx := context.Background()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	// Prompt for input
	logGroupName := promptInput("Enter Log Group name or press return to search all groups:")
	if logGroupName == "" {
		logGroupName = "ALL"
	}
	filterPattern := promptInput("Enter filter pattern to search logs:")
	if filterPattern == "" {
		log.Fatalf("Filter pattern must be set.")
	}

	// Prompt for timeframe
	timeframe := promptTimeframe()

	// Calculate start and end times based on selected timeframe
	endTime := time.Now()
	var startTime time.Time
	switch timeframe {
	case "5 minutes":
		startTime = endTime.Add(-5 * time.Minute)
	case "30 minutes":
		startTime = endTime.Add(-30 * time.Minute)
	case "1 hour":
		startTime = endTime.Add(-1 * time.Hour)
	case "6 hours":
		startTime = endTime.Add(-6 * time.Hour)
	case "12 hours":
		startTime = endTime.Add(-12 * time.Hour)
	case "1 day":
		startTime = endTime.Add(-24 * time.Hour)
	case "3 days":
		startTime = endTime.Add(-72 * time.Hour)
	case "7 days":
		startTime = endTime.Add(-168 * time.Hour)
	case "custom":
		// prompt user for input for year month date our minute minute second
		// and parse the input to time
		yearInput := promptInput("Enter year in format YYYY:")
		monthInput := promptInput("Enter month in format MM:")
		dayInput := promptInput("Enter day in format DD:")
		hourInput := promptInput("Enter hour in format HH (default 00):")
		if hourInput == "" {
			hourInput = "00"
		}
		minuteInput := promptInput("Enter minute in format MM (default 00):")
		if minuteInput == "" {
			minuteInput = "00"
		}
		secondInput := promptInput("Enter second in format SS (default 00):")
		if secondInput == "" {
			secondInput = "00"
		}

		startTimeStr := fmt.Sprintf("%s-%s-%sT%s:%s:%s", yearInput, monthInput, dayInput, hourInput, minuteInput, secondInput)
		startTime, err = time.Parse("2006-01-02T15:04:05", startTimeStr)
		if err != nil {
			log.Fatalf("Invalid start time format: %v", err)
		}
		// do the same for end time
		// and parse the input to time
		yearInput = promptInput("Enter year in format YYYY:")
		monthInput = promptInput("Enter month in format MM:")
		dayInput = promptInput("Enter day in format DD:")
		hourInput = promptInput("Enter hour in format HH (default 00):")
		if hourInput == "" {
			hourInput = "00"
		}
		minuteInput = promptInput("Enter minute in format MM (default 00):")
		if minuteInput == "" {
			minuteInput = "00"
		}
		secondInput = promptInput("Enter second in format SS (default 00):")
		if secondInput == "" {
			secondInput = "00"
		}
		endTimeStr := fmt.Sprintf("%s-%s-%sT%s:%s:%s", yearInput, monthInput, dayInput, hourInput, minuteInput, secondInput)
		endTime, err = time.Parse("2006-01-02T15:04:05", endTimeStr)
		if err != nil {
			log.Fatalf("Invalid end time format: %v", err)
		}

	default:
		// select 6 hours as default
		startTime = endTime.Add(-6 * time.Hour)
	}

	// Fetch all regions
	regions, err := getAllRegions(ctx, cfg)
	if err != nil {
		log.Fatalf("Error retrieving regions: %v", err)
	}

	// Prompt user to select regions
	promptSelectRegions(regions)

	// Iterate over regions and log groups
	cwlClient := cloudwatchlogs.NewFromConfig(cfg)
	for _, reg := range regions {
		var logGroups []string
		if logGroupName == "ALL" {
			logGroups, err = listLogGroups(ctx, cwlClient)
			if err != nil {
				color.Red("Error listing log groups in region %s: %v\n", reg, err)
				continue
			}
		} else {
			logGroups = []string{logGroupName}
		}

		// Search logs within the specified timeframe
		for _, logGroup := range logGroups {
			searchLogs(ctx, cwlClient, logGroup, filterPattern, startTime, endTime)
		}
	}
}

func promptTimeframe() string {
	var timeframe string
	prompt := &survey.Select{
		Message:  "Select timeframe:",
		Options:  []string{"5 minutes", "30 minutes", "1 hour", "6 hours", "12 hours", "1 day", "3 days", "7 days", "custom"},
		Default:  "6 hours",
		PageSize: 10,
	}
	err := survey.AskOne(prompt, &timeframe)
	if err != nil {
		log.Fatalf("Error selecting timeframe: %v", err)
	}
	return timeframe
}

func searchLogs(ctx context.Context, client *cloudwatchlogs.Client, logGroupName, filterPattern string, startTime, endTime time.Time) {
	// Implement the logic to search logs within the specified timeframe
	// using the CloudWatch Logs client
	fmt.Printf("Searching logs in log group: %s\n", logGroupName)

	input := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		FilterPattern: aws.String(filterPattern),
		StartTime:     aws.Int64(startTime.Unix() * 1000),
		EndTime:       aws.Int64(endTime.Unix() * 1000),
	}

	paginator := cloudwatchlogs.NewFilterLogEventsPaginator(client, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			color.Red("Error fetching log events: %v\n", err)
			return
		}
		// When a match is found, then print a message with the green tech saying that a match has been found in Log group
		if len(output.Events) > 0 {
			color.Green("Match found in log group: %s\n", logGroupName)
		}

		for _, event := range output.Events {
			fmt.Printf("[%s] %s\n", time.Unix(0, *event.Timestamp*int64(time.Millisecond)).Format(time.RFC3339), *event.Message)
		}
	}
}
