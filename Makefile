# Sir I have no idea
build:
	go build -a -ldflags '-w -extldflags "-static"' -o bin/aliyun-cr-exporter

dockerbuild:
	docker build -t bohrasd/aliyun-cr-exporter .

dockerpush:
	docker push bohrasd/aliyun-cr-exporter

fmt:
	go mod tidy
	gofmt -s -w .
