<?php

namespace App\Enums;

enum GameTransferStatus: string
{
    case Pending = 'pending';
    case Unknown = 'unknown';
    case Completed = 'completed';
    case Failed = 'failed';
    case ManualReview = 'manual_review';

    public function isUnresolved(): bool
    {
        return in_array($this, [self::Pending, self::Unknown, self::ManualReview], true);
    }
}
