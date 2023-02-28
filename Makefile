.PHONY: image
image:
	docker build -t image.aithu.com/platform/common:dev .

.PHONY: run
run:
	docker run --rm image.aithu.com/platform/common:dev
	