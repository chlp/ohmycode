<?php

namespace app;

use mysqli;
use mysqli_stmt;

class Db
{
    public function __construct(
        public mysqli $conn,
    ) {
    }

    public static function get(): self
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

    public function select(string $query, ?array $params = null): array
    {
        $stmt = $this->conn->prepare($query);
        if (!$stmt) {
            die('wrong select stmt');
        }
        $this->bindParams($query, $stmt, $params);
        $stmt->execute();
        $stmtRes = $stmt->get_result();
        $result = [];
        while ($row = $stmtRes->fetch_row()) {
            $result[] = $row;
        }
        $stmt->close();
        return $result;
    }

    public function exec(string $query, ?array $params): void
    {
        $stmt = $this->conn->prepare($query);
        if (!$stmt) {
            die('wrong exec stmt');
        }
        $this->bindParams($query, $stmt, $params);
        $stmt->execute();
        $stmt->close();
    }

    private function bindParams(string $query, mysqli_stmt $stmt, ?array $params): void
    {
        if ($params === null || count($params) === 0) {
            return;
        }
        $types = '';
        $vars = [];
        foreach ($params as $i => $param) {
            if (is_string($param)) {
                $types .= 's';
            } elseif (is_int($param)) {
                $types .= 'i';
            } else {
                die("wrong type: $i:" . gettype($param) . ' | ' . substr($query, 0, 50));
            }
            $vars[] = $param;
        }
        $stmt->bind_param($types, ...$vars);
    }
}
