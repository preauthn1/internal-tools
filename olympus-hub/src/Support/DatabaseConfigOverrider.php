<?php

namespace OlympusHub\Support;

use Illuminate\Support\Facades\DB;

class DatabaseConfigOverrider
{
    public function apply(array $vars): void
    {
        config([
            'database.connections.mysql.read.host'      => [$vars['DB_HOST_READ']],
            'database.connections.mysql.read.username'  => $vars['DB_USERNAME_READ'],
            'database.connections.mysql.read.password'  => $vars['DB_PASSWORD_READ'],
            'database.connections.mysql.write.host'     => [$vars['DB_HOST_WRITE']],
            'database.connections.mysql.write.username' => $vars['DB_USERNAME_WRITE'],
            'database.connections.mysql.write.password' => $vars['DB_PASSWORD_WRITE'],
            'database.connections.mysql.database'       => $vars['DB_DATABASE'],
            'database.connections.mysql.port'           => $vars['DB_PORT'],
        ]);

        DB::purge('mysql');
    }
}
