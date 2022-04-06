build:
	docker build --rm --iidfile current-image-id.txt .
start:
	docker run -dt --env-file .env $(cat current-image-id.txt)
prune_builders:
	docker image prune --filter label=stage=builder