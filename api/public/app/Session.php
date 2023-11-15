<?php

class Session
{
    private const DEFAULT_LANG = 'php82';

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
        public string        $id,
        public string        $name,
        public string        $code,
        public string        $lang,
        public string        $executor,
        public DateTime|null $executorCheckedAt,
        public DateTime      $updatedAt,
        public string        $writer,
        public array         $users,
        public bool          $request,
        public ?string       $result,
    )
    {
        $this->db = Db::get();
    }

    static public function createNew(string $id): ?self
    {
        if (!Utils::isUuid($id)) {
            return null;
        }
        $name = date('Y-m-d');
        return new self($id, $name, '', self::DEFAULT_LANG, '', null, new DateTime(), '', [], false, null);
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
        return new self($id, $sessionName, $code, $lang, $executor, $executorCheckedAt, $updatedAt, $writer, [], false, null);
    }

    public function save(bool $new = false): self
    {
        if ($new) {
            $query = "INSERT INTO `sessions` SET `name` = ?, `code` = ?, `lang` = ?, `executor` = ?, `executor_checked_at` = ?, `writer` = ?, `id` = ?;";
        } else {
            $query = "UPDATE `sessions` SET `name` = ?, `code` = ?, `lang` = ?, `executor` = ?, `executor_checked_at` = ?, `writer` = ? WHERE `id` = ?;";
        }
        $this->db->exec($query, [$this->name, $this->code, $this->lang, $this->executor, $this->executorCheckedAt, $this->writer, $this->id]);
        return self::getById($this->id);
    }

    public function getJson(): string
    {
        return json_encode([
            'id' => $this->id,
            'name' => $this->name,
            'code' => $this->code,
            'lang' => $this->lang,
            'executor' => $this->executor,
            'executorCheckedAt' => $this->executorCheckedAt,
            'updatedAt' => $this->updatedAt,
            'writer' => $this->writer,
            'users' => $this->users,
            'request' => $this->request,
            'result' => $this->result,
        ]);
    }
}