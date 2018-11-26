package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

type ActionTrial struct {
	Name              string `json:"Name" xml:"Name"`
	OssBucketName     string `json:"OssBucketName" xml:"OssBucketName"`
	OssBucketLocation string `json:"OssBucketLocation" xml:"OssBucketLocation"`
	OssKeyPrefix      string `json:"OssKeyPrefix" xml:"OssKeyPrefix"`
	RoleName          string `json:"RoleName" xml:"RoleName"`
}

func CreateTrail(client *cms.Client, request *CreateTrailRequest) (response *CreateTrailResponse, err error) {
	response = CreateCreateTrailResponse()
	err = client.DoAction(request, response)
	return
}

func UpdateTrail(client *cms.Client, request *UpdateTrailRequest) (response *UpdateTrailResponse, err error) {
	response = CreateUpdateTrailResponse()
	err = client.DoAction(request, response)
	return
}

func DeleteTrail(client *cms.Client, request *DeleteTrailRequest) (response *DeleteTrailResponse, err error) {
	response = CreateDeleteTrailResponse()
	err = client.DoAction(request, response)
	return
}

func DescribeTrails(client *cms.Client, request *DescribeTrailsRequest) (response *DescribeTrailsResponse, err error) {
	response = CreateDescribeTrailsResponse()
	err = client.DoAction(request, response)
	return
}

func StartLogging(client *cms.Client, request *StartLoggingRequest) (response *StartLoggingResponse, err error) {
	response = CreateStartLoggingResponse()
	err = client.DoAction(request, response)
	return
}

func StopLogging(client *cms.Client, request *StopLoggingRequest) (response *StopLoggingResponse, err error) {
	response = CreateStopLoggingResponse()
	err = client.DoAction(request, response)
	return
}

// CreateTrailRequest is the request struct for api CreateTrail
type CreateTrailRequest struct {
	*requests.RpcRequest
	Name          string `position:"Query" name:"Name"`
	OssBucketName string `position:"Query" name:"OssBucketName"`
	RoleName      string `position:"Query" name:"RoleName"`
	OssKeyPrefix  string `position:"Query" name:"OssKeyPrefix"`
}

// CreateTrailResponse is the response struct for api CreateTrail
type CreateTrailResponse struct {
	*responses.BaseResponse
	Name          string `json:"Name" xml:"Name"`
	OssBucketName string `json:"OssBucketName" xml:"OssBucketName"`
	OssKeyPrefix  string `json:"OssKeyPrefix" xml:"OssKeyPrefix"`
	RoleName      string `json:"RoleName" xml:"RoleName"`
	RequestId     string `json:"RequestId" xml:"RequestId"`
	HostId        string `json:"HostId" xml:"HostId"`
	Code          string `json:"Code" xml:"Code"`
	Message       string `json:"Message" xml:"Message"`
}

// CreateCreateTrailRequest creates a request to invoke CreateTrail API
func CreateCreateTrailRequest() (request *CreateTrailRequest) {
	request = &CreateTrailRequest{
		RpcRequest: &requests.RpcRequest{},
	}

	request.InitWithApiInfo("", "2015-09-28", "CreateTrail", "", "")
	request.SetDomain("actiontrail.cn-hangzhou.aliyuncs.com")
	request.Method = requests.GET

	return
}

// CreateCreateTrailResponse creates a response to parse from CreateTrail response
func CreateCreateTrailResponse() (response *CreateTrailResponse) {
	response = &CreateTrailResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

// UpdateTrailRequest is the request struct for api UpdateTrail
type UpdateTrailRequest struct {
	*requests.RpcRequest
	Name          string `position:"Query" name:"Name"`
	OssBucketName string `position:"Query" name:"OssBucketName"`
	RoleName      string `position:"Query" name:"RoleName"`
	OssKeyPrefix  string `position:"Query" name:"OssKeyPrefix"`
}

// UpdateTrailResponse is the response struct for api UpdateTrail
type UpdateTrailResponse struct {
	*responses.BaseResponse
	Name          string `json:"Name" xml:"Name"`
	OssBucketName string `json:"OssBucketName" xml:"OssBucketName"`
	OssKeyPrefix  string `json:"OssKeyPrefix" xml:"OssKeyPrefix"`
	RoleName      string `json:"RoleName" xml:"RoleName"`
	HomeRegion    string `json:"HomeRegion" xml:"HomeRegion"`
	RequestId     string `json:"RequestId" xml:"RequestId"`
	HostId        string `json:"HostId" xml:"HostId"`
	Code          string `json:"Code" xml:"Code"`
	Message       string `json:"Message" xml:"Message"`
}

// CreateUpdateTrailRequest creates a request to invoke UpdateTrail API
func CreateUpdateTrailRequest() (request *UpdateTrailRequest) {
	request = &UpdateTrailRequest{
		RpcRequest: &requests.RpcRequest{},
	}

	request.InitWithApiInfo("", "2015-09-28", "UpdateTrail", "", "")
	request.SetDomain("actiontrail.cn-hangzhou.aliyuncs.com")
	request.Method = requests.GET

	return
}

// CreatelUpdateTrailResponse creates a response to parse from UpdateTrail response
func CreateUpdateTrailResponse() (response *UpdateTrailResponse) {
	response = &UpdateTrailResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

// DeleteTrailRequest is the request struct for api DeleteTrail
type DeleteTrailRequest struct {
	*requests.RpcRequest
	Name string `position:"Query" name:"Name"`
}

// DeleteTrailResponse is the response struct for api DeleteTrail
type DeleteTrailResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	HostId    string `json:"HostId" xml:"HostId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

// CreateDeleteTrailRequest creates a request to invoke DeleteTrail API
func CreateDeleteTrailRequest() (request *DeleteTrailRequest) {
	request = &DeleteTrailRequest{
		RpcRequest: &requests.RpcRequest{},
	}

	request.InitWithApiInfo("", "2015-09-28", "DeleteTrail", "", "")
	request.SetDomain("actiontrail.cn-hangzhou.aliyuncs.com")
	request.Method = requests.GET

	return
}

// CreateDeleteTrailResponse creates a response to parse from DeleteTrail response
func CreateDeleteTrailResponse() (response *DeleteTrailResponse) {
	response = &DeleteTrailResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

// DescribeTrailsRequest is the request struct for api DescribeTrails
type DescribeTrailsRequest struct {
	*requests.RpcRequest
	NameList            string `position:"Query" name:"NameList"`
	IncludeShadowTrails string `position:"IncludeShadowTrails" name:"IncludeShadowTrails"`
}

// DescribeTrailsResponse is the response struct for api DescribeTrails
type DescribeTrailsResponse struct {
	*responses.BaseResponse
	TrailList []ActionTrial `json:"TrailList" xml:"TrailList"`
}

// CreateDescribeTrailsRequest creates a request to invoke DescribeTrails API
func CreateDescribeTrailsRequest() (request *DescribeTrailsRequest) {
	request = &DescribeTrailsRequest{
		RpcRequest: &requests.RpcRequest{},
	}

	request.InitWithApiInfo("", "2015-09-28", "DescribeTrails", "", "")
	request.SetDomain("actiontrail.cn-hangzhou.aliyuncs.com")
	request.Method = requests.GET
	request.AcceptFormat = "JSON"

	return
}

// CreateDescribeTrailsResponse creates a response to parse from DescribeTrails response
func CreateDescribeTrailsResponse() (response *DescribeTrailsResponse) {
	response = &DescribeTrailsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

// StartLoggingRequest is the request struct for api StartLogging
type StartLoggingRequest struct {
	*requests.RpcRequest
	Name string `position:"Query" name:"Name"`
}

// StartLoggingResponse is the response struct for api StartLogging
type StartLoggingResponse struct {
	*responses.BaseResponse
}

// CreateStartLoggingRequest creates a request to invoke StartLogging API
func CreateStartLoggingRequest() (request *StartLoggingRequest) {
	request = &StartLoggingRequest{
		RpcRequest: &requests.RpcRequest{},
	}

	request.InitWithApiInfo("", "2015-09-28", "StartLogging", "", "")
	request.SetDomain("actiontrail.cn-hangzhou.aliyuncs.com")
	request.Method = requests.GET
	request.AcceptFormat = "JSON"

	return
}

// CreateStartLoggingResponse creates a response to parse from StartLogging response
func CreateStartLoggingResponse() (response *StartLoggingResponse) {
	response = &StartLoggingResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

// StopLoggingRequest is the request struct for api StopLogging
type StopLoggingRequest struct {
	*requests.RpcRequest
	Name string `position:"Query" name:"Name"`
}

// StopLoggingResponse is the response struct for api StopLogging
type StopLoggingResponse struct {
	*responses.BaseResponse
}

// CreateStopLoggingRequest creates a request to invoke StopLogging API
func CreateStopLoggingRequest() (request *StopLoggingRequest) {
	request = &StopLoggingRequest{
		RpcRequest: &requests.RpcRequest{},
	}

	request.InitWithApiInfo("", "2015-09-28", "StopLogging", "", "")
	request.SetDomain("actiontrail.cn-hangzhou.aliyuncs.com")
	request.Method = requests.GET
	request.AcceptFormat = "JSON"

	return
}

// CreateStopLoggingResponse creates a response to parse from StopLogging response
func CreateStopLoggingResponse() (response *StopLoggingResponse) {
	response = &StopLoggingResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
