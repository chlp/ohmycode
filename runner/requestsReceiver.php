<?php

namespace runner;

require __DIR__ . '/tools.php';

echo "Runner. Requests receiver starting\n";

$conf = Conf::load();
if ($conf === null) {
    echo "problem with conf\n";
    return;
}
$api = new Api($conf->runnerId, $conf->isPublic, $conf->apiUrl);

usleep(500000); // 0.5 sec
echo "requests receiver initiating. id: $conf->runnerId\n";

while (true) {
    $requests = $api->request('/run/get_tasks', ['is_keep_alive' => true], true);

    if (!$requests->isOk()) {
        echo json_encode([date('Y-m-d H:i:s'), 'get requests', $requests->code, $requests->data]);
        sleep(2);
        continue;
    }

    foreach ($requests->data as $request) {
        $lang = $request['lang'];
        $hash = $request['hash'];
        if (in_array($request['lang'], $conf->languages)) {
            $res = $api->request('/run/ack_task', [
                'lang' => $lang,
                'hash' => (int)$hash,
            ]);
            if ($res->code !== 404) {
                $filePath = __DIR__ . "/$lang/requests/$hash";
                file_put_contents($filePath, $request['content']);
                chmod($filePath, 0700);
            }
        } else {
            $api->result('/result/set', [
                'lang' => $lang,
                'hash' => (int)$hash,
                'result' => "No runner for $lang",
            ]);
        }
    }

    usleep(200000); // 0.2 sec
}
