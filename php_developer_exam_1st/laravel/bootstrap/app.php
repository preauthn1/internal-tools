<?php

use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__ . '/../routes/web.php',
        commands: __DIR__ . '/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        $middleware->alias([
            'validate.vendor' => \App\Http\Middleware\GlobalParamValidator::class,
        ]);
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        // Global exception logging - all exceptions are automatically logged by Laravel
        // This callback allows you to customize the logging behavior
        $exceptions->report(function (\Throwable $e) {
            // Laravel logs exceptions automatically, but you can add custom logic here
            // For example: send to external service, add additional context, etc.

            \Illuminate\Support\Facades\Log::error('Exception occurred', [
                'message' => $e->getMessage(),
                'file' => $e->getFile(),
                'line' => $e->getLine(),
                'exception_class' => get_class($e),
                'trace' => $e->getTraceAsString(),
            ]);
        });
    })->create();
