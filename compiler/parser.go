package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)
var mem map[int]string = make(map[int]string, 16367)
var rmem map[string]int = make(map[string]int, 16367)
var jmp map[string]int = make(map[string]int, 100)

const starM int = 16
var countM int = 0

const data string = "@33\nM=1 //this is a comment\nD=A+1;JEQ\nM=0;JEQ\n@R0\n(LOOP)\n@i\nM=1\n@LOOP"

func main() {

    mem[16384] = "SCREEN"
    mem[24576] = "KBD"

    fc, err := os.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println("error when reading file contents " + err.Error())
        return
    }
    fc_str := string(fc)
    str := getLines(fc_str)
    result := handleInstructions(sanitizeSymbols(str[:(len(str)-1)]))
    output_file := giveOutputFile(os.Args[1])
    fo, err := os.Create(output_file)
    if err != nil {
        panic(err)
    }
    defer fo.Close()
    for i := range result {
        fmt.Fprintln(fo, result[i])
    }

}

func giveOutputFile(data string) string {

    path := strings.Split(data, "/")
    file := path[len(path)-1]
    fmt.Println(path)
    file_name := strings.Split(file, ".")[0]
    return file_name + ".asma"
}

func handleInstructions(data []string) []string {
    var result []string
    n := len(data)
    for i:=0; i<n; i++ {
        if strings.HasPrefix(data[i], "@") {
            fmt.Println("A instruction")
            result = append(result, handleAInstruction(data[i]))
        } else {
            fmt.Println("C instruction")
            result = append(result, handleCInstructions(data[i]))
        }
    }

    return result
}

func handleAInstruction(data string) string {
    var a string
    if data[1] == 'R' {
        a = data[2:]
    } else {
        a = data[1:]
    }
    var result strings.Builder
    result.WriteString("0")
    i, _ := strconv.Atoi(a)
    if i > 15 {
        //throw err and advice to use a symbol
    }
    str := strconv.FormatInt(int64(i),2)
    result.WriteString(fmt.Sprintf("%015s", str))
    return result.String()
}

func handleCInstructions(data string) string {
    var result strings.Builder
    result.WriteString("111") // c instruction padding
    if strings.Contains(data, "=") { //destination exists

        lhs := strings.Split(data, "=")[0]
        rhs := strings.Split(data, "=")[1]

        if strings.Contains(rhs, ";") {
            comp := strings.Split(rhs, ";")[0]
            jmp := strings.Split(rhs, ";")[1]
            if strings.Contains(comp, "A") && !strings.Contains(comp, "M") {
                result.WriteString("0") // a == 0 
            } else {
                result.WriteString("1") // a == 1
            }
            d := handleLhs(lhs)
            c := handleComputation(comp)
            j := handleJump(jmp)
            result.WriteString(c) //computation control bits
            result.WriteString(d) //destination
            result.WriteString(j) //jump to
        } else {
            if strings.Contains(rhs, "A") && !strings.Contains(rhs, "M") {
                result.WriteString("0") // a == 0 
            } else {
                result.WriteString("1") // a == 1
            }
            result.WriteString(handleComputation(rhs))
            result.WriteString(handleLhs(lhs))
            result.WriteString("000")
        }
    } else { //destination does not exist
        if strings.Contains(data, ";") {
            comp := strings.Split(data, ";")[0]
            jmp := strings.Split(data, ";")[1]
            if strings.Contains(comp, "A") && !strings.Contains(comp, "M") {
                result.WriteString("0")
            } else {
                result.WriteString("1")
            }
            c := handleComputation(comp)
            j := handleJump(jmp)
            result.WriteString(c)
            result.WriteString("000") //destination
            result.WriteString(j)
        }
    }
    return result.String()
}

