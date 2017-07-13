GCC Compiler:
-------------------
```shell
sudo apt-get install gcc
```

Make:
-------------------
```shell
sudo apt-get install make
```

TCL:
-------------------
```shell
sudo apt-get install tcl
```

Apache Server:
--------------------
```shell
sudo apt-get install apache2
```

PHP:
-----------------------------
```shell
$ sudo apt install php libapache2-mod-php php7.0-dev
```

Install Redis Server:
---------------------------
```shell
$ wget http://download.redis.io/redis-stable.tar.gz
$ tar xvzf redis-stable.tar.gz
$ cd redis-stable
$ make
$ make test
$ sudo make install
```

#This will cause redis-server to start on boot

$ sudo ./utils/install_server.sh

	Port           : 9999
	Config file    : /etc/redis/9999.conf
	Log file       : /var/log/redis_9999.log
	Data dir       : /var/lib/redis/9999
	Executable     : /usr/local/bin/redis-server
	Cli Executable : /usr/local/bin/redis-cli

Start and Stop Redis Server manually:
------------------------------------
$ sudo service redis_9999 start
$ sudo service redis_9999 stop
$ /etc/init.d/redis_9999 start
$ /etc/init.d/redis_9999 stop

Set the password:
------------------------------------
```shell
$ redis-cli -p 9999
$ redis 127.0.0.1:9999> CONFIG SET requirepass <secret_password>
$ redis 127.0.0.1:9999> quit
$ sudo service redis_9999 restart
$ redis-cli -p 9999
$ redis 127.0.0.1:9999> auth <secret_password>
```

To set the password, edit your redis.conf file, Uncomment the line and set the password

#requirepass foobared

```shell
$ nano /etc/redis/9999.conf
	requirepass redis_p@ssw0rd
```

Install GIT:
-----------------------
```shell
$ sudo apt-get install git
```

Install Redis PHP driver:
--------------------------------
```shell
$ git clone https://github.com/nicolasff/phpredis
$ git clone https://github.com/phpredis/phpredis
$ sudo apt-get install php7.0-dev
$ cd phpredis 
$ sudo phpize 
$ sudo ./configure 
$ sudo make 
$ sudo make install 
```

Locate php-extension directory:
------------------------------------
```shell
$ php-config --extension-dir
```

	/usr/lib/php/20151012

Copy the redis.so extension into php-extension directory
--------------------------------------------------------
```shell
$ cp phpredis/modules/redis.so /usr/lib/php/20151012
```

Edit the php.ini to include the extension
-------------------------------------------
```shell
$ nano /etc/php/7.0/apache2/php.ini
```
	# Add the line
	extension=redis.so

Installing GO Tools:
-------------------------
```shell
$ wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz
```

Configure GOROOT and GOPATH env variables
------------------------------------------
```shell
$ vi /etc/profile
```
#Add these line
```shell
 export GOROOT=/usr/local/go
 export GOPATH=/root/go
 export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Refresh the Profile:
--------------------
```shell
$ source /etc/profile
```

Redis Client for GO lang:
------------------------------------
```shell
$ sudo apt-get install git
$ go get github.com/garyburd/redigo/redis
```