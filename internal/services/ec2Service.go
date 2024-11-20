package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2Service encapsulates EC2 operations
type EC2Service struct {
	Client *ec2.Client
}

type LaunchInstanceInput struct {
	AMI            string   `json:"ami"`
	InstanceType   string   `json:"instance_type"`
	KeyPair        string   `json:"key_pair"`
	SecurityGroups []string `json:"security_groups"`
	InstanceName   string   `json:"instance_name"`
	MinCount       int32    `json:"min_count"`
	MaxCount       int32    `json:"max_count"`
}

type InstanceStatus struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	State     string `json:"state"`
	PublicIP  string `json:"public_ip"`
	PrivateIP string `json:"private_ip"`
}

func (s *EC2Service) ListRegions() ([]types.Region, error) {
	output, err := s.Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}
	return output.Regions, nil
}

func (s *EC2Service) LaunchInstance(input LaunchInstanceInput) (string, error) {
	runInput := &ec2.RunInstancesInput{
		ImageId:          aws.String(input.AMI),
		InstanceType:     types.InstanceType(input.InstanceType),
		KeyName:          aws.String(input.KeyPair),
		MinCount:         aws.Int32(input.MinCount),
		MaxCount:         aws.Int32(input.MaxCount),
		SecurityGroupIds: input.SecurityGroups,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(input.InstanceName),
					},
				},
			},
		},
	}

	output, err := s.Client.RunInstances(context.TODO(), runInput)
	if err != nil {
		return "", fmt.Errorf("failed to launch instance: %w", err)
	}

	if len(output.Instances) > 0 {
		return aws.ToString(output.Instances[0].InstanceId), nil
	}

	return "", fmt.Errorf("no instances were launched")

}

func (s *EC2Service) StopInstanceById(instanceID string) (string, error) {

	input := &ec2.StopInstancesInput{
		InstanceIds: []string{instanceID},
	}

	output, err := s.Client.StopInstances(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("unable to stop instance: %w", err)
	}

	if len(output.StoppingInstances) > 0 {
		return aws.ToString(output.StoppingInstances[0].InstanceId), nil
	}

	return "", fmt.Errorf("instance not found or failed to stop")
}

func (s *EC2Service) GetAllRunningInstancesStatus() ([]InstanceStatus, error) {
	input := &ec2.DescribeInstancesInput{}

	output, err := s.Client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %v", err)
	}

	var runningInstances []InstanceStatus

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State != nil && instance.State.Name == types.InstanceStateNameRunning {
				instanceName := "N/A"
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						instanceName = *tag.Value
						break
					}
				}

				// Append the instance info to the runningInstances slice
				runningInstances = append(runningInstances, InstanceStatus{
					Name:      instanceName,
					ID:        aws.ToString(instance.InstanceId),
					State:     string(instance.State.Name),
					PublicIP:  aws.ToString(instance.PublicIpAddress),
					PrivateIP: aws.ToString(instance.PrivateIpAddress),
				})
			}
		}
	}

	return runningInstances, nil
}
