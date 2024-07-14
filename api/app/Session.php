<?php

namespace app;

use DateTime;

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
        'postgres13' => [
            'name' => 'PostgreSQL 13',
            'highlighter' => 'sql',
        ],
        'java' => [
            'name' => 'Java',
            'highlighter' => 'text/x-java',
        ],
        'go' => [
            'name' => 'GoLang',
            'highlighter' => 'go',
        ],
    ];

    public function __construct(
        public string $id,
        public string $name,
        public string $code,
        public string $lang,
        public string $runner,
        public ?DateTime $runnerCheckedAt,
        public ?DateTime $updatedAt,
        public string $writer,
        public array $users,
        public bool $isWaitingForResult,
        public string $result,
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
            'runner' => $this->runner,
            'isRunnerOnline' => $this->isRunnerOnline(),
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

    static public function get(string $id, ?string $updatedAfter = null): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        $query = "
            SELECT `name`, `sessions`.`code`, `sessions`.`lang`, `sessions`.`runner`, `runner_checked_at`,
                `sessions`.`updated_at`, `writer`, `requests`.`session` IS NOT NULL AS `isWaitingForResult`, `results`.`result`
            FROM `sessions`
            LEFT JOIN `requests` ON `requests`.`session` = `sessions`.`id`
            LEFT JOIN `results` ON `results`.`session` = `sessions`.`id`
            WHERE `id` = ?
        ";
        $params = [$id];
        if ($updatedAfter !== null) {
            $query .= "AND `updated_at` > ?";
            $params[] = $updatedAfter;
        }
        $res = Db::get()->select($query, $params);
        if (count($res) === 0) {
            return null;
        }
        [$sessionName, $code, $lang, $runner, $runnerCheckedAtStr, $updatedAtStr, $writer, $isWaitingForResult, $result] = $res[0];
        $runnerCheckedAt = null;
        if ($runnerCheckedAtStr !== null) {
            $runnerCheckedAt = DateTime::createFromFormat('Y-m-d H:i:s', $runnerCheckedAtStr);
        }
        $updatedAt = DateTime::createFromFormat('Y-m-d H:i:s.u', $updatedAtStr);

        if ($result === '') {
            $result = '_';
        }
        $session = new self($id, $sessionName, $code, $lang, $runner, $runnerCheckedAt, $updatedAt, $writer, [], $isWaitingForResult, $result ?? '');
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
        $query = "INSERT INTO `sessions` SET `name` = ?, `code` = ?, `lang` = ?, `runner` = ?, `writer` = ?, `id` = ?;";
        $this->db->exec($query, [$this->name, $this->code, $this->lang, $this->runner, $this->writer, $this->id]);
        return self::get($this->id);
    }

    public function setSessionName(string $name): bool
    {
        if (!Utils::isValidString($name)) {
            return false;
        }
        $query = "UPDATE `sessions` SET `name` = ?, `updated_at` = NOW(3) WHERE `id` = ?";
        $this->db->exec($query, [$name, $this->id]);
        return true;
    }

    public function setLang(string $lang): bool
    {
        if (!isset(self::LANGS[$lang])) {
            return false;
        }
        $query = "UPDATE `sessions` SET `lang` = ?, `updated_at` = NOW(3) WHERE `id` = ?";
        $this->db->exec($query, [$lang, $this->id]);
        return true;
    }

    public function setCode(string $code): bool
    {
        if (strlen($code) > self::CODE_MAX_LENGTH) {
            return false;
        }
        $query = "UPDATE `sessions` SET `code` = ?, `updated_at` = NOW(3) WHERE `id` = ?";
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

    public function setRunner(string $runner): bool
    {
        if (!Utils::isUuid($runner)) {
            return false;
        }
        $query = "UPDATE `sessions` SET `runner` = ?, `updated_at` = NOW(3) WHERE `id` = ?";
        $this->db->exec($query, [$runner, $this->id]);
        return true;
    }

    static public function setCheckedByRunner(string $runner): void
    {
        if (!Utils::isUuid($runner)) {
            return;
        }
        $query = "UPDATE `sessions` SET `runner_checked_at` = NOW() WHERE `runner` = ?";
        Db::get()->exec($query, [$runner]);
    }

    public function setWriter(string $userId): bool
    {
        if (!Utils::isUuid($userId)) {
            return false;
        }
        $query = "UPDATE `sessions` SET `writer` = ?, `updated_at` = NOW(3) WHERE `id` = ?";
        $this->db->exec($query, [$userId, $this->id]);
        return true;
    }

    public function isRunnerOnline(): bool
    {
        if ($this->runnerCheckedAt === null) {
            return false;
        }
        return time() - $this->runnerCheckedAt->getTimestamp() < 10;
    }
}