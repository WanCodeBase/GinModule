# GinModule
SimpleBank Based on https://github.com/techschool/simplebank
## Docker Setting
https://hub.docker.com/_/postgres
1. **get image**: docker pull {image}:{tag}  
> docker pull postgres:12-alpine 
2. **start a container**: docker run --name {container_name} -e {environment_variable} -p {host_ports:container_ports} -d {image}:{tag}
> docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine
3. **connect to container**: docker exec -it {container_name_or_id} {command} [args]
> docker exec -it postgres12 /bin/sh
> docker exec -it postgres12 psql -U root

## database migrate