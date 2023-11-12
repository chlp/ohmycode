<?php

if (!isset($_GET['id']) && !isset($_POST['id'])) {
    http_response_code(400);
    return;
}

require 'db.php';
$conn = dbConn();

if (isset($_GET['id'])) {
    $id = $_GET['id'];
    $sql = "SELECT * FROM `sessions` WHERE `id` = $id"; // $id as $1
    $result = mysqli_query($conn, $sql);
    while ($row = mysqli_fetch_assoc($result)) {
        var_dump($row);
    }
    // get session
    // get request
    // get result
    // if not exist -> create
} else {
    $id = $_POST['id'];
    switch ($_POST['type']) {
        case 'lang':
            $lang = $_POST['lang'];
            $sql = "UPDATE `sessions` SET `` WHERE `id` = $id";
            $result = mysqli_execute_query($conn, $sql, []);
            break;
        case 'executor':
            $executor = $_POST['executor'];
            $sql = "UPDATE `sessions` SET `` WHERE `id` = $id";
            $result = mysqli_execute_query($conn, $sql, []);
            break;
        case 'code':
            $code = $_POST['code'];
            $sql = "UPDATE `sessions` SET `` WHERE `id` = $id";
            $result = mysqli_execute_query($conn, $sql, []);
            break;
        case 'executor_check':
            $executor = $_POST['executor'];
            $sql = "UPDATE `sessions` SET `executor_checked_at` = NOW() WHERE `id` = $id";
            $result = mysqli_execute_query($conn, $sql, []);
            break;
        default:
            http_response_code(400);
            return;
    }
}