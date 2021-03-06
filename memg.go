#!/usr/bin/env gorun
package main;

import (
    "fmt"
    "net"
    "bufio"
    "strings"
    "os"
    "strconv"
    //"runtime/pprof"   // Uncomment to profile
)

var CACHE map[string] string;

func main() {

    // Uncomment these three lines to profile
    /*
    handle, _ := os.Create("memg.prof")
    pprof.StartCPUProfile(handle)
    defer pprof.StopCPUProfile()
    */

    listener, err := net.Listen("tcp", "127.0.0.1:11211")
    if err != nil {
        panic("Error listening on 11211: " + err.String())
    }

    CACHE = make(map[string] string)

    if isSingle() {
        netconn, err := listener.Accept()
        if err != nil {
            panic("Accept error: " + err.String())
        }

        handleConn(netconn)

    } else {
        for {
            netconn, err := listener.Accept()
            if err != nil {
                panic("Accept error: " + err.String())
            }

            go handleConn(netconn)
        }
    }

}

func isSingle() bool {
    for _, arg := range os.Args  {
        if arg == "--single" {
            return true
        }
    }
    return false
}

/*
 * Networking
 */
func handleConn(conn net.Conn) {

    reader := bufio.NewReader(conn)
    for {

        // Fetch

        content, err := reader.ReadString('\n')
        if err == os.EOF {
            break
        } else if err != nil {
            fmt.Println(err)
            return
        }

        content = content[:len(content) - 2]    // Chop \r\n

        // Handle

        parts := strings.Split(content, " ")
        cmd := parts[0]
        switch cmd {

        case "get":
            key := parts[1]
            val, ok := CACHE[key]
            if ok {
                length := strconv.Itoa(len(val))
                conn.Write([]uint8("VALUE " + key + " 0 " + length + "\r\n"))
                conn.Write([]uint8(val + "\r\n"))
            }
            conn.Write([]uint8("END\r\n"))

        case "set":
            key := parts[1]
            //exp := parts[2]
            //flags := parts[3]
            length, _ := strconv.Atoi(parts[4])
            // Really we should read exactly 'length' bytes + \r\n
            val := make([]byte, length);
            reader.Read(val)
            CACHE[key] = string(val);
            conn.Write([]uint8("STORED\r\n"))
        }
    }
}
