# \ItemsHealthApi

All URIs are relative to *https://api.mercadolibre.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ItemsIdHealthActionsGet**](ItemsHealthApi.md#ItemsIdHealthActionsGet) | **Get** /items/{id}/health/actions | Return item health actions by id.
[**ItemsIdHealthGet**](ItemsHealthApi.md#ItemsIdHealthGet) | **Get** /items/{id}/health | Return health by id.
[**SitesSiteIdHealthLevelsGet**](ItemsHealthApi.md#SitesSiteIdHealthLevelsGet) | **Get** /sites/{site_id}/health_levels | Return health levels.



## ItemsIdHealthActionsGet

> ItemsIdHealthActionsGet(ctx, id, accessToken)

Return item health actions by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**|  | 
**accessToken** | **string**|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ItemsIdHealthGet

> ItemsIdHealthGet(ctx, id, accessToken)

Return health by id.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**|  | 
**accessToken** | **string**|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SitesSiteIdHealthLevelsGet

> SitesSiteIdHealthLevelsGet(ctx, siteId)

Return health levels.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**siteId** | **string**|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

