package main

import (
    "fmt"
    "os"
    "os/exec"
    "flag"
    "io"
    "time"
    "log"
)

const ALIVECHR string = "@ "
const DEADCHR string = "- "
const SLEEPDUR time.Duration = 200

/* returns positive modulo */
func posModulo(val, mod int) int {
    result := val % mod
    if result < 0 {
        return result + mod
    }
    return result
}

/* set cell to alive or dead */
func setAlive(isalive bool, board [][]string, row, col int) {
    if isalive {
        board[row][col] = ALIVECHR
    } else {
        board[row][col] = DEADCHR
    }
}

/* Returns a new square board with length specified by size */
func createBoard(size int) (board [][]string) {
    board = [][]string{}
    for i := 0; i < size; i++ {
        row := make([]string, size)
        for j := 0; j < size; j++ {
            row[j] = DEADCHR
        }
        board = append(board, row)
    }
    return
}

/* Creates a copy of a board */
func copyBoard(board [][]string) (duplicate [][]string) {
    duplicate = make([][]string, len(board))
    for i := range board {
        duplicate[i] = make([]string, len(board))
        copy(duplicate[i], board[i])
    }
    return duplicate
}

/* Print out the given board*/
func printBoard(board [][]string) {
    for i, _ := range board {
        for j, _ := range board[i] {
            fmt.Print(board[i][j])
        }
        fmt.Println()
    }
}

/* inititalizes and returns boards and threadinfo */
func initGameVars(fptr io.Reader) (board1, board2 [][]string, iterations int){
    var size, cellstoinit, row, col int
    // get board and file parsing details
    _, err := fmt.Fscanf(fptr, "%d", &size)
    _, err = fmt.Fscanf(fptr, "%d", &iterations)
    _, err = fmt.Fscanf(fptr, "%d", &cellstoinit)
    if err != nil {
        log.Fatal("File Format error")
    }

    // create boards
    board1 = createBoard(size)
    board2 = createBoard(size)

    // init cells
    for i := 0; i < cellstoinit; i++ {
        _, err := fmt.Fscanf(fptr, "%d %d", &row, &col)
        if err != nil {
            log.Fatal("File Format error")
        }
        fmt.Printf("row %d, col %d\n", row, col)
        setAlive(true, board1, row, col)
     }

    return
}

/* get number of alive neighbors around specified cell */
func getNumAliveNeighbors(board [][]string, row, col int) int {
    size := len(board)
    count := 0
    // check cells above
    for i:= 0; i < 3; i++ {
        r := posModulo(row-1, size)
        c := posModulo(col+1-i, size)
        if board[r][c] == ALIVECHR {
            count++
        }
    }
    // check cells on side
    if board[row][posModulo(col-1,size)] == ALIVECHR {
        count++
    }
    if board[row][posModulo(col+1,size)] == ALIVECHR {
        count++
    }
    // check cells below
    for i := 0; i < 3; i++ {
        r := posModulo(row+1, size)
        c := posModulo(col+1-i, size)
        if board[r][c] == ALIVECHR {
            count++
        }
    }
    return count
}

/* check if specified cell should be alive or dead in the next iteration */
func updateCell(mod, orig [][]string, row, col int) {
    // check if cell is alive
    alive := orig[row][col] == ALIVECHR
    // if alive,
    //      < 2 neighbors: dies
    //      2 or 3 neighbors: stay alive
    //      > 3 neighbors: dies
    // else,
    //      3 neighbors: lives
    aliveNeighbors := getNumAliveNeighbors(orig, row, col)
    if alive {
        if aliveNeighbors < 2 || aliveNeighbors > 3 {
            mod[row][col] = DEADCHR
        } else {
            mod[row][col] = orig[row][col]
        }
    } else {
        if aliveNeighbors == 3 {
            mod[row][col] = ALIVECHR
        }
    }
}

/* Run one step of the gol simulation */
func golstep(board1, board2 [][]string) ([][]string, [][]string) {
    // loop through cells and update
    for r := 0; r < len(board1); r++ {
        for c := 0; c < len(board1); c++ {
            updateCell(board2, board1, r, c)
        }
    }
    duplicate := copyBoard(board1)
    board1 = board2

    return board2, duplicate
}

/* Displays the gol simulation for the specified number of iterations */
func rungol(board1, board2 [][]string, iterations int) {
    // loop for iterations
    //      run a step of the gol
    // Initial print
    printBoard(board1)
    fmt.Println()

    for i := 0; i < iterations; i++ {
        // run a step of the gol
        board1, board2 = golstep(board1, board2)
        // print out the result
        printBoard(board1)
        fmt.Println()
        time.Sleep(SLEEPDUR * time.Millisecond)
        cmd := exec.Command("clear")
        cmd.Run()
    }
}

func main() {
    // Get filename from cmdl
    filename := flag.String("fpath", "test.txt", "file path of board init")
    flag.Parse()

    // Try to open file
    f, err := os.Open(*filename)
    if err != nil {
        flag.PrintDefaults()
        log.Fatal(err)
    }
    defer func() {
        if err = f.Close(); err != nil {
        log.Fatal(err)
    }
    }()

    board1, board2, iterations := initGameVars(f)

    rungol(board1, board2, iterations)
}
