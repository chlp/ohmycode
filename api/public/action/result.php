<?php

$input = require __DIR__ . '/actions.php';

$action = (string)($input['action'] ?? '');
switch ($action) {
    case 'set':
        $requests = Request::get((string)($input['executor'] ?? ''), (string)($input['lang'] ?? ''), (string)($input['hash'] ?? ''));
        if (count($requests) === 0) {
            // no more need for result
            return;
        }
        Result::set($requests[0], (string)($input['result'] ?? ''));
        break;
    default:
        error('wrong action', 404);
}
