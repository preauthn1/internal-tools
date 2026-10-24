<?php

namespace OlympusHub\Tests\Support;

use OlympusHub\Exceptions\InvalidTetherDataException;
use OlympusHub\Support\EnvParser;
use PHPUnit\Framework\TestCase;

/**
 * ⚡ 聖域考驗：環境變數解析器測試 (Divine Trial: Env Parser Test)
 * 驗證從觀測位面獲取的原始能量 (Raw Env String) 是否能被正確解構。
 */
class EnvParserTest extends TestCase
{
    /** @var EnvParser 解析器實例 */
    private EnvParser $parser;

    /**
     * 儀式準備：初始化解析器。
     */
    protected function setUp(): void
    {
        $this->parser = new EnvParser();
    }

    /**
     * 測試：將 Env 格式轉化為關聯陣列 (Divine Transmutation: Env to Array)
     * 確保各項能量指標 (DB_HOST, DB_PORT...) 皆能被精確歸類。
     */
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

    /**
     * 測試：當輸入為空字串時拋出異常 (Calamity: Empty Energy)
     * 驗證當獲取到的能量為空時，聖域會判定資料異常。
     */
    public function test_parse_throws_when_empty_string(): void
    {
        $this->expectException(InvalidTetherDataException::class);
        $this->expectExceptionMessage('Tether value is empty');

        $this->parser->parse('');
    }

    /**
     * 測試：當輸入僅包含空白字元時拋出異常 (Calamity: Void Energy)
     * 驗證無意義的能量波動會引發災厄。
     */
    public function test_parse_throws_when_whitespace_only(): void
    {
        $this->expectException(InvalidTetherDataException::class);
        $this->expectExceptionMessage('Tether value is empty');

        $this->parser->parse("   \n  ");
    }

    /**
     * 測試：忽略不符合格式的行 (Divine Sieve: Ignore Invalid Lines)
     * 確保解析器能自動過濾雜訊，僅保留具備等號 (=) 的有效神諭。
     */
    public function test_parse_ignores_lines_without_equals(): void
    {
        $raw = "DB_HOST_READ=0.0.0.0\ninvalid_line\nDB_PORT=3306";

        $result = $this->parser->parse($raw);

        $this->assertSame('0.0.0.0', $result['DB_HOST_READ']);
        $this->assertSame('3306',     $result['DB_PORT']);
        $this->assertArrayNotHasKey('invalid_line', $result);
    }

    /**
     * 測試：處理包含特殊字元的數值 (Divine Tolerance: Special Characters)
     * 驗證當能量數值 (如密碼) 包含混沌符號時，仍能被原樣保留。
     */
    public function test_parse_handles_value_containing_special_chars(): void
    {
        $raw = "DB_PASSWORD_WRITE=!QAZ2wsxForXSGames_online";

        $result = $this->parser->parse($raw);

        $this->assertSame('!QAZ2wsxForXSGames_online', $result['DB_PASSWORD_WRITE']);
    }
}
