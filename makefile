NAME = endocode
TAG = latest
#INSTANCE = endocode

.PHONY: default build copy debug clean push buildrelease

default: run

build:
	docker build -t $(NAME):$(TAG) -f Dockerfile .

clean:
	docker rm $(NAME)

debug:
	docker run --rm -it --name $(NAME) /bin/bash

run:
	docker run --rm -p 8080:8080 --name $(NAME) $(NAME):$(TAG)

push:
	docker push $(NAME):$(TAG)