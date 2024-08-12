<?php

namespace app;

use DateTime;

class Session
{
    private const DEFAULT_LANG = 'php82';
    private const CODE_MAX_LENGTH = 32768;
    private const IS_ACTIVE_FROM_LAST_UPDATE_SEC = 10;
    private const IS_WRITER_STILL_WRITING_SEC = 2;
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
        public ?DateTime $codeUpdatedAt,
        public string $writer,
        public array $users,
        public bool $isWaitingForResult,
        public string $result,
    ) {
        $this->db = Db::get();
    }

    public function getJson(): string
    {
        return json_encode([
            'id' => $this->id,
            'name' => $this->name,
            'code' => $this->code,
            'codeHash' => Utils::ohMySimpleHash($this->code),
            'lang' => $this->lang,
            'runner' => $this->runner,
            'isRunnerOnline' => $this->isRunnerOnline(),
            'updatedAt' => $this->updatedAt,
            'codeUpdatedAt' => $this->codeUpdatedAt,
            'writer' => $this->writer,
            'users' => $this->users,
            'isWaitingForResult' => $this->isWaitingForResult,
            'result' => $this->result,
        ]);
    }

    public static function createNew(string $id): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        $name = date('Y-m-d');
        $runner = self::getRandomActiveRunner();
        $runnerCheckedAt = null;
        if ($runner !== '') {
            $runnerCheckedAt = new DateTime();
        }
        return new self($id, $name, '', self::DEFAULT_LANG, $runner, $runnerCheckedAt, null, null, '', [], false, '');
    }

    public static function get(string $id, ?string $updatedAfter = null): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        $query = "
            SELECT `name`, `sessions`.`code`, `sessions`.`lang`, `sessions`.`runner`, `runner_checked_at`,
                `sessions`.`updated_at`, `sessions`.`code_updated_at`, `writer`,
                `requests`.`session` IS NOT NULL AS `isWaitingForResult`, `results`.`result`
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
        [$sessionName, $code, $lang, $runner, $runnerCheckedAtStr, $updatedAtStr, $codeUpdatedAtStr, $writer, $isWaitingForResult, $result] = $res[0];
        $runnerCheckedAt = null;
        if ($runnerCheckedAtStr !== null) {
            $runnerCheckedAt = DateTime::createFromFormat('Y-m-d H:i:s', $runnerCheckedAtStr);
        }
        $updatedAt = DateTime::createFromFormat('Y-m-d H:i:s.u', $updatedAtStr);
        $codeUpdatedAt = DateTime::createFromFormat('Y-m-d H:i:s.u', $codeUpdatedAtStr);

        if ($result === '') {
            $result = '_';
        }
        $session = new self($id, $sessionName, $code, $lang, $runner, $runnerCheckedAt, $updatedAt, $codeUpdatedAt, $writer, [], $isWaitingForResult, $result ?? '');
        $session->loadUsers();

        return $session;
    }

    private function loadUsers(): void
    {
        $res = $this->db->select("select `user`, `name` from `session_users` where session = ? ORDER BY `name`", [$this->id]);
        $users = [];
        foreach ($res as $row) {
            $users[$row[0]] = [
                'id' => $row[0],
                'name' => $row[1],
            ];
        }
        $this->users = $users;
    }

    public static function updateUserOnline(string $sessionId, string $userId): void
    {
        if (Utils::isUuid($sessionId) && Utils::isUuid($userId)) {
            Db::get()->exec("update `session_users` set updated_at = NOW(3) where session = ? and user = ?", [$sessionId, $userId]);
        }
    }

    public static function removeOldUsers(string $sessionId): void
    {
        if (Utils::isUuid($sessionId)) {
            $query = "delete from `session_users` where `session` = ? and `updated_at` < NOW(3) - INTERVAL " . self::IS_ACTIVE_FROM_LAST_UPDATE_SEC . " second";
            if (Db::get()->exec($query, [$sessionId]) > 0) {
                self::updateTime($sessionId);
            }
        }
    }

    public static function updateWriter(string $sessionId): void
    {
        if (Utils::isUuid($sessionId)) {
            $query = "update `sessions` set writer = '', updated_at = NOW(3) where `id` = ? and `code_updated_at` < NOW(3) - INTERVAL " . self::IS_WRITER_STILL_WRITING_SEC . " second";
            Db::get()->exec($query, [$sessionId]);
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
        $query = "UPDATE `sessions` SET `code` = ?, `updated_at` = NOW(3), `code_updated_at` = NOW(3) WHERE `id` = ?";
        $this->db->exec($query, [$code, $this->id]);
        return true;
    }

    public static function updateTime(string $sessionId): void
    {
        if (Utils::isUuid($sessionId)) {
            return;
        }
        $query = "UPDATE `sessions` SET `updated_at` = NOW(3) WHERE `id` = ?;";
        Db::get()->exec($query, [$sessionId]);
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
        self::updateTime($this->id);
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

    public static function setCheckedByRunner(string $runner): void
    {
        if (!Utils::isUuid($runner)) {
            return;
        }

        $query = "UPDATE `sessions` SET `runner_checked_at` = NOW() WHERE `runner` = ?";
        Db::get()->exec($query, [$runner]);

        self::setActiveRunner($runner);
    }

    public function setWriter(string $userId): bool
    {
        if (!Utils::isUuid($userId)) {
            return false;
        }
        $this->writer = $userId;
        $query = "UPDATE `sessions` SET `writer` = ?, `updated_at` = NOW(3), `code_updated_at` = NOW(3) WHERE `id` = ?";
        $this->db->exec($query, [$userId, $this->id]);
        return true;
    }

    public function isRunnerOnline(): bool
    {
        if ($this->runnerCheckedAt === null) {
            return false;
        }
        return time() - $this->runnerCheckedAt->getTimestamp() < self::IS_ACTIVE_FROM_LAST_UPDATE_SEC;
    }

    private static function getRandomActiveRunner(): string
    {
        $query = "SELECT `id` FROM `runners` WHERE checked_at >= NOW() - INTERVAL " . self::IS_ACTIVE_FROM_LAST_UPDATE_SEC . " SECOND ORDER BY RAND() LIMIT 1;";
        $runners = Db::get()->select($query);
        if (count($runners) === 1) {
            return $runners[0][0];
        }
        return '';
    }

    private static function setActiveRunner(string $runner): void
    {
        $setActiveRunnersQuery = "INSERT INTO runners (id, checked_at) VALUES (?, NOW()) ON DUPLICATE KEY UPDATE checked_at = NOW()";
        Db::get()->exec($setActiveRunnersQuery, [$runner]);
    }
}
