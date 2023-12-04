<?php

echo "executor manager starting\n";

$conf = loadConf();
if ($conf === null) {
    echo "problem with conf\n";
    return;
}

$runnerId = $conf['id'];
$runnerLanguages = $conf['id'];
$runnerApiUrl = $conf['api'];

echo "id: $runnerId\n";

while (true) {
    echo "next cycle\n";
    foreach ($runnerLanguages as $lang) {
        $resultsDir = __DIR__ . "/$lang/results";
        $files = preg_grep('/^([^.])/', scandir($resultsDir));
        foreach ($files as $file) {
            $result = substr(file_get_contents($resultsDir . '/' . $file), 0, 16384);
            [$code, $response] = post($runnerApiUrl . '/action/result.php', [
                'action' => 'set',
                'executor' => $runnerId,
                'lang' => $lang,
                'hash' => $file,
                'result' => $result,
            ]);
            if ($code !== 200) {
                var_dump('set result', $code, $response);
            } else {
                unlink($resultsDir . '/' . $file);
            }
        }
    }

    echo "2\n";

    [$code, $requests] = post($runnerApiUrl . '/action/request.php', [
        'action' => 'get',
        'executor' => $runnerId,
    ]);
    echo "3\n";
    if ($code !== 200) {
        var_dump('get requests', $code, $requests);
    } else {
        foreach ($requests as $request) {
            $lang = $request['lang'];
            $hash = $request['hash'];
            if (in_array($request['lang'], $runnerLanguages)) {
                file_put_contents(__DIR__ . "/$lang/requests/$hash", $request['code']);
                post($runnerApiUrl . '/action/request.php', [
                    'action' => 'markReceived',
                    'executor' => $runnerId,
                    'lang' => $lang,
                    'hash' => $hash,
                ]);
            } else {
                post($runnerApiUrl . '/action/result.php', [
                    'action' => 'set',
                    'executor' => $runnerId,
                    'lang' => $lang,
                    'hash' => $hash,
                    'result' => "There is no necessary runner for $lang",
                ]);
            }
        }
    }
    echo "4\n";
    sleep(1);
}

function loadConf(): ?array
{
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

function genUuid(): string
{
    return sprintf('%04x%04x%04x%04x%04x%04x%04x%04x',
        mt_rand(0, 0xffff), mt_rand(0, 0xffff),
        mt_rand(0, 0xffff),
        mt_rand(0, 0x0fff) | 0x4000,
        mt_rand(0, 0x3fff) | 0x8000,
        mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff)
    );
}


function post($url, $data): array
{
    $curl = curl_init($url);
    curl_setopt($curl, CURLOPT_URL, $url);
    curl_setopt($curl, CURLOPT_POST, true);
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($curl, CURLOPT_POSTFIELDS, json_encode($data));
    $resp = curl_exec($curl);
    $json = json_decode((string)$resp, true);
    $code = curl_getinfo($curl, CURLINFO_HTTP_CODE);
    curl_close($curl);
    if ($code === 200 && $resp === '') {
        return [200, []];
    }
    if (!is_array($json)) {
        return ['e' . $code, $resp];
    }
    if ($code !== 200) {
        return [$code, $resp];
    }
    return [200, $json];
}