# R1wallet

This project is my code challenge for ArvanCloud company interview.<br>
This is a voucher platform that can handle credit vouchers without ACL.<br>
It uses redis pub/sub as a queue for communicate with wallet to apply credit.<br>
I tried to design this project based on SOLID principles.<br>
To run the project make sure you have installed docker and docker-compose and it's running.<br>
Use following command to check docker installed and it's running:
``` shell script
$docker -v
```
You should see something like:
``` shell script
 Docker version 19.03.13, build 4484c46d9d
```

After that check that you have installed docker-compose as well with following command:
```shell script
$docker-compose -v
```
You should see something like:
``` shell script
docker-compose version 1.27.4, build 40524192
```

Make a copy from `env.yaml.default` and rename it to `env.yml`:
```shell script
$cp env.yaml.default env.yaml
```

After that you have to use following command to run project:
```shell script
$docker-compose up -d
```

Wait till application starts after that use following command to see forwarded ports:
```shell script
$docker-compose ps
```
You will see something like:
```shell script
       Name                     Command               State            Ports         
-------------------------------------------------------------------------------------
r1wallet_app         ./entrypoint.sh bash -c ./ ...   Up      0.0.0.0:32852->1323/tcp
r1wallet_mariadb_1   docker-entrypoint.sh mysqld      Up      0.0.0.0:32849->3306/tcp
r1wallet_redis_1     docker-entrypoint.sh redis ...   Up      0.0.0.0:32848->6379/tcp

```

It shows that you application is running on port 1323 and forward to port 32852 of your local machine<br>
