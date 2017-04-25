# gocr

OCR easily from CIL. Backed [Google Cloud Vision API](https://cloud.google.com/vision/).

## Dependencies

- Go
- [GoogleCloudPlatform/google-cloud-go: Google Cloud APIs Go Client Library](https://github.com/GoogleCloudPlatform/google-cloud-go)

## Installation

```bash
$ go get github.com/acomagu/gocr
```

## Usage

Before of all you must set [Google Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials). The simplest way is use `gcloud` command.

```bash
$ gcloud auth login
```

And now you can use it!

```bash
$ gocr file1.jpg file2.jpg --concurrency=2
```

## License

This software includes the work that is distributed in the Apache License 2.0.

This software is released under the MIT License, see [LICENSE.txt](https://github.com/acomagu/gocr/blob/master/LICENSE).
