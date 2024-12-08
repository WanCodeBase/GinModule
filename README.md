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

## database 

### ACID Property
Atomicity  
Consistency  
Isolation  
Durabilty

### CURD
**database/sql**: fast but mistakes cannot be caught until runtime  
**gorm**: low code but slow on high load (gorm的运行速度比标准库慢3-5倍)  
**sqlx**: fast & easy but mistakes cannot be caught until runtime  
**sqlc**: automatic code generation, especially for Postgres

### Isolation Level 事务隔离级别
> In Mysql:   
> select @@transaction_isolation  
> set session transaction isolation level {isolation_level}  

> In Postgres:  
> show transaction isolation level;  
> // **only can set transaction level within trasaction**  
> set transaction isolation level {isolation level}

**isolation level**:  
- read uncommitted  
- read committed: stop dirty read  
- repeatable read(Mysql default): stop dirty read and unrepeatable read  
- serializable: stop dirty read, unrepeated read and serialization anomaly(序列化异常)

**In Postgres**
- read uncommitted mode behaves like read committed. 
- Using dependence detection  
- default level: RC

**In Mysql**  
- default level: RR
- Using locking mechanism  

## Github Workflow

1. Actions on the website
2. Coding .yml file to decide steps

## gomock
1. install: https://github.com/golang/mock
2. vi ~/.zshrc & add export PATH=$PATH:~/go/bin
3. source ~/.zshrc
4. using which mockgen to make sure mockgen is valid.