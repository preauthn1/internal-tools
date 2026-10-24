<?php

namespace App\Http\Middleware;

use App\Rules\ValidGameId;
use App\Rules\ValidVendorId;
use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\ValidationException;
use Symfony\Component\HttpFoundation\Response;

class GlobalParamValidator
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        $vendorId = $request->route('vendorId');
        $gameId = $request->route('gameId');

        // Validate vendorId
        if (!$vendorId) {
            abort(404, 'Vendor ID is required');
        }

        try {
            $rules = [
                'vendorId' => ['required', new ValidVendorId()],
            ];

            // Validate gameId if present in route
            if ($gameId !== null) {
                $rules['gameId'] = ['required', new ValidGameId()];
            }

            Validator::make(
                [
                    'vendorId' => $vendorId,
                    'gameId' => $gameId,
                ],
                $rules
            )->validate();
        } catch (ValidationException $e) {
            abort(404, 'Invalid route parameter');
        }

        return $next($request);
    }
}
