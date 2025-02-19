package mocks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/stretchr/testify/mock"
)

type FakeMetricsAPI struct {
	cloudwatchiface.CloudWatchAPI
	cloudwatch.GetMetricDataOutput

	Metrics        []*cloudwatch.Metric
	MetricsPerPage int

	CallsGetMetricDataWithContext []*cloudwatch.GetMetricDataInput
}

func (c *FakeMetricsAPI) GetMetricDataWithContext(ctx aws.Context, input *cloudwatch.GetMetricDataInput, opts ...request.Option) (*cloudwatch.GetMetricDataOutput, error) {
	c.CallsGetMetricDataWithContext = append(c.CallsGetMetricDataWithContext, input)

	return &c.GetMetricDataOutput, nil
}

func (c *FakeMetricsAPI) ListMetricsPages(input *cloudwatch.ListMetricsInput, fn func(*cloudwatch.ListMetricsOutput, bool) bool) error {
	if c.MetricsPerPage == 0 {
		c.MetricsPerPage = 1000
	}
	chunks := chunkSlice(c.Metrics, c.MetricsPerPage)

	for i, metrics := range chunks {
		response := fn(&cloudwatch.ListMetricsOutput{
			Metrics: metrics,
		}, i+1 == len(chunks))
		if !response {
			break
		}
	}
	return nil
}

func chunkSlice(slice []*cloudwatch.Metric, chunkSize int) [][]*cloudwatch.Metric {
	var chunks [][]*cloudwatch.Metric
	for {
		if len(slice) == 0 {
			break
		}
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

type MetricsClient struct {
	cloudwatchiface.CloudWatchAPI
	mock.Mock
}

func (m *MetricsClient) GetMetricDataWithContext(ctx aws.Context, input *cloudwatch.GetMetricDataInput, opts ...request.Option) (*cloudwatch.GetMetricDataOutput, error) {
	args := m.Called(ctx, input, opts)

	return args.Get(0).(*cloudwatch.GetMetricDataOutput), args.Error(1)
}
