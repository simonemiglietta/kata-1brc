<?php

namespace App\ValueObjects;

use Stringable;

readonly class Detection implements Stringable
{
    public function __construct(
        public string $station,
        public float  $temperature
    ) {}

    public function __toString(): string
    {
        return "{$this->station};{$this->temperature}";
    }

    static public function fromRow(string $row): static {
        $split = explode(';', $row);

        return new Detection($split[0], $split[1]);
    }
}
