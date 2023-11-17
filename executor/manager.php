<?php

$conf = loadConf();
if ($conf === null) {
    return;
}

while (true) {
    foreach ($conf['languages'] as $lang) {
        $resultsDir = $lang . '/results';
        $files = preg_grep('/^([^.])/', scandir($resultsDir));
        foreach ($files as $file) {
            $result = file_get_contents($resultsDir . '/' . $file);
            post($conf['api'] . '/action/result.php', [
                'executor' => $conf['id'],
                'lang' => $lang,
                'hash' => $file,
            ]);
            unlink($resultsDir . '/' . $file);
        }
    }

    $requests = post($conf['api'] . '/action/request.php', [
        'executor' => $conf['id'],
    ]);
    foreach ($requests as $request) {
        file_put_contents(__DIR__. "/{$request['lang']}/{$request['hash']}", $request['code']);
    }
    sleep(1);
}

function loadConf(): ?array {
    define('CONF_PATH', 'conf.json');
    define('CONF_EXAMPLE_PATH', 'conf-example.json');
    if (!file_exists(CONF_PATH)) {
        $conf = json_decode(file_get_contents(CONF_EXAMPLE_PATH), true);
        if (!is_array($conf)) {
            echo 'conf: wrong conf-example';
            return null;
        }
        $conf['id'] = genUuid();
        file_put_contents(CONF_PATH, json_encode($conf, JSON_PRETTY_PRINT));
        return $conf;
    }
    $conf = json_decode(file_get_contents(CONF_PATH), true);
    if (!is_array($conf)) {
        echo 'conf: can not parse file: ' . CONF_PATH;
        return null;
    }
    if (!isset($conf['id']) || strlen($conf['id']) !== 32) {
        $conf['id'] = genUuid();
        file_put_contents(CONF_PATH, json_encode($conf, JSON_PRETTY_PRINT));
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
    return $conf;
}

function genUuid(): string {
    return sprintf( '%04x%04x%04x%04x%04x%04x%04x%04x',
        mt_rand( 0, 0xffff ), mt_rand( 0, 0xffff ),
        mt_rand( 0, 0xffff ),
        mt_rand( 0, 0x0fff ) | 0x4000,
        mt_rand( 0, 0x3fff ) | 0x8000,
        mt_rand( 0, 0xffff ), mt_rand( 0, 0xffff ), mt_rand( 0, 0xffff )
    );
}


function post($url, $data): array {
    $curl = curl_init($url);
    curl_setopt($curl, CURLOPT_URL, $url);
    curl_setopt($curl, CURLOPT_POST, true);
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($curl, CURLOPT_POSTFIELDS, json_encode($data));
    $resp = json_decode((string)curl_exec($curl), true);
    curl_close($curl);
    if (!is_array($resp)) {
        return [];
    }
    return $resp;
}