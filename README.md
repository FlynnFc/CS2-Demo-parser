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
