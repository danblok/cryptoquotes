all:
	@./init.sh

up:
	docker compose up -d --build
