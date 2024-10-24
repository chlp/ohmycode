<?php

use app\Request;
use app\Session;
use app\Utils;

$input = require __DIR__ . '/actions.php';

$action = (string)($input['action'] ?? '');
switch ($action) {
    case 'set':
        $session = Session::get((string)($input['session'] ?? ''));
        if ($session === null) {
            error('No session');
        }
        if (!$session->runnerIsOnline()) {
            error('Runner is not ready');
        }
        Request::set($session);
        break;
    case 'mark_received':
        Request::markReceived((string)($input['runner'] ?? ''), (bool)($input['is_public'] ?? false), (string)($input['lang'] ?? ''), (string)($input['hash'] ?? ''));
        break;
    case 'get':
        $isKeepAlive = (bool)$input['is_keep_alive'];
        $keepAliveRequestTimeSec = 10;
        if ($isKeepAlive) {
            ini_set('max_execution_time', $keepAliveRequestTimeSec + 3);
        }
        $runner = (string)($input['runner'] ?? '');
        $isPublic = (bool)($input['is_public'] ?? false);
        if (!Utils::isUuid($runner)) {
            error('not valid runner', 404);
        }
        $requests = [];
        $lastInCycleUpdateTime = 0;
        while (true) {
            $currentTime = microtime(true);
            if ($currentTime - $lastInCycleUpdateTime >= 1) {
                // updating max one time per second
                $lastInCycleUpdateTime = $currentTime;
                Session::updateRunnerCheckedAt($runner, $isPublic);
            }
            $requests = Request::get($runner, $isPublic);
            if (!$isKeepAlive) {
                break;
            }
            if (count($requests) > 0) {
                break;
            }
            if (Utils::timer() > $keepAliveRequestTimeSec) {
                break;
            }
            if (connection_status() !== CONNECTION_NORMAL) {
                break;
            }
            usleep(100000); // 0.1 sec
        }
        $output = [];
        foreach ($requests as $request) {
            $output[] = [
                'code' => $request->code,
                'lang' => $request->lang,
                'hash' => $request->hash,
            ];
        }
        echo json_encode($output);
        break;
    default:
        error('wrong action', 404);
}
