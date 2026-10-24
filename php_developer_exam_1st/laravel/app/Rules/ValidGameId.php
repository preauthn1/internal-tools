<?php

namespace App\Rules;

use Closure;
use Illuminate\Contracts\Validation\ValidationRule;

class ValidGameId implements ValidationRule
{
    /**
     * Run the validation rule.
     *
     * @param  \Closure(string, ?string=): \Illuminate\Translation\PotentiallyTranslatedString  $fail
     */
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

        if (!preg_match('/^[a-zA-Z0-9]+$/', $value)) {
            $fail('The :attribute may only contain letters and numbers.');
            return;
        }
    }
}
