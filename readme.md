# S3 Object Management CLI (Go)

This Go application provides a command-line interface (CLI) for managing objects in an Amazon S3 bucket. Users can upload, download, delete, and list objects in a specified S3 bucket.

## Features

- List contents of an S3 bucket
- Upload files to an S3 bucket
- Download files from an S3 bucket
- Delete objects from an S3 bucket

## Prerequisites

- Go 1.22 or newer
- AWS account (IAM) with S3 access
- AWS credentials (Access Key ID and Secret Access Key)
- `.env` file in the project root with the following variables:
```
AWS_ACCESS_KEY_ID=your_access_key_id
AWS_SECRET_ACCESS_KEY=your_secret_access_key
```

## Usage

Run the application with the S3 bucket name as an argument:

```
go run main.go <bucket-name>
```


### Commands

Once the application is running, you can use the following commands:

- **upload**: Enter `upload` to upload a file. You will be prompted to enter the file path.
- **download**: Enter `download` to download an object. You will be prompted to enter the object key.
- **delete**: Enter `delete` to delete an object. You will be prompted to enter the object key.
- **exit**: Enter `exit` to terminate the application.

