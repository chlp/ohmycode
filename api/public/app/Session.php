<?php

class Session
{
    private const DEFAULT_LANG = 'php82';
    private const CODE_MAX_LENGTH = 32768;
    private Db $db;
    public const LANGS = [
        'php82' => [
            'name' => 'PHP 8.2',
            'highlighter' => 'php',
        ],
        'mysql8' => [
            'name' => 'MySQL 8',
            'highlighter' => 'sql',
        ],
        'go' => [
            'name' => 'GoLang',
            'highlighter' => 'go',
        ],
    ];

    public function __construct(
        public string    $id,
        public string    $name,
        public string    $code,
        public string    $lang,
        public string    $executor,
        public ?DateTime $executorCheckedAt,
        public ?DateTime $updatedAt,
        public string    $writer,
        public array     $users,
        public bool      $isWaitingForResult,
        public string    $result,
    )
    {
        $this->db = Db::get();
    }

    public function getJson(): string
    {
        return json_encode([
            'id' => $this->id,
            'name' => $this->name,
            'code' => $this->code,
            'lang' => $this->lang,
            'executor' => $this->executor,
            'isExecutorOnline' => $this->isExecutorOnline(),
            'updatedAt' => $this->updatedAt,
            'writer' => $this->writer,
            'users' => $this->users,
            'isWaitingForResult' => $this->isWaitingForResult,
            'result' => $this->result,
        ]);
    }

    static public function createNew(string $id): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        $name = date('Y-m-d');
        return new self($id, $name, '', self::DEFAULT_LANG, '', null, null, '', [], false, '');
    }

    static public function getById(string $id, ?string $updatedAfter = null): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        // todo: join request, result and users
        $query = "SELECT `name`, `code`, `lang`, `executor`, `executor_checked_at`, `updated_at`, `writer` FROM `sessions` WHERE `id` = ?";
        $params = [$id];
        if ($updatedAfter !== null) {
            $query .= "updated_at > ?";
            $params[] = $updatedAfter;
        }
        $res = Db::get()->select($query, $params);
        if (count($res) === 0) {
            return null;
        }
        [$sessionName, $code, $lang, $executor, $executorCheckedAtStr, $updatedAtStr, $writer] = $res[0];
        $executorCheckedAt = null;
        if ($executorCheckedAtStr !== null) {
            $executorCheckedAt = DateTime::createFromFormat('Y-m-d H:i:s', $executorCheckedAtStr);
        }
        $updatedAt = DateTime::createFromFormat('Y-m-d H:i:s.u', $updatedAtStr);

        $session = new self($id, $sessionName, $code, $lang, $executor, $executorCheckedAt, $updatedAt, $writer, [], false, '');
        $session->loadUsers();

        return $session;
    }

    private function loadUsers(): void
    {
        $res = $this->db->select("select `user`, `name` from `session_users` where session = ?", [$this->id]);
        $users = [];
        foreach ($res as $row) {
            $users[] = [
                'id' => $row[0],
                'name' => $row[1],
            ];
        }
        $this->users = $users;
    }

    static public function updateUserOnline(string $sessionId, string $userId): void
    {
        if (Utils::isUuid($sessionId) && Utils::isUuid($userId)) {
            Db::get()->exec("update `session_users` set updated_at = NOW(3) where session = ? and user = ?", [$sessionId, $userId]);
        }
    }

    static public function removeOldUsers(string $sessionId): void
    {
        if (Utils::isUuid($sessionId)) {
            Db::get()->exec("delete from `session_users` where `session` = ? and `updated_at` < NOW(3) - INTERVAL 20 second", [$sessionId]);
        }
    }

    public function insert(): self
    {
        $query = "INSERT INTO `sessions` SET `name` = ?, `code` = ?, `lang` = ?, `executor` = ?, `writer` = ?, `id` = ?;";
        $this->db->exec($query, [$this->name, $this->code, $this->lang, $this->executor, $this->writer, $this->id]);
        return self::getById($this->id);
    }

    public function setSessionName(string $name): bool
    {
        if (!Utils::isValidString($name)) {
            return false;
        }
        $query = "UPDATE `sessions` SET `name` = ? WHERE `id` = ?";
        $this->db->exec($query, [$name, $this->id]);
        return true;
    }

    public function setLang(string $lang): bool
    {
        if (!isset(self::LANGS[$lang])) {
            return false;
        }
        $query = "UPDATE `sessions` SET `lang` = ? WHERE `id` = ?";
        $this->db->exec($query, [$lang, $this->id]);
        return true;
    }

    public function setCode(string $code): bool
    {
        if (strlen($code) > self::CODE_MAX_LENGTH) {
            return false;
        }
        $query = "UPDATE `sessions` SET `code` = ? WHERE `id` = ?";
        $this->db->exec($query, [$code, $this->id]);
        return true;
    }

    public function updateTime(): void
    {
        $query = "UPDATE `sessions` SET `updated_at` = NOW(3) WHERE `id` = ?;";
        $this->db->exec($query, [$this->id]);
    }

    public function setUserName(string $userId, string $name): bool
    {
        if (!Utils::isUuid($userId)) {
            return false;
        }
        if (!Utils::isValidString($name)) {
            return false;
        }
        $query = "INSERT INTO `session_users` SET `session` = ?, `user` = ?, `name` = ? ON DUPLICATE KEY UPDATE `name` = ?";
        $this->db->exec($query, [$this->id, $userId, $name, $name]);
        $this->updateTime();
        return true;
    }

    public function setExecutor(string $executor): bool
    {
        if (!Utils::isUuid($executor)) {
            return false;
        }
        $query = "UPDATE `sessions` SET `executor` = ? WHERE `id` = ?";
        $this->db->exec($query, [$executor, $this->id]);
        return true;
    }

    public function setWriter(string $userId): bool
    {
        if (!Utils::isUuid($userId)) {
            return false;
        }
        $query = "UPDATE `sessions` SET `writer` = ? WHERE `id` = ?";
        $this->db->exec($query, [$userId, $this->id]);
        return true;
    }

    public function isExecutorOnline(): bool
    {
        if ($this->executorCheckedAt === null) {
            return false;
        }
        return time() - $this->executorCheckedAt->getTimestamp() < 10;
    }
}