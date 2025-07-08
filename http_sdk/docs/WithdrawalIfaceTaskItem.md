# WithdrawalIfaceTaskItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TeamId** | Pointer to **string** |  | [optional] 
**MpId** | Pointer to **string** |  | [optional] 
**Filename** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**WithdrawalIfaceTaskStatus**](WithdrawalIfaceTaskStatus.md) |  | [optional] [default to TASK_UNKNOWN]
**Source** | Pointer to [**WithdrawalIfaceImporterSource**](WithdrawalIfaceImporterSource.md) |  | [optional] [default to SOURCE_UNKNOWN]
**MpType** | Pointer to [**WithdrawalIfaceOrderMpType**](WithdrawalIfaceOrderMpType.md) |  | [optional] [default to CUSTOM]
**ResourceUri** | Pointer to **string** |  | [optional] 
**ErrMessage** | Pointer to **string** |  | [optional] 
**IsErr** | Pointer to **bool** |  | [optional] 

## Methods

### NewWithdrawalIfaceTaskItem

`func NewWithdrawalIfaceTaskItem() *WithdrawalIfaceTaskItem`

NewWithdrawalIfaceTaskItem instantiates a new WithdrawalIfaceTaskItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWithdrawalIfaceTaskItemWithDefaults

`func NewWithdrawalIfaceTaskItemWithDefaults() *WithdrawalIfaceTaskItem`

NewWithdrawalIfaceTaskItemWithDefaults instantiates a new WithdrawalIfaceTaskItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTeamId

`func (o *WithdrawalIfaceTaskItem) GetTeamId() string`

GetTeamId returns the TeamId field if non-nil, zero value otherwise.

### GetTeamIdOk

`func (o *WithdrawalIfaceTaskItem) GetTeamIdOk() (*string, bool)`

GetTeamIdOk returns a tuple with the TeamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTeamId

`func (o *WithdrawalIfaceTaskItem) SetTeamId(v string)`

SetTeamId sets TeamId field to given value.

### HasTeamId

`func (o *WithdrawalIfaceTaskItem) HasTeamId() bool`

HasTeamId returns a boolean if a field has been set.

### GetMpId

`func (o *WithdrawalIfaceTaskItem) GetMpId() string`

GetMpId returns the MpId field if non-nil, zero value otherwise.

### GetMpIdOk

`func (o *WithdrawalIfaceTaskItem) GetMpIdOk() (*string, bool)`

GetMpIdOk returns a tuple with the MpId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMpId

`func (o *WithdrawalIfaceTaskItem) SetMpId(v string)`

SetMpId sets MpId field to given value.

### HasMpId

`func (o *WithdrawalIfaceTaskItem) HasMpId() bool`

HasMpId returns a boolean if a field has been set.

### GetFilename

`func (o *WithdrawalIfaceTaskItem) GetFilename() string`

GetFilename returns the Filename field if non-nil, zero value otherwise.

### GetFilenameOk

`func (o *WithdrawalIfaceTaskItem) GetFilenameOk() (*string, bool)`

GetFilenameOk returns a tuple with the Filename field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilename

`func (o *WithdrawalIfaceTaskItem) SetFilename(v string)`

SetFilename sets Filename field to given value.

### HasFilename

`func (o *WithdrawalIfaceTaskItem) HasFilename() bool`

HasFilename returns a boolean if a field has been set.

### GetStatus

`func (o *WithdrawalIfaceTaskItem) GetStatus() WithdrawalIfaceTaskStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *WithdrawalIfaceTaskItem) GetStatusOk() (*WithdrawalIfaceTaskStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *WithdrawalIfaceTaskItem) SetStatus(v WithdrawalIfaceTaskStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *WithdrawalIfaceTaskItem) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetSource

`func (o *WithdrawalIfaceTaskItem) GetSource() WithdrawalIfaceImporterSource`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *WithdrawalIfaceTaskItem) GetSourceOk() (*WithdrawalIfaceImporterSource, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *WithdrawalIfaceTaskItem) SetSource(v WithdrawalIfaceImporterSource)`

SetSource sets Source field to given value.

### HasSource

`func (o *WithdrawalIfaceTaskItem) HasSource() bool`

HasSource returns a boolean if a field has been set.

### GetMpType

`func (o *WithdrawalIfaceTaskItem) GetMpType() WithdrawalIfaceOrderMpType`

GetMpType returns the MpType field if non-nil, zero value otherwise.

### GetMpTypeOk

`func (o *WithdrawalIfaceTaskItem) GetMpTypeOk() (*WithdrawalIfaceOrderMpType, bool)`

GetMpTypeOk returns a tuple with the MpType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMpType

`func (o *WithdrawalIfaceTaskItem) SetMpType(v WithdrawalIfaceOrderMpType)`

SetMpType sets MpType field to given value.

### HasMpType

`func (o *WithdrawalIfaceTaskItem) HasMpType() bool`

HasMpType returns a boolean if a field has been set.

### GetResourceUri

`func (o *WithdrawalIfaceTaskItem) GetResourceUri() string`

GetResourceUri returns the ResourceUri field if non-nil, zero value otherwise.

### GetResourceUriOk

`func (o *WithdrawalIfaceTaskItem) GetResourceUriOk() (*string, bool)`

GetResourceUriOk returns a tuple with the ResourceUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceUri

`func (o *WithdrawalIfaceTaskItem) SetResourceUri(v string)`

SetResourceUri sets ResourceUri field to given value.

### HasResourceUri

`func (o *WithdrawalIfaceTaskItem) HasResourceUri() bool`

HasResourceUri returns a boolean if a field has been set.

### GetErrMessage

`func (o *WithdrawalIfaceTaskItem) GetErrMessage() string`

GetErrMessage returns the ErrMessage field if non-nil, zero value otherwise.

### GetErrMessageOk

`func (o *WithdrawalIfaceTaskItem) GetErrMessageOk() (*string, bool)`

GetErrMessageOk returns a tuple with the ErrMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrMessage

`func (o *WithdrawalIfaceTaskItem) SetErrMessage(v string)`

SetErrMessage sets ErrMessage field to given value.

### HasErrMessage

`func (o *WithdrawalIfaceTaskItem) HasErrMessage() bool`

HasErrMessage returns a boolean if a field has been set.

### GetIsErr

`func (o *WithdrawalIfaceTaskItem) GetIsErr() bool`

GetIsErr returns the IsErr field if non-nil, zero value otherwise.

### GetIsErrOk

`func (o *WithdrawalIfaceTaskItem) GetIsErrOk() (*bool, bool)`

GetIsErrOk returns a tuple with the IsErr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsErr

`func (o *WithdrawalIfaceTaskItem) SetIsErr(v bool)`

SetIsErr sets IsErr field to given value.

### HasIsErr

`func (o *WithdrawalIfaceTaskItem) HasIsErr() bool`

HasIsErr returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


