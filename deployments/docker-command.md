## Команды для тестирования контейнеров


	docker build -t puzanovma/parser-bank -f ./deployments/scratch.Dockerfile .
	docker images
	docker run -d -p 7777:7777 -e TZ=Europe/Samara --name parser-bank puzanovma/parser-bank 
	docker run --rm -it -p 7777:7777 -e TZ=Europe/Samara --name parser-bank puzanovma/parser-bank 
	docker ps
	docker ps -a
	docker stop parser-bank


* Удаление недействительных образов  
`docker rmi $(docker images -f dangling=true -q)`
* Удаление всех остановленных контейнеров  
`docker rm $(docker ps -a -f status=exited -q)`
* Удаление недействительных томов (Docker 1.9 +)  
`docker volume rm $(docker volume ls -f dangling=true -q)`
