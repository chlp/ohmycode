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
    }
}