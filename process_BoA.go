package main

import (
	"os"
	"fmt"
	"log"
	"regexp"
)
//////////////////////////// Declarations /////////////////////////////
var source_path, destination_path string
var source_row[] string
/////////////////////////// Get Source Function ///////////////////////
func get_source () (string){
	var file string
	for !input_check_path(file) {
		fmt.Printf("Source File (include path from root directory):  ")
		fmt.Scanf("%s", &file)
		fmt.Println()
		if !input_is_valid(source_dir) {
			fmt.Printf("Source directory is invalid.  A-Z, a-z, 0-9, /-_ only")
		}
	}
	return(file)
}
////////////////////////// Get Destination Function ///////////////////
func get_destination () (string){
	var file string
	for !input_check_path(file) {
		fmt.Printf("Destination File (include path from root directory):  ")
		fmt.Scanf("%s", &file)
		fmt.Println()
		if !input_is_valid(file) {
			fmt.Printf("Source directory is invalid.  A-Z, a-z, 0-9, /-_ only")
		}
	}
	return(file)
}
///////////////////////// Check the input path for valid characters ////
func input_check_path (absolute_path string) (bool) {
	var d_regex = regexp.MustCompile(`[a-zA-Z0-9/_-]`)
	answer := d_regex.MatchString(absolute_path)
	return answer
}
//////////////////////// Verify the Source file exists ////////////////
func verify_source_path (absolute_path string) {
	if _, err := os.Stat(absolute_path); err != nil {
		log.Print("No such file ", absolute_path)
		os.Exit(1)
	}
}
//////////////////////// Read source file into data structure //////////
func read_row (absolute_path string) {
}
/////////////////////// Rewrite the name into standard format //////////
func process_names () {
	// takes and row from the CSV file as input does the translation function and calls write_row for outputting results
}
/////////////////////// Write a line of CSV to the output file /////////
func write_row ( file) {
}

///////////////////////       Main      ////////////////////////////////
func main() {
	source_path = get_source()
	destination_path = get_destination ()
	verify_source_path(source_dir)
//while file end hasn't been reached //
	source_row = read_row()
	process_names (source_row)
//end while //
}
