<?php

class Db
{
    static public function dbConn()
    {
        static $conn = null;
        if ($conn !== null) {
            return $conn;
        }
        $conf = Conf::loadApiConf()['db'];
        $conn = mysqli_connect($conf['servername'], $conf['username'], $conf['password'], $conf['dbname'], $conf['port']);
        if (!$conn) {
            die("mysql connection failed: " . mysqli_connect_error());
        }
        return $conn;
    }
}
//
//---
//
//$stmt = mysqli_prepare($link, "INSERT INTO CountryLanguage VALUES (?, ?, ?, ?)");
//mysqli_stmt_bind_param($stmt, 'sssd', $code, $language, $official, $percent);
//
//$code = 'DEU';
//$language = 'Bavarian';
//$official = "F";
//$percent = 11.2;
//
//mysqli_stmt_execute($stmt);
//
//printf("%d row inserted.\n", mysqli_stmt_affected_rows($stmt));
//
//---
//
//$mysqli = new mysqli('localhost', 'my_user', 'my_password', 'world');
//
//$stmt = $mysqli->prepare("SELECT Language FROM CountryLanguage WHERE CountryCode IN (?, ?)");
///* Using ... to provide arguments */
//$stmt->bind_param('ss', ...['DEU', 'POL']);
//$stmt->execute();
//$stmt->store_result();
//
