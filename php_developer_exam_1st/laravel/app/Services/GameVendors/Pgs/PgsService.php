<?php

namespace App\Services\GameVendors\target;

use App\Services\GameVendors\AbstractGameVendorService;
use App\Services\GameVendors\AuthStrategies\AuthStrategyInterface;
use App\Services\GameVendors\AuthStrategies\PayloadAuthStrategy;
use App\Services\GameVendors\target\PgsResponse;
use App\Services\GameVendors\ResponseInterface;
use Illuminate\Support\Facades\Cache;

class PgsService extends AbstractGameVendorService
{
    const FIELD_USER_NAME = "playerName";
    const FIELD_GAME_ID = "gameID";

    protected function authenticate(): array
    {
        return [
            "merchantId" => $this->getConfig("merchant_id"),
            "merchantSecretKey" => $this->getConfig("merchant_secret"),
        ];
    }

    protected function getAuthStrategy(): AuthStrategyInterface
    {
        return new PayloadAuthStrategy();
    }

    protected function getVendorResponse(): ResponseInterface
    {
        return new PgsResponse();
    }

    protected function getEndpoint(string $action): ?string
    {
        return match ($action) {
            "getGameList" => "/games",
            "isUserExist" => "/player/exist",
            "createUser" => "/player",
            "getBalance" => "/wallet/amount",
            "deposit" => "/wallet/transfer-in",
            "withdraw" => "/wallet/transfer-out",
            "getTransactionById" => "/transfer",
            "getGameRecords" => "/game-records",
            "getGameRecordsByTransactionNo" => "/gamer-records/transactionNo",
            "getMerchantSummaryReport" => "/game-records/report",
            "getGameUrl" => "/game/link",
            "getEventAwards" => "/event/prize-records",
            "forceQuitUser" => "/player/logout",
            default => throw new \InvalidArgumentException("Unsupported action {$action} for target vendor")
        };
    }

    public function getGameList(): array
    {
        $cacheKey = "target:game_list";
        #can get from crontab
        $cacheTtl = $this->getConfig("game_list_cache_ttl"); // default to 24 hours
        $redisCache = Cache::store('redis');

        // Use cache if it exists
        if ($redisCache->has($cacheKey)) {
            return $redisCache->get($cacheKey);
        }

        // Cache doesn't exist, make API request and store in cache
        $response = $this->makeGetRequest("getGameList");
        $gameList = $this->getVendorResponse()->init($response)->getGameList();
        $redisCache->put($cacheKey, $gameList, $cacheTtl);

        return $gameList;
    }

    public function getGameUrl(string $username, string $gameId, array $options = []): string
    {
        $payload = array_merge([
            self::FIELD_USER_NAME => $username,
            self::FIELD_GAME_ID => $gameId,
        ], $options);
        $response = $this->makeGetRequest("getGameUrl", $payload);

        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    public function isUserExist(string $username): bool
    {
        $response = $this->makeGetRequest(
            "isUserExist",
            [self::FIELD_USER_NAME => $username]
        );

        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    public function createUser(string $username): bool
    {
        $response = $this->makePostRequest(
            "createUser",
            [self::FIELD_USER_NAME => $username]
        );

        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    public function createUserIfNotExist(string $username): bool
    {
        if (!$this->isUserExist($username)) {
            return $this->createUser($username);
        } else {
            return true;
        }
    }

    public function getBalance(string $username): float
    {
        $response = $this->makeGetRequest("getBalance", [self::FIELD_USER_NAME => $username]);
        $data = $this->getVendorResponse()->init($response)->{__FUNCTION__}();
        return $data;
    }

    #需新增回應物件包裝回傳結果
    public function deposit(string $username, float $amount, string $accountTransactionNo): array
    {
        $response = $this->makePostRequest("deposit", [
            self::FIELD_USER_NAME => $username,
            "amount" => $amount,
            "merchantOrderNo" => $accountTransactionNo,
        ]);
        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    public function withdraw(string $username, float $amount, string $accountTransactionNo): array
    {
        $response = $this->makePostRequest("withdraw", [
            self::FIELD_USER_NAME => $username,
            "amount" => $amount,
            "merchantOrderNo" => $accountTransactionNo,
        ]);
        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    public function getTransactionById(array $filters = []): array
    {
        $response = $this->makeGetRequest("getTransactionById", $filters);
        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    /**
     * Get game records within a time range with optional username filter
     * @param string $startTime date("Y-m-d H:i:s")
     * @param string $endTime  date("Y-m-d H:i:s")
     * @param int $page start from 1
     * @param int $pageSize max 10000
     * @param string $username
     */
    public function getGameRecords(array $filters = []): array
    {
        $fixed_filters = $this->validateGameRecordsFilters($filters);
        $response = $this->makeGetRequest("getGameRecords", $filters);

        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    //OdinTodo: use a class instead of this method to validate filters
    private function validateGameRecordsFilters(array $filters): array
    {
        $requiredFields = ["startTime", "endTime", "pageIndex", "pageSize", self::FIELD_USER_NAME];
        $fixed_filters = [];
        $isStartTimeEmpty = true;
        $isEndTimeEmpty = true;
        $isPageIndexEmpty = true;
        $isPageSizeEmpty = true;

        foreach ($requiredFields as $field) {
            switch ($field) {
            //case "vendorId":
                //if (!isset($filters[$field]) or !is_string($filters[$field]) or empty($filters[$field])) {
                    //throw new \InvalidArgumentException("Invalid or missing {$field} in filters for getGameRecords");
                //}
                //$fixed_filters[$field] = $filters[$field];
                //break;

            case "startTime":
                if (isset($filters[$field]) and strtotime($filters[$field])) {
                    $isStartTimeEmpty = false;
                    $fixed_filters[$field] = $filters[$field];
                }
                break;

            case "endTime":
                if (isset($filters[$field]) and strtotime($filters[$field])) {
                    $isEndTimeEmpty = false;
                    $fixed_filters[$field] = $filters[$field];
                }

                if ($isStartTimeEmpty or $isEndTimeEmpty) {
                    throw new \InvalidArgumentException("Either startTime or endTime must be provided in filters for getGameRecords");
                }

                $startTimestamp = strtotime($fixed_filters["startTime"]);
                $endTimestamp = strtotime($fixed_filters["endTime"]);

                if (($startTimestamp > $endTimestamp) or ($endTimestamp - $startTimestamp) > 8 * 3600) {
                    throw new \InvalidArgumentException(
                        "Invalid time range in filters for getGameRecords,
                        startTime must be less than endTime and the time range must not exceed 8 hours
                    ");
                }
                break;

            case "pageIndex":
                if (!isset($filters[$field]) or !is_int($filters[$field]) or $filters[$field] <= 0) {
                    $fixed_filters[$field] = 1;
                } else {
                    $fixed_filters[$field] = $filters[$field];
                }
                break;

            case "pageSize":
                if (!isset($filters[$field]) or !is_int($filters[$field]) or $filters[$field] <= 0) {
                    $fixed_filters[$field] = $this->getConfig("game_records_page_size");
                } else{
                    $fixed_filters[$field] = $filters[$field];
                }
                break;

            case self::FIELD_USER_NAME:
                if (isset($filters[$field]) and is_string($filters[$field])) {
                    $fixed_filters[$field] = $filters[$field];
                }
                break;
            }
        }

        return $fixed_filters;
    }

    public function getGameRecordsByTransactionNo(array $filters = []): array
    {
        $response = $this->makeGetRequest("getGameRecordsByTransactionNo", $filters);
        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    public function forceQuitUser(string $username): bool
    {
        $response = $this->makePostRequest("forceQuitUser", [self::FIELD_USER_NAME => $username]);
        return $this->getVendorResponse()->init($response)->{__FUNCTION__}();
    }

    //OdinTodo: add this to interface
    public function getVendorFieldUserName(): string
    {
        return self::FIELD_USER_NAME;
    }

    protected function processMerchantSummaryReportResponse(array $response): array
    {
        return $response["data"] ?? [];
    }

    protected function processEventAwardsResponse(array $response): array
    {
        return $response["data"] ?? [];
    }
}
