# \WithdrawalServiceAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**WithdrawalServiceGetTaskList**](WithdrawalServiceAPI.md#WithdrawalServiceGetTaskList) | **Get** /v4/withdrawal/task/list | 
[**WithdrawalServiceHealthCheck**](WithdrawalServiceAPI.md#WithdrawalServiceHealthCheck) | **Get** /v4/withdrawal/health | 
[**WithdrawalServiceRun**](WithdrawalServiceAPI.md#WithdrawalServiceRun) | **Get** /v4/withdrawal/run | 
[**WithdrawalServiceStop**](WithdrawalServiceAPI.md#WithdrawalServiceStop) | **Get** /v4/withdrawal/stop | 
[**WithdrawalServiceSubmitWithdrawal**](WithdrawalServiceAPI.md#WithdrawalServiceSubmitWithdrawal) | **Post** /v4/withdrawal/task/submit | 



## WithdrawalServiceGetTaskList

> WithdrawalIfaceTaskListResponse WithdrawalServiceGetTaskList(ctx).TeamId(teamId).Status(status).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	teamId := "teamId_example" // string |  (optional)
	status := "status_example" // string |  (optional) (default to "TASK_UNKNOWN")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WithdrawalServiceAPI.WithdrawalServiceGetTaskList(context.Background()).TeamId(teamId).Status(status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WithdrawalServiceAPI.WithdrawalServiceGetTaskList``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WithdrawalServiceGetTaskList`: WithdrawalIfaceTaskListResponse
	fmt.Fprintf(os.Stdout, "Response from `WithdrawalServiceAPI.WithdrawalServiceGetTaskList`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWithdrawalServiceGetTaskListRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **teamId** | **string** |  | 
 **status** | **string** |  | [default to &quot;TASK_UNKNOWN&quot;]

### Return type

[**WithdrawalIfaceTaskListResponse**](WithdrawalIfaceTaskListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WithdrawalServiceHealthCheck

> WithdrawalIfaceCommonResponse WithdrawalServiceHealthCheck(ctx).Id(id).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WithdrawalServiceAPI.WithdrawalServiceHealthCheck(context.Background()).Id(id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WithdrawalServiceAPI.WithdrawalServiceHealthCheck``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WithdrawalServiceHealthCheck`: WithdrawalIfaceCommonResponse
	fmt.Fprintf(os.Stdout, "Response from `WithdrawalServiceAPI.WithdrawalServiceHealthCheck`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWithdrawalServiceHealthCheckRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string** |  | 

### Return type

[**WithdrawalIfaceCommonResponse**](WithdrawalIfaceCommonResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WithdrawalServiceRun

> WithdrawalIfaceCommonResponse WithdrawalServiceRun(ctx).Id(id).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WithdrawalServiceAPI.WithdrawalServiceRun(context.Background()).Id(id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WithdrawalServiceAPI.WithdrawalServiceRun``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WithdrawalServiceRun`: WithdrawalIfaceCommonResponse
	fmt.Fprintf(os.Stdout, "Response from `WithdrawalServiceAPI.WithdrawalServiceRun`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWithdrawalServiceRunRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string** |  | 

### Return type

[**WithdrawalIfaceCommonResponse**](WithdrawalIfaceCommonResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WithdrawalServiceStop

> WithdrawalIfaceCommonResponse WithdrawalServiceStop(ctx).Id(id).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WithdrawalServiceAPI.WithdrawalServiceStop(context.Background()).Id(id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WithdrawalServiceAPI.WithdrawalServiceStop``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WithdrawalServiceStop`: WithdrawalIfaceCommonResponse
	fmt.Fprintf(os.Stdout, "Response from `WithdrawalServiceAPI.WithdrawalServiceStop`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWithdrawalServiceStopRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string** |  | 

### Return type

[**WithdrawalIfaceCommonResponse**](WithdrawalIfaceCommonResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WithdrawalServiceSubmitWithdrawal

> WithdrawalIfaceCommonResponse WithdrawalServiceSubmitWithdrawal(ctx).Body(body).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	body := *openapiclient.NewWithdrawalIfaceSubmitWdRequest() // WithdrawalIfaceSubmitWdRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WithdrawalServiceAPI.WithdrawalServiceSubmitWithdrawal(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WithdrawalServiceAPI.WithdrawalServiceSubmitWithdrawal``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WithdrawalServiceSubmitWithdrawal`: WithdrawalIfaceCommonResponse
	fmt.Fprintf(os.Stdout, "Response from `WithdrawalServiceAPI.WithdrawalServiceSubmitWithdrawal`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWithdrawalServiceSubmitWithdrawalRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WithdrawalIfaceSubmitWdRequest**](WithdrawalIfaceSubmitWdRequest.md) |  | 

### Return type

[**WithdrawalIfaceCommonResponse**](WithdrawalIfaceCommonResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

