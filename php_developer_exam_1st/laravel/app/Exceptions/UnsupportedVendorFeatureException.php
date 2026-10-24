<?php

namespace App\Exceptions;

use Exception;

class UnsupportedVendorFeatureException extends Exception
{
    public function __construct(string $vendorName, string $feature)
    {
        parent::__construct(
            "Feature '{$feature}' is not supported by vendor '{$vendorName}'"
        );
    }
}
