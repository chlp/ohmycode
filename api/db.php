<?php

// docker run --name local-mysql -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=rootpass --restart unless-stopped mysql:8

function dbConn()
{
    require 'conf.php';
    $conf = loadApiConf();
    $conn = mysqli_connect($conf['servername'], $conf['username'], $conf['password'], $conf['dbname']);
    if (!$conn) {
        die("mysql connection failed: " . mysqli_connect_error());
    }
    return $conn;
}
