# CS2-Demo-parser

Takes a folder of CS2 demos and creates a spreadsheet with detailed stats

Credit:
Demo file parsing - https://github.com/markus-wa/demoinfocs-golang

## Spreadsheet generation times

| # of demos | Total file size | Single-core  | Multi-core     |
| ---------- | --------------- | ------------ | -------------- |
| 200        | ~15gbs          | 700+ seconds | 90-120 seconds |
| 60         | ~5gbs           | 500+ seconds | 50-55 seconds  |
| 14         | ~1gbs           | 90+ seconds  | 10-17 seconds  |
| 5          | ~500mbs         | 20 seconds   | 5 seconds      |

### Getting Started

1. Clone the repository to your local machine:

   ```sh
   git clone https://github.com/FlynnFc/CS2-Demo-parser.git
   ```

2. Navigate to the cloned directory:

   ```sh
   cd CS2-Demo-parser
   ```

3. Run the project:

   ```sh
   go run main
   ```

## Usage

Run the program from the command line. You will be prompted for a folder path. enter the FULL path

### Example

```sh
Please enter the path of the demo folder:  C:\Users\user\demos
```

## Example output

### Basic Information and Performance
| ID                | Name          | Team Name      | Matches | Rounds | Kills | Assists | Deaths | Damage | ADR  |
|-------------------|---------------|----------------|---------|--------|-------|---------|--------|--------|------|
| 76561198016577250 | TAVARES       | team_androxZ-  | 1       | 20     | 14    | 0       | 18     | 1797   | 89   |
| 76561198857416779 | sneppyyyyyyy  | team_androxZ-  | 1       | 20     | 22    | 0       | 10     | 1926   | 96   |
| 76561199004526824 | Ourob         | Shaman         | 1       | 36     | 20    | 0       | 24     | 2244   | 62   |
| 76561199082282831 | JJaredd       | Shaman         | 1       | 36     | 27    | 0       | 21     | 2981   | 82   |
| 76561198157156096 | dell-w        | TYREECESIMPSON | 2       | 53     | 42    | 0       | 33     | 3982   | 75   |


### All other columns

| Total Entry Attempts | CT Entry Kills | CT Entry Deaths | T Entry Kills | T Entry Deaths |
|----------------------|----------------|-----------------|---------------|----------------|

| CT Traded Kills | CT Failed Trades | CT Traded Deaths | T Trade Kills | T Failed Trades | T Traded Deaths |
|------------------|-----------------|-----------------|-----------------|-----------------|-----------------|

| 1v1 Attempts | 1v1 Wins | 1v2 Attempts | 1v2 Wins | 1v3 Attempts | 1v3 Wins | 1v4 Attempts | 1v4 Wins | 1v5 Attempts | 1v5 Wins |
|--------------|----------|--------------|----------|--------------|----------|--------------|----------|--------------|----------|





