<?php

namespace runner;

class Conf
{
    private const filePath = 'conf.json';
    private const exampleFilePath = 'conf-example.json';

    public function __construct(
        public readonly string $runnerId,
        public readonly bool $isPublic,
        public readonly string $runnerName,
        public readonly string $apiUrl,
        public readonly array $languages
    )
    {
    }

    public static function load(): ?self
    {
        if (!file_exists(self::filePath)) {
            $conf = json_decode(file_get_contents(self::exampleFilePath), true);
            if (!is_array($conf)) {
                echo 'conf: wrong conf-example';
                return null;
            }
            $conf['id'] = self::genUuid();
            file_put_contents(self::filePath, json_encode($conf, JSON_PRETTY_PRINT));
            chmod(self::filePath, 0700);
        }
        $conf = json_decode(file_get_contents(self::filePath), true);
        if (!is_array($conf)) {
            echo 'conf: can not parse file: ' . self::filePath;
            return null;
        }
        if (!isset($conf['id']) || strlen($conf['id']) !== 32) {
            $conf['id'] = self::genUuid();
            file_put_contents(self::filePath, json_encode($conf, JSON_PRETTY_PRINT));
            chmod(self::filePath, 0700);
            return null;
        }
        if (!isset($conf['name']) || !isset($conf['languages'])) {
            echo 'conf: incomplete file';
            return null;
        }
        if (!is_string($conf['id'])) {
            echo 'conf: wrong id format';
            return null;
        }
        if (!is_bool($conf['is_public'])) {
            echo 'conf: wrong is_public format';
            return null;
        }
        if (!is_string($conf['api'])) {
            echo 'conf: wrong api format';
            return null;
        }
        if (!is_string($conf['name'])) {
            echo 'conf: wrong name format';
            return null;
        }
        if (!is_array($conf['languages'])) {
            echo 'conf: wrong id format';
            return null;
        }
        return new self($conf['id'], $conf['is_public'], $conf['name'], $conf['api'], $conf['languages']);
    }

    private static function genUuid(): string
    {
        return sprintf('%04x%04x%04x%04x%04x%04x%04x%04x',
            mt_rand(0, 0xffff), mt_rand(0, 0xffff),
            mt_rand(0, 0xffff),
            mt_rand(0, 0x0fff) | 0x4000,
            mt_rand(0, 0x3fff) | 0x8000,
            mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff)
        );
    }
}

readonly class ApiResponse
{
    public function __construct(
        public int $code,
        public mixed $data,
    )
    {
    }

    public function isOk(): bool
    {
        return $this->code === 200;
    }
}

readonly class Api
{
    private const DEFAULT_TIMEOUT_SEC = 2;
    private const KEEPALIVE_TIMEOUT_SEC = 15;

    public function __construct(
        private string $runnerId,
        private bool $isPublic,
        private string $url,
    )
    {
    }

    public function request(string $action, ?array $moreData = null, bool $isKeepAlive = false): ApiResponse
    {
        $data = [
            'action' => $action,
            'isPublic' => $this->isPublic,
            'runner' => $this->runnerId,
        ];
        if ($moreData !== null) {
            $data = array_merge($data, $moreData);
        }
        [$code, $requests] = self::post($this->url . '/action/request.php', $data, $isKeepAlive);
        return new ApiResponse($code, $requests);
    }

    public function result(string $action, ?array $moreData = null): ApiResponse
    {
        $data = [
            'action' => $action,
            'isPublic' => $this->isPublic,
            'runner' => $this->runnerId,
        ];
        if ($moreData !== null) {
            $data = array_merge($data, $moreData);
        }
        [$code, $requests] = self::post($this->url . '/action/result.php', $data);
        return new ApiResponse($code, $requests);
    }

    private static function post(string $url, array $data, bool $isKeepAlive = false): array
    {
        $curl = curl_init($url);
        curl_setopt_array($curl, [
            CURLOPT_URL => $url,
            CURLOPT_POST => true,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_POSTFIELDS => json_encode($data),
            CURLOPT_ENCODING => '',
            CURLOPT_MAXREDIRS => 10,
            CURLOPT_TIMEOUT => $isKeepAlive ? self::KEEPALIVE_TIMEOUT_SEC : self::DEFAULT_TIMEOUT_SEC,
            CURLOPT_CONNECTTIMEOUT => self::DEFAULT_TIMEOUT_SEC,
            CURLOPT_FOLLOWLOCATION => true,
            CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
            CURLOPT_HTTPHEADER => ['Content-Type: application/json'],
            CURLOPT_SSL_VERIFYHOST => 0,
            CURLOPT_SSL_VERIFYPEER => 0,
            CURLOPT_IPRESOLVE => CURL_IPRESOLVE_V4,
        ]);
        $resp = curl_exec($curl);
        $json = json_decode((string)$resp, true);
        $code = curl_getinfo($curl, CURLINFO_HTTP_CODE);
        $err = curl_error($curl);
        curl_close($curl);
        if ($code === 200 && $resp === '') {
            return [200, []];
        }
        if (!is_array($json)) {
            if ((string)$resp === '') {
                $resp = $err;
            }
            return [1000 + (int)$code, $resp];
        }
        if ($code !== 200) {
            return [(int)$code, $resp];
        }
        return [200, $json];
    }
}