# \DoubleWDServiceAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DoubleWDServiceHealthCheck**](DoubleWDServiceAPI.md#DoubleWDServiceHealthCheck) | **Get** /v4/double/health | 



## DoubleWDServiceHealthCheck

> WithdrawalIfaceCommonResponse DoubleWDServiceHealthCheck(ctx).Id(id).Execute()



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
	resp, r, err := apiClient.DoubleWDServiceAPI.DoubleWDServiceHealthCheck(context.Background()).Id(id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DoubleWDServiceAPI.DoubleWDServiceHealthCheck``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DoubleWDServiceHealthCheck`: WithdrawalIfaceCommonResponse
	fmt.Fprintf(os.Stdout, "Response from `DoubleWDServiceAPI.DoubleWDServiceHealthCheck`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDoubleWDServiceHealthCheckRequest struct via the builder pattern


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

