package costexplorer

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	ceTypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCostExplorerClient is a mock implementation of the AWS Cost Explorer client
type MockCostExplorerClient struct {
	mock.Mock
}

func (m *MockCostExplorerClient) GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*costexplorer.GetCostAndUsageOutput), args.Error(1)
}

func TestGetCosts(t *testing.T) {
	mockClient := new(MockCostExplorerClient)
	client := &Client{svc: mockClient}

	// Test data
	startDate := "2024-01-01"
	endDate := "2024-02-01"
	costMetric := "NetAmortizedCost"
	grouping := "SERVICE"
	serviceFilter := ""
	tagFilters := []string{}

	// Mock response
	mockResponse := &costexplorer.GetCostAndUsageOutput{
		ResultsByTime: []ceTypes.ResultByTime{
			{
				Groups: []ceTypes.Group{
					{
						Keys: []string{"Amazon EC2"},
						Metrics: map[string]ceTypes.MetricValue{
							"NetAmortizedCost": {
								Amount: aws.String("100.00"),
								Unit:   aws.String("USD"),
							},
						},
					},
				},
			},
		},
	}

	mockClient.On("GetCostAndUsage", mock.Anything, mock.Anything).Return(mockResponse, nil)

	// Test the function
	costs, err := client.GetCosts(context.Background(), startDate, endDate, costMetric, grouping, serviceFilter, tagFilters)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, costs)
	assert.Equal(t, 100.00, costs["Amazon EC2"])
	mockClient.AssertExpectations(t)
}

func TestGetTagExpression(t *testing.T) {
	tag := "Environment"
	value := "Production"
	expr := GetTagExpression(tag, value)

	assert.NotNil(t, expr.Tags)
	assert.Equal(t, tag, *expr.Tags.Key)
	assert.Equal(t, []string{value}, expr.Tags.Values)
}

func TestGetDimensionExpression(t *testing.T) {
	dimension := "SERVICE"
	value := "Amazon EC2"
	expr := GetDimensionExpression(dimension, value)

	assert.NotNil(t, expr.Dimensions)
	assert.Equal(t, ceTypes.Dimension(dimension), expr.Dimensions.Key)
	assert.Equal(t, []string{value}, expr.Dimensions.Values)
}
