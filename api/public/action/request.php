<?php

require __DIR__ . '/../app/bootstrap.php';

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    error('Method not allowed', 405);
}

$input = json_decode(file_get_contents('php://input'), true);

$action = (string)($input['action'] ?? '');
switch ($action) {
    case 'set':
        $session = Session::getById((string)($input['session'] ?? ''));
        if ($session === null) {
            error('No session');
        }
        if (!$session->isExecutorOnline()) {
            error('Executor is not ready');
        }
        Request::set($session);
        break;
    case 'get':
        $executor = (string)($input['executor'] ?? '');
        echo json_encode(Request::get($executor));
        break;
    default:
        error('wrong action', 404);
}

function error($str, $code = 400): void
{
    http_response_code($code);
    die(json_encode(['error' => $str]));
}