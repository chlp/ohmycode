<?php

use app\Session;
use app\Utils;

$input = require __DIR__ . '/actions.php';

$sessionId = (string)($input['session'] ?? '');
if (!Utils::isUuid($sessionId)) {
    error('Invalid: session');
}

$userId = (string)($input['user'] ?? '');
if (!Utils::isUuid($userId)) {
    error('Invalid: user');
}

$userName = (string)($input['userName'] ?? '');

$action = (string)($input['action'] ?? '');
switch ($action) {
    case 'getUpdate':
        $isKeepAlive = (bool)($input['isKeepAlive'] ?? false);
        $keepAliveRequestTimeSec = 30;
        if ($isKeepAlive) {
            ini_set('max_execution_time', $keepAliveRequestTimeSec + 3);
        }
        $lastUpdate = (isset($input['lastUpdate']) && is_string($input['lastUpdate'])) ? $input['lastUpdate'] : null;
        $lastInCycleUpdateTime = 0;
        while (true) {
            $session = Session::get($sessionId, $lastUpdate);
            if ($session !== null) {
                break;
            } else {
                if ($lastUpdate !== null) {
                    $currentTime = microtime(true);
                    if ($currentTime - $lastInCycleUpdateTime >= 1) {
                        // updating max one time per second
                        $lastInCycleUpdateTime = $currentTime;
                        Session::updateUserOnline($sessionId, $userId);
                        Session::cleanupUsers($sessionId);
                        Session::cleanupWriter($sessionId);
                    }
                }
                if (Utils::timer() > $keepAliveRequestTimeSec) {
                    return;
                }
                if (connection_status() !== CONNECTION_NORMAL) {
                    return;
                }
                usleep(200000); // 0.2 sec
            }
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
            Session::updateUserOnline($sessionId, $userId);
        }
        Session::cleanupUsers($sessionId);
        Session::cleanupWriter($sessionId);
        echo $session->getJson();
        break;
    case 'setSessionName':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setSessionName((string)($input['sessionName'] ?? ''))) {
            error('Wrong session name');
        }
        break;
    case 'setUserName':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setUserName($userId, $userName)) {
            error('Wrong user name');
        }
        break;
    case 'setLang':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setLang((string)($input['lang'] ?? ''))) {
            error('Wrong lang');
        }
        break;
    case 'setRunner':
        $session = getSession($sessionId, $userId, $userName);
        if (!$session->setRunner((string)($input['runner'] ?? ''))) {
            error('Wrong runner');
        }
        break;
    case 'setCode':
        $session = getSession($sessionId, $userId, $userName);
        if ($session->writer !== '' && $session->writer !== $userId) {
            error('Temporary forbidden 1', 403);
        }
        if (!$session->setCode((string)($input['code'] ?? ''), $userId)) {
            error('Temporary forbidden 2', 403);
        }
        break;
    default:
        error('wrong action', 404);
}

function getSession(string $sessionId, string $userId, string $userName): Session
{
    $session = Session::get($sessionId);
    if ($session === null) {
        $session = Session::createNew($sessionId);
        $session->writer = $userId;
        $session->insert();
        $session->setUserName($userId, $userName);
    }
    return $session;
}
