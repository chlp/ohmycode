<?php

namespace app;

class Result
{
    public static function set(Request $request, string $result): void
    {
        $query = "INSERT INTO `results` SET `session` = ?, `code` = ?, `result` = ?, `lang` = ?
                       ON DUPLICATE KEY UPDATE `code` = ?, `result` = ?, `lang` = ?";
        Db::get()->exec($query, [
            $request->session, $request->code, $result, $request->lang,
            $request->code, $result, $request->lang,
        ]);
        Session::updateTime($request->session);
        Request::removeByRunner($request->runner, $request->lang, $request->hash);
    }

    public static function removeBySession(string $session): void
    {
        if (!Utils::isUuid($session)) {
            return;
        }
        $query = "DELETE FROM `results` WHERE `session` = ?";
        Db::get()->exec($query, [$session]);
        Session::updateTime($session);
    }
}
