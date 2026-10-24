<?php

namespace App\Services\GameVendors\target;

use App\Services\GameVendors\ResponseInterface;
use Illuminate\Http\Client\Response as HttpResponse;

class PgsResponse implements ResponseInterface
{
    private HttpResponse $rawResponse;
    private mixed $bodyData = [];

    public function __construct() {}

    public function init(HttpResponse $response): ResponseInterface
    {
        $this->rawResponse = $response;
        $decoded_json_res = $response->json();

        if ($decoded_json_res === false) {
            throw new \RuntimeException("Invalid API response structure from vendor.");
        } elseif (!isset($decoded_json_res["data"])) {
            throw new \RuntimeException("Response missing required fields from vendor.");
        }

        $this->bodyData = $decoded_json_res["data"];
        return $this;
    }

    public function getGameList(): array
    {
        if (empty($this->bodyData)) {
            throw new \RuntimeException("Response empty");
        }

        return $this->bodyData;
    }

    public function isUserExist(): bool
    {
        if (!isset($this->bodyData["isExist"])) {
            LOG::error("target isUserExist response missing isExist field", ['response_body' => $this->rawResponse->body()]);
            throw new \RuntimeException("Response missing required fields from vendor.");
        }

        return (bool) $this->bodyData["isExist"];
    }

    public function createUser(): bool
    {
        if (!empty($this->bodyData)) {
            throw new \RuntimeException("Create user response an non-empty response which is not compatable with doc");
        }
        return true;
    }

    public function getBalance(): float
    {
        if (!isset($this->bodyData["amountWithConvert"])) {
            LOG::error("target getBalance response missing amountWithConvert field", ['response_body' => $this->rawResponse->body()]);
            throw new \RuntimeException("Response missing required fields from vendor.");
        }

        return floatval($this->bodyData["amountWithConvert"]);
    }

    public function deposit(): array
    {
        if (!isset($this->bodyData["transactionNo"]) or empty($this->bodyData["transactionNo"])) {
            throw new \RuntimeException("empty transactionNo in deposit response");
        }

        return $this->bodyData;
    }

    public function withdraw(): array
    {
        if (!isset($this->bodyData["transactionNo"]) or empty($this->bodyData["transactionNo"])) {
            throw new \RuntimeException("empty merchantOrderNo in withdraw response");
        }

        return $this->bodyData;
    }

    public function getTransactionById(): array
    {
        return $this->bodyData;
    }

    public function getGameRecords(): array
    {
        return is_array($this->bodyData) ? $this->bodyData : [];
    }

    public function getGameUrl(): string
    {
        if (!isset($this->bodyData["link"])) {
            throw new \RuntimeException("No game url given in response");
        }

        return $this->bodyData["link"];
    }

    public function getMerchantSummaryReport(): array
    {
    }

    public function getEventAwards(): array
    {
        return is_array($this->bodyData) ? $this->bodyData : [];
    }

    public function forceQuitUser(): bool
    {
        return (bool) ($this->bodyData['success'] ?? false);
    }
}
