# Mini-Alt

### S3-Compatible API
[Amazon Docs](https://docs.aws.amazon.com/AmazonS3/latest/API/API_Operations_Amazon_Simple_Storage_Service.html)
#### Bucket Operations:
- CreateBucket [x]
- DeleteBucket [x]
- ListBuckets [x]

#### Object Operations:
- PutObject [x]
- GetObject [x]
- DeleteObject [x]
- CopyObject [x]
- MultipartUpload (optional)
- ListObjectsV2 [x]
- HeadObject []

## Build for MacOS
```bash
GOOS=darwin GOARCH=arm64 go build -o build/mini-alt-arm .
GOOS=darwin GOARCH=amd64 go build -o build/mini-alt-amd .

lipo -create -output build/mini-alt build/mini-alt-arm build/mini-alt-amd
```