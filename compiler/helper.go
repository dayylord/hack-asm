package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

//elegance
func getLines(data string) []string {
    elems := strings.Split(removeComments(removeWhiteSpaces(data)), "\n")
    return elems
}

func removeWhiteSpaces(data string) string {
    n := len(data)
    var result strings.Builder
    for i:=0; i<n; i++{
        switch data[i] {
            case ' ':
                break
        default:
            result.WriteByte(data[i])
        }
    }
    return result.String()
}

func removeComments(data string) string {
    var result strings.Builder
    n := len(data)
    for i:=0; i<n; i++ {
        switch data[i] {
        case '/':
            if data[i+1] == '/' {

                for k:=i; k<n; k++ {

                    if data[k] == '\n' {
                        i = k-1
                        break
                    }
                }
            } 
        default:
            result.WriteByte(data[i])
        }
    }
    return result.String()
}

func handleComputation(data string) string {
    var c string
    switch data {
    case "0"  : c = "101010"; break
    case "1"  : c = "111111"; break
    case "-1" : c = "111010"; break
    case "D"  : c = "001100"; break
    case "A"  : c = "110000"; break
    case "!D" : c = "001101"; break
    case "!A" : c = "110001"; break
    case "-D" : c = "001111"; break
    case "-A" : c = "110011"; break
    case "D+1": c = "011111"; break
    case "A+1": c = "110111"; break
    case "D-1": c = "001110"; break
    case "A-1": c = "110010"; break
    case "D+A": c = "000010"; break
    case "D-A": c = "010011"; break
    case "A-D": c = "000111"; break
    case "D&A": c = "000000"; break
    case "D|A": c = "010101"; break

    case "M"  : c = "110000"; break
    case "!M" : c = "110001"; break
    case "-M" : c = "110011"; break
    case "M+1": c = "110111"; break
    case "M-1": c = "110010"; break
    case "D+M": c = "000010"; break
    case "D-M": c = "010011"; break
    case "M-D": c = "000111"; break
    case "D&M": c = "000000"; break
    case "D|M": c = "010101"; break

    case "A+D": c = "000010"; break
    case "M+D": c = "000010"; break
    }
    return c
}

func handleLhs(data string) string {
    var d string
    switch data {
    case "M":
        d = "001"
        break
    case "A":
        d = "100"
        break
    case "D":
        d = "010"
        break
    case "AM":
        d = "101"
        break
    case "MD":
        d = "011"
        break
    case "AD":
        d = "110"
        break
    case "AMD":
        d = "111"
        break
    default:
        d = "000"
        break
    }
    return d
}

func handleJump(data string) string {
    var d string
    switch data {
    case "JGT":
        d = "001"
        break
    case "JEQ":
        d = "010"
        break
    case "JGE":
        d = "011"
        break
    case "JLT":
        d = "100"
        break
    case "JNE":
        d = "101"
        break
    case "JLE":
        d = "110"
        break
    case "JMP":
        d = "111"
        break
    default:
        d = "000"
    }
    return d
}

func customRegisters(data string) string {

    var d string
    switch data {

    }
    return d
}

func alloc(data string) {

    d := starM + countM
    countM++
    _, ok := mem[d]

    if !ok {
        mem[d] = data
        rmem[data] = d
    } else {
        fmt.Println("not unique cannot alloc again")
    }
    fmt.Println(d)
    fmt.Println(mem)
}

func sanitizeSymbols(d []string) []string {

    data := d
    var dlr []string

    for i:=0; i<len(d); i++ {
        reg := regexp.MustCompile(`^@([^\dR]\w*)$`)
        bol := reg.Match([]byte(data[i]))
        if strings.HasPrefix(data[i], "@") && bol {
            alloc(data[i][1:])
            fmt.Println(data[i] + " matches " + strconv.FormatBool(bol))
        }
    }
    readSymbols(data)
    convertSymbols(data)
    convertAddress(data)
    for i:=0; i<len(data); i++ {
        if !strings.HasPrefix(data[i], "(") {

            dlr = append(dlr, data[i])
        }
    }
    return dlr
}

func readSymbols(data []string) {

    n := len(data)
    for i:=0; i<n; i++ {
        if strings.HasPrefix(data[i], "(") && strings.HasSuffix(data[i], ")") {
            k := len(data[i])
            jmp[data[i][1:k-1]] = i;
        }
    }
    fmt.Println(jmp)
    fmt.Println(data)
}

func convertSymbols(data []string) {
    for i:=0; i<len(data); i++ {
        if len(data[i]) > 1 {
            val, ok := jmp[data[i][1:]]
            if  ok {
                data[i] = "@" + strconv.Itoa(val);
            }
        }
    }
    fmt.Println("data mutation from converSymbols")
    fmt.Println(data)
    fmt.Println("data mutation still in converSymbols")
}

func convertAddress(data []string) {

    list := make([]string, 0)
    for k := range rmem {
        list = append(list, k)
    }
    fmt.Println(list)

    for i:=0; i<len(data); i++ {
        if len(data[i]) > 1 {
            val, _ := rmem[data[i][1:]]
            if strings.HasPrefix(data[i], "@") && slices.Contains(list, data[i][1:]) {
                data[i] = "@" + strconv.Itoa(val)
            }
        }
    }
}

