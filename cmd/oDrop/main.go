package main

import (
	"compress/gzip"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net"
	"oDrop/core"
	"oDrop/utils"
	"os"
	"strconv"
)

func main() {
	var mode string
	var file string
	var n int
	mode, file, n = PromptUser()

	if utils.ModeToSimple(mode) == "s" {

		// generates a random number to verify receiver
		r := utils.GetRandomNumber()
		fmt.Printf("Passcode is %s\n", r)

		// broadcasts the file
		fmt.Println("Waiting for connections")
		err := core.Send(core.SendDataCallback{
			SentCallback: func(c net.Conn) {
				fmt.Printf("Sent file to %v", c.RemoteAddr())
				os.Exit(0)
			},
			DataBroker: func(c net.Conn, reader io.Reader, size int64) {
				gzr, err := gzip.NewReader(reader)
				if err != nil {
					log.Fatal(err)
				}

				bar := progressbar.DefaultBytes(size, "Sending")
				bar.Describe("Sending")
				io.Copy(io.MultiWriter(c, bar), gzr)
			},
		}, file, r)
		if err != nil {
			log.Fatalf("cant send file %v", err)
		}
	} else {
		err := core.Receive(file, strconv.Itoa(n), func(d io.Reader, f io.Writer) {
			bar := progressbar.DefaultBytes(
				-1,
				"downloading",
			)
			// copy contents of data to the file
			wb, err := io.Copy(io.MultiWriter(f, bar), d)

			if err != nil {
				log.Fatalln(err)
			}
			if wb == 0 {
				fmt.Printf("got %d bytes passcode might be wrong", wb)
			} else {
				fmt.Printf("%d B written in %s", wb, file)
			}
		})
		if err != nil {
			log.Fatalf("cant receive file %v", err)
		}
	}

}

// this function return the mode and filename if the mode is send the filename is the name of the file to send
// else it is the name of the file to save as
func PromptUser() (string, string, int) {
	var (
		mode   string
		file   string
		number int
	)

	for {
		fmt.Print("Do you want to send/receive: ")
		_, _ = fmt.Scanln(&mode)
		mode = utils.RemoveWhitespace(mode)

		if mode == "send" || mode == "receive" || mode == "r" || mode == "s" {
			break
		} else if mode == "exit" {
			os.Exit(0)
		}
		fmt.Print("\n")

		fmt.Println("wrong input")
	}
	if utils.ModeToSimple(mode) == "s" {
		for {
			fmt.Print("enter the location of your file: ")
			_, _ = fmt.Scanln(&file)
			file = utils.RemoveWhitespace(file)

			if file != "" && utils.DoesFileExist(file) {
				break
			}
			fmt.Print("\n")
			fmt.Println("file doesnt exit")
		}
	} else {
		for {
			fmt.Print("enter the location to save file: ")
			_, _ = fmt.Scanln(&file)
			file = utils.RemoveWhitespace(file)

			if file != "" && utils.DoesFileExist(file) == false {
				break
			} else if file == "exit" {
				os.Exit(0)
			}
			fmt.Print("\n")

			fmt.Println("file exists")
		}
		for {
			fmt.Print("enter the pass code: ")
			_, err := fmt.Scanln(&number)
			fmt.Print("\n")
			if err != nil {
				log.Fatalln(err)
			} else {
				break
			}
		}
	}
	return mode, file, number
}
