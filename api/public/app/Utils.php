<?php

class Utils
{
    private const ID_LENGTH = 32;

    static public function genUuid(): string
    {
        return sprintf('%04x%04x%04x%04x%04x%04x%04x%04x',
            mt_rand(0, 0xffff), mt_rand(0, 0xffff),
            mt_rand(0, 0xffff),
            mt_rand(0, 0x0fff) | 0x4000,
            mt_rand(0, 0x3fff) | 0x8000,
            mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff)
        );
    }

    static public function isUuid(string $id): bool
    {
        return preg_match('/^[a-z0-9]{' . self::ID_LENGTH . '}$/', $id) === 1;
    }
}