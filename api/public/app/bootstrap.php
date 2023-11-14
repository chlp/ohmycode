<?php

static $loaded = false;
if ($loaded) {
    return true;
}

require __DIR__ . '/Utils.php';
require __DIR__ . '/Conf.php';
require __DIR__ . '/Db.php';
require __DIR__ . '/Session.php';