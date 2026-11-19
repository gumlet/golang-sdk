---
name: gumlet-go-sdk
description: "Go SDK for Gumlet API. Use when writing Go code that calls Gumlet API with the nodejssdk package: installing it, constructing and authenticating the client, and calling API operations."
---

# Gumlet Go SDK

Generated Go client for Gumlet API, published as `nodejssdk`. Use the generated client instead of hand-writing HTTP requests.

## Install

```sh
go get nodejssdk
```

## Client setup and authentication

```go
import (
	"context"
	"fmt"

	sdk "nodejssdk"
)

client := sdk.NewClient()
```

Provide credentials using the options below. Environment variables are read automatically when the target runtime supports them:

- `option.WithAPIKey` (env: `API_KEY`) — Credential for the API_KEY client option.

## Calling operations

```go
package main

import (
	"context"
	"fmt"
	"os"

	sdk "nodejssdk"
	"nodejssdk/option"
)

func main() {
	client := sdk.NewClient(
		option.WithAPIKey(os.Getenv("API_KEY")),
	)

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
}
```

Method names, parameter shapes, and response types are generated from the API description — do not guess them. Look up the exact call signature in [api.md](../../../api.md) before writing a call.

## Error handling

Non-success responses return generated API errors. Error objects expose status, headers, response body, and request metadata where the target runtime supports it.

```go
videoAsset, err := client.VideoAssets.New(context.Background(), sdk.VideoAssetNewParams{
	Input:       sdk.F[string]("http://devimages.apple.com/iphone/samples/bipbop/bipbopall.m3u8"),
	WorkspaceID: sdk.F[string]("<your workspace id>"),
	Format:      sdk.F[sdk.VideoAssetNewParamsFormat](sdk.VideoAssetNewParamsFormat("ABR")),
	Title:       sdk.F[string]("Example Title"),
})
if err != nil {
	var apiErr *sdk.Error
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode, apiErr.JSON.RawJSON())
	}
	panic(err)
}

// imports: "context", "errors", "fmt", sdk "nodejssdk"
```

## Requirements

- Go 1.22 or newer

## Reference files

- [README.md](../../../README.md) — full feature tour: client options, request options, retries and timeouts, webhooks, logging.
- [api.md](../../../api.md) — complete catalogue of every operation with request and response types.
