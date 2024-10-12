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
        public bool $runnerIsPublic,
        public string $runner,
        public ?DateTime $runnerCheckedAt,
        public ?DateTime $updatedAt,
        public ?DateTime $codeUpdatedAt,
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
            'codeHash' => Utils::ohMySimpleHash($this->code),
            'lang' => $this->lang,
            'runnerIsPublic' => $this->runnerIsPublic,
            'runner' => $this->runner,
            'runnerCA' => $this->runnerCheckedAt,
            'runnerIsOnline' => $this->runnerIsOnline(),
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
        $runnerCheckedAt = self::getNewestPublicRunnerCheckedAt();
        return new self($id, $name, '', self::DEFAULT_LANG, true, '', $runnerCheckedAt, null, null, '', [], false, '');
    }

    public static function get(string $id, ?string $updatedAfter = null): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        $query = "
            SELECT `name`, `sessions`.`code`, `sessions`.`lang`,
                `sessions`.`runner_is_public`, `sessions`.`runner`, `runner_checked_at`,
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
        [$sessionName, $code, $lang, $runnerIsPublic, $runner, $runnerCheckedAtStr, $updatedAtStr, $codeUpdatedAtStr, $writer, $isWaitingForResult, $result] = $res[0];
        $runnerCheckedAt = null;
        if ($runnerCheckedAtStr !== null) {
            $runnerCheckedAt = DateTime::createFromFormat('Y-m-d H:i:s', $runnerCheckedAtStr);
        }
        $updatedAt = DateTime::createFromFormat('Y-m-d H:i:s.u', $updatedAtStr);
        $codeUpdatedAt = DateTime::createFromFormat('Y-m-d H:i:s.u', $codeUpdatedAtStr);

        if ($result === '') {
            $result = '_';
        }
        $session = new self($id, $sessionName, $code, $lang, $runnerIsPublic, $runner, $runnerCheckedAt, $updatedAt, $codeUpdatedAt, $writer, [], $isWaitingForResult, $result ?? '');
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

    public static function cleanupUsers(string $sessionId): void
    {
        if (Utils::isUuid($sessionId)) {
            $query = "delete from `session_users` where `session` = ? and `updated_at` < NOW(3) - INTERVAL " . self::IS_ACTIVE_FROM_LAST_UPDATE_SEC . " second";
            if (Db::get()->exec($query, [$sessionId]) > 0) {
                error_log('removeOldUsers');
                self::updateTime($sessionId);
            }
        }
    }

    public function insert(): self
    {
        $query = "INSERT INTO `sessions` SET `name` = ?, `code` = ?, `lang` = ?, `runner_is_public` = ?, `runner` = ?, `runner_checked_at` = ?, `writer` = ?, `id` = ?;";
        $this->db->exec($query, [$this->name, $this->code, $this->lang, $this->runnerIsPublic, $this->runner, $this->runnerCheckedAt, $this->writer, $this->id]);
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

    public function setCode(string $code, string $userId): bool
    {
        if (strlen($code) > self::CODE_MAX_LENGTH) {
            return false;
        }
        $query = "
            UPDATE `sessions` SET
                `code` = ?,
                `updated_at` = NOW(3),
                `code_updated_at` = NOW(3),
                `writer` = ?
            WHERE `id` = ? AND (`writer` = ? OR `writer` = '');
        ";
        return 1 === $this->db->exec($query, [$code, $userId, $this->id, $userId]);
    }

    public static function updateTime(string $sessionId): void
    {
        if (!Utils::isUuid($sessionId)) {
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

    public static function setCheckedByRunner(string $runner, bool $isPublic): void
    {
        if (!Utils::isUuid($runner)) {
            return;
        }

        if ($isPublic) {
            $query = "UPDATE `sessions` SET `runner_checked_at` = NOW() WHERE `runner` = ''";
            Db::get()->exec($query, []);
        } else {
            $query = "UPDATE `sessions` SET `runner_checked_at` = NOW() WHERE `runner` = ?";
            Db::get()->exec($query, [$runner]);
        }

        self::setActiveRunner($runner, $isPublic);
    }

    public static function cleanupWriter(string $sessionId): void
    {
        if (Utils::isUuid($sessionId)) {
            $query = "update `sessions` set `writer` = '', `updated_at` = NOW(3) where `id` = ? and `writer` != '' and `code_updated_at` < NOW(3) - INTERVAL " . self::IS_WRITER_STILL_WRITING_SEC . " second";
            Db::get()->exec($query, [$sessionId]);
        }
    }

    public function runnerIsOnline(): bool
    {
        if ($this->runnerCheckedAt === null) {
            return false;
        }
        return time() - $this->runnerCheckedAt->getTimestamp() < self::IS_ACTIVE_FROM_LAST_UPDATE_SEC;
    }

    private static function getNewestPublicRunnerCheckedAt(): ?DateTime
    {
        $checkedAt = null;
        $query = "SELECT `checked_at` FROM `runners` WHERE `is_public` = true ORDER BY `checked_at` DESC LIMIT 1;";
        $res = Db::get()->select($query);
        if (count($res) === 1) {
            $checkedAtStr = $res[0][0];
            if ($checkedAtStr !== null) {
                $checkedAt = DateTime::createFromFormat('Y-m-d H:i:s', $checkedAtStr);
            }
        }
        return $checkedAt;
    }

    private static function setActiveRunner(string $runner, bool $isPublic): void
    {
        $setActiveRunnersQuery = "INSERT INTO runners (id, checked_at) VALUES (?, NOW()) ON DUPLICATE KEY UPDATE checked_at = NOW()";
        Db::get()->exec($setActiveRunnersQuery, [$runner]);
    }
}
