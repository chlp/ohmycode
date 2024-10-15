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

    public static function set(Session $session): void
    {
        $query = "INSERT INTO `requests` SET `session` = ?, `runner` = ?, `code` = ?, `lang` = ?
                       ON DUPLICATE KEY UPDATE `runner` = ?, `code` = ?, `lang` = ?, `received` = 0";
        Db::get()->exec($query, [
            $session->id, $session->runner, $session->code, $session->lang ?? Session::DEFAULT_LANG,
            $session->runner, $session->code, $session->lang ?? Session::DEFAULT_LANG,
        ]);
        Session::updateTime($session->id);
        Result::remove($session->id);
    }

    /**
     * @param string $runner
     * @param bool $isPublic
     * @param string|null $lang
     * @param string|null $hash
     * @return self[]
     */
    public static function get(string $runner, bool $isPublic, ?string $lang = null, ?string $hash = null): array
    {
        if (!Utils::isUuid($runner)) {
            return [];
        }
        $query = "SELECT `session`, `code`, `lang`, md5(`code`) as `hash` FROM `requests` WHERE `runner` = ?";
        if ($isPublic) {
            $params = [''];
        } else {
            $params = [$runner];
        }
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

    public static function markReceived(string $runner, bool $isPublic, string $lang, string $hash): void
    {
        if (!Utils::isUuid($runner)) {
            return;
        }
        if ($isPublic) {
            $query = "UPDATE `requests` SET `received` = 1, `runner` = ? WHERE `runner` = '' and `lang` = ? and md5(`code`) = ?";
        } else {
            $query = "UPDATE `requests` SET `received` = 1 WHERE `runner` = ? and `lang` = ? and md5(`code`) = ?";
        }
        Db::get()->exec($query, [$runner, $lang, $hash]);
    }

    public static function remove(string $runner, string $lang, string $hash): void
    {
        if (!Utils::isUuid($runner)) {
            return;
        }
        $query = "DELETE FROM `requests` WHERE `runner` = ? and `lang` = ? and md5(`code`) = ?";
        Db::get()->exec($query, [$runner, $lang, $hash]);
    }
}
