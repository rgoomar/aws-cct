package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/rgoomar/aws-cct/pkg/costexplorer"
	"github.com/rgoomar/aws-cct/pkg/display"
	"github.com/rgoomar/aws-cct/pkg/models"
	"github.com/urfave/cli/v2"
)

func main() {
	var dateFormat = "2006-01-02"
	var firstMonthStart string
	var secondMonthStart string
	var costMetric string
	var serviceFilter string
	var sortColumn string
	var sortOrder string
	var output string

	currentDate := time.Now()
	thisMonthFirst := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	previousMonthFirst := thisMonthFirst.AddDate(0, -1, 0)
	nextMonthFirst := thisMonthFirst.AddDate(0, 1, 0)
	lastDayOfThisMonth := nextMonthFirst.AddDate(0, 0, -1).Day()

	app := &cli.App{
		Name:  "aws-cct",
		Usage: "AWS Cost Comparison Tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "start",
				Value:       previousMonthFirst.Format(dateFormat),
				Usage:       "First month to compare (2020-01-01)",
				Destination: &firstMonthStart,
			},
			&cli.StringFlag{
				Name:        "end",
				Value:       thisMonthFirst.Format(dateFormat),
				Usage:       "Second month to compare (2020-02-01)",
				Destination: &secondMonthStart,
			},
			&cli.StringFlag{
				Name:        "cost-metric",
				Value:       "NetAmortizedCost",
				Usage:       "Cost Metric to compare (NetAmortizedCost, UnblendedCost, etc.)",
				Destination: &costMetric,
			},
			&cli.StringFlag{
				Name:        "service",
				Value:       "",
				Usage:       "Define a service to dig into",
				Destination: &serviceFilter,
			},
			&cli.StringSliceFlag{
				Name:  "tag",
				Usage: "Tag value to filter results (app=web, env=prod, etc.)",
			},
			&cli.StringFlag{
				Name:        "sort",
				Value:       "name",
				Usage:       "Column to sort results on (name, start, end, delta, deltapercent)",
				Destination: &sortColumn,
			},
			&cli.StringFlag{
				Name:        "sort-order",
				Value:       "asc",
				Usage:       "Order to sort in (asc or desc)",
				Destination: &sortOrder,
			},
			&cli.StringFlag{
				Name:        "output",
				Value:       "table",
				Usage:       "Output format (supported formats: table, csv)",
				Destination: &output,
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			client, err := costexplorer.NewClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create AWS client: %w", err)
			}

			start, _ := time.Parse(dateFormat, firstMonthStart)
			firstMonthEnd := start.AddDate(0, 1, 0).Format(dateFormat)
			end, _ := time.Parse(dateFormat, secondMonthStart)
			isProjection := false
			multiplier := 1.0
			var secondMonthEnd string

			if currentDate.Month() == end.Month() {
				isProjection = true
				secondMonthEndDate := currentDate.AddDate(0, 0, -1)
				multiplier = float64(lastDayOfThisMonth) / float64(secondMonthEndDate.Day())
				secondMonthEnd = secondMonthEndDate.Format(dateFormat)
			} else {
				secondMonthEnd = end.AddDate(0, 1, 0).Format(dateFormat)
			}

			var grouping = "SERVICE"
			if serviceFilter != "" {
				grouping = "USAGE_TYPE"
			}

			tagFilters := c.StringSlice("tag")

			firstResultsCosts, err := client.GetCosts(ctx, firstMonthStart, firstMonthEnd, costMetric, grouping, serviceFilter, tagFilters)
			if err != nil {
				return fmt.Errorf("failed to get first month costs: %w", err)
			}

			secondResultsCosts, err := client.GetCosts(ctx, secondMonthStart, secondMonthEnd, costMetric, grouping, serviceFilter, tagFilters)
			if err != nil {
				return fmt.Errorf("failed to get second month costs: %w", err)
			}

			allServiceNames := extractAllServiceNames(firstResultsCosts, secondResultsCosts)
			var serviceCosts []models.ServiceCosts
			var totalAmount, totalSecondAmount, totalDelta float64

			for _, service := range allServiceNames {
				amount := firstResultsCosts[service]
				secondAmount := secondResultsCosts[service] * multiplier
				delta := secondAmount - amount
				deltaPercent := 0.0
				if amount != 0 {
					deltaPercent = delta / amount * 100
				}

				serviceCosts = append(serviceCosts, models.ServiceCosts{
					ServiceName:  service,
					Amount:       amount,
					SecondAmount: secondAmount,
					Delta:        delta,
					DeltaPercent: deltaPercent,
				})

				totalAmount += amount
				totalSecondAmount += secondAmount
				totalDelta += delta
			}

			// Sort results
			sort.Slice(serviceCosts, func(i, j int) bool {
				var retVal bool
				switch sortColumn {
				case "start":
					retVal = serviceCosts[i].Amount < serviceCosts[j].Amount
				case "end":
					retVal = serviceCosts[i].SecondAmount < serviceCosts[j].SecondAmount
				case "delta":
					retVal = serviceCosts[i].Delta < serviceCosts[j].Delta
				case "deltapercent":
					retVal = serviceCosts[i].DeltaPercent < serviceCosts[j].DeltaPercent
				default: // default to service name
					retVal = serviceCosts[i].ServiceName < serviceCosts[j].ServiceName
				}
				if sortOrder == "desc" {
					retVal = !retVal
				}
				return retVal
			})

			comparison := &models.CostComparison{
				FirstMonthStart:   firstMonthStart,
				SecondMonthStart:  secondMonthStart,
				IsProjection:      isProjection,
				Multiplier:        multiplier,
				ServiceCosts:      serviceCosts,
				TotalAmount:       totalAmount,
				TotalSecondAmount: totalSecondAmount,
				TotalDelta:        totalDelta,
				TotalDeltaPercent: totalDelta / totalAmount * 100,
			}

			tableDisplay := display.NewTableDisplay()
			fmt.Print(tableDisplay.Render(comparison, output))

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func extractAllServiceNames(firstResultsCosts, secondResultsCosts map[string]float64) []string {
	var allServiceNames []string
	for serviceName := range firstResultsCosts {
		allServiceNames = append(allServiceNames, serviceName)
	}
	for serviceName := range secondResultsCosts {
		if _, ok := firstResultsCosts[serviceName]; !ok {
			allServiceNames = append(allServiceNames, serviceName)
		}
	}
	sort.Strings(allServiceNames)
	return allServiceNames
}
