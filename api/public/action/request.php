<?php

$input = require __DIR__ . '/actions.php';

$action = (string)($input['action'] ?? '');
switch ($action) {
    case 'set':
        $session = Session::get((string)($input['session'] ?? ''));
        if ($session === null) {
            error('No session');
        }
        if (!$session->isExecutorOnline()) {
            error('Executor is not ready');
        }
        Request::set($session);
        break;
    case 'markReceived':
        Request::markReceived((string)($input['executor'] ?? ''), (string)($input['lang'] ?? ''), (string)($input['hash'] ?? ''));
        break;
    case 'get':
        $executor = (string)($input['executor'] ?? '');
        Session::setCheckedByExecutor($executor);
        $requests = Request::get($executor);
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
