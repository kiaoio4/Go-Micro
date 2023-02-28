GROUP_NAME:=platform
PROJECT_NAME:=go-micro
VERSION:=dev

.PHONY: image data run

image:
	docker build -t ${GROUP_NAME}/${PROJECT_NAME}:${VERSION} .

data:
	docker run --rm \
		-p 80:80 \
		-v /root/go-micro/simulate:/simulate \
		${GROUP_NAME}/${PROJECT_NAME}:${VERSION} \
		public -w /simulate
run:
	docker run --rm \
		-p 82:82 \
		${GROUP_NAME}/${PROJECT_NAME}:${VERSION}