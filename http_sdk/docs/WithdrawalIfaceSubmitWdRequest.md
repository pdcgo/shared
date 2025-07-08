# WithdrawalIfaceSubmitWdRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TeamId** | Pointer to **string** |  | [optional] 
**MpId** | Pointer to **string** |  | [optional] 
**Source** | Pointer to [**WithdrawalIfaceImporterSource**](WithdrawalIfaceImporterSource.md) |  | [optional] [default to SOURCE_UNKNOWN]
**MpType** | Pointer to [**WithdrawalIfaceOrderMpType**](WithdrawalIfaceOrderMpType.md) |  | [optional] [default to CUSTOM]
**ResourceUri** | Pointer to **string** |  | [optional] 

## Methods

### NewWithdrawalIfaceSubmitWdRequest

`func NewWithdrawalIfaceSubmitWdRequest() *WithdrawalIfaceSubmitWdRequest`

NewWithdrawalIfaceSubmitWdRequest instantiates a new WithdrawalIfaceSubmitWdRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWithdrawalIfaceSubmitWdRequestWithDefaults

`func NewWithdrawalIfaceSubmitWdRequestWithDefaults() *WithdrawalIfaceSubmitWdRequest`

NewWithdrawalIfaceSubmitWdRequestWithDefaults instantiates a new WithdrawalIfaceSubmitWdRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTeamId

`func (o *WithdrawalIfaceSubmitWdRequest) GetTeamId() string`

GetTeamId returns the TeamId field if non-nil, zero value otherwise.

### GetTeamIdOk

`func (o *WithdrawalIfaceSubmitWdRequest) GetTeamIdOk() (*string, bool)`

GetTeamIdOk returns a tuple with the TeamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTeamId

`func (o *WithdrawalIfaceSubmitWdRequest) SetTeamId(v string)`

SetTeamId sets TeamId field to given value.

### HasTeamId

`func (o *WithdrawalIfaceSubmitWdRequest) HasTeamId() bool`

HasTeamId returns a boolean if a field has been set.

### GetMpId

`func (o *WithdrawalIfaceSubmitWdRequest) GetMpId() string`

GetMpId returns the MpId field if non-nil, zero value otherwise.

### GetMpIdOk

`func (o *WithdrawalIfaceSubmitWdRequest) GetMpIdOk() (*string, bool)`

GetMpIdOk returns a tuple with the MpId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMpId

`func (o *WithdrawalIfaceSubmitWdRequest) SetMpId(v string)`

SetMpId sets MpId field to given value.

### HasMpId

`func (o *WithdrawalIfaceSubmitWdRequest) HasMpId() bool`

HasMpId returns a boolean if a field has been set.

### GetSource

`func (o *WithdrawalIfaceSubmitWdRequest) GetSource() WithdrawalIfaceImporterSource`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *WithdrawalIfaceSubmitWdRequest) GetSourceOk() (*WithdrawalIfaceImporterSource, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *WithdrawalIfaceSubmitWdRequest) SetSource(v WithdrawalIfaceImporterSource)`

SetSource sets Source field to given value.

### HasSource

`func (o *WithdrawalIfaceSubmitWdRequest) HasSource() bool`

HasSource returns a boolean if a field has been set.

### GetMpType

`func (o *WithdrawalIfaceSubmitWdRequest) GetMpType() WithdrawalIfaceOrderMpType`

GetMpType returns the MpType field if non-nil, zero value otherwise.

### GetMpTypeOk

`func (o *WithdrawalIfaceSubmitWdRequest) GetMpTypeOk() (*WithdrawalIfaceOrderMpType, bool)`

GetMpTypeOk returns a tuple with the MpType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMpType

`func (o *WithdrawalIfaceSubmitWdRequest) SetMpType(v WithdrawalIfaceOrderMpType)`

SetMpType sets MpType field to given value.

### HasMpType

`func (o *WithdrawalIfaceSubmitWdRequest) HasMpType() bool`

HasMpType returns a boolean if a field has been set.

### GetResourceUri

`func (o *WithdrawalIfaceSubmitWdRequest) GetResourceUri() string`

GetResourceUri returns the ResourceUri field if non-nil, zero value otherwise.

### GetResourceUriOk

`func (o *WithdrawalIfaceSubmitWdRequest) GetResourceUriOk() (*string, bool)`

GetResourceUriOk returns a tuple with the ResourceUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceUri

`func (o *WithdrawalIfaceSubmitWdRequest) SetResourceUri(v string)`

SetResourceUri sets ResourceUri field to given value.

### HasResourceUri

`func (o *WithdrawalIfaceSubmitWdRequest) HasResourceUri() bool`

HasResourceUri returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


