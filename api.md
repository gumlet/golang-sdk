# Gumlet Go API

Complete reference of every operation, grouped by resource. See [the README](./README.md) for usage and configuration.

## Contents

- [`VideoAssets`](#videoassets)
  - [Create Asset](#create-asset)
  - [Create Asset Direct Upload](#create-asset-direct-upload)
  - [Asset Details](#asset-details)
  - [Delete Asset](#delete-asset)
  - [Update Asset](#update-asset)
  - [Update thumbnail from video](#update-thumbnail-from-video)
  - [Update thumbnail via upload](#update-thumbnail-via-upload)
  - [Create/Update Video Asset Chapters](#createupdate-video-asset-chapters)
  - [List Assets (Deprecated)](#list-assets-deprecated)
  - [List Assets](#list-assets)
  - [Bulk Delete](#bulk-delete)
  - [Bulk Tag](#bulk-tag)
  - [Asset Analytics](#asset-analytics)
  - [Create Asset Direct Upload](#create-asset-direct-upload-1)
  - [Upload a thumbnail file](#upload-a-thumbnail-file)
- [`SubtitleUpload`](#subtitleupload)
  - [Upload Subtitles](#upload-subtitles)
  - [Complete Subtitle Upload](#complete-subtitle-upload)
- [`AudioUpload`](#audioupload)
  - [Add Audio](#add-audio)
  - [Complete Audio Upload](#complete-audio-upload)
- [`VideoUsageAnalytics`](#videousageanalytics)
  - [Video Usage Analytics](#video-usage-analytics)
  - [Top Streamed Assets](#top-streamed-assets)
- [`MultipartUpload`](#multipartupload)
  - [Get Part Upload URL](#get-part-upload-url)
  - [Complete Multipart Upload](#complete-multipart-upload)
  - [Abort Upload](#abort-upload)
  - [List Uploads](#list-uploads)
- [`VideoProfiles`](#videoprofiles)
  - [Create Profile](#create-profile)
  - [List Profiles](#list-profiles)
  - [Update Profile](#update-profile)
  - [Get Profile](#get-profile)
  - [Delete Profile](#delete-profile)
- [`VideoPlaylists`](#videoplaylists)
  - [Create Playlist](#create-playlist)
  - [Get all playlists](#get-all-playlists)
  - [Add asset to playlist](#add-asset-to-playlist)
  - [Remove asset from playlist](#remove-asset-from-playlist)
  - [Update Playlist](#update-playlist)
  - [Delete Playlist](#delete-playlist)
  - [Get playlist assets](#get-playlist-assets)
  - [Arrange Videos In Playlist](#arrange-videos-in-playlist)
- [`WebhookEndpoints`](#webhookendpoints)
  - [Create Webhook](#create-webhook)
  - [List Webhooks](#list-webhooks)
  - [Update Webhook](#update-webhook)
  - [Delete Webhook](#delete-webhook)
  - [Get History](#get-history)
- [`ImageSources`](#imagesources)
  - [Create Source](#create-source)
  - [List Sources](#list-sources)
  - [Get Image Source](#get-image-source)
  - [Update Source](#update-source)
  - [Delete Source](#delete-source)
  - [Purge Cache](#purge-cache)
  - [Purge Image Cache](#purge-image-cache)
- [`ImageUsageAnalytics`](#imageusageanalytics)
  - [Image Usage Analytics](#image-usage-analytics)
- [`LiveStreamAssets`](#livestreamassets)
  - [Create Live Asset](#create-live-asset)
  - [Update Live Asset](#update-live-asset)
  - [Get Live Asset Status](#get-live-asset-status)
  - [Delete Live Asset](#delete-live-asset)
  - [Complete Live Stream](#complete-live-stream)
  - [Filter Live Assets](#filter-live-assets)
  - [`start`](#start)
  - [Upload Live Thumbnails](#upload-live-thumbnails)
  - [Get Live Asset Status History](#get-live-asset-status-history)
- [`RecycleBin`](#recyclebin)
  - [Recover Deleted Asset](#recover-deleted-asset)
  - [List Recycle Bin](#list-recycle-bin)
- [`VideoWorkspaces`](#videoworkspaces)
  - [List Workspaces](#list-workspaces)
  - [Create Workspace](#create-workspace)
  - [Update Workspace](#update-workspace)
  - [Get Workspace](#get-workspace)
  - [Delete Workspace](#delete-workspace)
- [`Folders`](#folders)
  - [Create Folder](#create-folder)
  - [List Folders](#list-folders)
  - [Get Folder](#get-folder)
  - [Update Folder](#update-folder)
  - [Delete Folder](#delete-folder)
  - [Remove Assets From Folder](#remove-assets-from-folder)
- [`ChannelViewers`](#channelviewers)
  - [Invite Channel Viewers](#invite-channel-viewers)
  - [Remove Channel Viewers](#remove-channel-viewers)
  - [Invite Channel Viewers via CSV](#invite-channel-viewers-via-csv)
  - [List Subscribers](#list-subscribers)
- [`VideoAnalytics`](#videoanalytics)
  - [Viewer Analytics](#viewer-analytics)
  - [Breakdown Data](#breakdown-data)
  - [Aggregated Data](#aggregated-data)
- [`OrganizationData`](#organizationdata)
  - [Get Organization Details](#get-organization-details)
- [`UserData`](#userdata)
  - [Get User](#get-user)
- [`AuditLogs`](#auditlogs)
  - [Fetch Audit Logs](#fetch-audit-logs)
- [`Billing`](#billing)
  - [List Invoices](#list-invoices)
  - [Get Billing Details](#get-billing-details)
  - [Update Billing Details](#update-billing-details)
  - [Upcoming Invoice](#upcoming-invoice)
- [`LiveStreamWorkspaces`](#livestreamworkspaces)
  - [List Workspaces](#list-workspaces-1)
  - [Create Workspace](#create-workspace-1)
  - [Update Workspace](#update-workspace-1)
  - [Delete Workspace](#delete-workspace-1)
- [`LiveStreamAnalytics`](#livestreamanalytics)
  - [Usage Analytics](#usage-analytics)
- [`GlobalSearch`](#globalsearch)
  - [Global Search](#global-search)

## Setup

```go
import (
	"context"
	"fmt"

	sdk "github.com/gumlet/golang-sdk"
)

client := sdk.NewClient()
```

## `VideoAssets`

Upload, manage, and retrieve video assets, including thumbnails, subtitles, audio tracks, and chapters.

### Create Asset

An asset refers to media content/video that is processed, stored, and delivered through Gumlet. This endpoint creates an asset allowing users to ingest media content into the Gumlet system for processing and delivery.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetNewParams`](./videoasset.go) |
| Response | [`VideoAssetNewResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.New(context.Background(), sdk.VideoAssetNewParams{
	Input:       sdk.F[string]("http://devimages.apple.com/iphone/samples/bipbop/bipbopall.m3u8"),
	WorkspaceID: sdk.F[string]("<your workspace id>"),
	Format:      sdk.F[sdk.VideoAssetNewParamsFormat](sdk.VideoAssetNewParamsFormat("ABR")),
	Title:       sdk.F[string]("Example Title"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Create Asset Direct Upload

This endpoint creates a video asset allowing upload of a video from the local file system and ingest media content into the Gumlet system for processing and delivery. Body parameters are the same as the Create Asset body parameters except for the `input` parameter, which this endpoint does not take. A successful response will be returned with `upload_url` field. You can make `PUT` request to that URL to upload video. To upload video using `upload_url` refer to [Direct Upload](https://docs.gumlet.com/developers/vod-developers-section/direct-upload#2.-use-the-url-to-upload-a-file) guide.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetUploadParams`](./videoasset.go) |
| Response | [`VideoAssetUploadResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.Upload(context.Background(), sdk.VideoAssetUploadParams{
	WorkspaceID: sdk.F[string]("646df1c9173a4a2fcac180b4"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Asset Details

This endpoint retrieves the details of an asset that has previously been created.

| Direction | Type |
| --- | --- |
| Response | [`VideoAssetGetDetailsResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.GetDetails(context.Background(), "assetId")
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Delete Asset

This endpoint removes an asset given its unique asset id. The asset will be removed from storage as well, associated URLs will be inaccessible.

```go
err := client.VideoAssets.Delete(context.Background(), "assetId")
if err != nil {
	panic(err)
}
```

### Update Asset

This endpoint allows users to update video asset that has previously been created.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetUpdateParams`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.Update(context.Background(), sdk.VideoAssetUpdateParams{
	AssetID: sdk.F[string]("<YOUR_ASSET_ID>"),
	Title:   sdk.F[string]("Updated Title"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Update thumbnail from video

Select frame from video to use as thumbnail.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetThumbnailSelectParams`](./videoasset.go) |
| Response | [`VideoAssetThumbnailSelectResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.ThumbnailSelect(context.Background(), "assetId", sdk.VideoAssetThumbnailSelectParams{
	FrameAtSecond: sdk.F[float64](2),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Update thumbnail via upload

Use any image file to use as thumbnail. Once you use the API, you will get `upload_url` in the response, and that can be used to upload the image file.

Here is the sample curl request.

```bash
curl --location --request PUT '<upload_url>' \
--data '<YOUR_FILE_PATH>'
```

| Direction | Type |
| --- | --- |
| Response | [`VideoAssetThumbnailUploadResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.ThumbnailUpload(context.Background(), "assetId")
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Create/Update Video Asset Chapters

This endpoint will create/update video asset chapters.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetNewUpdateChapterParams`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.NewUpdateChapter(context.Background(), "assetId", sdk.VideoAssetNewUpdateChapterParams{
	Chapters: sdk.F[[]sdk.VideoAssetNewUpdateChapterParamsChapter]([]sdk.VideoAssetNewUpdateChapterParamsChapter{sdk.VideoAssetNewUpdateChapterParamsChapter{
		Label:     sdk.F[string]("Chapter 1"),
		StartTime: sdk.F[int64](0),
	}, sdk.VideoAssetNewUpdateChapterParamsChapter{
		Label:     sdk.F[string]("Chapter 2"),
		StartTime: sdk.F[int64](10),
	}}),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### List Assets (Deprecated)

[Deprecated] Combined folders + assets list for a workspace. Use `GET /video/assets/list/{workspace_id}` for assets and `GET /video/workspaces/{workspace_id}/folders` for folders instead.

List folders and assets for a workspace in a single response. Use `parent_id` to browse a specific folder, or filters like `title`, `status`, and `playlist_id` to search assets.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetListWorkspaceFoldersAndDeprecatedParams`](./videoasset.go) |
| Response | [`VideoAssetListWorkspaceFoldersAndDeprecatedResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.ListWorkspaceFoldersAndDeprecated(context.Background(), "workspaceId", sdk.VideoAssetListWorkspaceFoldersAndDeprecatedParams{
	Offset:      sdk.F[int64](0),
	SignedToken: sdk.F[sdk.VideoAssetListWorkspaceFoldersAndDeprecatedParamsSignedToken](sdk.VideoAssetListWorkspaceFoldersAndDeprecatedParamsSignedToken("false")),
	Size:        sdk.F[int64](20),
	Type:        sdk.F[sdk.VideoAssetListWorkspaceFoldersAndDeprecatedParamsType](sdk.VideoAssetListWorkspaceFoldersAndDeprecatedParamsType("all")),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### List Assets

List assets in a video workspace (local Mongo list; not Atlas Search).

Use `page_number` and `page_size` for page-based pagination. Use `next_page_token` for cursor-based pagination where supported.

`next_page_token` is supported only with the **default** sort or `sortBy=created_at` (and without a `title` search). It is **not** supported with other sort fields (for example `title`, `duration`, `status`, `playlist_order`, `custom`) or when `title` is set; in those cases use `page_number` / `page_size` only.

For folders, use `GET /video/workspaces/{workspace_id}/folders` and combine results in your client. The combined endpoint `GET /video/workspaces/{workspace_id}/list` is deprecated.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetListParams`](./videoasset.go) |
| Response | [`VideoAssetListResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.List(context.Background(), "workspaceId", sdk.VideoAssetListParams{
	OrderBy:     sdk.F[sdk.VideoAssetListParamsOrderBy](sdk.VideoAssetListParamsOrderBy("desc")),
	PageNumber:  sdk.F[int64](1),
	PageSize:    sdk.F[int64](20),
	SignedToken: sdk.F[sdk.VideoAssetListParamsSignedToken](sdk.VideoAssetListParamsSignedToken("false")),
	SortBy:      sdk.F[sdk.VideoAssetListParamsSortBy](sdk.VideoAssetListParamsSortBy("created_at")),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Bulk Delete

Delete multiple VOD assets at once.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetDeleteManyParams`](./videoasset.go) |
| Response | [`VideoAssetDeleteManyResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.DeleteMany(context.Background(), sdk.VideoAssetDeleteManyParams{
	AssetList:   sdk.F[[]string]([]string{"64249a8858fd3a208b987702", "64784bae843b155b829bbf84"}),
	WorkspaceID: sdk.F[string]("60bd2ba353ff754d28179ee6"),
	SourceID:    sdk.F[string]("60bd2ba353ff754d28179ee6"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Bulk Tag

Add / remove tags from multiple assets at once.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetTagManyParams`](./videoasset.go) |
| Response | [`VideoAssetTagManyResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.TagMany(context.Background(), sdk.VideoAssetTagManyParams{
	AddTags:     sdk.F[[]string]([]string{"tag-1"}),
	AssetList:   sdk.F[[]string]([]string{"6221db301c8b821b0519fba0", "61e8f2726ec832ab2ac4fa6e"}),
	RemoveTags:  sdk.F[[]string]([]string{"playlist-1"}),
	WorkspaceID: sdk.F[string]("60bd2ba353ff754d28179ee6"),
	SourceID:    sdk.F[string]("60bd2ba353ff754d28179ee6"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Asset Analytics

Get video analytics for a single asset.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetAnalyticsParams`](./videoasset.go) |
| Response | [`VideoAssetAnalyticsResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.Analytics(context.Background(), "assetId", sdk.VideoAssetAnalyticsParams{
	DateRange: sdk.F[sdk.VideoAssetAnalyticsParamsDateRange](sdk.VideoAssetAnalyticsParamsDateRange{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	GroupBy: sdk.F[sdk.VideoAssetAnalyticsParamsGroupBy](sdk.VideoAssetAnalyticsParamsGroupBy("daily")),
	Metrics: sdk.F[[]sdk.VideoAssetAnalyticsParamsMetric]([]sdk.VideoAssetAnalyticsParamsMetric{"impressions"}),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Create Asset Direct Upload

This endpoint creates a video asset allowing upload of a video from the local file system and ingest media content into the Gumlet system for processing and delivery. Body parameters are the same as the Create Asset body parameters except for the `input` parameter, which this endpoint does not take. A successful response will be returned with `upload_url` field. You can make `PUT` request to that URL to upload video. To upload video using `upload_url` refer to [Direct Upload](https://docs.gumlet.com/developers/vod-developers-section/direct-upload#2.-use-the-url-to-upload-a-file) guide.

| Direction | Type |
| --- | --- |
| Request | [`VideoAssetUploadFileParams`](./videoasset.go) |
| Response | [`VideoAssetUploadFileResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.UploadFile(context.Background(), strings.NewReader("file contents"), sdk.VideoAssetUploadFileParams{
	WorkspaceID: sdk.F[string]("646df1c9173a4a2fcac180b4"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

### Upload a thumbnail file

Requests a thumbnail upload URL for the asset and PUTs the file to it.

| Direction | Type |
| --- | --- |
| Response | [`VideoAssetUploadThumbnailFileResponse`](./videoasset.go) |

```go
videoAsset, err := client.VideoAssets.UploadThumbnailFile(context.Background(), "assetId", strings.NewReader("file contents"))
if err != nil {
	panic(err)
}

fmt.Println(videoAsset)
```

## `SubtitleUpload`

Add subtitles to an existing asset.

### Upload Subtitles

Upload `.srt` or `.vtt` file to the video asset. The response of this API call gives `upload_url` for each language specified. You need to send a `PUT` request of the subtitle files to those URLs. Once that's done, you need to call the subtitle upload complete API. Only after that, Gumlet will add subtitles to asset.

| Direction | Type |
| --- | --- |
| Request | [`SubtitleUploadUploadParams`](./subtitleupload.go) |
| Response | [`SubtitleUploadUploadResponse`](./subtitleupload.go) |

```go
subtitleUpload, err := client.SubtitleUpload.Upload(context.Background(), "assetId", sdk.SubtitleUploadUploadParams{
	LanguageCodes: sdk.F[[]string]([]string{"en"}),
})
if err != nil {
	panic(err)
}

fmt.Println(subtitleUpload)
```

### Complete Subtitle Upload

This API must be called after adding subtitles; the add subtitle call gives you URLs to upload, and you complete a `PUT` request to those URLs. 
Once that is done, calling this initiates the process to actually add the subtitle to the video.

| Direction | Type |
| --- | --- |
| Request | [`SubtitleUploadCompleteParams`](./subtitleupload.go) |

```go
subtitleUpload, err := client.SubtitleUpload.Complete(context.Background(), "assetId", sdk.SubtitleUploadCompleteParams{
	UploadResponses: sdk.F[[]sdk.SubtitleUploadCompleteParamsUploadResponse]([]sdk.SubtitleUploadCompleteParamsUploadResponse{sdk.SubtitleUploadCompleteParamsUploadResponse{
		LanguageCode: sdk.F[string]("en"),
		Uploaded:     sdk.F[bool](true),
	}}),
})
if err != nil {
	panic(err)
}

fmt.Println(subtitleUpload)
```

## `AudioUpload`

Add additional audio tracks to an existing asset.

### Add Audio

Add any audio file to the video asset. 
The response of this API call gives `upload_url` for each language specified. You need to send a `PUT` request of the audio files to those URLs. Once that's done, you need to call the audio upload complete API. Only after that will Gumlet add audio to the asset.

| Direction | Type |
| --- | --- |
| Request | [`AudioUploadUploadParams`](./audioupload.go) |
| Response | [`AudioUploadUploadResponse`](./audioupload.go) |

```go
audioUpload, err := client.AudioUpload.Upload(context.Background(), "assetId", sdk.AudioUploadUploadParams{
	LanguageCodes: sdk.F[[]string]([]string{"en"}),
})
if err != nil {
	panic(err)
}

fmt.Println(audioUpload)
```

### Complete Audio Upload

Finish adding audio tracks to an existing video asset. Call this only after the audio files have been PUT to the upload URLs from the add-audio call. This attaches those tracks to the asset.

| Direction | Type |
| --- | --- |
| Request | [`AudioUploadCompleteParams`](./audioupload.go) |

```go
audioUpload, err := client.AudioUpload.Complete(context.Background(), "assetId", sdk.AudioUploadCompleteParams{
	UploadResponses: sdk.F[[]sdk.AudioUploadCompleteParamsUploadResponse]([]sdk.AudioUploadCompleteParamsUploadResponse{sdk.AudioUploadCompleteParamsUploadResponse{
		Uploaded: sdk.F[bool](true),
	}}),
})
if err != nil {
	panic(err)
}

fmt.Println(audioUpload)
```

## `VideoUsageAnalytics`

Query video analytics and streaming duration usage data.

### Video Usage Analytics

This endpoint gives usage analytics data of your videos. Ex - top assets, bandwidth consumption.

Pagination applies only to `top_assets` and `widget_data` metrics.

| Direction | Type |
| --- | --- |
| Request | [`VideoUsageAnalyticGetParams`](./videousageanalytic.go) |
| Response | [`VideoUsageAnalyticGetResponse`](./videousageanalytic.go) |

```go
videoUsageAnalytic, err := client.VideoUsageAnalytics.Get(context.Background(), sdk.VideoUsageAnalyticGetParams{
	DateRange: sdk.F[sdk.VideoUsageAnalyticGetParamsDateRange](sdk.VideoUsageAnalyticGetParamsDateRange{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	Metrics:        sdk.F[[]sdk.VideoUsageAnalyticGetParamsMetric]([]sdk.VideoUsageAnalyticGetParamsMetric{"bandwidth_consumption", "asset_duration", "storage_unit", "top_assets", "drm_requests"}),
	GroupBy:        sdk.F[sdk.VideoUsageAnalyticGetParamsGroupBy](sdk.VideoUsageAnalyticGetParamsGroupBy("hourly")),
	PageNumber:     sdk.F[int64](1),
	PageSize:       sdk.F[int64](20),
	TopAssetsCount: sdk.F[string]("5"),
	TopAssetsPage:  sdk.F[string]("0"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoUsageAnalytic)
```

### Top Streamed Assets

This endpoint lists top streamed assets in a video collection

| Direction | Type |
| --- | --- |
| Request | [`VideoUsageAnalyticTopAssetsParams`](./videousageanalytic.go) |
| Response | [`VideoUsageAnalyticTopAssetsResponse`](./videousageanalytic.go) |

```go
videoUsageAnalytic, err := client.VideoUsageAnalytics.TopAssets(context.Background(), sdk.VideoUsageAnalyticTopAssetsParams{
	EndAt:    sdk.F[time.Time](time.Now()),
	StartAt:  sdk.F[time.Time](time.Now()),
	Page:     sdk.F[string]("1"),
	PageSize: sdk.F[string]("1000"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoUsageAnalytic)
```

## `MultipartUpload`

Upload large video files in parts and complete the multipart upload.

### Get Part Upload URL

Use this endpoint to retrieve a pre-signed upload URL for the given part number.

| Direction | Type |
| --- | --- |
| Response | [`MultipartUploadGetPartURLResponse`](./multipartupload.go) |

```go
multipartUpload, err := client.MultipartUpload.GetPartURL(context.Background(), "assetId", "partNumber")
if err != nil {
	panic(err)
}

fmt.Println(multipartUpload)
```

### Complete Multipart Upload

Once you upload all parts to S3 bucket via pre-signed URL, use this endpoint to complete the multipart upload.

| Direction | Type |
| --- | --- |
| Request | [`MultipartUploadCompleteParams`](./multipartupload.go) |

```go
multipartUpload, err := client.MultipartUpload.Complete(context.Background(), "assetId", sdk.MultipartUploadCompleteParams{})
if err != nil {
	panic(err)
}

fmt.Println(multipartUpload)
```

### Abort Upload

This call aborts multi-part upload and deletes the already uploaded parts from the storage.

| Direction | Type |
| --- | --- |
| Request | [`MultipartUploadAbortParams`](./multipartupload.go) |

```go
multipartUpload, err := client.MultipartUpload.Abort(context.Background(), "assetId", sdk.MultipartUploadAbortParams{
	Body: map[string]interface{}{},
})
if err != nil {
	panic(err)
}

fmt.Println(multipartUpload)
```

### List Uploads

Lists all parts uploaded so far.

| Direction | Type |
| --- | --- |
| Request | [`MultipartUploadListParams`](./multipartupload.go) |
| Response | [`MultipartUploadListResponse`](./multipartupload.go) |

```go
multipartUpload, err := client.MultipartUpload.List(context.Background(), "assetId", sdk.MultipartUploadListParams{
	Body: map[string]interface{}{},
})
if err != nil {
	panic(err)
}

fmt.Println(multipartUpload)
```

## `VideoProfiles`

Create and manage encoding/output profiles for video assets.

### Create Profile

Gumlet provides the functionality of creating multiple video assets using the same set of parameters. A Video profile is a set of parameters that can be referenced/used while creating a video as a single parameter.

| Direction | Type |
| --- | --- |
| Request | [`VideoProfileNewParams`](./videoprofile.go) |
| Response | [`VideoProfileNewResponse`](./videoprofile.go) |

```go
videoProfile, err := client.VideoProfiles.New(context.Background(), sdk.VideoProfileNewParams{
	Format:     sdk.F[sdk.VideoProfileNewParamsFormat](sdk.VideoProfileNewParamsFormat("ABR")),
	Name:       sdk.F[string]("Gumlet-Profile-1"),
	Transcribe: sdk.F[bool](false),
})
if err != nil {
	panic(err)
}

fmt.Println(videoProfile)
```

### List Profiles

This endpoint retrieves the details of all profiles that have previously been created.

Pagination uses `page_number` / `page_size`. Use `next_page_token` for cursor-based pagination.

| Direction | Type |
| --- | --- |
| Request | [`VideoProfileListParams`](./videoprofile.go) |
| Response | [`VideoProfileListResponse`](./videoprofile.go) |

```go
videoProfile, err := client.VideoProfiles.List(context.Background(), sdk.VideoProfileListParams{
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(videoProfile)
```

### Update Profile

Update an existing profile. Settings provided in body parameters will only be updated in the existing profile.

| Direction | Type |
| --- | --- |
| Request | [`VideoProfileUpdateParams`](./videoprofile.go) |
| Response | [`VideoProfileUpdateResponse`](./videoprofile.go) |

```go
videoProfile, err := client.VideoProfiles.Update(context.Background(), "profileId", sdk.VideoProfileUpdateParams{
	ProfileID:  sdk.F[string](""),
	Format:     sdk.F[sdk.VideoProfileUpdateParamsFormat](sdk.VideoProfileUpdateParamsFormat("ABR")),
	Transcribe: sdk.F[bool](false),
	Vc:         sdk.F[[]sdk.VideoProfileUpdateParamsVc]([]sdk.VideoProfileUpdateParamsVc{"libx264"}),
})
if err != nil {
	panic(err)
}

fmt.Println(videoProfile)
```

### Get Profile

This endpoint retrieves the details of a video profile that has previously been created.

| Direction | Type |
| --- | --- |
| Response | [`VideoProfileGetResponse`](./videoprofile.go) |

```go
videoProfile, err := client.VideoProfiles.Get(context.Background(), "profileId")
if err != nil {
	panic(err)
}

fmt.Println(videoProfile)
```

### Delete Profile

This endpoint removes a profile given its unique `profile_id`. The profile will be removed but video assets created using the profile will remain as it is.

```go
videoProfile, err := client.VideoProfiles.Delete(context.Background(), "profileId")
if err != nil {
	panic(err)
}

fmt.Println(videoProfile)
```

## `VideoPlaylists`

Create and manage playlists, and control which assets belong to them.

### Create Playlist

Create new playlist inside video wprkspace

| Direction | Type |
| --- | --- |
| Request | [`VideoPlaylistNewParams`](./videoplaylist.go) |
| Response | [`VideoPlaylistNewResponse`](./videoplaylist.go) |

```go
videoPlaylist, err := client.VideoPlaylists.New(context.Background(), sdk.VideoPlaylistNewParams{
	Title:        sdk.F[string]("Playlist-Title"),
	WorkspaceID:  sdk.F[string]("{{video-source-id}}"),
	CollectionID: sdk.F[string]("{{video-source-id}}"),
	Description:  sdk.F[string]("This is description for playlist."),
})
if err != nil {
	panic(err)
}

fmt.Println(videoPlaylist)
```

### Get all playlists

Get all playlists for given workspace

| Direction | Type |
| --- | --- |
| Request | [`VideoPlaylistListAllParams`](./videoplaylist.go) |
| Response | [`[]VideoPlaylistListAllResponse`](./videoplaylist.go) |

```go
videoPlaylist, err := client.VideoPlaylists.ListAll(context.Background(), sdk.VideoPlaylistListAllParams{})
if err != nil {
	panic(err)
}

fmt.Println(videoPlaylist)
```

### Add asset to playlist

This operation adds a single asset or a list of assets to a playlist.

| Direction | Type |
| --- | --- |
| Request | [`VideoPlaylistNewAssetParams`](./videoplaylist.go) |
| Response | [`VideoPlaylistNewAssetResponse`](./videoplaylist.go) |

```go
videoPlaylist, err := client.VideoPlaylists.NewAsset(context.Background(), "playlistId", sdk.VideoPlaylistNewAssetParams{
	AssetList: sdk.F[[]sdk.VideoPlaylistNewAssetParamsAssetList]([]sdk.VideoPlaylistNewAssetParamsAssetList{sdk.VideoPlaylistNewAssetParamsAssetList{
		AssetID: sdk.F[string]("6508790283e4d60611846790"),
	}, sdk.VideoPlaylistNewAssetParamsAssetList{
		AssetID:  sdk.F[string]("650878f883e4d6061184677d"),
		Position: sdk.F[int64](1),
	}, sdk.VideoPlaylistNewAssetParamsAssetList{
		AssetID: sdk.F[string]("650878de83e4d6061184676a"),
	}, sdk.VideoPlaylistNewAssetParamsAssetList{
		AssetID:  sdk.F[string]("650878d883e4d60611846757"),
		Position: sdk.F[int64](2),
	}, sdk.VideoPlaylistNewAssetParamsAssetList{
		AssetID:  sdk.F[string]("65578dd87eebc22dcdd549a2"),
		Position: sdk.F[int64](3),
	}}),
})
if err != nil {
	panic(err)
}

fmt.Println(videoPlaylist)
```

### Remove asset from playlist

Removed an asset or list of assets from a given playlist.

| Direction | Type |
| --- | --- |
| Request | [`VideoPlaylistDeleteAssetParams`](./videoplaylist.go) |
| Response | [`VideoPlaylistDeleteAssetResponse`](./videoplaylist.go) |

```go
videoPlaylist, err := client.VideoPlaylists.DeleteAsset(context.Background(), "playlistId", sdk.VideoPlaylistDeleteAssetParams{
	DeleteList: sdk.F[[]string]([]string{"6508790783e4d606118467a3"}),
})
if err != nil {
	panic(err)
}

fmt.Println(videoPlaylist)
```

### Update Playlist

This endpoint allows you to update playlist name, channel visibility, or playlist order on a channel page.

| Direction | Type |
| --- | --- |
| Request | [`VideoPlaylistUpdateParams`](./videoplaylist.go) |
| Response | [`VideoPlaylistUpdateResponse`](./videoplaylist.go) |

```go
videoPlaylist, err := client.VideoPlaylists.Update(context.Background(), "playlistId", sdk.VideoPlaylistUpdateParams{
	ChannelVisibility: sdk.F[bool](false),
	Description:       sdk.F[string]("This is updated description"),
	PlayerConfig: sdk.F[sdk.VideoPlaylistUpdateParamsPlayerConfig](sdk.VideoPlaylistUpdateParamsPlayerConfig{
		Preload:                 sdk.F[bool](true),
		Autoplay:                sdk.F[bool](false),
		DisableSeek:             sdk.F[bool](true),
		DisablePlayerControls:   sdk.F[bool](false),
		PoweredByGumletOverlay:  sdk.F[bool](false),
		AllowDrmProtectedVideos: sdk.F[bool](false),
		Loop:                    sdk.F[bool](false),
		PlayerColor:             sdk.F[string]("#6658ea"),
		IncludeSeo:              sdk.F[bool](true),
		SubtitleEnabled:         sdk.F[bool](true),
		PixelTags:               sdk.F[interface{}](map[string]interface{}{}),
		LogoWidth:               sdk.F[int64](51),
		LogoHeight:              sdk.F[int64](100),
		DynamicWatermark:        sdk.F[bool](false),
		WatermarkFontSize:       sdk.F[int64](1),
		WatermarkFontColor:      sdk.F[string]("#ff0000"),
		WatermarkBgColor:        sdk.F[string]("transparent"),
		WatermarkInterval:       sdk.F[int64](1000),
	}),
	Position: sdk.F[int64](6),
	Title:    sdk.F[string]("Playlist-Title-Updated"),
})
if err != nil {
	panic(err)
}

fmt.Println(videoPlaylist)
```

### Delete Playlist

Deletes a playlist by plalist ID.

```go
err := client.VideoPlaylists.Delete(context.Background(), "playlistId")
if err != nil {
	panic(err)
}
```

### Get playlist assets

Get a list of all assets inside playlist. You can choose in which order are assets returned.

| Direction | Type |
| --- | --- |
| Request | [`VideoPlaylistListAssetsParams`](./videoplaylist.go) |
| Response | [`VideoPlaylistListAssetsResponse`](./videoplaylist.go) |

```go
videoPlaylist, err := client.VideoPlaylists.ListAssets(context.Background(), "playlistId", sdk.VideoPlaylistListAssetsParams{
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[string]("10"),
	SortOrder:  sdk.F[int64](1),
})
if err != nil {
	panic(err)
}

fmt.Println(videoPlaylist)
```

### Arrange Videos In Playlist

Reorder videos inside a playlist either by moving a single asset to a position or by sorting the playlist by title or created date.

| Direction | Type |
| --- | --- |
| Request | [`VideoPlaylistReorderAssetParams`](./videoplaylist.go) |
| Response | [`VideoPlaylistReorderAssetResponse`](./videoplaylist.go) |

```go
videoPlaylist, err := client.VideoPlaylists.ReorderAsset(context.Background(), "playlistId", sdk.VideoPlaylistReorderAssetParams{
	Body: sdk.VideoPlaylistReorderAssetParamsBodyVariant0{
		AssetID:       sdk.F[string]("6e82bf783e88be000ab45ed2"),
		PageNumber:    sdk.F[int64](1),
		PageSize:      sdk.F[int64](10),
		AssetPosition: sdk.F[int64](0),
	},
})
if err != nil {
	panic(err)
}

fmt.Println(videoPlaylist)
```

## `WebhookEndpoints`

Create webhook listeners and receive video status, live video status, and product events. Gumlet POSTs JSON to your URL with the `x-gumlet-token` header.

### Create Webhook

Creates a new webhook listener. Gumlet POSTs JSON to `url` for each matching event and sends `secret_token` in the `x-gumlet-token` header. Payload schemas are documented in the webhooks section.

| Direction | Type |
| --- | --- |
| Request | [`WebhookEndpointNewParams`](./webhookendpoint.go) |
| Response | [`WebhookEndpointNewResponse`](./webhookendpoint.go) |

```go
webhookEndpoint, err := client.WebhookEndpoints.New(context.Background(), sdk.WebhookEndpointNewParams{
	SecretToken: sdk.F[string](""),
	Sources:     sdk.F[[]string]([]string{""}),
	Triggers:    sdk.F[[]sdk.WebhookEndpointNewParamsTrigger]([]sdk.WebhookEndpointNewParamsTrigger{"status"}),
	URL:         sdk.F[string](""),
})
if err != nil {
	panic(err)
}

fmt.Println(webhookEndpoint)
```

### List Webhooks

List all webhooks.

| Direction | Type |
| --- | --- |
| Response | [`WebhookEndpointListResponse`](./webhookendpoint.go) |

```go
webhookEndpoint, err := client.WebhookEndpoints.List(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(webhookEndpoint)
```

### Update Webhook

Update a webhook listener.

| Direction | Type |
| --- | --- |
| Request | [`WebhookEndpointUpdateParams`](./webhookendpoint.go) |
| Response | [`WebhookEndpointUpdateResponse`](./webhookendpoint.go) |

```go
webhookEndpoint, err := client.WebhookEndpoints.Update(context.Background(), "webhookId", sdk.WebhookEndpointUpdateParams{})
if err != nil {
	panic(err)
}

fmt.Println(webhookEndpoint)
```

### Delete Webhook

Delete webhook listener endpoint.

```go
webhookEndpoint, err := client.WebhookEndpoints.Delete(context.Background(), "webhookId")
if err != nil {
	panic(err)
}

fmt.Println(webhookEndpoint)
```

### Get History

Get logs history for a given webhook.

Pagination uses `page_number` / `page_size`. Use `next_page_token` for cursor-based pagination.

| Direction | Type |
| --- | --- |
| Request | [`WebhookEndpointHistoryParams`](./webhookendpoint.go) |
| Response | `[]map[string]interface{}` |

```go
webhookEndpoint, err := client.WebhookEndpoints.History(context.Background(), "webhookId", sdk.WebhookEndpointHistoryParams{
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(webhookEndpoint)
```

## `ImageSources`

Manage image sources, view image analytics, and purge the image cache.

### Create Source

This endpoint allows users to create image source.

| Direction | Type |
| --- | --- |
| Request | [`ImageSourceNewParams`](./imagesource.go) |
| Response | [`ImageSourceNewResponse`](./imagesource.go) |

```go
imageSource, err := client.ImageSources.New(context.Background(), sdk.ImageSourceNewParams{
	Namespace: sdk.F[string]("google-demo"),
	Type:      sdk.F[sdk.ImageSourceNewParamsType](sdk.ImageSourceNewParamsType("webfolder")),
	Webfolder: sdk.F[sdk.ImageSourceNewParamsWebfolder](sdk.ImageSourceNewParamsWebfolder{
		BaseURL: sdk.F[string]("https://www.google.com"),
	}),
})
if err != nil {
	panic(err)
}

fmt.Println(imageSource)
```

### List Sources

This endpoint list image sources which are assigned to the user or token.

Pagination uses `page_number` / `page_size`. Use `next_page_token` for cursor-based pagination.

| Direction | Type |
| --- | --- |
| Request | [`ImageSourceListParams`](./imagesource.go) |
| Response | [`ImageSourceListResponse`](./imagesource.go) |

```go
imageSource, err := client.ImageSources.List(context.Background(), sdk.ImageSourceListParams{
	OrderBy:    sdk.F[sdk.ImageSourceListParamsOrderBy](sdk.ImageSourceListParamsOrderBy("asc")),
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(imageSource)
```

### Get Image Source

Get all details about image source.

| Direction | Type |
| --- | --- |
| Response | [`ImageSourceGetResponse`](./imagesource.go) |

```go
imageSource, err := client.ImageSources.Get(context.Background(), "imageSourceId")
if err != nil {
	panic(err)
}

fmt.Println(imageSource)
```

### Update Source

Update settings on an image source that already exists. This does not create a source or purge cached images.

| Direction | Type |
| --- | --- |
| Request | [`ImageSourceUpdateParams`](./imagesource.go) |
| Response | [`ImageSourceUpdateResponse`](./imagesource.go) |

```go
imageSource, err := client.ImageSources.Update(context.Background(), "imageSourceId", sdk.ImageSourceUpdateParams{
	Aws: sdk.F[sdk.ImageSourceUpdateParamsAws](sdk.ImageSourceUpdateParamsAws{
		BucketName:   sdk.F[string]("my-bucket-test"),
		BucketRegion: sdk.F[string]("ap-southeast-1"),
		AccessKey:    sdk.F[string]("BQUA6QFXVWHAAB6IO2X1"),
		Secret:       sdk.F[string]("aws_secret"),
	}),
	Type: sdk.F[sdk.ImageSourceUpdateParamsType](sdk.ImageSourceUpdateParamsType("aws")),
})
if err != nil {
	panic(err)
}

fmt.Println(imageSource)
```

### Delete Source

This endpoint removes a image source. All image delivery using this subdomain will be stopped.

```go
imageSource, err := client.ImageSources.Delete(context.Background(), "imageSourceId")
if err != nil {
	panic(err)
}

fmt.Println(imageSource)
```

### Purge Cache

You can purge cache for any image by using our cache purge API.

| Direction | Type |
| --- | --- |
| Request | [`ImageSourcePurgeCacheParams`](./imagesource.go) |
| Response | [`ImageSourcePurgeCacheResponse`](./imagesource.go) |

```go
imageSource, err := client.ImageSources.PurgeCache(context.Background(), "subdomain", sdk.ImageSourcePurgeCacheParams{
	Paths: sdk.F[[]string]([]string{"image.jpeg", "image2.png"}),
})
if err != nil {
	panic(err)
}

fmt.Println(imageSource)
```

### Purge Image Cache

You can purge the cache for any image path by using this cache purge API.

| Direction | Type |
| --- | --- |
| Request | [`ImageSourcePurgeParams`](./imagesource.go) |
| Response | [`ImageSourcePurgeResponse`](./imagesource.go) |

```go
imageSource, err := client.ImageSources.Purge(context.Background(), "sourceId", sdk.ImageSourcePurgeParams{
	Paths: sdk.F[[]string]([]string{"image.jpeg", "image2.png"}),
})
if err != nil {
	panic(err)
}

fmt.Println(imageSource)
```

## `ImageUsageAnalytics`

Query aggregated and chart-ready image usage analytics data.

### Image Usage Analytics

This endpoint helps you get image analytics data like bandwidth consumption, request count, CDN hit ratio, etc.

| Direction | Type |
| --- | --- |
| Request | [`ImageUsageAnalyticGetParams`](./imageusageanalytic.go) |
| Response | [`ImageUsageAnalyticGetResponse`](./imageusageanalytic.go) |

```go
imageUsageAnalytic, err := client.ImageUsageAnalytics.Get(context.Background(), sdk.ImageUsageAnalyticGetParams{
	DateRange: sdk.F[sdk.ImageUsageAnalyticGetParamsDateRange](sdk.ImageUsageAnalyticGetParamsDateRange{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	Metrics: sdk.F[[]sdk.ImageUsageAnalyticGetParamsMetric]([]sdk.ImageUsageAnalyticGetParamsMetric{"bandwidth_consumption"}),
	GroupBy: sdk.F[sdk.ImageUsageAnalyticGetParamsGroupBy](sdk.ImageUsageAnalyticGetParamsGroupBy("daily")),
})
if err != nil {
	panic(err)
}

fmt.Println(imageUsageAnalytic)
```

## `LiveStreamAssets`

Create, control, and monitor live stream assets.

### Create Live Asset

A live asset refers to media content/video that is live-streamed through Gumlet. This endpoint creates a live streaming asset allowing users to live stream a video that will be pushed to Gumlet.

| Direction | Type |
| --- | --- |
| Request | [`LiveStreamAssetNewParams`](./livestreamasset.go) |
| Response | [`LiveStreamAssetNewResponse`](./livestreamasset.go) |

```go
liveStreamAsset, err := client.LiveStreamAssets.New(context.Background(), sdk.LiveStreamAssetNewParams{
	LiveSourceID: sdk.F[string](""),
	Resolution:   sdk.F[string](""),
})
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

### Update Live Asset

A live asset refers to media content/video that is live-streamed through Gumlet. This endpoint allows users to update a live streaming asset.

| Direction | Type |
| --- | --- |
| Request | [`LiveStreamAssetUpdateParams`](./livestreamasset.go) |
| Response | [`LiveStreamAssetUpdateResponse`](./livestreamasset.go) |

```go
liveStreamAsset, err := client.LiveStreamAssets.Update(context.Background(), sdk.LiveStreamAssetUpdateParams{
	LiveAssetID: sdk.F[string](""),
})
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

### Get Live Asset Status

This endpoint retrieves the details of a live video asset that has previously been created.

| Direction | Type |
| --- | --- |
| Response | [`LiveStreamAssetGetStatusResponse`](./livestreamasset.go) |

```go
liveStreamAsset, err := client.LiveStreamAssets.GetStatus(context.Background(), "liveAssetId")
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

### Delete Live Asset

This endpoint removes a live asset given its unique live asset id. The live asset will be removed from storage as well, associated URLs will be inaccessible.

```go
liveStreamAsset, err := client.LiveStreamAssets.Delete(context.Background(), "liveAssetId")
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

### Complete Live Stream

This endpoint allows marking live assets complete. Once the live asset is marked complete, it can no longer be used to ingest the live stream on Gumlet.

```go
liveStreamAsset, err := client.LiveStreamAssets.Complete(context.Background(), "liveAssetId")
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

### Filter Live Assets

This endpoint lists live assets on the basis of `status` for the given `live_source_id`.

Pagination uses `page_number` / `page_size`. Use `next_page_token` for cursor-based pagination.

| Direction | Type |
| --- | --- |
| Request | [`LiveStreamAssetFilterParams`](./livestreamasset.go) |
| Response | [`LiveStreamAssetFilterResponse`](./livestreamasset.go) |

```go
liveStreamAsset, err := client.LiveStreamAssets.Filter(context.Background(), "liveSourceId", sdk.LiveStreamAssetFilterParams{
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

### `start`

Start a live stream.

```go
err := client.LiveStreamAssets.Start(context.Background(), "liveAssetId")
if err != nil {
	panic(err)
}
```

### Upload Live Thumbnails

Generate presigned upload URLs for live stream thumbnails. Supported thumbnail states are `preparing`, `disconnected`, and `end`.

| Direction | Type |
| --- | --- |
| Request | [`LiveStreamAssetUploadParams`](./livestreamasset.go) |
| Response | [`LiveStreamAssetUploadResponse`](./livestreamasset.go) |

```go
liveStreamAsset, err := client.LiveStreamAssets.Upload(context.Background(), sdk.LiveStreamAssetUploadParams{
	LiveAssetID: sdk.F[string]("68c406b147f9ad0c0d584ce2"),
	Statuses:    sdk.F[interface{}]("preparing"),
})
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

### Get Live Asset Status History

This endpoint retrieves the history of a live video asset that has previously been created.

| Direction | Type |
| --- | --- |
| Response | [`LiveStreamAssetStatusHistoryResponse`](./livestreamasset.go) |

```go
liveStreamAsset, err := client.LiveStreamAssets.StatusHistory(context.Background(), "liveAssetId")
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAsset)
```

## `RecycleBin`

Endpoints to list deleted assets and restore them.

### Recover Deleted Asset

Recovers a deleted asset from the recycle bin.

| Direction | Type |
| --- | --- |
| Request | [`RecycleBinRecoverParams`](./recyclebin.go) |

```go
err := client.RecycleBin.Recover(context.Background(), sdk.RecycleBinRecoverParams{
	AssetID: sdk.F[string](""),
})
if err != nil {
	panic(err)
}
```

### List Recycle Bin

List all assets in a recycle bin for a given workspace. The deleted assets are available for 30 days. After that, assets are permanently deleted.

Pagination uses `page_number` / `page_size` only (no `next_page_token`).

| Direction | Type |
| --- | --- |
| Request | [`RecycleBinListParams`](./recyclebin.go) |
| Response | [`RecycleBinListResponse`](./recyclebin.go) |

```go
recycleBin, err := client.RecycleBin.List(context.Background(), sdk.RecycleBinListParams{
	WorkspaceID: sdk.F[string]("workspaceId"),
	PageNumber:  sdk.F[int64](1),
	PageSize:    sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(recycleBin)
```

## `VideoWorkspaces`

Create and manage video workspaces.

### List Workspaces

This endpoint list video workspace which are assigned to the user or token.

Pagination uses `page_number` / `page_size` (max page size 200). Use `next_page_token` for cursor-based pagination.

| Direction | Type |
| --- | --- |
| Request | [`VideoWorkspaceListParams`](./videoworkspace.go) |
| Response | [`VideoWorkspaceListResponse`](./videoworkspace.go) |

```go
videoWorkspace, err := client.VideoWorkspaces.List(context.Background(), sdk.VideoWorkspaceListParams{
	OrderBy:    sdk.F[sdk.VideoWorkspaceListParamsOrderBy](sdk.VideoWorkspaceListParamsOrderBy("asc")),
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(videoWorkspace)
```

### Create Workspace

Video workspaces are top-level entities in Gumlet. You can use them to organize videos for different teams/departments or use cases.

| Direction | Type |
| --- | --- |
| Request | [`VideoWorkspaceNewParams`](./videoworkspace.go) |
| Response | [`VideoWorkspaceNewResponse`](./videoworkspace.go) |

```go
videoWorkspace, err := client.VideoWorkspaces.New(context.Background(), sdk.VideoWorkspaceNewParams{
	Name:             sdk.F[string]("My Workspace"),
	Type:             sdk.F[sdk.VideoWorkspaceNewParamsType](sdk.VideoWorkspaceNewParamsType("direct-upload")),
	DataLocalization: sdk.F[sdk.VideoWorkspaceNewParamsDataLocalization](sdk.VideoWorkspaceNewParamsDataLocalization("us")),
})
if err != nil {
	panic(err)
}

fmt.Println(videoWorkspace)
```

### Update Workspace

This endpoint allows users to update video workspace that has previously been created.

| Direction | Type |
| --- | --- |
| Request | [`VideoWorkspaceUpdateParams`](./videoworkspace.go) |
| Response | [`VideoWorkspaceUpdateResponse`](./videoworkspace.go) |

```go
videoWorkspace, err := client.VideoWorkspaces.Update(context.Background(), "workspaceId", sdk.VideoWorkspaceUpdateParams{
	Aws: sdk.F[sdk.VideoWorkspaceUpdateParamsAws](sdk.VideoWorkspaceUpdateParamsAws{
		BucketName:   sdk.F[string]("my-bucket-test"),
		BucketRegion: sdk.F[string]("ap-southeast-1"),
		AccessKey:    sdk.F[string]("BQUA6QFXVWHAAB6IO2X1"),
		Secret:       sdk.F[string]("aws_secret"),
	}),
	DefaultProfileID: sdk.F[string]("646df1c9173a4a2fcac180b7"),
	Name:             sdk.F[string]("awsrename"),
	Type:             sdk.F[sdk.VideoWorkspaceUpdateParamsType](sdk.VideoWorkspaceUpdateParamsType("aws")),
})
if err != nil {
	panic(err)
}

fmt.Println(videoWorkspace)
```

### Get Workspace

This endpoint get all the data of video workspace that has previously been created.

| Direction | Type |
| --- | --- |
| Response | [`VideoWorkspaceGetResponse`](./videoworkspace.go) |

```go
videoWorkspace, err := client.VideoWorkspaces.Get(context.Background(), "workspaceId")
if err != nil {
	panic(err)
}

fmt.Println(videoWorkspace)
```

### Delete Workspace

This endpoint removes a video workspace given its unique asset id. All the asset in workspace will be removed from storage as well, associated URLs will be inaccessible.

```go
videoWorkspace, err := client.VideoWorkspaces.Delete(context.Background(), "workspaceId")
if err != nil {
	panic(err)
}

fmt.Println(videoWorkspace)
```

## `Folders`

Organize video assets into folders within a workspace.

### Create Folder

Create a folder inside a video workspace. Optionally provide `parent_id` to create a nested folder.

| Direction | Type |
| --- | --- |
| Request | [`FolderNewParams`](./folder.go) |
| Response | [`FolderNewResponse`](./folder.go) |

```go
folder, err := client.Folders.New(context.Background(), "workspaceId", sdk.FolderNewParams{
	Name:     sdk.F[string]("Course Assets"),
	ParentID: sdk.F[string](""),
})
if err != nil {
	panic(err)
}

fmt.Println(folder)
```

### List Folders

List folders for a video workspace. Use `parent_id` to list only folders inside a specific parent folder.

| Direction | Type |
| --- | --- |
| Request | [`FolderListParams`](./folder.go) |
| Response | [`[]FolderListResponse`](./folder.go) |

```go
folder, err := client.Folders.List(context.Background(), "workspaceId", sdk.FolderListParams{})
if err != nil {
	panic(err)
}

fmt.Println(folder)
```

### Get Folder

Get a single folder by id.

| Direction | Type |
| --- | --- |
| Response | [`FolderGetResponse`](./folder.go) |

```go
folder, err := client.Folders.Get(context.Background(), "workspaceId", "folderId")
if err != nil {
	panic(err)
}

fmt.Println(folder)
```

### Update Folder

Rename a folder, move it to another parent folder, or move assets into the folder by sending `asset_ids`.

| Direction | Type |
| --- | --- |
| Request | [`FolderUpdateParams`](./folder.go) |
| Response | [`FolderUpdateResponse`](./folder.go) |

```go
folder, err := client.Folders.Update(context.Background(), "workspaceId", "folderId", sdk.FolderUpdateParams{
	Body: sdk.FolderUpdateParamsBodyVariant0{
		Name: sdk.F[string]("Course Assets Updated"),
	},
})
if err != nil {
	panic(err)
}

fmt.Println(folder)
```

### Delete Folder

Delete a folder. Descendant folders and assets inside them are deleted by the backend workflow.

| Direction | Type |
| --- | --- |
| Response | [`FolderDeleteResponse`](./folder.go) |

```go
folder, err := client.Folders.Delete(context.Background(), "workspaceId", "folderId")
if err != nil {
	panic(err)
}

fmt.Println(folder)
```

### Remove Assets From Folder

Remove one or more assets from their current folder assignment inside the workspace.

| Direction | Type |
| --- | --- |
| Request | [`FolderRemoveAssetsParams`](./folder.go) |
| Response | [`FolderRemoveAssetsResponse`](./folder.go) |

```go
folder, err := client.Folders.RemoveAssets(context.Background(), "workspaceId", sdk.FolderRemoveAssetsParams{
	AssetIDs: sdk.F[[]string]([]string{"67e4f2b4403562dbea654301", "67e4f2bb403562dbea654302"}),
})
if err != nil {
	panic(err)
}

fmt.Println(folder)
```

## `ChannelViewers`

Invite and remove viewers on a private video channel.

### Invite Channel Viewers

Invite one or more viewers to a members-only channel.

| Direction | Type |
| --- | --- |
| Request | [`ChannelViewerInviteParams`](./channelviewer.go) |
| Response | [`ChannelViewerInviteResponse`](./channelviewer.go) |

```go
channelViewer, err := client.ChannelViewers.Invite(context.Background(), "videoWorkspaceId", sdk.ChannelViewerInviteParams{
	Users: sdk.F[[]sdk.ChannelViewerInviteParamsUser]([]sdk.ChannelViewerInviteParamsUser{sdk.ChannelViewerInviteParamsUser{
		Email: sdk.F[string]("test@gumlet.com"),
		Name:  sdk.F[string]("Test User-0"),
	}, sdk.ChannelViewerInviteParamsUser{
		Email: sdk.F[string]("test+1@gumlet.com"),
		Name:  sdk.F[string]("Test User-1"),
	}, sdk.ChannelViewerInviteParamsUser{
		Email: sdk.F[string]("test+2@gumlet.com"),
		Name:  sdk.F[string]("Test User-2"),
	}}),
})
if err != nil {
	panic(err)
}

fmt.Println(channelViewer)
```

### Remove Channel Viewers

Remove one or more viewers from a channel by email address.

| Direction | Type |
| --- | --- |
| Request | [`ChannelViewerRemoveParams`](./channelviewer.go) |
| Response | [`ChannelViewerRemoveResponse`](./channelviewer.go) |

```go
channelViewer, err := client.ChannelViewers.Remove(context.Background(), "videoWorkspaceId", sdk.ChannelViewerRemoveParams{
	Emails: sdk.F[[]string]([]string{"test@gumlet.com", "test+2@gumlet.com"}),
})
if err != nil {
	panic(err)
}

fmt.Println(channelViewer)
```

### Invite Channel Viewers via CSV

Invite viewers to a channel by uploading a CSV file.

| Direction | Type |
| --- | --- |
| Request | [`ChannelViewerInviteCsvParams`](./channelviewer.go) |
| Response | [`ChannelViewerInviteCsvResponse`](./channelviewer.go) |

```go
channelViewer, err := client.ChannelViewers.InviteCsv(context.Background(), "videoWorkspaceId", sdk.ChannelViewerInviteCsvParams{
	ViewersCsv: sdk.F[io.Reader](strings.NewReader("")),
})
if err != nil {
	panic(err)
}

fmt.Println(channelViewer)
```

### List Subscribers

List all channel subscribers.

Pagination uses `page_number` / `page_size`. Use `next_page_token` for cursor-based pagination.

| Direction | Type |
| --- | --- |
| Request | [`ChannelViewerListSubscribersParams`](./channelviewer.go) |
| Response | [`ChannelViewerListSubscribersResponse`](./channelviewer.go) |

```go
channelViewer, err := client.ChannelViewers.ListSubscribers(context.Background(), "workspaceId", sdk.ChannelViewerListSubscribersParams{
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(channelViewer)
```

## `VideoAnalytics`

Query aggregated and chart-ready video analytics data.

### Viewer Analytics

This endpoint retrieves viewer analytics data. This endpoint is use for deep insights on the analytics data.

| Direction | Type |
| --- | --- |
| Request | [`VideoAnalyticViewerDataParams`](./videoanalytic.go) |
| Response | [`VideoAnalyticViewerDataResponse`](./videoanalytic.go) |

```go
videoAnalytic, err := client.VideoAnalytics.ViewerData(context.Background(), sdk.VideoAnalyticViewerDataParams{
	DateRange: sdk.F[sdk.VideoAnalyticViewerDataParamsDateRange](sdk.VideoAnalyticViewerDataParamsDateRange{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	Metrics: sdk.F[[]sdk.VideoAnalyticViewerDataParamsMetric]([]sdk.VideoAnalyticViewerDataParamsMetric{"views"}),
	GroupBy: sdk.F[sdk.VideoAnalyticViewerDataParamsGroupBy](sdk.VideoAnalyticViewerDataParamsGroupBy("daily")),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAnalytic)
```

### Breakdown Data

Read viewer metrics for one workspace, grouped by a breakdown field, for a date range of at most 60 days. This does not change assets or delivery.

| Direction | Type |
| --- | --- |
| Request | [`VideoAnalyticBreakdownDataParams`](./videoanalytic.go) |
| Response | [`VideoAnalyticBreakdownDataResponse`](./videoanalytic.go) |

```go
videoAnalytic, err := client.VideoAnalytics.BreakdownData(context.Background(), sdk.VideoAnalyticBreakdownDataParams{
	Breakdowns: sdk.F[[]sdk.VideoAnalyticBreakdownDataParamsBreakdown]([]sdk.VideoAnalyticBreakdownDataParamsBreakdown{sdk.VideoAnalyticBreakdownDataParamsBreakdown{
		Name:     sdk.F[sdk.VideoAnalyticBreakdownDataParamsBreakdownsName](sdk.VideoAnalyticBreakdownDataParamsBreakdownsName("custom_video_id")),
		Metric:   sdk.F[sdk.VideoAnalyticBreakdownDataParamsBreakdownsMetric](sdk.VideoAnalyticBreakdownDataParamsBreakdownsMetric("views")),
		Page:     sdk.F[int64](1),
		PageSize: sdk.F[int64](100),
	}, sdk.VideoAnalyticBreakdownDataParamsBreakdown{
		Name:     sdk.F[sdk.VideoAnalyticBreakdownDataParamsBreakdownsName](sdk.VideoAnalyticBreakdownDataParamsBreakdownsName("custom_video_title")),
		Metric:   sdk.F[sdk.VideoAnalyticBreakdownDataParamsBreakdownsMetric](sdk.VideoAnalyticBreakdownDataParamsBreakdownsMetric("completion_percent_by_views")),
		Page:     sdk.F[int64](1),
		PageSize: sdk.F[int64](100),
	}}),
	DateRange: sdk.F[sdk.VideoAnalyticBreakdownDataParamsDateRange](sdk.VideoAnalyticBreakdownDataParamsDateRange{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	WorkspaceID: sdk.F[string]("6694c405e63913eecf3cf5fb"),
	Filters:     sdk.F[[]sdk.VideoAnalyticBreakdownDataParamsFilter]([]sdk.VideoAnalyticBreakdownDataParamsFilter{}),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAnalytic)
```

### Aggregated Data

This endpoint retrieves aggregated data of the given metrics.

| Direction | Type |
| --- | --- |
| Request | [`VideoAnalyticAggregatedDataParams`](./videoanalytic.go) |
| Response | [`VideoAnalyticAggregatedDataResponse`](./videoanalytic.go) |

```go
videoAnalytic, err := client.VideoAnalytics.AggregatedData(context.Background(), sdk.VideoAnalyticAggregatedDataParams{
	Aggregate: sdk.F[[]sdk.VideoAnalyticAggregatedDataParamsAggregate]([]sdk.VideoAnalyticAggregatedDataParamsAggregate{sdk.VideoAnalyticAggregatedDataParamsAggregate{
		Metric:   sdk.F[sdk.VideoAnalyticAggregatedDataParamsAggregateMetric](sdk.VideoAnalyticAggregatedDataParamsAggregateMetric("views")),
		Function: sdk.F[sdk.VideoAnalyticAggregatedDataParamsAggregateFunction](sdk.VideoAnalyticAggregatedDataParamsAggregateFunction("sum")),
	}}),
	Timeframe: sdk.F[sdk.VideoAnalyticAggregatedDataParamsTimeframe](sdk.VideoAnalyticAggregatedDataParamsTimeframe{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	WorkspaceID: sdk.F[string](""),
})
if err != nil {
	panic(err)
}

fmt.Println(videoAnalytic)
```

## `OrganizationData`

Endpoints to get organization details.

### Get Organization Details

You can get organization data using this API.

| Direction | Type |
| --- | --- |
| Response | [`OrganizationDataFetchOrgResponse`](./organizationdata.go) |

```go
organizationData, err := client.OrganizationData.FetchOrg(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(organizationData)
```

## `UserData`

Endpoints to get user details.

### Get User

This endpoint gives information about the user account.

| Direction | Type |
| --- | --- |
| Response | [`UserDataFetchResponse`](./userdata.go) |

```go
userData, err := client.UserData.Fetch(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(userData)
```

## `AuditLogs`

Get detailed user activity logs for the entire organisation.

### Fetch Audit Logs

Get audit logs for the user activity in your organisation. Please note that this endpoint can only be accessed by `owner` and `admin` role users.

| Direction | Type |
| --- | --- |
| Request | [`AuditLogFetchParams`](./auditlog.go) |
| Response | [`AuditLogFetchResponse`](./auditlog.go) |

```go
auditLog, err := client.AuditLogs.Fetch(context.Background(), sdk.AuditLogFetchParams{
	DateRange: sdk.F[sdk.AuditLogFetchParamsDateRange](sdk.AuditLogFetchParamsDateRange{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	PageNumber: sdk.F[int64](1),
	PageSize:   sdk.F[int64](100),
})
if err != nil {
	panic(err)
}

fmt.Println(auditLog)
```

## `Billing`

Get / change all details about billing and invoices.

### List Invoices

List all invoices that are generated so far.

| Direction | Type |
| --- | --- |
| Request | [`BillingListInvoicesParams`](./billing.go) |
| Response | [`BillingListInvoicesResponse`](./billing.go) |

```go
billing, err := client.Billing.ListInvoices(context.Background(), sdk.BillingListInvoicesParams{
	PageSize: sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(billing)
```

### Get Billing Details

Get billing details for this organization.

| Direction | Type |
| --- | --- |
| Response | [`BillingFetchDetailsResponse`](./billing.go) |

```go
billing, err := client.Billing.FetchDetails(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(billing)
```

### Update Billing Details

Update billing details

| Direction | Type |
| --- | --- |
| Request | [`BillingUpdateDetailsParams`](./billing.go) |
| Response | [`BillingUpdateDetailsResponse`](./billing.go) |

```go
billing, err := client.Billing.UpdateDetails(context.Background(), sdk.BillingUpdateDetailsParams{
	AddressLine: sdk.F[string](""),
	City:        sdk.F[string](""),
	CompanyName: sdk.F[string](""),
	CountryCode: sdk.F[string](""),
	GstNumber:   sdk.F[string](""),
	Postal:      sdk.F[string](""),
	StateCode:   sdk.F[string](""),
})
if err != nil {
	panic(err)
}

fmt.Println(billing)
```

### Upcoming Invoice

Get details about upcoming invoice.

| Direction | Type |
| --- | --- |
| Response | [`BillingFetchUpcomingInvoiceResponse`](./billing.go) |

```go
billing, err := client.Billing.FetchUpcomingInvoice(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(billing)
```

## `LiveStreamWorkspaces`

Create and manage live stream workspaces.

### List Workspaces

List all live stream workspaces.

| Direction | Type |
| --- | --- |
| Response | [`LiveStreamWorkspaceListResponse`](./livestreamworkspace.go) |

```go
liveStreamWorkspace, err := client.LiveStreamWorkspaces.List(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(liveStreamWorkspace)
```

### Create Workspace

Create live stream workspace.

| Direction | Type |
| --- | --- |
| Request | [`LiveStreamWorkspaceNewParams`](./livestreamworkspace.go) |
| Response | [`LiveStreamWorkspaceNewResponse`](./livestreamworkspace.go) |

```go
liveStreamWorkspace, err := client.LiveStreamWorkspaces.New(context.Background(), sdk.LiveStreamWorkspaceNewParams{
	Name: sdk.F[string](""),
})
if err != nil {
	panic(err)
}

fmt.Println(liveStreamWorkspace)
```

### Update Workspace

Update live stream workspace.

| Direction | Type |
| --- | --- |
| Request | [`LiveStreamWorkspaceUpdateParams`](./livestreamworkspace.go) |
| Response | [`LiveStreamWorkspaceUpdateResponse`](./livestreamworkspace.go) |

```go
liveStreamWorkspace, err := client.LiveStreamWorkspaces.Update(context.Background(), "liveWorkspaceId", sdk.LiveStreamWorkspaceUpdateParams{
	Name:          sdk.F[string]("live-stream-collections"),
	VideoSourceID: sdk.F[string]("67bea1d66ca0059a95bf7de9"),
})
if err != nil {
	panic(err)
}

fmt.Println(liveStreamWorkspace)
```

### Delete Workspace

Delete the live stream workspace.

```go
liveStreamWorkspace, err := client.LiveStreamWorkspaces.Delete(context.Background(), "liveWorkspaceId")
if err != nil {
	panic(err)
}

fmt.Println(liveStreamWorkspace)
```

## `LiveStreamAnalytics`

Get usage analytics for live streams.

### Usage Analytics

Get usage analytics for your live streams.

| Direction | Type |
| --- | --- |
| Request | [`LiveStreamAnalyticUsageParams`](./livestreamanalytic.go) |
| Response | [`LiveStreamAnalyticUsageResponse`](./livestreamanalytic.go) |

```go
liveStreamAnalytic, err := client.LiveStreamAnalytics.Usage(context.Background(), sdk.LiveStreamAnalyticUsageParams{
	DateRange: sdk.F[sdk.LiveStreamAnalyticUsageParamsDateRange](sdk.LiveStreamAnalyticUsageParamsDateRange{
		StartAt: sdk.F[time.Time](time.Now()),
		EndAt:   sdk.F[time.Time](time.Now()),
	}),
	GroupBy: sdk.F[sdk.LiveStreamAnalyticUsageParamsGroupBy](sdk.LiveStreamAnalyticUsageParamsGroupBy("daily")),
	Metrics: sdk.F[[]sdk.LiveStreamAnalyticUsageParamsMetric]([]sdk.LiveStreamAnalyticUsageParamsMetric{"bandwidth_consumption"}),
})
if err != nil {
	panic(err)
}

fmt.Println(liveStreamAnalytic)
```

## `GlobalSearch`

Search all video assets / playlists / folders etc across workspaces.

### Global Search

Search all video assets / playlists / folders etc across workspaces.

| Direction | Type |
| --- | --- |
| Request | [`GlobalSearchSearchParams`](./globalsearch.go) |
| Response | [`GlobalSearchSearchResponse`](./globalsearch.go) |

```go
globalSearch, err := client.GlobalSearch.Search(context.Background(), sdk.GlobalSearchSearchParams{
	SearchQuery:     sdk.F[string]("searchQuery"),
	AssetsOffset:    sdk.F[int64](0),
	ChannelsOffset:  sdk.F[int64](0),
	FoldersOffset:   sdk.F[int64](0),
	PlaylistsOffset: sdk.F[int64](0),
	Size:            sdk.F[int64](20),
})
if err != nil {
	panic(err)
}

fmt.Println(globalSearch)
```
