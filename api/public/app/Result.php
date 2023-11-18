<?php

class Result
{
    static public function set(Request $request, string $result): void
    {
        $query = "INSERT INTO `results` SET `session` = ?, `code` = ?, `result` = ?, `lang` = ?
                       ON DUPLICATE KEY UPDATE `code` = ?, `result` = ?, `lang` = ?";
        Db::get()->exec($query, [
            $request->session, $request->code, $result, $request->lang,
            $request->code, $result, $request->lang,
        ]);
        Session::get($request->session)?->updateTime();
        Request::remove($request->executor, $request->lang, $request->hash);
    }

    static public function remove(string $session): void
    {
        if (!Utils::isUuid($session)) {
            return;
        }
        $query = "DELETE FROM `results` WHERE `session` = ?";
        Db::get()->exec($query, [$session]);
    }
}