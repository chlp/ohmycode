<?php

require __DIR__ . '/app/bootstrap.php';

if (!isset($_POST['session'])) {
    http_response_code(400);
    echo 'Not found: session';
    return;
}
$sessionId = (string)$_POST['session'];
if (!Utils::isUuid($sessionId)) {
    http_response_code(400);
    echo 'Invalid: session';
    return;
}

if (!isset($_POST['user'])) {
    http_response_code(400);
    return;
}
$userId = (string)$_POST['user'];
if (!Utils::isUuid($userId)) {
    http_response_code(400);
    echo 'Invalid: user';
    return;
}

$lastUpdate = isset($_POST['lastUpdate']) ? (string)$_POST['lastUpdate'] : null;

$session = Session::get($sessionId, $userId, $lastUpdate);
