<?php

namespace App\Enums;

enum GameTransferAuditStatus: string
{
    case PendingRecords = 'pending_records';
    case Matched = 'matched';
    case ManualReview = 'manual_review';
}
