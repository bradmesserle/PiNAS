rest_package_path=./cmd/rest-api
build_dir=./bin
binary_name=rest-api

.REST_API: build

build:
	@mkdir -p ${build_dir}
	GOARCH=arm64 GOOS=linux go build -o ${build_dir}/${binary_name}-arm64 ${rest_package_path}
