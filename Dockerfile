FROM image.aithu.com/dev/builder-go:1.18-bullseye as builder
COPY . .
RUN bash build.sh common-dev /common-dev
