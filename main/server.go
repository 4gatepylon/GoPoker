package main


import "fmt"
import "bufio"
import "strings"
import "log"
import "flag"

// copy pasted
func request(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if i == 0 {
			m := strings.Fields(line)[0]
			fmt.Println("Methods", m)
		}
		if line == "" {
			break
		}
		i++
	}
}

// copy pasted
func response(conn net.Conn) {
	body := `<h1>This is Go Http Server using TCP</h1>`

	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}
// copy pasted
func handleConnection(conn net.Conn) {
	defer conn.Close()
	request(conn)
	response(conn)
}

func main() {
	var port string
	var serverName string

	flag.StringVar(&port, "port", "8888", "Specify the port to host from.")
	flag.StringVar(&port, "server", "Poker Server", "Specify the name of the server.")
	flag.Parse()

	// start the tcp server
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer listen.Close() // this is called when main() exits

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConnection(conn)
	}
}