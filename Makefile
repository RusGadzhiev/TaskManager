.PHONY: start
start:
	mysql_pass="pass" server_pass="pass_serv" docker compose up --build

.DEFAULT_GOAL := start
