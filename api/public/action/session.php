<?php

require __DIR__ . '/../app/bootstrap.php';

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

$userName = (string)$_POST['userName'];

$action = (string)$_POST['action'] ?? '';
switch ($action) {
    case 'getUpdate':
        $lastUpdate = isset($_POST['lastUpdate']) ? (string)$_POST['lastUpdate'] : null;
        $session = Session::getById($sessionId, $lastUpdate);
        if ($session === null) {
            http_response_code(400);
            echo 'Not found';
            return;
        }
        $userFound = false;
        foreach ($session->users as $user) {
            if ($user['id'] === $userId) {
                $userFound = true;
                break;
            }
        }
        if (!$userFound) {
            $session->setUserName($userId, $userName);
        } else {
            $session->updateUserOnline($userId);
        }
        $session->removeOldUsers();
        echo $session->getJson();
        break;
    case 'setSessionName':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setSessionName((string)$_POST['sessionName'] ?? '')) {
            http_response_code(400);
            echo 'Wrong session name';
            return;
        }
        break;
    case 'setUserName':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setUserName($userId, $userName)) {
            http_response_code(400);
            echo 'Wrong user name';
            return;
        }
        break;
    case 'setLang':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setLang((string)$_POST['lang'] ?? '')) {
            http_response_code(400);
            echo 'Wrong lang';
            return;
        }
        break;
    case 'setExecutor':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setExecutor((string)$_POST['executor'] ?? '')) {
            http_response_code(400);
            echo 'Wrong executor';
            return;
        }
        break;
    case 'setCode':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setCode((string)$_POST['code'] ?? '')) {
            http_response_code(400);
            echo 'Wrong lang';
            return;
        }
        break;
}

function getSession(int $sessionId, int $userId, string $userName): Session
{
    $session = Session::getById($sessionId);
    if ($session === null) {
        $session = Session::createNew($sessionId);
        $session->writer = $userId;
        $session->insert();
        $session->setUserName($userId, $userName);
    }
    return $session;
}
