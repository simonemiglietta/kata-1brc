<?php

namespace App\Tools;

use App\ValueObjects\Detection;
use App\ValueObjects\StationAggregate;
use Generator;

class FileTool
{

    /**
     * @return Generator<void, void, Detection, void>
     */
    public static function fileParser(string $fileName): Generator
    {
        $file = fopen($fileName, 'r');

        while (($row = fgets($file)) !== false) {
            yield Detection::fromRow($row);
        }

        fclose($file);
    }

    /**
     * @param StationAggregate[] $aggregates
     */
    public static function writeAggregates(string $fileName, array $aggregates): void {
        $file = fopen($fileName, 'w');

        fwrite( $file,'{');
        fwrite($file, implode(', ', $aggregates));
        fwrite( $file,'}' . PHP_EOL);

        fclose($file);
    }
}
