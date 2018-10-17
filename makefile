
.PHONY: version build build_linux docker_login docker_build docker_push_dev docker_push_pro
.PHONY: rm_stop test_rm test_clear docker_test

Version = v0.3.0
GOBIN = $(pwd)
VersionFile = version/version.go
GitCommit = `git rev-parse --short=8 HEAD`
BuildVersion = "`date +%FT%T%z`"
GOBIN = $(shell pwd)

ImagesPrefix = "registry.cn-hangzhou.aliyuncs.com/yuanben/"
ImageName = "usmint"
ImageTestName = "$(ImageName):test"
ImageCommitName = "$(ImageName):$(GitCommit)"

version:
	@echo "项目版本处理"
	@echo "package version" > $(VersionFile)
	@echo "const Version = "\"$(Version)\" >> $(VersionFile)
	@echo "const BuildVersion = "\"$(BuildVersion)\" >> $(VersionFile)
	@echo "const GitCommit = "\"$(GitCommit)\" >> $(VersionFile)

build:
	@echo "开始编译"
	GOBIN=$(GOBIN) go install main.go

build_linux: version
	@echo "交叉编译成linux应用"
	GOBIN=`pwd` CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install main.go

docker_login:
	@echo "登陆镜像仓库"
	sudo docker login -u baiyunhui@yuanben -p ybl12345 registry.cn-hangzhou.aliyuncs.com

test_rm:
	@echo "删除文件"
	@rm -rf test/d1
	@rm -rf test/d2
	@rm -rf test/d3
	@rm -rf test/d4
	@rm -rf test/d5
	@rm -rf test/d6

test_create:
	@echo "创建文件"
	./main --home test/d1 init
	./main --home test/d2 init
	./main --home test/d3 init
	./main --home test/d4 init
	./main --home test/d5 init
	./main --home test/d6 init


test_clear:
	@echo "reset文件"
	./main --home test/d1 unsafe_reset_all
	./main --home test/d2 unsafe_reset_all
	./main --home test/d3 unsafe_reset_all
	./main --home test/d4 unsafe_reset_all
	./main --home test/d5 unsafe_reset_all
	./main --home test/d6 unsafe_reset_all

docker_test:
	@echo "kchain docker test"
	@ls * | grep example_data || mkdir example_data
	sudo docker run --rm -it -v `pwd`/example_data:/kdata -p 46656:46656 -p 46657:46657 kchain init
	sudo docker run --rm -it -v `pwd`/example_data:/kdata -p 46656:46656 -p 46657:46657 kchain

rm_stop:
	@echo "删除所有的的容器"
	sudo docker rm -f $(sudo docker ps -qa)
	sudo docker ps -a

rm_none:
	@echo "删除所为none的image"
	sudo docker images  | grep none | awk '{print $3}' | xargs docker rmi

docker_push_pro: docker_build
	@echo "docker push pro"
	sudo docker tag $(ImageName) $(ImagesPrefix)$(ImageName)
	sudo docker tag $(ImageName) $(ImagesPrefix)$(ImageCommitName)

	sudo docker push $(ImagesPrefix)$(ImageCommitName)
	sudo docker push $(ImagesPrefix)$(ImageName)

docker_push_dev: docker_build
	@echo "docker push test"
	sudo docker tag {{.ImageName}} {{.ImagesPrefix}}{{.ImageTestName}}
	sudo docker push $(ImagesPrefix)$(ImageTestName)

docker_build: build_linux
	@echo "构建docker镜像"
	sudo docker build -t $(ImageName) .

