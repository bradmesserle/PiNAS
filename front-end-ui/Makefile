main_go_path=./cmd/main
build_dir=./bin
binary_name=nas-ui
version=0.0.1

all: compile-templ build

compile-templ:
	go tool templ generate


build:
	@mkdir -p ${build_dir}
	GOARCH=arm64 GOOS=linux go build -o ${build_dir}/${binary_name}-arm64 ${main_go_path}


package-deb:
	@mkdir -p /tmp/PINAS-UI-Services-${version}_arm64/DEBIAN
	@mkdir -p /tmp/PINAS-UI-Services-${version}_arm64/opt/pinas/ui-services
	cp ${build_dir}/${binary_name}-arm64 /tmp/PINAS-UI-Services-${version}_arm64/opt/pinas/ui-services
	cp unix-scripts/systemd/pinas-ui-services.service /tmp/PINAS-UI-Services-${version}_arm64/opt/pinas/ui-services
	chmod +x /tmp/PINAS-UI-Services-${version}_arm64/opt/pinas/ui-services
	cp deb-package-files/control /tmp/PINAS-UI-Services-${version}_arm64/DEBIAN
	cp deb-package-files/postinst /tmp/PINAS-UI-Services-${version}_arm64/DEBIAN
	cp deb-package-files/postrm /tmp/PINAS-UI-Services-${version}_arm64/DEBIAN
	chmod +x  /tmp/PINAS-UI-Services-${version}_arm64/DEBIAN/postinst
	chmod +x  /tmp/PINAS-UI-Services-${version}_arm64/DEBIAN/postrm
	dpkg-deb --build --root-owner-group /tmp/PINAS-UI-Services-${version}_arm64
