<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="image/png" href="/favicon.png">

    <title>OhMyCode</title>
</head>

<body>

<?php
require 'db.php';
$conn = dbConn();
$sql = "SELECT * FROM `sessions`";
$result = mysqli_query($conn, $sql);
while ($row = mysqli_fetch_assoc($result)) {
    var_dump($row);
}
?>
</body>

</html>
