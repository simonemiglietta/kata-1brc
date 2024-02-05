# One Billion Rows Challenge

This kata is based on [gunnarmorling/1brc](https://github.com/gunnarmorling/1brc) challenge.

Differently from the original project, the main goal is to explore tools and performances of different languages and/or
algorithms.

## Initial Setup

Working data are included in `data` folder. But this repository cannot include the full set of data because 1 billion
rows take around 15 GB of disk space.

Instead, `data` folder includes a sample list of distinct weather stations and a several tools that can be used to
generate the full samples list over 10 thousand stations.

In folder `data/generators` you can find several tools. At this moment, the only one fully working is the Python one.

## Scripts prerequisites

Each script must have these characteristics:

* fully contained in the folder `src/<script-name>`
* work as a shell command, without a GUI
* has not explicit error handling (job boundaries are predetermined and well defined)
* explain build and run instructions in a `README.md` file
* represent the work in progress updated every second (or more, in order to not interfere significantly with total
  execution time)
* a progress bar will be very appreciated
* read data from `data/measurements.txt` and write results in its own folder (script folder)
* test against data samples provided in `data/testcases` folder

## Results

All the tests have been performed on a Dell Inspiron 5502 PC (CPU  i7-1165G7 4 core 8 thread, RAM 12 GB)

| Script     | Language | Algorithm notes | Exec Time | Memory Used |
|------------|----------|-----------------|-----------|-------------|
| ScriptName | Italiano | Placeholder row | Tanto     | Poco        |

Each script is evaluated by the shell command `/usr/bin/time -f "time: %E - memory: %M" <command>`. Pay attention: you
have to use the `time` linux program, not the shell command!

All the tests used to fill the above table are listed in `data/rinnings.csv` file. 
