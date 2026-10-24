<?php

namespace App\Http\Controllers;

use Illuminate\View\View;

class DashboardController extends Controller
{
    public function index(): View
    {
        $allVendors = config('game_vendors', []);

        $vendors = [];
        foreach ($allVendors as $vendorId => $vendorConfig) {
            $vendors[$vendorId] = [
                'icon' => $vendorConfig['icon'] ?? null,
            ];
        }

        return view('dashboard', [
            'vendors' => $vendors,
        ]);
    }
}
