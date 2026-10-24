<?php

namespace App\Rules;

use App\Services\GameVendors\GameVendorFactory;
use Closure;
use Illuminate\Contracts\Validation\ValidationRule;

class ValidVendorId implements ValidationRule
{
    public function validate(string $attribute, mixed $value, Closure $fail): void
    {
        if (!is_string($value)) {
            $fail('The :attribute must be a string.');
            return;
        }

        if (strlen($value) > 20) {
            $fail('The :attribute must not exceed 20 characters.');
            return;
        }

        if (!preg_match('/^[a-zA-Z0-9_-]+$/', $value)) {
            $fail('The :attribute may only contain letters, numbers, dashes, and underscores.');
            return;
        }

        if (!in_array($value, GameVendorFactory::getAvailableVendors(), true)) {
            $fail('The selected :attribute is invalid.');
            return;
        }
    }
}

