package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"os"
)
var amount_re = regexp.MustCompile(`\d+\.\d+`)
var date_re = regexp.MustCompile(`\s\d+\s`)
var name_re = regexp.MustCompile(`[A-Za-z _]+`)
var dr_list []os.FileInfo
var destination, source, amount, date, name, s, d, sf, df string
var v os.FileInfo
var fm int
/////////////////////////////////////////////////////////////////
func get_input() (string, string) {
	// Get the name of the source directory
	fmt.Printf("Source Directory:  ")
	fmt.Scanf("%s", &source)
	fmt.Println()
	// Get the name of the destination directory
	fmt.Printf("Destination Directory: ")
	fmt.Scanf("%s", &destination)
	return source, destination
}
/////////////////////////////////////////////////////////////////
func input_is_valid(d string) (bool){
	var d_regex = regexp.MustCompile(`[a-zA-Z0-9/_-]`)
	answer:= d_regex.MatchString(d)
	return answer
}
////////////////////////////////////////////////////////////////
func get_file_data(v os.FileInfo) (string, string, string) {
	amount = amount_re.FindString(v.Name())
	date = date_re.FindString(v.Name())
	var names = name_re.FindAllString(v.Name(), 2)
	name = names[len(names)-1]
	return amount, date, name
}
////////////////////////////////////////////////////////////////
func move_the_file (source, sf, destintion, df string) (n int) {
/// if the destination directory does not exist, make one
	 err := os.MkdirAll(destination , 0744)
	if err != nil {
			fmt.Printf("Unable to create the directory \n")
			os.Exit(1)
	}
	os.Rename(source+sf, destination+df)
	if err == nil {
		n = 1
	} else {
		fmt.Printf("Unable to write the file. %s", source+sf)
		n=0
	}
/// if error on writing a file (non-duplicate), write error message and move on
/// next step is to file these into their proper directories
	return n
}
////////////////////////////////////////////////////////////////
func main() {
	get_input()
// Validate source and destination input.  Exit with error code 1 if not valid.
	if !input_is_valid(source) {
		fmt.Printf("Source direcgory is invalid.  A-Z, a-z, 0-9, /-_ only")
		os.Exit(1)
	}
	if !input_is_valid(destination) {
		fmt.Printf("Destination directory is invalid.  A-Z, a-z, 0-9, /-_ only")
		os.Exit(1)
	}
// Read the directory.  Exit with error code 1 if unable.
	dr_list, err := ioutil.ReadDir(source)
	if err != nil {
		fmt.Println("Could not read the source directory")
		os.Exit(1)
	}
//iterate over the file names, call the parsing function and if possible rewrite the file names
	for _, v = range (dr_list) {
		amount, date, name = get_file_data(v)
		switch {
		case (amount != "") && (date != "") && (name != ""): {
			sf := "/" + v.Name()
			df := "/" + date + "($" + amount + ") " + name + ".pdf"
			fm = fm+move_the_file(source, sf, destination, df)
		}
		case (amount != "") && (date != "") && (name == ""): {
			// missing name only
		}
		case (amount != "") && (date == "") && (name != ""): {
			// missing date only
		}
		case (amount != "") && (date == "") && (name == ""): {
			// missing date and name
		}
		case (amount == "") && (date != "") && (name != ""): {
			// missing everything
		}
		case (amount == "") && (date != "") && (name == ""): {
			//missing amount and name
		}
		case (amount == "") && (date == "") && (name != ""): {
			// missing amount and date
		}
		case (amount == "") && (date != "") && (name != ""): {
			// missing amount only
		}


			// A	D	N
			// 0	0	0 - Everything missing 				- done
			// 0	0	1 - Missing amount and date			- done
			// 0	1	0 = Missing  amount and name		- done
			// 0	1	1 = Missing amount only				- done
			// 1	0	0 = Missing date and name			- done
			// 1	0	1 = Missing date only				- done
			// 1	1	0 = Missing name only				- done
			// 1	1	1 = Nothing missing					- done



			//		if (amount != "") && (date != "") && (name != "") {
			//			sf := "/" + v.Name()
			//			df := "/" + date + "($" + amount + ") " + name + ".pdf"
			//			fm = fm + move_the_file(source, sf, destination, df)
			//		}
		}
	}
	fmt.Printf("%d" +  " files renamed and moved", fm)
}


