<?php

namespace App\ValueObjects;

use Stringable;

class StationAggregate implements Stringable
{
    public string $station;
    public int $itemCount;
    public float $temperatureMinimum;
    public float $temperatureMaximum;
    public float $temperatureSum;

    public function __construct(Detection $detection)
    {
        $this->station = $detection->station;
        $this->itemCount = 1;
        $this->temperatureMinimum = $detection->temperature;
        $this->temperatureMaximum = $detection->temperature;
        $this->temperatureSum = $detection->temperature;
    }

    public function addDetection(Detection $detection): static
    {
        $this->temperatureMinimum = min($this->temperatureMinimum, $detection->temperature);
        $this->temperatureMaximum = max($this->temperatureMaximum, $detection->temperature);
        $this->itemCount++;
        $this->temperatureSum += $detection->temperature;

        return $this;
    }

    private function temperatureMean(): float
    {
        return $this->temperatureSum / $this->itemCount;
    }

    public function __toString(): string
    {
        $roundedMean = round($this->temperatureMean(),1);

        return sprintf(
            '%s=%.1f/%.1f/%.1f',
            $this->station, $this->temperatureMinimum, $roundedMean, $this->temperatureMaximum
        );
    }
}
