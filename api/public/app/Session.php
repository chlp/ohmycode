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

    static public function createNew(string $id, string $writer): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        if (!Utils::isUuid($writer)) {
            return null;
        }
        $name = date('Y-m-d');
        return (new self($id, $name, '', self::DEFAULT_LANG, '', null, new DateTime(), $writer))->save(true);
    }

    static public function get(string $id, ?string $user = null, ?string $updatedAfter = null): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        if ($user !== null && !Utils::isUuid($id)) {
            return null;
        }
        $query = "SELECT `name`, `code`, `lang`, `executor`, `executor_checked_at`, `updated_at`, `writer` FROM `sessions` WHERE `id` = ?";
        $params = [$id];
        if ($updatedAfter !== null) {
            $query .= "updated_at > ?";
            $params[] = $updatedAfter;
        }
        $res = Db::get()->select($query, $params);
        if (count($res) === 0) {
            if ($updatedAfter === null) {
                return self::createNew($id, $user);
            }
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

    public function save(bool $new = false): self
    {
        if ($new) {
            $query = "INSERT INTO `sessions` SET `name` = ?, `code` = ?, `lang` = ?, `executor` = ?, `executor_checked_at` = ?, `writer` = ?, `id` = ?;";
        } else {
            $query = "UPDATE `sessions` SET `name` = ?, `code` = ?, `lang` = ?, `executor` = ?, `executor_checked_at` = ?, `writer` = ? WHERE `id` = ?;";
        }
        $this->db->exec($query, [$this->name, $this->code, $this->lang, $this->executor, $this->executorCheckedAt, $this->writer, $this->id]);
        return self::get($this->id);
    }
}