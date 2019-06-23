# Youtube Comments Crawler

This function is youtube comment crawler to use AWS SAM.

## Description

This function has the following configuration.

![AWS Diagram](https://github.com/gtongy/youtube-comments-crawler/blob/master/images/aws-youtube-comments-crawler.png)

This function is use under AWS Service

- Compute
  - Lambda
    - Function Exec
- Database
  - DynamoDB
    - Resource Storage
- Administration & Security
  - CloudWatch
    - Schedule Event, and Logging
- Deployment & Management
  - CloudFormation
    - Create deploy stack
  - X-Ray
    - trace func exec and aws services

## Requirement

- aws-sam-cli
  - 0.14.2
- go
  - 1.12.1
- docker
  - 18.09.2
- docker-compose
  - 1.23.2
- aws-cli
  - 1.16.130

## Usage

### Development

####

Development use aws-sam-cli.
So, Please install aws-sam-cli.

- use Mac or Linux

```
$ brew tap aws/tap
$ brew install aws-sam-cli
```

- using pip

```
$ pip install --user aws-sam-cli
```

Before execute commands, you need to get google service account with YouTube Data API certification enabled.
After clone this repository, the following command will be executed.

- container up

```sh
$ cd /path/to/youtube-comments-crawler
$ make create-network && docker-compose up -d
```

- put item

```sh
$  make put-item TABLE_NAME='YoutubeCommentsCrawlerYoutubers' ITEM='{ "id": { "S": "unique xid insert" }, "name": { "S": "Please Input Youtuber Name" }, "channel_id": { "S": "Please Input Youtuber Channel ID" }}'
```

- create dummy tables

```sh
 $ make create-table TABLE_NAME="YoutubeCommentsCrawlerVideos"
 $ make create-table TABLE_NAME="YoutubeCommentsCrawlerComments"
 $ make create-table TABLE_NAME="YoutubeCommentsCrawlerYoutubers"
```

- execute lambda local

```sh
$ cd /path/to/youtube-comments-crawler
$ aws s3 cp /path/to/service-account.json s3://google-service-accounts-dev/youtube-comments-crawler \
--endpoint-url=http://localhost:9001 \
--region ap-northeast-1 --profile $MINIO_PROFILE
```

### Deploy

Deploy is direct create croudformation stack.
Inside calls sam package, sam deploy. Please Look at Makefile.

```sh
$ make create-package
$ make deploy-package
```

## Install

`git clone https://github.com/gtongy/youtube-comments-crawler`

## Author

[gtongy](https://github.com/gtongy)
