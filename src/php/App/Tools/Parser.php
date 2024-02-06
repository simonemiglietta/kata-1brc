<?php

namespace App\Tools;

use App\ValueObjects\StationAggregate;
use Generator;

class Parser
{
    /**
     * @return Generator<void, void, int, StationAggregate[]>
     */
    public static function aggregateGenerator(string $sourceFile, string $destFile): Generator
    {
        /** @var array<string,StationAggregate> $aggregates */
        $aggregates = [];
        $i = 0;

        foreach (FileTool::fileParser($sourceFile) as $detection) {
            $aggregate = $aggregates[$detection->station] ?? null;

            if ($aggregate) {
                $aggregate->addDetection($detection);
            } else {
                $aggregates[$detection->station] = new StationAggregate($detection);
            }

            yield ++$i;
        }

        usort($aggregates, fn($a, $b) => $a->station <=> $b->station);
        FileTool::writeAggregates($destFile, $aggregates);

        return $aggregates;
    }
}
