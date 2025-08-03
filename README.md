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

## Build for macOS
```bash
GOOS=darwin GOARCH=arm64 go build -o build/mini-alt-arm .
GOOS=darwin GOARCH=amd64 go build -o build/mini-alt-amd .

lipo -create -output build/mini-alt build/mini-alt-arm build/mini-alt-amd
```
## Initial Data Directory Setup

Mini-Alt comes with a method to load a directory on startup :)

Here I explain how to prepare your local directories so they can be automatically loaded as buckets and objects into the system.

> Run the binary with the `-load-initial-data` flag to trigger this loading process.

---

### Directory Structure

The root directory is used to represent your **bucket storage**.  
Each **top-level folder** becomes a **bucket**.  
All files and subfolders inside a bucket are treated as **objects**.

#### Example

```
data/
├── bucket-one/
│   ├── file1.txt
│   └── nested/
│       └── file2.png
├── bucket-two/
│   └── video.mp4
└── bucket-three/
    └── logs/
        └── 2025/
            └── report.json
```

- `bucket-one`, `bucket-two`, and `bucket-three` → registered as buckets
- Files inside (e.g. `file1.txt`, `video.mp4`, `report.json`) → registered as objects
- Folders inside (e.g. `nested/`, `logs/2025/`) → also registered as empty "directory" objects

---

### Metadata Handling

Each file will be automatically scanned to collect:
- **Content Length**
- **Content Type** (based on file content or extension)
- **MD5 Hash** (used as the ETag)

---

### Usage

1. Prepare your `data/` directory with the structure above.
2. Run your binary with the flag:

```bash
./minialt -load-initial-data
```

The system will load all buckets and objects into the database, including metadata.

---
