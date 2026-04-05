main_go_path=./cmd/main
build_dir=./bin
binary_name=nas-ui

.REST_API: build

build:
	@mkdir -p ${build_dir}
	GOARCH=arm64 GOOS=linux go build -o ${build_dir}/${binary_name}-arm64 ${main_go_path}