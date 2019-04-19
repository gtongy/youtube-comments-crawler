build:
	GOOS=linux go build -o main
main-zip: build
	zip main.zip main
clean:
	rm -f main main.zip

generate-event:
	sam local generate-event cloudwatch scheduled-event \
	--region ap-northeast-1 > event.json

create-network:
	docker network inspect youtube-comments-crawler-network &>/dev/null || \
	docker network create --driver bridge youtube-comments-crawler-network

# dynamodb operate
list-tables:
	aws dynamodb list-tables \
	--region ap-northeast-1 \
	--endpoint-url http://localhost:8000

ATTRIBUTE_DEFINITIONS ?= '[{"AttributeName":"key","AttributeType": "S"}]'
KEY_SCHEMA ?= '[{"AttributeName":"key","KeyType": "HASH"}]'
PROVISIONED_THROUGHPUT ?= '{"ReadCapacityUnits": 5,"WriteCapacityUnits": 5}'

create-table:
	aws dynamodb create-table --table-name $(TABLE_NAME) \
	--region ap-northeast-1 \
	--attribute-definitions $(ATTRIBUTE_DEFINITIONS) \
	--key-schema $(KEY_SCHEMA) \
	--provisioned-throughput $(PROVISIONED_THROUGHPUT) \
	--endpoint-url http://localhost:8000

local-exec:
	sam local invoke $(FUNCTION_NAME) \
	--region ap-northeast-1 \
	--env-vars env.json \
	--event event.json \
	--template=deploy/template/staging.yml \
	--docker-network minio-network
