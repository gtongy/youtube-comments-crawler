# create handler
build:
	GOOS=linux go build -o main
main-zip: build
	zip main.zip main
clean:
	rm -f main main.zip

# deploy commands
create-package:
	sam package \
	--template-file ./deploy/template/production.yml \
	--output-template-file package.template.yml \
	--s3-bucket $(CREATE_PACKAGE_BUCKET_NAME)
deploy-package:
	sam deploy \
	--template-file ./package.template.yml \
	--region ap-northeast-1 \
	--stack-name youtube-comments-crawler \
	--capabilities CAPABILITY_IAM

# local network create
create-network:
	docker network inspect youtube-comments-crawler-network &>/dev/null || \
	docker network create --driver bridge youtube-comments-crawler-network

# dynamodb operate
list-tables:
	aws dynamodb list-tables \
	--region ap-northeast-1 \
	--endpoint-url http://localhost:8000

ATTRIBUTE_DEFINITIONS ?= '[{"AttributeName":"id","AttributeType": "S"}]'
KEY_SCHEMA ?= '[{"AttributeName":"id","KeyType": "HASH"}]'
PROVISIONED_THROUGHPUT ?= '{"ReadCapacityUnits": 5,"WriteCapacityUnits": 5}'
create-table:
	aws dynamodb create-table --table-name $(TABLE_NAME) \
	--region ap-northeast-1 \
	--attribute-definitions $(ATTRIBUTE_DEFINITIONS) \
	--key-schema $(KEY_SCHEMA) \
	--provisioned-throughput $(PROVISIONED_THROUGHPUT) \
	--endpoint-url http://localhost:8000
delete-table:
	aws dynamodb delete-table --table-name $(TABLE_NAME) \
	--region ap-northeast-1 \
	--endpoint-url http://localhost:8000
put-item:
	aws dynamodb put-item --table-name $(TABLE_NAME) \
	--item '$(ITEM)' \
	--region ap-northeast-1 \
	--endpoint-url http://localhost:8000
scan-items:
	aws dynamodb scan --table-name $(TABLE_NAME) \
	--region ap-northeast-1 \
	--endpoint-url http://localhost:8000

# docker lambda local exec
local-exec:
	sam local invoke $(FUNCTION_NAME) \
	--region ap-northeast-1 \
	--env-vars env.json \
	--event event.json \
	--template=deploy/template/production.yml \
	--docker-network youtube-comments-crawler-network
generate-event:
	sam local generate-event cloudwatch scheduled-event \
	--region ap-northeast-1 > event.json