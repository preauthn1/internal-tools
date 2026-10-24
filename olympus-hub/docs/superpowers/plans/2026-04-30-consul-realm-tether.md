# Consul Realm Tether Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Before executing artisan commands in `HubController::invoke()`, fetch DB connection info from Consul KV for the given `realmId` and override Laravel's runtime database config so the artisan command connects to the correct production realm.

**Architecture:** A `ConsulRealmTetherService` orchestrates three focused components: `ConsulKvClient` fetches raw KV data via Consul's HTTP API (`/v1/kv/DB/{REALM}?raw`), `EnvParser` converts the `.env`-format string to an associative array, and `DatabaseConfigOverrider` applies the result to Laravel's runtime `config('database.connections.mysql.*')` and purges the DB connection. `HubController::invoke()` calls the service before `Artisan::call()`, catching all exceptions at the controller level. Game realms skip tether activation.

**Tech Stack:** PHP 8.1+, Laravel 5.5–13 (illuminate/support), Consul HTTP API v1, friendsofphp/consul-php-sdk ^5.3, orchestra/testbench ^8|^9, PHPUnit ^6|^10, Mockery (via testbench)

---

## File Map

### New Files

| File | Responsibility |
|---|---|
| `src/Exceptions/ConsulConnectionException.php` | Thrown when Consul is unreachable or returns non-200/404 |
| `src/Exceptions/RealmTetherNotFoundException.php` | Thrown when Consul key does not exist (404) |
| `src/Exceptions/InvalidTetherDataException.php` | Thrown when KV value is empty or unparseable |
| `src/Consul/ConsulKvClient.php` | Wraps `friendsofphp/consul-php-sdk` KV service, maps SDK exceptions to domain exceptions, returns raw string |
| `src/Support/EnvParser.php` | Parses `KEY=VALUE\n...` format string into associative array |
| `src/Support/DatabaseConfigOverrider.php` | Applies parsed vars to Laravel config and calls `DB::purge('mysql')` |
| `src/Services/ConsulRealmTetherService.php` | Orchestrates client → parser → overrider for a given realmId |
| `tests/Consul/ConsulKvClientTest.php` | Unit tests for ConsulKvClient |
| `tests/Support/EnvParserTest.php` | Unit tests for EnvParser |
| `tests/Services/ConsulRealmTetherServiceTest.php` | Unit tests for ConsulRealmTetherService |
| `phpunit.xml` | PHPUnit configuration |

### Modified Files

| File | Change |
|---|---|
| `config/olympus.php` | Add `consul.base_url` key |
| `src/OlympusHubServiceProvider.php` | Add `registerServices()`, bind `ConsulRealmTetherService` |
| `src/Http/Controllers/HubController.php` | Inject service, call `activate()` before `Artisan::call()`, catch tether exceptions |

---

## Task 1: Test Infrastructure

**Files:**
- Create: `phpunit.xml`
- Create: `tests/TestCase.php`

- [X] **Step 1: Add sensiolabs/consul-php-sdk and run composer install**

```bash
cd /home/rhys/projects/www/olympus-hub && composer require sensiolabs/consul-php-sdk && composer install
```

Expected: `sensiolabs/consul-php-sdk` added to `composer.json` require block, packages installed under `vendor/`

- [X] **Step 2: Create `phpunit.xml`**

Note: `orchestra/testbench` is in `require-dev` but not used — all tests extend `PHPUnit\Framework\TestCase` directly.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<phpunit xmlns:xsi="http://example.invalid/2001/XMLSchema-instance"
         xsi:noNamespaceSchemaLocation="vendor/phpunit/phpunit/phpunit.xsd"
         bootstrap="vendor/autoload.php"
         colors="true">
    <testsuites>
        <testsuite name="Unit">
            <directory suffix="Test.php">./tests</directory>
        </testsuite>
    </testsuites>
</phpunit>
```

- [X] **Step 3: Verify test runner works**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit --version
```

Expected: `PHPUnit 10.x.x ...`

- [X] **Step 4: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add phpunit.xml && git commit -m "chore: add test infrastructure"
```

---

## Task 2: Exception Classes

**Files:**
- Create: `src/Exceptions/ConsulConnectionException.php`
- Create: `src/Exceptions/RealmTetherNotFoundException.php`
- Create: `src/Exceptions/InvalidTetherDataException.php`

- [X] **Step 1: Create `src/Exceptions/ConsulConnectionException.php`**

```php
<?php

namespace OlympusHub\Exceptions;

use RuntimeException;

class ConsulConnectionException extends RuntimeException {}
```

- [X] **Step 2: Create `src/Exceptions/RealmTetherNotFoundException.php`**

```php
<?php

namespace OlympusHub\Exceptions;

use RuntimeException;

class RealmTetherNotFoundException extends RuntimeException {}
```

- [X] **Step 3: Create `src/Exceptions/InvalidTetherDataException.php`**

```php
<?php

namespace OlympusHub\Exceptions;

use RuntimeException;

class InvalidTetherDataException extends RuntimeException {}
```

- [X] **Step 4: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add src/Exceptions/ && git commit -m "feat: add Consul tether exception classes"
```

---

## Task 3: EnvParser

**Files:**
- Create: `src/Support/EnvParser.php`
- Create: `tests/Support/EnvParserTest.php`

- [X] **Step 1: Create `tests/Support/EnvParserTest.php`**

```php
<?php

namespace OlympusHub\Tests\Support;

use OlympusHub\Exceptions\InvalidTetherDataException;
use OlympusHub\Support\EnvParser;
use PHPUnit\Framework\TestCase;

class EnvParserTest extends TestCase
{
    private EnvParser $parser;

    protected function setUp(): void
    {
        $this->parser = new EnvParser();
    }

    public function test_parse_returns_associative_array_from_env_format(): void
    {
        $raw = "DB_CONNECTION=mysql\nDB_HOST_READ=0.0.0.0\nDB_HOST_WRITE=0.0.0.0\nDB_PORT=3306\nDB_DATABASE=xsgcasino\nDB_USERNAME_READ=reader\nDB_PASSWORD_READ=rpass\nDB_USERNAME_WRITE=writer\nDB_PASSWORD_WRITE=wpass";

        $result = $this->parser->parse($raw);

        $this->assertSame('mysql',      $result['DB_CONNECTION']);
        $this->assertSame('0.0.0.0',   $result['DB_HOST_READ']);
        $this->assertSame('0.0.0.0',   $result['DB_HOST_WRITE']);
        $this->assertSame('3306',       $result['DB_PORT']);
        $this->assertSame('xsgcasino',  $result['DB_DATABASE']);
        $this->assertSame('reader',     $result['DB_USERNAME_READ']);
        $this->assertSame('rpass',      $result['DB_PASSWORD_READ']);
        $this->assertSame('writer',     $result['DB_USERNAME_WRITE']);
        $this->assertSame('wpass',      $result['DB_PASSWORD_WRITE']);
    }

    public function test_parse_throws_when_empty_string(): void
    {
        $this->expectException(InvalidTetherDataException::class);
        $this->expectExceptionMessage('Tether value is empty');

        $this->parser->parse('');
    }

    public function test_parse_throws_when_whitespace_only(): void
    {
        $this->expectException(InvalidTetherDataException::class);
        $this->expectExceptionMessage('Tether value is empty');

        $this->parser->parse("   \n  ");
    }

    public function test_parse_ignores_lines_without_equals(): void
    {
        $raw = "DB_HOST_READ=0.0.0.0\ninvalid_line\nDB_PORT=3306";

        $result = $this->parser->parse($raw);

        $this->assertSame('0.0.0.0', $result['DB_HOST_READ']);
        $this->assertSame('3306',     $result['DB_PORT']);
        $this->assertArrayNotHasKey('invalid_line', $result);
    }

    public function test_parse_handles_value_containing_special_chars(): void
    {
        $raw = "DB_PASSWORD_WRITE=!QAZ2wsxForXSGames_online";

        $result = $this->parser->parse($raw);

        $this->assertSame('!QAZ2wsxForXSGames_online', $result['DB_PASSWORD_WRITE']);
    }
}
```

- [X] **Step 2: Run test to verify it fails**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit tests/Support/EnvParserTest.php --colors=always
```

Expected: FAIL — `Class "OlympusHub\Support\EnvParser" not found`

- [X] **Step 3: Create `src/Support/EnvParser.php`**

```php
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
```

- [X] **Step 4: Run tests to verify they pass**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit tests/Support/EnvParserTest.php --colors=always
```

Expected: `OK (5 tests, 9 assertions)`

- [X] **Step 5: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add src/Support/EnvParser.php tests/Support/EnvParserTest.php && git commit -m "feat: add EnvParser with tests"
```

---

## Task 4: ConsulKvClient

**Files:**
- Create: `src/Consul/ConsulKvClient.php`
- Create: `tests/Consul/ConsulKvClientTest.php`

`ConsulKvClient` 是 SDK `KV` 服務的 thin wrapper，負責把 SDK exceptions 轉換成我們的 domain exceptions，使上層不依賴 SDK 細節。

- [X] **Step 1: Create `tests/Consul/ConsulKvClientTest.php`**

用 `createMock(ClientInterface::class)` 注入假的 Client，模擬各種情況。

```php
<?php

namespace OlympusHub\Tests\Consul;

use Consul\ClientInterface;
use Consul\ConsulResponse;
use Consul\Exception\ClientException;
use Consul\Exception\ServerException;
use Consul\Services\KV;
use OlympusHub\Consul\ConsulKvClient;
use OlympusHub\Exceptions\ConsulConnectionException;
use OlympusHub\Exceptions\RealmTetherNotFoundException;
use PHPUnit\Framework\TestCase;

class ConsulKvClientTest extends TestCase
{
    public function test_fetch_returns_raw_body_on_success(): void
    {
        $expected = "DB_HOST_READ=0.0.0.0\nDB_PORT=3306";

        $client = $this->createMock(ClientInterface::class);
        $client->expects($this->once())
            ->method('get')
            ->with('v1/kv/DB/TPOWER', ['query' => ['raw' => true]])
            ->willReturn(new ConsulResponse([], $expected, 200));

        $kv = new KV($client);
        $this->assertSame($expected, (new ConsulKvClient($kv))->fetch('DB/TPOWER'));
    }

    public function test_fetch_throws_realm_tether_not_found_on_404(): void
    {
        $this->expectException(RealmTetherNotFoundException::class);

        $client = $this->createMock(ClientInterface::class);
        $client->method('get')->willThrowException(
            new ClientException('Not Found', 404)
        );

        $kv = new KV($client);
        (new ConsulKvClient($kv))->fetch('DB/UNKNOWN');
    }

    public function test_fetch_throws_consul_connection_exception_on_server_error(): void
    {
        $this->expectException(ConsulConnectionException::class);

        $client = $this->createMock(ClientInterface::class);
        $client->method('get')->willThrowException(
            new ServerException('Internal Server Error', 500)
        );

        $kv = new KV($client);
        (new ConsulKvClient($kv))->fetch('DB/TPOWER');
    }

    public function test_fetch_throws_consul_connection_exception_on_non_404_client_error(): void
    {
        $this->expectException(ConsulConnectionException::class);

        $client = $this->createMock(ClientInterface::class);
        $client->method('get')->willThrowException(
            new ClientException('Forbidden', 403)
        );

        $kv = new KV($client);
        (new ConsulKvClient($kv))->fetch('DB/TPOWER');
    }
}
```

- [X] **Step 2: Run test to verify it fails**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit tests/Consul/ConsulKvClientTest.php --colors=always
```

Expected: FAIL — `Class "OlympusHub\Consul\ConsulKvClient" not found`

- [X] **Step 3: Create `src/Consul/ConsulKvClient.php`**

```php
<?php

namespace OlympusHub\Consul;

use Consul\Exception\ClientException;
use Consul\Exception\ServerException;
use Consul\Services\KV;
use OlympusHub\Exceptions\ConsulConnectionException;
use OlympusHub\Exceptions\RealmTetherNotFoundException;

class ConsulKvClient
{
    /** @var KV */
    private $kv;

    public function __construct(KV $kv)
    {
        $this->kv = $kv;
    }

    public function fetch(string $key): string
    {
        try {
            $response = $this->kv->get($key, ['raw' => true]);
            return $response->getBody();
        } catch (ClientException $e) {
            if ($e->getCode() === 404) {
                throw new RealmTetherNotFoundException("Consul key '{$key}' not found");
            }
            throw new ConsulConnectionException($e->getMessage());
        } catch (ServerException $e) {
            throw new ConsulConnectionException("Consul server error: {$e->getMessage()}");
        }
    }
}
```

- [X] **Step 4: Run tests to verify they pass**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit tests/Consul/ConsulKvClientTest.php --colors=always
```

Expected: `OK (4 tests, 5 assertions)`

- [X] **Step 5: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add src/Consul/ConsulKvClient.php tests/Consul/ConsulKvClientTest.php && git commit -m "feat: add ConsulKvClient using friendsofphp/consul-php-sdk with tests"
```

---

## Task 5: DatabaseConfigOverrider

`apply()` 只包裝 Laravel 的 `config()` 與 `DB::purge()`，屬於 framework adapter，不寫單元測試。

**Files:**
- Create: `src/Support/DatabaseConfigOverrider.php`

- [X] **Step 1: Create `src/Support/DatabaseConfigOverrider.php`**

```php
<?php

namespace OlympusHub\Support;

use Illuminate\Support\Facades\DB;

class DatabaseConfigOverrider
{
    public function apply(array $vars): void
    {
        config([
            'database.connections.mysql.read.host'      => [$vars['DB_HOST_READ']],
            'database.connections.mysql.read.username'  => $vars['DB_USERNAME_READ'],
            'database.connections.mysql.read.password'  => $vars['DB_PASSWORD_READ'],
            'database.connections.mysql.write.host'     => [$vars['DB_HOST_WRITE']],
            'database.connections.mysql.write.username' => $vars['DB_USERNAME_WRITE'],
            'database.connections.mysql.write.password' => $vars['DB_PASSWORD_WRITE'],
            'database.connections.mysql.database'       => $vars['DB_DATABASE'],
            'database.connections.mysql.port'           => $vars['DB_PORT'],
        ]);

        DB::purge('mysql');
    }
}
```

- [X] **Step 2: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add src/Support/DatabaseConfigOverrider.php && git commit -m "feat: add DatabaseConfigOverrider"
```

---

## Task 6: ConsulRealmTetherService

**Files:**
- Create: `src/Services/ConsulRealmTetherService.php`
- Create: `tests/Services/ConsulRealmTetherServiceTest.php`

- [X] **Step 1: Create `tests/Services/ConsulRealmTetherServiceTest.php`**

```php
<?php

namespace OlympusHub\Tests\Services;

use OlympusHub\Consul\ConsulKvClient;
use OlympusHub\Exceptions\RealmTetherNotFoundException;
use OlympusHub\Services\ConsulRealmTetherService;
use OlympusHub\Support\DatabaseConfigOverrider;
use OlympusHub\Support\EnvParser;
use PHPUnit\Framework\TestCase;

class ConsulRealmTetherServiceTest extends TestCase
{
    private string $rawValue = "DB_CONNECTION=mysql\nDB_HOST_READ=0.0.0.0\nDB_HOST_WRITE=0.0.0.0\nDB_PORT=3306\nDB_DATABASE=xsgcasino\nDB_USERNAME_READ=reader\nDB_PASSWORD_READ=rpass\nDB_USERNAME_WRITE=writer\nDB_PASSWORD_WRITE=wpass";

    private array $parsedVars = [
        'DB_CONNECTION'     => 'mysql',
        'DB_HOST_READ'      => '0.0.0.0',
        'DB_HOST_WRITE'     => '0.0.0.0',
        'DB_PORT'           => '3306',
        'DB_DATABASE'       => 'xsgcasino',
        'DB_USERNAME_READ'  => 'reader',
        'DB_PASSWORD_READ'  => 'rpass',
        'DB_USERNAME_WRITE' => 'writer',
        'DB_PASSWORD_WRITE' => 'wpass',
    ];

    public function test_activate_fetches_parses_and_applies_config(): void
    {
        $client   = $this->createMock(ConsulKvClient::class);
        $parser   = $this->createMock(EnvParser::class);
        $overrider = $this->createMock(DatabaseConfigOverrider::class);

        $client->expects($this->once())->method('fetch')->with('DB/TPOWER')->willReturn($this->rawValue);
        $parser->expects($this->once())->method('parse')->with($this->rawValue)->willReturn($this->parsedVars);
        $overrider->expects($this->once())->method('apply')->with($this->parsedVars);

        (new ConsulRealmTetherService($client, $parser, $overrider))->activate('tpower');
    }

    public function test_activate_uppercases_realm_id_for_consul_key(): void
    {
        $client   = $this->createMock(ConsulKvClient::class);
        $parser   = $this->createMock(EnvParser::class);
        $overrider = $this->createMock(DatabaseConfigOverrider::class);

        $client->expects($this->once())->method('fetch')->with('DB/ALFAWIN77')->willReturn($this->rawValue);
        $parser->method('parse')->willReturn($this->parsedVars);
        $overrider->method('apply');

        (new ConsulRealmTetherService($client, $parser, $overrider))->activate('alfawin77');
    }

    public function test_activate_propagates_realm_tether_not_found_exception(): void
    {
        $this->expectException(RealmTetherNotFoundException::class);

        $client   = $this->createMock(ConsulKvClient::class);
        $parser   = $this->createMock(EnvParser::class);
        $overrider = $this->createMock(DatabaseConfigOverrider::class);

        $client->method('fetch')->willThrowException(new RealmTetherNotFoundException('not found'));

        (new ConsulRealmTetherService($client, $parser, $overrider))->activate('tpower');
    }
}
```

- [X] **Step 2: Run test to verify it fails**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit tests/Services/ConsulRealmTetherServiceTest.php --colors=always
```

Expected: FAIL — `Class "OlympusHub\Services\ConsulRealmTetherService" not found`

- [X] **Step 3: Create `src/Services/ConsulRealmTetherService.php`**

```php
<?php

namespace OlympusHub\Services;

use OlympusHub\Consul\ConsulKvClient;
use OlympusHub\Support\DatabaseConfigOverrider;
use OlympusHub\Support\EnvParser;

class ConsulRealmTetherService
{
    /** @var ConsulKvClient */
    private $client;

    /** @var EnvParser */
    private $parser;

    /** @var DatabaseConfigOverrider */
    private $overrider;

    public function __construct(
        ConsulKvClient $client,
        EnvParser $parser,
        DatabaseConfigOverrider $overrider
    ) {
        $this->client   = $client;
        $this->parser   = $parser;
        $this->overrider = $overrider;
    }

    public function activate(string $realmId): void
    {
        $key  = 'DB/' . strtoupper($realmId);
        $raw  = $this->client->fetch($key);
        $vars = $this->parser->parse($raw);
        $this->overrider->apply($vars);
    }
}
```

- [X] **Step 4: Run tests to verify they pass**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit tests/Services/ConsulRealmTetherServiceTest.php --colors=always
```

Expected: `OK (3 tests, 5 assertions)`

- [X] **Step 5: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add src/Services/ConsulRealmTetherService.php tests/Services/ConsulRealmTetherServiceTest.php && git commit -m "feat: add ConsulRealmTetherService with tests"
```

---

## Task 7: Config & ServiceProvider Binding

**Files:**
- Modify: `config/olympus.php`
- Modify: `src/OlympusHubServiceProvider.php`

- [X] **Step 1: Update `config/olympus.php`** — add `consul` block at the end:

- [X] **Step 2: Update `src/OlympusHubServiceProvider.php`** — add `registerServices()` and its call:

- [X] **Step 3: Run full test suite to verify nothing breaks**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit --colors=always
```

Expected: All existing tests still pass.

- [X] **Step 4: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add config/olympus.php src/OlympusHubServiceProvider.php && git commit -m "feat: add Consul config and bind ConsulRealmTetherService"
```

---

## Task 8: HubController::invoke() Modification

`HubController` 屬於 web layer（HTTP routing + Artisan facade），不寫單元測試。

**Files:**
- Modify: `src/Http/Controllers/HubController.php`

- [X] **Step 1: Replace `src/Http/Controllers/HubController.php`**

- [X] **Step 2: Run full test suite**

```bash
cd /home/rhys/projects/www/olympus-hub && vendor/bin/phpunit --colors=always
```

Expected: All tests pass (12 tests across EnvParser, ConsulKvClient, ConsulRealmTetherService).

- [X] **Step 3: Commit**

```bash
cd /home/rhys/projects/www/olympus-hub && git add src/Http/Controllers/HubController.php && git commit -m "feat: activate Consul realm tether in HubController::invoke()"
```
