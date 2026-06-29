rest_package_path=./cmd/rest-api
build_dir=./bin
binary_name=rest-api
version=0.0.1

.REST_API: build

build:
	@mkdir -p ${build_dir}
	GOARCH=arm64 GOOS=linux go build -o ${build_dir}/${binary_name}-arm64 ${rest_package_path}

package-deb:
	@mkdir -p /tmp/PI-NAS-REST-Services-${version}_arm64/DEBIAN
	@mkdir -p /tmp/PI-NAS-REST-Services-${version}_arm64/opt/pinas/rest-services
	cp ${build_dir}/${binary_name}-arm64 /tmp/PI-NAS-REST-Services-${version}_arm64/opt/pinas/rest-services
	chmod + /tmp/PI-NAS-REST-Services-${version}_arm64/opt/pinas/rest-services
	cp deb-package-files/control /tmp/PI-NAS-REST-Services-${version}_arm64/DEBIAN
	cp deb-package-files/postinst /tmp/PI-NAS-REST-Services-${version}_arm64/DEBIAN
	chmod +  /tmp/PI-NAS-REST-Services-${version}_arm64/DEBIAN/postinst
	dpkg-deb --build --root-owner-group /tmp/PI-NAS-REST-Services_0.0-1_amd64