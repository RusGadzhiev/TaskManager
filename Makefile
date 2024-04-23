.PHONY: start
start:
	mysql_pass="pass" docker compose up --build

.DEFAULT_GOAL := start
