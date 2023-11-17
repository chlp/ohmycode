<?php

class Request
{
    public function __construct(
        public string $session,
        public string $executor,
        public string $code,
        public string $lang,
        public string $hash,
    )
    {
    }

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

    /**
     * @param string $executor
     * @param string|null $lang
     * @param string|null $hash
     * @return self[]
     */
    static public function get(string $executor, ?string $lang = null, ?string $hash = null): array
    {
        if (!Utils::isUuid($executor)) {
            return [];
        }
        $query = "SELECT `session`, `code`, `lang`, md5(`code`) as `hash` FROM `requests` WHERE `executor` = ?";
        $params = [$executor];
        if ($lang !== null) {
            $query = " AND `lang` = ?";
            $params[] = $lang;
        }
        if ($hash !== null) {
            $query = " AND md5(`code`) = ?";
            $params[] = $hash;
        }
        $query .= " LIMIT ?";
        $params[] = self::MAX_REQUESTS_FOR_EXECUTOR_PER_REQUEST;
        $res = Db::get()->select($query, $params);
        $requests = [];
        foreach ($res as $row) {
            $requests[] = new self($row[0], $executor, $row[1], $row[2], $row[3]);
        }
        return $requests;
    }

    static public function remove(string $executor, string $lang, string $hash): void
    {
        $query = "DELETE FROM `requests` WHERE `executor` = ? and `lang` = ? and md5(`code`) = ?";
        Db::get()->exec($query, [$executor, $lang, $hash]);
    }
}