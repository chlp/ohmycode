<?php

namespace runner;

require __DIR__ . '/tools.php';

echo "Runner. Results sender starting\n";

$conf = Conf::load();
if ($conf === null) {
    echo "problem with conf\n";
    return;
}
$api = new Api($conf->runnerId, $conf->apiUrl);

echo "results sender. id: $conf->runnerId\n";

while (true) {
    $isEmpty = true;
    foreach ($conf->languages as $lang) {
        $resultsDir = __DIR__ . "/$lang/results";
        $files = preg_grep('/^([^.])/', scandir($resultsDir));
        foreach ($files as $file) {
            $isEmpty = false;
            $newResultData = substr(file_get_contents($resultsDir . '/' . $file), 0, 16384);
            $setter = $api->result('set', [
                'lang' => $lang,
                'hash' => $file,
                'result' => $newResultData,
            ]);
            if (!$setter->isOk()) {
                echo json_encode([date('Y-m-d H:i:s'), 'set result', $setter->code, $setter->data]);
                sleep(2);
                continue;
            }
            unlink($resultsDir . '/' . $file);
        }
    }
    if ($isEmpty) {
        usleep(200000); // 0.2 sec
    }
}
