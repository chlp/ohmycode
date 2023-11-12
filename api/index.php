<?php

require 'db.php';
$conn = dbConn();

$sql = "SELECT * FROM `sessions`";
$result = mysqli_query($conn, $sql);
while ($row = mysqli_fetch_assoc($result)) {
    var_dump($row);
}