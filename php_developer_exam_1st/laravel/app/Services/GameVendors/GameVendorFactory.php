<?php

namespace App\Services\GameVendors;

use App\Services\GameVendors\GameVendorServiceInterface;
use InvalidArgumentException;

class GameVendorFactory
{
    public static function make(string $vendorName): GameVendorServiceInterface
    {
        $config = config("game_vendors.{$vendorName}");

        if (!$config) {
            throw new InvalidArgumentException(
                "Vendor '{$vendorName}' is not configured in game_vendors.php"
            );
        }

        $serviceClass = self::getServiceClassName($vendorName);

        if (!class_exists($serviceClass)) {
            throw new InvalidArgumentException(
                "Service class '{$serviceClass}' does not exist for vendor '{$vendorName}'"
            );
        }

        return new $serviceClass($vendorName, $config);
    }

    public static function getAvailableVendors(): array
    {
        return array_keys(config('game_vendors', []));
    }

    public static function isVendorAvailable(string $vendorName): bool
    {
        return config("game_vendors.{$vendorName}") !== null;
    }

    protected static function getServiceClassName(string $vendorName): string
    {
        $className = str_replace(
            ' ',
            '',
            ucwords(str_replace(['-', '_' ], ' ', $vendorName))
        ) . 'Service';

        $baseClass = "App\\Services\\GameVendors\\{$className}";
        $vendorDir = str_replace(' ', '', ucwords(str_replace(['-', '_' ], ' ', $vendorName)));
        $vendorClass = "App\\Services\\GameVendors\\{$vendorDir}\\{$className}";

        if (class_exists($vendorClass)) {
            return $vendorClass;
        }

        return $baseClass;
    }
}
