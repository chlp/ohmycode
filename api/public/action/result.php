<?php

use app\Request;
use app\Result;

$input = require __DIR__ . '/actions.php';

$action = (string)($input['action'] ?? '');
switch ($action) {
    case 'set':
        $requests = Request::get((string)($input['runner'] ?? ''), false, (string)($input['lang'] ?? ''), (string)($input['hash'] ?? ''));
        if (count($requests) === 0) {
            // no more need for result
            return;
        }
        $result = (string)($input['result'] ?? '');
        $result = substr($result, 0, 16384);
        Result::set($requests[0], $result);
        break;
    default:
        error('wrong action', 404);
}
