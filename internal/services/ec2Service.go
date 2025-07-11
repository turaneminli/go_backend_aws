package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2Service encapsulates EC2 operations
type EC2Service struct {
	Client *ec2.Client
}

type LaunchInstanceInput struct {
	AMI            string   `json:"ami"`
	InstanceType   string   `json:"instanceType"`
	KeyPair        string   `json:"keyPair"`
	SecurityGroups []string `json:"securityGroups"`
	InstanceName   string   `json:"instanceName"`
	MinCount       int32    `json:"minCount"`
	MaxCount       int32    `json:"maxCount"`
	Region         string   `json:"region"`
}

type InstanceStatus struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	State     string `json:"state"`
	PublicIP  string `json:"public_ip"`
	PrivateIP string `json:"private_ip"`
}

type InstanceDetail struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	State          string   `json:"state"`
	PrivateIP      string   `json:"privateIp"`
	PublicIP       string   `json:"publicIp"`
	InstanceType   string   `json:"instanceType"`
	LaunchTime     string   `json:"launchTime"`
	Tags           []string `json:"tags"`
	SecurityGroups []string `json:"securityGroups"`
	Volumes        []string `json:"volumes"`
}

func (s *EC2Service) ListRegions() ([]types.Region, error) {
	output, err := s.Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}
	return output.Regions, nil
}

func (s *EC2Service) ListSecurityGroups() ([]map[string]string, error) {
	req := &ec2.DescribeSecurityGroupsInput{}
	resp, err := s.Client.DescribeSecurityGroups(context.TODO(), req)
	if err != nil {
		return nil, fmt.Errorf("unable to describe security groups: %v", err)
	}

	var result []map[string]string
	for _, sg := range resp.SecurityGroups {
		result = append(result, map[string]string{
			"GroupId":   *sg.GroupId,
			"GroupName": *sg.GroupName,
		})
	}

	return result, nil
}

func (s *EC2Service) LaunchInstance(input LaunchInstanceInput) (string, error) {
	region := input.Region

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}
	s.Client = ec2.NewFromConfig(cfg)

	// Prepare the EC2 run instance input
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

	// Run the instance
	output, err := s.Client.RunInstances(context.TODO(), runInput)
	if err != nil {
		return "", fmt.Errorf("failed to launch instance: %w", err)
	}

	// Return the instance ID if successful
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

func (s *EC2Service) StartInstanceById(instanceID string) (string, error) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []string{instanceID},
	}

	output, err := s.Client.StartInstances(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("unable to start instance: %w", err)
	}

	if len(output.StartingInstances) > 0 {
		return aws.ToString(output.StartingInstances[0].InstanceId), nil
	}

	return "", fmt.Errorf("instance not found or failed to start")
}

func (s *EC2Service) RebootInstanceById(instanceID string) (string, error) {
	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{instanceID},
	}

	_, err := s.Client.RebootInstances(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("unable to reboot instance: %w", err)
	}

	return instanceID, nil
}

func (s *EC2Service) TerminateInstanceById(instanceID string) (string, error) {
	// Create input for terminating the instance
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	// Call TerminateInstances method
	_, err := s.Client.TerminateInstances(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("unable to terminate instance: %w", err)
	}

	// Return the instance ID if successful
	return instanceID, nil
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
			if instance.State != nil && (instance.State.Name == types.InstanceStateNameRunning || instance.State.Name == types.InstanceStateNameStopped) {
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

func (s *EC2Service) GetInstanceDetails(instanceId string) (*InstanceDetail, error) {
	// Create the request to describe the instance
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}

	output, err := s.Client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %v", err)
	}

	// Check if the instance is found
	if len(output.Reservations) == 0 || len(output.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance with ID %s not found", instanceId)
	}

	instance := output.Reservations[0].Instances[0]
	instanceDetail := &InstanceDetail{
		ID:           aws.ToString(instance.InstanceId),
		State:        string(instance.State.Name),
		PrivateIP:    aws.ToString(instance.PrivateIpAddress),
		PublicIP:     aws.ToString(instance.PublicIpAddress),
		InstanceType: string(instance.InstanceType),
	}

	// Convert LaunchTime (*time.Time) to string
	if instance.LaunchTime != nil {
		instanceDetail.LaunchTime = instance.LaunchTime.Format(time.RFC3339) // or any other format
	}

	// Get the name tag of the instance (if available)
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			instanceDetail.Name = *tag.Value
			break
		}
	}

	// Fetch security groups associated with the instance
	for _, sg := range instance.SecurityGroups {
		instanceDetail.SecurityGroups = append(instanceDetail.SecurityGroups, *sg.GroupName)
	}

	// Fetch attached volumes (if any)
	for _, blockDevice := range instance.BlockDeviceMappings {
		instanceDetail.Volumes = append(instanceDetail.Volumes, *blockDevice.Ebs.VolumeId)
	}

	// Fetch instance tags
	for _, tag := range instance.Tags {
		instanceDetail.Tags = append(instanceDetail.Tags, fmt.Sprintf("%s: %s", *tag.Key, *tag.Value))
	}

	return instanceDetail, nil
}
