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
$lang = (string)($input['lang'] ?? '');

$action = (string)($input['action'] ?? '');
switch ($action) {
    case 'get_update':
        $isKeepAlive = (bool)($input['is_keep_alive'] ?? false);
        $keepAliveRequestTimeSec = 30;
        if ($isKeepAlive) {
            ini_set('max_execution_time', $keepAliveRequestTimeSec + 3);
        }
        $lastUpdate = (isset($input['lastUpdate']) && is_string($input['lastUpdate'])) ? $input['lastUpdate'] : null;
        $lastInCycleUpdateTime = 0;
        while (true) {
            if ($lastUpdate !== null) {
                $session = Session::get($sessionId, $lastUpdate);
            } else {
                $session = getSession($sessionId, $userId, $userName, $lang);
            }
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
                if (connection_aborted()) {
                    return;
                }
                echo ' '; // flush-hack: to work with connection_status() and connection_aborted()
                flush();
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
    case 'set_session_name':
        $session = getSession($sessionId, $userId, $userName, $lang);
        if (!$session->setSessionName((string)($input['sessionName'] ?? ''))) {
            error('Wrong session name');
        }
        break;
    case 'set_user_name':
        $session = getSession($sessionId, $userId, $userName, $lang);
        if (!$session->setUserName($userId, $userName)) {
            error('Wrong user name');
        }
        break;
    case 'set_lang':
        $session = getSession($sessionId, $userId, $userName, $lang);
        if (!$session->setLang($lang)) {
            error('Wrong lang');
        }
        break;
    case 'set_runner':
        $session = getSession($sessionId, $userId, $userName, $lang);
        if (!$session->setRunner((string)($input['runner'] ?? ''))) {
            error('Wrong runner');
        }
        break;
    case 'set_code':
        $session = getSession($sessionId, $userId, $userName, $lang);
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

function getSession(string $sessionId, string $userId, string $userName, string $lang): Session
{
    $session = Session::get($sessionId);
    if ($session === null) {
        $session = Session::createNew($sessionId);
        $session->writer = $userId;
        if ($lang !== '') {
            $session->lang = $lang;
        } else {
            $session->lang = Session::DEFAULT_LANG;
        }
        $session->insert();
        $session->setUserName($userId, $userName);
    }
    return $session;
}
