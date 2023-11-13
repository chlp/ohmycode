<?php

function dbConn()
{
    require_once 'conf.php';
    $conf = loadApiConf()['db'];
    $conn = mysqli_connect($conf['servername'], $conf['username'], $conf['password'], $conf['dbname'], $conf['port']);
    if (!$conn) {
        die("mysql connection failed: " . mysqli_connect_error());
    }
    return $conn;
}
