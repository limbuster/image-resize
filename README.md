# Serverless with Go

## Before running anything

Please change the `serverless.yml` file to reflect the correct s3 bucket.

## Background

When an image is uploaded to `s3://bucket/uploads/file.jpg` it will be copied as original to `s3://bucket/uploads/<md5_checksum>.jpg` and the file will be removed from the original. It will also create a thumbnail with max dimension of 500 pixels.

This project uses [resize](https://github.com/nfnt/resize) to generate the thumbnail.

## Build

```bash
make
```

## Deploy

```bash
sls deploy
```

## Upload to s3 bucket

```bash
aws s3 cp resources/test.jpg s3://<bucket>/uploads/
```

## Improvements

- Add unit tests
