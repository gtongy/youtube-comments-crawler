# util commands
build:
	GOOS=linux go build -o main
main-zip: build
	zip main.zip main
clean:
	rm -f main main.zip

# local exec commands
generate-event:
	sam local generate-event cloudwatch scheduled-event \
	--region ap-northeast-1 > event.json

local-exec:
	sam local invoke $(FUNCTION_NAME) \
	--region ap-northeast-1 \
	--env-vars env.json \
	--event event.json \
	--template=deploy/template/staging.yml
