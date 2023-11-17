<?php

class Request
{
    const MAX_REQUESTS_FOR_EXECUTOR_PER_REQUEST = 5;

    static public function set(Session $session): void
    {
        $query = "INSERT INTO `requests` SET `session` = ?, `executor` = ?, `code` = ?, `lang` = ?
                       ON DUPLICATE KEY UPDATE `executor` = ?, `code` = ?, `lang` = ?";
        Db::get()->exec($query, [
            $session->id, $session->executor, $session->code, $session->lang,
            $session->executor, $session->code, $session->lang,
        ]);
        $session->updateTime();
    }

    static public function get(string $executor): array
    {
        $query = "SELECT `code`, md5(`code`) as `hash`, `lang` FROM `requests` WHERE `executor` = ? LIMIT ?";
        $res = Db::get()->select($query, [$executor, self::MAX_REQUESTS_FOR_EXECUTOR_PER_REQUEST]);
        $requests = [];
        foreach ($res as $row) {
            $requests[] = [
                'code' => $row[0],
                'hash' => $row[1],
                'lang' => $row[2],
            ];
        }
        return $requests;
    }

    static public function remove(string $executor, string $lang, string $hash): void
    {
        $query = "DELETE FROM `requests` WHERE `executor` = ? and `lang` = ? and md5(`code`) = ?";
        Db::get()->exec($query, [$executor, $lang, $hash]);
    }
}