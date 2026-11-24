# Changelog

## [1.1.1](https://github.com/gumlet/golang-sdk/compare/v1.1.0...v1.1.1) (2026-11-24)


### ⚠ BREAKING CHANGES

* **api:** 19 breaking changes to the SDK surface.
    - method-name POST video/assets/upload — method names differ (breaking: renamed; the old spelling no longer resolves)
    - request-shape POST video/assets/upload — extra positional arg count differs (breaking: changed in a way existing callers can observe)
    - response-shape POST video/assets/upload→ — holds another named type — after: VideoAssetUploadFileResponse, before: VideoAssetUploadResponse (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - response-shape POST video/assets/upload→.input — holds another named type — after: VideoAssetUploadFileResponseInput, before: VideoAssetUploadResponseInput (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - response-shape POST video/assets/upload→.input.call_to_actions[] — holds another named type — after: VideoAssetUploadFileResponseInputCallToAction, before: VideoAssetUploadResponseInputCallToAction (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - response-shape POST video/assets/upload→.input.metadata — holds another named type — after: VideoAssetUploadFileResponseInputMetadata, before: VideoAssetUploadResponseInputMetadata (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - response-shape POST video/assets/upload→.input.transformations — holds another named type — after: VideoAssetUploadFileResponseInputTransformations, before: VideoAssetUploadResponseInputTransformations (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - response-shape POST video/assets/upload→.output — holds another named type — after: VideoAssetUploadFileResponseOutput, before: VideoAssetUploadResponseOutput (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams — holds another named type — after: VideoAssetUploadFileParams, before: VideoAssetUploadParams (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.additional_tracks[] — holds another named type — after: VideoAssetUploadFileParamsAdditionalTrack, before: VideoAssetUploadParamsAdditionalTrack (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.animated_gif — holds another named type — after: VideoAssetUploadFileParamsAnimatedGif, before: VideoAssetUploadParamsAnimatedGif (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.call_to_actions[] — holds another named type — after: VideoAssetUploadFileParamsCallToAction, before: VideoAssetUploadParamsCallToAction (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.crop — holds another named type — after: VideoAssetUploadFileParamsCrop, before: VideoAssetUploadParamsCrop (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - enum-values VideoAssetUploadFileParams.format — enum const names differ (wire values match) (breaking: changed in a way existing callers can observe)
    - request-shape VideoAssetUploadFileParams.generate_subtitles — holds another named type — after: VideoAssetUploadFileParamsGenerateSubtitles, before: VideoAssetUploadParamsGenerateSubtitles (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.image_overlay — holds another named type — after: VideoAssetUploadFileParamsImageOverlay, before: VideoAssetUploadParamsImageOverlay (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.pad — holds another named type — after: VideoAssetUploadFileParamsPad, before: VideoAssetUploadParamsPad (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.text_overlay — holds another named type — after: VideoAssetUploadFileParamsTextOverlay, before: VideoAssetUploadParamsTextOverlay (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - request-shape VideoAssetUploadFileParams.trim — holds another named type — after: VideoAssetUploadFileParamsTrim, before: VideoAssetUploadParamsTrim (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
* **api:** 3 breaking changes to the SDK surface.
    - method-name POST video/assets/{p}/thumbnail — method names differ (breaking: renamed; the old spelling no longer resolves)
    - response-shape POST video/assets/{p}/thumbnail→ — holds another named type — after: VideoAssetUploadThumbnailResponse, before: VideoAssetUploadThumbnailFileResponse (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - types-classes videoasset.go:VideoAssetUploadThumbnailFileResponse — struct only in before (breaking: removed; VideoAssetUploadThumbnailResponse appeared in its place, likely a rename)
* **api:** 4 breaking changes to the SDK surface.
    - method-name POST video/assets/{p}/thumbnail — method names differ (breaking: renamed; the old spelling no longer resolves)
    - request-shape POST video/assets/{p}/thumbnail — extra positional arg count differs (breaking: changed in a way existing callers can observe)
    - response-shape POST video/assets/{p}/thumbnail→ — holds another named type — after: VideoAssetUploadThumbnailFileResponse, before: VideoAssetThumbnailUploadResponse (breaking: retyped to another named struct; Go types are nominal, so values of the old type no longer assign)
    - types-classes videoasset.go:VideoAssetThumbnailUploadResponse — struct only in before (breaking: removed; VideoAssetUploadThumbnailFileResponse appeared in its place, likely a rename)

### Features

* **api:** update SDK surface (1 change) ([35b65bb](https://github.com/gumlet/golang-sdk/commit/35b65bb2725552f3e0514b65a6adde73668ed27e))
* **api:** update SDK surface (36 changes) ([b0b8f4e](https://github.com/gumlet/golang-sdk/commit/b0b8f4ef930e8497d29ec00218fd2027c29fc9f5))
* **api:** update SDK surface (4 changes) ([9117422](https://github.com/gumlet/golang-sdk/commit/9117422e215436426eb383afd224ec465ffcc82d))
* **api:** update SDK surface (6 changes) ([109b870](https://github.com/gumlet/golang-sdk/commit/109b8705a1816c25fbe487aed49c213027489055))


### Chores

* **api:** update generated SDK content ([6bcd37a](https://github.com/gumlet/golang-sdk/commit/6bcd37a14e3880706ef361485b9e76e697f289b3))
* release 1.1.1 ([782d595](https://github.com/gumlet/golang-sdk/commit/782d5950cf5d1c167e25560ec6a669082477d5b2))

## [1.1.0](https://github.com/gumlet/golang-sdk/compare/v1.0.1...v1.1.0) (2026-11-23)


### Features

* **api:** update SDK surface (2 changes) ([4d56fbe](https://github.com/gumlet/golang-sdk/commit/4d56fbe28322cb488214b34a8a5caea8c0990aa9))


### Chores

* **api:** update generated SDK content ([5655292](https://github.com/gumlet/golang-sdk/commit/5655292d814f3c4b4226a7e46df1e3d99957228f))
* **api:** update generated SDK content ([e41820f](https://github.com/gumlet/golang-sdk/commit/e41820ff30c75b781fe27f33d73297fa21b7aac2))
* **api:** update generated SDK content ([655a638](https://github.com/gumlet/golang-sdk/commit/655a63832d751b22213e02ab2b8fa054c8711bd7))

## [1.0.1](https://github.com/gumlet/golang-sdk/compare/v1.0.0...v1.0.1) (2026-11-19)


### Chores

* **api:** regenerate SDK ([e4bfb7b](https://github.com/gumlet/golang-sdk/commit/e4bfb7b135b63a1df27fcc1e17b81cbd3499578d))

## [1.0.0](https://github.com/gumlet/golang-sdk/compare/v0.1.0...v1.0.0) (2026-11-19)


### Features

* **api:** initial SDK generation ([33d0dec](https://github.com/gumlet/golang-sdk/commit/33d0dec61b07e62798b54be53eb2eab4baf8fb16))


### Chores

* **api:** update generated SDK content ([3a8cdc7](https://github.com/gumlet/golang-sdk/commit/3a8cdc7e8e9c66fcf40979b58b0f82c7ba816d9c))
* **api:** update generated SDK content ([1674c27](https://github.com/gumlet/golang-sdk/commit/1674c27c7bd7f6f74d4394e0485fff46e5cba0c4))
* **api:** update generated SDK content ([1500d56](https://github.com/gumlet/golang-sdk/commit/1500d56eb03d40f8769354cd47acb3994842f52b))
* release 1.0.0 ([4c27bb4](https://github.com/gumlet/golang-sdk/commit/4c27bb4530cd9d83c4f365144c43aebbf776f242))
