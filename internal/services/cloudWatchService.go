package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type CloudWatchService struct {
	Client *cloudwatch.Client
}

type EC2Metrics struct {
	InstanceID string  `json:"instance_id"`
	MetricName string  `json:"metric_name"`
	Timestamp  string  `json:"timestamp"`
	Value      float64 `json:"value"`
}

func (s *CloudWatchService) GetEC2Metrics(instanceID string) ([]EC2Metrics, error) {
	metrics := []types.MetricDataQuery{
		{
			Id: aws.String("cpu_utilization"),
			MetricStat: &types.MetricStat{
				Metric: &types.Metric{
					Namespace:  aws.String("AWS/EC2"),
					MetricName: aws.String("CPUUtilization"),
					Dimensions: []types.Dimension{
						{
							Name:  aws.String("InstanceId"),
							Value: aws.String(instanceID),
						},
					},
				},
				Period: aws.Int32(60), // Period in seconds
				Stat:   aws.String("Average"),
			},
			ReturnData: aws.Bool(true),
		},
		{
			Id: aws.String("disk_read_bytes"),
			MetricStat: &types.MetricStat{
				Metric: &types.Metric{
					Namespace:  aws.String("AWS/EC2"),
					MetricName: aws.String("DiskReadBytes"),
					Dimensions: []types.Dimension{
						{
							Name:  aws.String("InstanceId"),
							Value: aws.String(instanceID),
						},
					},
				},
				Period: aws.Int32(60),
				Stat:   aws.String("Sum"),
			},
			ReturnData: aws.Bool(true),
		},
	}

	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: metrics,
		StartTime:         aws.Time(time.Now().Add(-1 * time.Hour)),
		EndTime:           aws.Time(time.Now()),
	}

	output, err := s.Client.GetMetricData(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric data: %v", err)
	}

	var ec2Metrics []EC2Metrics
	for _, result := range output.MetricDataResults {
		for i, timestamp := range result.Timestamps {
			ec2Metrics = append(ec2Metrics, EC2Metrics{
				InstanceID: instanceID,
				MetricName: *result.Label,
				Timestamp:  timestamp.String(),
				Value:      result.Values[i],
			})
		}
	}

	return ec2Metrics, nil
}
