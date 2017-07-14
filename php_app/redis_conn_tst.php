<?php
   //Connecting to Redis server on localhost
   try{
         $redis = new Redis();
         $redis->connect('127.0.0.1', 9999);
         $redis->auth("redis_p@ssw0rd");
   }
   catch(Exception $e){
        die($e->getMessage());
   }
   echo "Starting Redis Connection".PHP_EOL;
   echo "Connection to server sucessfully".PHP_EOL;
   //check whether server is running or not
   echo "Server is running: ".$redis->ping().PHP_EOL;
?>

