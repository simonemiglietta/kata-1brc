<?php

require __DIR__ . '/vendor/autoload.php';

use App\Tools\Parser;
use Dariuszp\CliProgressBar;

const MAX_ROWS = 1_000_000_000;
const SRC_FILE = __DIR__ . '/../../data/measurements.txt';
const DST_FILE = __DIR__ . '/measurements.out';

$bar = new CliProgressBar(MAX_ROWS);

$lastEpoch = 0;
$generator = Parser::aggregateGenerator(SRC_FILE, DST_FILE);
foreach ($generator as $counter) {
    $currentEpoch = time();

    $bar->setCurrentStep($counter);

    if ($currentEpoch != $lastEpoch) {
        $bar->display();
        $lastEpoch = $currentEpoch;
    }
}

$bar->setCurrentStep(MAX_ROWS);
$bar->display();
$bar->end();

$aggregates = $generator->getReturn();

echo 'Extracted ' . count($aggregates) . ' stations' . PHP_EOL;
