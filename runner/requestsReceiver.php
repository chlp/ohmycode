<?php

namespace runner;

require __DIR__ . '/tools.php';

echo "Runner. Requests receiver starting\n";

$conf = Conf::load();
if ($conf === null) {
    echo "problem with conf\n";
    return;
}
$api = new Api($conf->runnerId, $conf->apiUrl);

usleep(500000);
echo "requests receiver initiating. id: $conf->runnerId\n";

while (true) {
    $requests = $api->request('get');

    if (!$requests->isOk()) {
        echo json_encode([date('Y-m-d H:i:s'), 'get requests', $requests->code, $requests->data]);
        sleep(2);
        continue;
    }

    foreach ($requests->data as $request) {
        $lang = $request['lang'];
        $hash = $request['hash'];
        if (in_array($request['lang'], $conf->languages)) {
            $filePath = __DIR__ . "/$lang/requests/$hash";
            file_put_contents($filePath, $request['code']);
            chmod($filePath, 0777);
            $api->request('markReceived', [
                'lang' => $lang,
                'hash' => $hash,
            ]);
        } else {
            $api->result('set', [
                'lang' => $lang,
                'hash' => $hash,
                'result' => "No runner for $lang",
            ]);
        }
    }

    usleep(500000); // 0.5 sec
}
