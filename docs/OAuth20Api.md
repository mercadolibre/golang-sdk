# \OAuth20Api

All URIs are relative to *https://api.mercadolibre.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Auth**](OAuth20Api.md#Auth) | **Get** /authorization | Authentication Endpoint
[**GetToken**](OAuth20Api.md#GetToken) | **Post** /oauth/token | Request Access Token



## Auth

> Auth(ctx, responseType, clientId, redirectUri)

Authentication Endpoint

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**responseType** | **string**|  | [default to code]
**clientId** | **string**|  | 
**redirectUri** | **string**|  | 

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


## GetToken

> GetToken(ctx, optional)

Request Access Token

Partner makes a request to the token endpoint by adding the following parameters described below

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetTokenOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a GetTokenOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **grantType** | **optional.String**|  | 
 **clientId** | **optional.String**|  | 
 **clientSecret** | **optional.String**|  | 
 **redirectUri** | **optional.String**|  | 
 **code** | **optional.String**|  | 
 **refreshToken** | **optional.String**|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

