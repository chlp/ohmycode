<?php

class Session
{
    private const DEFAULT_LANG = 'php82';

    private Db $db;

    public function __construct(
        public string        $id,
        public string        $name,
        public string        $code,
        public string        $lang,
        public string        $executor,
        public DateTime|null $executorCheckedAt,
        public DateTime      $updatedAt,
        public string        $writer,
    )
    {
        $this->db = Db::get();
    }

    static public function createNew(string $writer): ?self
    {
        if (!Utils::isUuid($writer)) {
            return null;
        }
        $id = Utils::genUuid();
        return new self($id, '', '', self::DEFAULT_LANG, '', null, new DateTime(), $writer);
    }

    static public function getById(string $id): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        $res = Db::get()->select(
            "SELECT `name`, `code`, `lang`, `executor`, `executor_checked_at`, `updated_at`, `writer` FROM `sessions` WHERE `id` = ?",
            [$id]
        );
        if (count($res) === 0) {
            return null;
        }
        [$sessionName, $code, $lang, $executor, $executorCheckedAtStr, $updatedAtStr, $writer] = $res[0];
        $executorCheckedAt = null;
        if ($executorCheckedAtStr !== null) {
            $executorCheckedAt = DateTime::createFromFormat('Y-m-d H:i:s', $executorCheckedAtStr);
        }
        $updatedAt = DateTime::createFromFormat('Y-m-d H:i:s.u', $updatedAtStr);
        return new self($id, $sessionName, $code, $lang, $executor, $executorCheckedAt, $updatedAt, $writer);
    }
}