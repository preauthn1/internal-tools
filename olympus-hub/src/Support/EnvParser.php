<?php

namespace OlympusHub\Support;

use OlympusHub\Exceptions\InvalidTetherDataException;

class EnvParser
{
    public function parse(string $raw): array
    {
        if (trim($raw) === '') {
            throw new InvalidTetherDataException('Tether value is empty');
        }

        $result = [];
        foreach (explode("\n", $raw) as $line) {
            $line = trim($line);
            if ($line === '' || strpos($line, '=') === false) {
                continue;
            }
            [$key, $value] = explode('=', $line, 2);
            $result[trim($key)] = trim($value);
        }

        return $result;
    }
}
