NAME = endocode
#INSTANCE = endocode

.PHONY: default build copy debug clean push buildrelease

default: buildrelease

build:
	docker build -t $(NAME) .

copy:
	docker create $(NAME)
	docker cp $(NAME):latest ./app
	docker rm $(NAME)

release:
	docker build -t $(NAME):$(TAG) -f Dockerfile .
	docker tag $(NAME):$(TAG) $(NAME):latest

buildrelease: build copy release

clean:
	docker rm $(NAME)

debug:
	docker run --rm -it --name $(NAME) /bin/bash

run:
	docker run --rm -p 8080:8080 --name $(NAME)

dev:
	docker run -it --rm -p 8080:8080 -w /app/$(NAME)

push:
	docker push $(NAME):$(TAG)
	docker push $(NAME):latest