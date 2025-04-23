package costexplorer

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	ceTypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

// CostExplorerAPI defines the interface for AWS Cost Explorer operations
type CostExplorerAPI interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

// Client represents the AWS Cost Explorer client
type Client struct {
	svc CostExplorerAPI
}

// NewClient creates a new AWS Cost Explorer client
func NewClient(ctx context.Context) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		svc: costexplorer.NewFromConfig(cfg),
	}, nil
}

// GetCosts retrieves costs for a given time period
func (c *Client) GetCosts(ctx context.Context, startDate, endDate, costMetric, grouping string, serviceFilter string, tagFilters []string) (map[string]float64, error) {
	var expressions []ceTypes.Expression

	if serviceFilter != "" {
		expressions = append(expressions, GetDimensionExpression("SERVICE", serviceFilter))
	}

	for _, tag := range tagFilters {
		parts := strings.Split(tag, "=")
		if len(parts) == 2 {
			expressions = append(expressions, GetTagExpression(parts[0], parts[1]))
		}
	}

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &ceTypes.DateInterval{
			Start: aws.String(startDate),
			End:   aws.String(endDate),
		},
		Granularity: ceTypes.GranularityMonthly,
		Metrics:     []string{costMetric},
		GroupBy: []ceTypes.GroupDefinition{
			{
				Type: ceTypes.GroupDefinitionTypeDimension,
				Key:  aws.String(grouping),
			},
		},
	}

	if len(expressions) > 0 {
		input.Filter = &ceTypes.Expression{
			And: expressions,
		}
	}

	result, err := c.svc.GetCostAndUsage(ctx, input)
	if err != nil {
		return nil, err
	}

	costs := make(map[string]float64)
	for _, result := range result.ResultsByTime {
		for _, group := range result.Groups {
			for _, key := range group.Keys {
				amount, _ := strconv.ParseFloat(*group.Metrics[costMetric].Amount, 64)
				costs[key] = amount
			}
		}
	}

	return costs, nil
}

// GetTagExpression creates a tag filter expression
func GetTagExpression(tag string, value string) ceTypes.Expression {
	return ceTypes.Expression{
		Tags: &ceTypes.TagValues{
			Key:    aws.String(tag),
			Values: []string{value},
		},
	}
}

// GetDimensionExpression creates a dimension filter expression
func GetDimensionExpression(dimension string, value string) ceTypes.Expression {
	return ceTypes.Expression{
		Dimensions: &ceTypes.DimensionValues{
			Key:    ceTypes.Dimension(dimension),
			Values: []string{value},
		},
	}
}
