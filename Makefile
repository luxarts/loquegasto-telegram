docker-up:
	docker build -t lqg-tg .
	docker run --name lqg-tg lqg-tg
docker-down:
	docker rm lqg-tg || true
	docker rmi lqg-tg || true

docker:
	make docker-down
	make docker-up