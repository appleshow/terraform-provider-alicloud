package cdn

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeLiveStreamFrameAndBitRateByDomain invokes the cdn.DescribeLiveStreamFrameAndBitRateByDomain API synchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreamframeandbitratebydomain.html
func (client *Client) DescribeLiveStreamFrameAndBitRateByDomain(request *DescribeLiveStreamFrameAndBitRateByDomainRequest) (response *DescribeLiveStreamFrameAndBitRateByDomainResponse, err error) {
	response = CreateDescribeLiveStreamFrameAndBitRateByDomainResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLiveStreamFrameAndBitRateByDomainWithChan invokes the cdn.DescribeLiveStreamFrameAndBitRateByDomain API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreamframeandbitratebydomain.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamFrameAndBitRateByDomainWithChan(request *DescribeLiveStreamFrameAndBitRateByDomainRequest) (<-chan *DescribeLiveStreamFrameAndBitRateByDomainResponse, <-chan error) {
	responseChan := make(chan *DescribeLiveStreamFrameAndBitRateByDomainResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLiveStreamFrameAndBitRateByDomain(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeLiveStreamFrameAndBitRateByDomainWithCallback invokes the cdn.DescribeLiveStreamFrameAndBitRateByDomain API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreamframeandbitratebydomain.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamFrameAndBitRateByDomainWithCallback(request *DescribeLiveStreamFrameAndBitRateByDomainRequest, callback func(response *DescribeLiveStreamFrameAndBitRateByDomainResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLiveStreamFrameAndBitRateByDomainResponse
		var err error
		defer close(result)
		response, err = client.DescribeLiveStreamFrameAndBitRateByDomain(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeLiveStreamFrameAndBitRateByDomainRequest is the request struct for api DescribeLiveStreamFrameAndBitRateByDomain
type DescribeLiveStreamFrameAndBitRateByDomainRequest struct {
	*requests.RpcRequest
	AppName       string           `position:"Query" name:"AppName"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	PageSize      requests.Integer `position:"Query" name:"PageSize"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	PageNumber    requests.Integer `position:"Query" name:"PageNumber"`
}

// DescribeLiveStreamFrameAndBitRateByDomainResponse is the response struct for api DescribeLiveStreamFrameAndBitRateByDomain
type DescribeLiveStreamFrameAndBitRateByDomainResponse struct {
	*responses.BaseResponse
	RequestId                string                                                              `json:"RequestId" xml:"RequestId"`
	Count                    int                                                                 `json:"Count" xml:"Count"`
	PageNumber               int                                                                 `json:"pageNumber" xml:"pageNumber"`
	PageSize                 int                                                                 `json:"pageSize" xml:"pageSize"`
	FrameRateAndBitRateInfos FrameRateAndBitRateInfosInDescribeLiveStreamFrameAndBitRateByDomain `json:"FrameRateAndBitRateInfos" xml:"FrameRateAndBitRateInfos"`
}

// CreateDescribeLiveStreamFrameAndBitRateByDomainRequest creates a request to invoke DescribeLiveStreamFrameAndBitRateByDomain API
func CreateDescribeLiveStreamFrameAndBitRateByDomainRequest() (request *DescribeLiveStreamFrameAndBitRateByDomainRequest) {
	request = &DescribeLiveStreamFrameAndBitRateByDomainRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeLiveStreamFrameAndBitRateByDomain", "", "")
	return
}

// CreateDescribeLiveStreamFrameAndBitRateByDomainResponse creates a response to parse from DescribeLiveStreamFrameAndBitRateByDomain response
func CreateDescribeLiveStreamFrameAndBitRateByDomainResponse() (response *DescribeLiveStreamFrameAndBitRateByDomainResponse) {
	response = &DescribeLiveStreamFrameAndBitRateByDomainResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
