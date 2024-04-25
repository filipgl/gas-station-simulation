# Simulation of gas station

## Idea
This shema represents how the gas station operates:
![idea](https://github.com/filipgl/gas-station-simulation/assets/72188289/9ef028fb-433a-436b-9444-f9ed0dfadabd)

## Solution schema
Below schema represents how upper idea is expressed in application. It is inspired by [this article](https://go.dev/blog/pipelines):
![solution_schema](https://github.com/filipgl/gas-station-simulation/assets/72188289/97fab127-b9a9-4197-8bf4-c828d0925c05)

## Detail of group block
Group block enables distribution of `Cars` to the best `Server` (`Server` = `Pump` or cash `Register`).
Best server is the least occupied.
![group_schema](https://github.com/filipgl/gas-station-simulation/assets/72188289/bfd08a62-50e7-4257-a9f4-20ba11afbd3b)
