# Go-Gol
Serialized Game of Life implementation in Go.

## Usage
.go file placed in src/go_gol. Use
```
go_gol -fpath <filepath> 
```
to run the simulation

## Config Structure
Sample configuration for the [glider pattern](https://en.wikipedia.org/wiki/Conway's_Game_of_Life#Examples_of_patterns).
```
10    // 10x10 grid
50    // run for 50 iterations
5     // 5 points in grid that are alive
5 5   // (row, col) coordinates
6 6 
6 7 
7 5
7 6
```
