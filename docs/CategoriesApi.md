# \CategoriesApi

All URIs are relative to *https://api.mercadolibre.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CategoriesCategoryIdGet**](CategoriesApi.md#CategoriesCategoryIdGet) | **Get** /categories/{category_id} | Return by category.
[**SitesSiteIdCategoriesGet**](CategoriesApi.md#SitesSiteIdCategoriesGet) | **Get** /sites/{site_id}/categories | Return a categories by site.
[**SitesSiteIdDomainDiscoverySearchGet**](CategoriesApi.md#SitesSiteIdDomainDiscoverySearchGet) | **Get** /sites/{site_id}/domain_discovery/search | Predictor



## CategoriesCategoryIdGet

> CategoriesCategoryIdGet(ctx, categoryId)

Return by category.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**categoryId** | **string**|  | 

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


## SitesSiteIdCategoriesGet

> SitesSiteIdCategoriesGet(ctx, siteId)

Return a categories by site.

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


## SitesSiteIdDomainDiscoverySearchGet

> SitesSiteIdDomainDiscoverySearchGet(ctx, siteId, q, limit)

Predictor

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**siteId** | **string**|  | 
**q** | **string**|  | 
**limit** | **string**|  | 

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

