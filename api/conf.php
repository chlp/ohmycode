<?php

function loadApiConf(): array
{
    define('CONF_PATH', 'api-conf.json');
    if (!file_exists(CONF_PATH)) {
        die('conf: please create conf file');
    }
    $conf = json_decode(file_get_contents(CONF_PATH), true);
    if (!is_array($conf)) {
        die('conf: can not parse file');
    }
    return $conf;
}