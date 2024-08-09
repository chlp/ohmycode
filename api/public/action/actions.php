<?php

require __DIR__ . '/../../app/bootstrap.php';

function error($str, $code = 400): void
{
    http_response_code($code);
    die(json_encode(['error' => $str]));
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    error('Method not allowed', 405);
}

session_write_close();

return json_decode(file_get_contents('php://input'), true);
