<?php

use app\Utils;

static $loaded = false;
if ($loaded) {
    return true;
}

require __DIR__ . '/Utils.php';
Utils::timer();
require __DIR__ . '/Conf.php';
require __DIR__ . '/Db.php';
require __DIR__ . '/Session.php';
require __DIR__ . '/Request.php';
require __DIR__ . '/Result.php';
