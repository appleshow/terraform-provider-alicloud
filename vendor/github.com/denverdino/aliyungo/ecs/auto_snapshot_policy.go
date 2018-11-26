package ecs

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type CreateAutoSnapshotPolocyArgs struct {
	RegionId           common.Region
	TimePoints             string
	RepeatWeekdays         string
	RetentionDays          requests.Integer
	AutoSnapshotPolicyName           string
	ClientToken  string
}

type CreateAutoSnapshotPolocyResponse struct {
	common.Response
	AutoSnapshotPolicyId string
}

type DescribeAutoSnapshotPolocyExArgs struct {
	RegionId           common.Region
	AutoSnapshotPolicyId             string
}

type ModifyAutoSnapshotPolicyExArgs struct {
	RegionId           common.Region
	AutoSnapshotPolicyId             string
	AutoSnapshotPolicyName             string
	TimePoints           string
	RepeatWeekdays        string
	RetentionDays         requests.Integer 
	ClientToken  string
}

type ModifyAutoSnapshotPolicyExResponse struct {
	common.Response
}

type DeleteAutoSnapshotPolicyArgs struct {
	RegionId           common.Region
	AutoSnapshotPolicyId             string
}

type DescribeAutoSnapshotPolicyExResponse struct {
	common.Response
	common.PaginationResult
	RegionId common.Region
	AutoSnapshotPolicies    struct {
		AutoSnapshotPolicy []AutoSnapshotPolicyType 
	}
}

type DeleteAutoSnapshotPolicyResponse struct {
	common.Response
}

type AutoSnapshotPolicyType struct {
	AutoSnapshotPolicyId             string
	RegionId           common.Region
	AutoSnapshotPolicyName             string
	TimePoints           string
	RepeatWeekdays        string
	RetentionDays               int
	DiskNums          int
	OperationLocks           OperationLocksType
	Status               string
	CreationTime		string
}

type AutoSnapshotPolicyApplicationArgs struct {
	RegionId           common.Region	
	AutoSnapshotPolicyId         string
	DiskIds             string
}

type AutoSnapshotPolicyApplicationResponse struct {
	common.Response
}

func (client *Client) ApplyAutoSnapshotPolicy(args *AutoSnapshotPolicyApplicationArgs) error {
	response := AutoSnapshotPolicyApplicationResponse{}
	err := client.Invoke("ApplyAutoSnapshotPolicy", args, &response)
	return err
}

func (client *Client) CancelAutoSnapshotPolicy(args *AutoSnapshotPolicyApplicationArgs) error {
	response := AutoSnapshotPolicyApplicationResponse{}
	err := client.Invoke("CancelAutoSnapshotPolicy", args, &response)
	return err
}

func (client *Client) CreateAutoSnapshotPolicy(args *CreateAutoSnapshotPolocyArgs) (autoSnapshotPolicyId string, err error) {
	response := CreateAutoSnapshotPolocyResponse{}
	err = client.Invoke("CreateAutoSnapshotPolicy", args, &response)
	if err != nil {
		return "", err
	}
	return response.AutoSnapshotPolicyId, err
}

func (client *Client) DescribeAutoSnapshotPolicyEx(args *DescribeAutoSnapshotPolocyExArgs) (autoSnapshotPolicies []AutoSnapshotPolicyType, pagination *common.PaginationResult, err error) {
	response, err := client.DescribeAutoSnapshotPolicyExWithRaw(args)
	if err != nil {
		return nil, nil, err
	}

	return response.AutoSnapshotPolicies.AutoSnapshotPolicy, &response.PaginationResult, err
}

func (client *Client) DescribeAutoSnapshotPolicyExWithRaw(args *DescribeAutoSnapshotPolocyExArgs) (response *DescribeAutoSnapshotPolicyExResponse, err error) {
	response = &DescribeAutoSnapshotPolicyExResponse{}

	err = client.Invoke("DescribeAutoSnapshotPolicyEx", args, response)

	if err != nil {
		return nil, err
	}

	return response, err
}

func (client *Client) ModifyAutoSnapshotPolicyEx(args *ModifyAutoSnapshotPolicyExArgs) error {
	response := ModifyAutoSnapshotPolicyExResponse{}
	err := client.Invoke("ModifyAutoSnapshotPolicyEx", args, &response)
	return err
}

func (client *Client) DeleteAutoSnapshotPolicy(args *DeleteAutoSnapshotPolicyArgs) error {
	response := DeleteAutoSnapshotPolicyResponse{}
	err := client.Invoke("DeleteAutoSnapshotPolicy", args, &response)
	return err
}

