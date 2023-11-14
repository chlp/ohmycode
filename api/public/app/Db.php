<?php

class Db
{
    public function __construct(
        public mysqli $conn,
    )
    {
    }

    static public function get(): self
    {
        static $db = null;
        if ($db !== null) {
            return $db;
        }
        $conf = Conf::loadApiConf()['db'];
        $conn = mysqli_connect($conf['servername'], $conf['username'], $conf['password'], $conf['dbname'], $conf['port']);
        if (!$conn) {
            die("mysql connection failed: " . mysqli_connect_error());
        }
        $db = new self($conn);
        return $db;
    }

    public function select(string $query, ?array $params): array
    {
        $stmt = $this->conn->prepare($query);
        if (!$stmt) {
            die('wrong stmt');
        }
        if ($params !== null) {
            $types = '';
            $vars = [];
            foreach ($params as $param) {
                if (is_string($param)) {
                    $types .= 's';
                } else {
                    die('wrong type: ' . gettype($param));
                }
                $vars[] = $param;
            }
            $stmt->bind_param($types, ...$vars);
        }
        $stmt->execute();
        $stmtRes = $stmt->get_result();
        $result = [];
        while ($row = $stmtRes->fetch_row()) {
            $result[] = $row;
        }
        $stmt->close();
        return $result;
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
