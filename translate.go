package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	// "path/filepath"

	"github.com/bregydoc/gtranslate"
)

var file *os.File
var startTime time.Time

// CmdFlag handles command line arguments
type CmdFlag struct {
	lang    *string
	inFile  *string
	outFile *string
	delim   *string
	addUrl  *bool
	url     *string
	col     *uint
	start   *uint
}

func openFile(path string) (*bufio.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return bufio.NewReader(nil), err
	}
	// defer file.Close()
	return bufio.NewReader(file), nil
}

func translate(word string, lang []string) string {
	translated, err := gtranslate.TranslateWithParams(
		word,
		gtranslate.TranslationParams{
			From: lang[0], //"it",
			To:   lang[1], //"ru",
		},
	)
	if err != nil {
		panic(err)
	}
	return translated
}

// type QuiteFlagSet struct {
// 	*flag.FlagSet
//   }

//   func (f *QuiteFlagSet) failf(format string, a ...interface{}) error {
// 	// err := fmt.Errorf(format, a...)
// 	// fmt.Fprintln(os.Stderr, err)
// 	// return err
// 	return nil
//   }

func printLink(url string, text string) (out string) {
	out = fmt.Sprintf("\x1b]8;;%+v\a%s\x1b]8;;\a", url, text)
	// echo -e '\e]8;;http://example.com\aThis is a link\e]8;;\a'
	return out
}

func main() {
	// fmt.Printf("\nhere\n")
	startTime = time.Now()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	go func() {
		for {
			<-killSignal
			fmt.Println("\x1b[1A\x1b[K\n..closing...")
			exit()
		}
	}()

	// contMode := flag.NewFlagSet("start", flag.ContinueOnError)
	fileMode := flag.NewFlagSet("file", flag.ExitOnError)
	mainMode := flag.NewFlagSet("", flag.ContinueOnError)

	var cmdFlag CmdFlag
	lng := mainMode.String("lang", "it-ru", "language translate")
	inv := mainMode.Bool("_", false, "invert language` `\n")

	cmdFlag.lang = fileMode.String("lang", "it-ru", "language translate")
	cmdFlag.inFile = fileMode.String("in", "", "input file path")
	cmdFlag.outFile = fileMode.String("out", "", "output file path")
	cmdFlag.delim = fileMode.String("delim", ";", "delimeter symbol (one byte)")
	cmdFlag.addUrl = fileMode.Bool("addlink", true, "option to add link into the next link")
	cmdFlag.url = fileMode.String("url", "https://context.reverso.net/traduzione/italiano-russo/", "url to add")
	cmdFlag.col = fileMode.Uint("col", 4, "column number in file")
	cmdFlag.start = fileMode.Uint("start", 0, "number of line to start with")

	// flag.ErrorHandling = 0
	// mainMode.Parse(os.Args[1:])
	if err := mainMode.Parse(os.Args[1:]); err != nil {
		// mainMode.PrintDefaults()
		fmt.Println("\nUsage of subcommand file:")
		fileMode.PrintDefaults()
		os.Exit(0)
	}

	/*
		// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(pwd, os.Args[0],"\n\n")

		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)//*/

	if len(os.Args) > 1 && os.Args[1] == "file" {
		fileMode.Parse(os.Args[2:])

		lang := strings.Split(*cmdFlag.lang, "-")
		delim := []byte(*cmdFlag.delim)
		col := *cmdFlag.col

		if len(*cmdFlag.inFile) < 1 {
			args := fileMode.Args()
			if len(args) >= 2 {
				*cmdFlag.inFile = args[0]
				*cmdFlag.outFile = args[1]
			} else {
				fmt.Println("Path not recognized")
				os.Exit(2)
			}
		}

		fmt.Printf("lang:  \t%#v\ncol:  \t%v\ninFile:  %#v\noutFile: %#v\ndelim:  %#v\naddUrl:  %#v\nurl:  \t%#v\nstart:  %v\n\n",
			*cmdFlag.lang,
			*cmdFlag.col,
			*cmdFlag.inFile,
			*cmdFlag.outFile,
			*cmdFlag.delim,
			*cmdFlag.addUrl,
			*cmdFlag.url,
			*cmdFlag.start)

		reader, err1 := openFile(*cmdFlag.inFile)
		if err1 != nil {
			fmt.Printf("Open file error: %s\n", err1)
		}

		outfile, err := os.Create(*cmdFlag.outFile)
		if err != nil {
			fmt.Printf("File create error: %s \n", err)
		}

		defer outfile.Close()

		rowNum := uint(0)
		for rowNum < *cmdFlag.start {
			rowNum++
			reader.ReadLine()
		}
		for {
			rowNum++
			wrd := ""
			line := ""
			for i := uint(0); i < col; i++ {
				n2, err := reader.ReadString(delim[0])
				if err != nil {
					if err.Error() != "EOF" {
						fmt.Printf("Read string error: %s %s\n", n2, err)
					}
					break
				}
				line += n2
				wrd = n2
			}
			n1, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			wrd = strings.TrimSuffix(wrd, string(delim[0]))
			fmt.Printf("%v:\t%s", rowNum, wrd)

			translated := translate(wrd, lang)

			fmt.Printf("\t%s: %s \t", lang[1], translated)

			line += translated + ";"
			if *cmdFlag.addUrl {
				txt := strings.ReplaceAll(wrd, " ", "+")
				url := *cmdFlag.url + txt
				fmt.Printf("%+v\n", printLink(url, "visit link"))
				line += url + ";"
			}
			line += string(n1)

			_, err2 := outfile.WriteString(line + "\n")

			if err2 != nil {
				fmt.Printf("Write string error: %s \n", err2)
			}
		}

		defer file.Close()

		exit()

	} else {

		// mainMode.Parse(os.Args[1:])

		lang := strings.Split(*lng, "-")
		if *inv {
			lang = []string{lang[1], lang[0]}
		}
		if len(mainMode.Args()) != 0 {
			wrd := strings.Join(mainMode.Args(), " ")
			translated := translate(wrd, lang)
			fmt.Printf("%v: %#v -> %v: %#v\n", lang[0], wrd, lang[1], translated)
		} else {
			reader := bufio.NewReader(os.Stdin)
			info, _ := os.Stdin.Stat()
			// go func() {
			// fmt.Print("Enter text: ")
			user := false
			if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
				user = true
			}
			// fmt.Println(user)
			fmt.Printf("%v: ", lang[0])
			for {
				wrd, _, ok := reader.ReadLine()
				if ok == nil {
					if len(wrd) != 0 {
						translated := translate(string(wrd), lang)
						url := *cmdFlag.url + strings.ReplaceAll(string(wrd), " ", "+")
						if !user {
							fmt.Printf("%v\n", string(wrd)) //\x1b[1A\x1b[K
						}
						fmt.Printf("%v: %+v\n", lang[1], printLink(url, translated))
						if user {
							fmt.Println()
						}
					} else if user {
						lang = []string{lang[1], lang[0]}
					}
					fmt.Printf("%v: ", lang[0])
				} else {
					if ok.Error() == "EOF" {
						exit()
					} else {
						panic(ok)
					}
				}
			}
			// }()
			exit()
		}

		exit()
	}

	os.Exit(0)
}

func exit() {
	elapsed := time.Since(startTime)
	fmt.Printf("\nTranslation took %s\n", elapsed)
	os.Exit(0)
}
