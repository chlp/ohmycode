<?php

namespace app;

class Request
{
    public function __construct(
        public string $session,
        public string $runner,
        public string $code,
        public string $lang,
        public string $hash,
    )
    {
    }

    const MAX_REQUESTS_FOR_RUNNER_PER_REQUEST = 5;

    static public function set(Session $session): void
    {
        $query = "INSERT INTO `requests` SET `session` = ?, `runner` = ?, `code` = ?, `lang` = ?
                       ON DUPLICATE KEY UPDATE `runner` = ?, `code` = ?, `lang` = ?, `received` = 0";
        Db::get()->exec($query, [
            $session->id, $session->runner, $session->code, $session->lang,
            $session->runner, $session->code, $session->lang,
        ]);
        $session->updateTime();
        Result::remove($session->id);
    }

    /**
     * @param string $runner
     * @param string|null $lang
     * @param string|null $hash
     * @return self[]
     */
    static public function get(string $runner, ?string $lang = null, ?string $hash = null): array
    {
        if (!Utils::isUuid($runner)) {
            return [];
        }
        $query = "SELECT `session`, `code`, `lang`, md5(`code`) as `hash` FROM `requests` WHERE `runner` = ?";
        $params = [$runner];
        if ($lang !== null && $hash !== null) {
            $query .= " AND `lang` = ? AND md5(`code`) = ?";
            $params[] = $lang;
            $params[] = $hash;
        } else {
            $query .= " AND `received` = 0";
        }
        $query .= " LIMIT ?";
        $params[] = self::MAX_REQUESTS_FOR_RUNNER_PER_REQUEST;
        $res = Db::get()->select($query, $params);
        $requests = [];
        foreach ($res as $row) {
            $requests[] = new self($row[0], $runner, $row[1], $row[2], $row[3]);
        }
        return $requests;
    }

    static public function markReceived(string $runner, string $lang, string $hash): void
    {
        if (!Utils::isUuid($runner)) {
            return;
        }
        $query = "UPDATE `requests` SET `received` = 1 WHERE `runner` = ? and `lang` = ? and md5(`code`) = ?";
        Db::get()->exec($query, [$runner, $lang, $hash]);
    }

    static public function remove(string $runner, string $lang, string $hash): void
    {
        if (!Utils::isUuid($runner)) {
            return;
        }
        $query = "DELETE FROM `requests` WHERE `runner` = ? and `lang` = ? and md5(`code`) = ?";
        Db::get()->exec($query, [$runner, $lang, $hash]);
    }
}