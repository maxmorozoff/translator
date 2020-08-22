# translator

## Installation:

```
$ make install
```
To rebuild project use:
```
$ make build
```

## Usage:

```
$ translate -help
Usage:
  -_  
    	invert language 
    	
  -lang string
    	language translate (default "it-ru")

Usage of subcommand file:
  -addlink
    	option to add link into the next link (default true)
  -col uint
    	column number in file (default 4)
  -delim string
    	delimeter symbol (one byte) (default ";")
  -in string
    	input file path
  -lang string
    	language translate (default "it-ru")
  -out string
    	output file path
  -start uint
    	number of line to start with
  -url string
    	url to add (default "https://context.reverso.net/traduzione/italiano-russo/")
```

## Examples:
```
$ translate
  it: // Enter a word here
      // To switch language in this mode just hit enter
$ translate -lang en-ru 
  en: // To modify language use flag -lang
```
```
$ translate ciao
  it: "ciao" -> ru: "Привет"

  Translation took 486.391326ms
// To invert language use -_
$ translate -_ привет
  ru: "привет" -> it: "ciao"
  
  Translation took 386.214161ms

$ translate -_ Привет
  ru: "Привет" -> it: "Ciao"

  Translation took 386.214161ms
```
```
$ cat text.txt | translate > output.txt
```
```
$ translate file text.txt output.txt -col 3 -delim ; -lang it-ru
```