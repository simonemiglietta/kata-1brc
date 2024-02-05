<?php

namespace Tests;

use App\Tools\Parser;
use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;

class ParserTest extends TestCase
{
    #[DataProvider('fileDataProvider')]
    public function testFileShouldMatch(string $sourceFile, string $expectedFile)
    {
        $actualFile = 'measurements.out';

        foreach (Parser::aggregateGenerator($sourceFile, $actualFile) as $count) {
            // do nothing
        }

        $this->assertTrue($this->filesCompare($actualFile, $expectedFile));
    }

    public static function fileDataProvider(): array
    {
        $files = glob(__DIR__ . '/../../../data/testcases/*.txt');

        $cases = [];

        foreach ($files as $srcFile) {
            $pathInfo = pathinfo($srcFile);
            $label = $pathInfo['filename'];
            $destFile = str_replace('.txt', '.out', $srcFile);
            $cases[$label] = [$srcFile, $destFile];
        }

        return $cases;
    }

    private function filesCompare(string $a, string $b): bool
    {
        // Check if filesize is different
        if (filesize($a) !== filesize($b))
            return false;

        // Check if content is different
        $ah = fopen($a, 'rb');
        $bh = fopen($b, 'rb');

        $result = true;
        while (!feof($ah)) {
            if (fread($ah, 8192) != fread($bh, 8192)) {
                $result = false;
                break;
            }
        }

        fclose($ah);
        fclose($bh);

        return $result;
    }
}
