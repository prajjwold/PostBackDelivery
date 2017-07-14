<?php
include_once("./redis_conn_params.php");

//Make sure that it is a POST request.
if(strcasecmp($_SERVER['REQUEST_METHOD'], 'POST') != 0){
    throw new Exception('Request method must be POST!');
}
 
//Make sure that the content type of the POST request has been set to application/json
$contentType = isset($_SERVER["CONTENT_TYPE"]) ? trim($_SERVER["CONTENT_TYPE"]) : '';
if(strcasecmp($contentType, 'application/json') != 0){
    throw new Exception('Content type must be: application/json');
}

//Receive RAW POST JSON Data
$content = trim(file_get_contents("php://input"));

//Decode the RAW POST JSON Data
$decoded = json_decode($content, true);
if(! is_array($decoded)){
    throw new Exception("Decoding JSON data failed");
}

//Connecting to Redis server on localhost
 try{
      $redis = new Redis();
      $redis->connect(REDIS_URL, REDIS_PORT);
      $redis->auth(REDIS_PASSWORD);
   }
catch(Exception $e){
     die($e->getMessage());
}
echo "Starting Redis Connection".PHP_EOL;
echo "Connection to server sucessfully".PHP_EOL;
//check whether server is running or not
echo "Server is running: ".$redis->ping().PHP_EOL;

if (isset($decoded['data']) && isset($decoded['endpoint'])) {
    foreach($decoded['data'] as &$data) {
        $postback = array(
            "endpoint" => $decoded['endpoint'],
            "data" => $data,
        );
        $redis->lPush(REDIS_MESSAGE_QUEUE, json_encode($postback));
    }
} else {
    echo 'No data received.' . PHP_EOL;
}

?>
