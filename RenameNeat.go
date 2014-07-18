// ** The next iteration of this program should move the resulting file into the file structure that already exists.
//
// Written by:  Bob Failla
//
// Reads files in a directory provided by the user
// Parses the files for Name, Amount and Date
// Renames the files according to YYYYMMDD ($amount) Name
// ** Places the files in a directory supplied by the user
// Leaves the files that could not be renamed in the source directory
// Prints a message if a file could not be moved
// Either prints what files were moved or logs an error
//


package main
////////////////////////////// Packages
import (
	"fmt"
	"io/ioutil"
	"regexp"
	"os"
	"strings"
	"log"
)
/////////////////////////////// Global Variables
var source_dir, destination_dir, amount, date, name string
var dr_list []os.FileInfo
//////////////////////////////////////////
func get_input() (string, string, []os.FileInfo) {

	// Get the name of the source directory
	for !input_is_valid(source_dir) {
		fmt.Printf("Source Directory:  ")
		fmt.Scanf("%s", &source_dir)
		fmt.Println()
		if !input_is_valid(source_dir) {
			fmt.Printf("Source directory is invalid.  A-Z, a-z, 0-9, /-_ only")
		}
	}
	// Get the name of the destination directory
	for !input_is_valid(destination_dir) {
		fmt.Printf("Destination Directory: ")
		fmt.Scanf("%s", &destination_dir)
		if !input_is_valid(destination_dir) {
			fmt.Printf("Destination directory is invalid.  A-Z, a-z, 0-9, /-_ only")
		}
	}
	// Read the destination directory.  Exit with error code 1 if unable.
	dr_list, err := ioutil.ReadDir(source_dir)
	if err != nil {
		fmt.Println("Could not read the source directory")
		os.Exit(1)
	}
	return source_dir, destination_dir, dr_list
}
/////////////////////////////////////////////////////////////////
func input_is_valid(d string) (bool) {
	var d_regex = regexp.MustCompile(`[a-zA-Z0-9/_-]`)
		answer := d_regex.MatchString(d)
	return answer
}
////////////////////////////////////////////////////////////////
func get_file_data(v os.FileInfo) (string, string, string) {
	var br bool
	amount_re 	:= regexp.MustCompile(`\d+\.\d+`)
	date_re 	:= regexp.MustCompile(`\d{8}`)
	name_re 	:= regexp.MustCompile(`[A-Za-z_ ]+`)
	// Use the regex defined above to find name, amount, date from the filename. (Move to here - these variables don't need global scope)
	amount 		:= amount_re.FindString(v.Name())
	date 		:= date_re.FindString(v.Name())
	names 		:= name_re.FindAllString(v.Name(), 3)
// Iterate through name, which could contain as many as 3 matches to find the correct name.  Consider including "cash and visa, as these provide no value.
	name :=""
	br = false
	for x:=len(names)-1; x>=0 && !br; x-- {
		if names[x] != "pdf" && names[x] != "Receipt" && names[x] !=" "{
			name = names[x]
			br = true
		}
	}
	return amount, date, name
}
////////////////////////////////////////////////////////////////
func rename_files () {

	var  amount, date, name, source_file, destination_file string
	var v os.FileInfo
	R:= strings.Replace

//iterate over the file names, call the parsing function and if possible rewrite the file names
	for _, v = range (dr_list) {
		amount, date, name = get_file_data(v)
		source_file = source_dir+"/"+v.Name()

// Keep calm and realize this is just a nested strings.replace function (shortened by the variable declaration for
// R above) replacing ". " with ".", "$" with "\$" and " " with "\ "
// The " ." gets replaced first to make sure it doesn't get escaped and become permanent.

		if date != "" {
			year := date[0:4]
			month := date[5:6]
			if _, err := os.Stat(destination_dir + "/" + year + "/" + month); err != nil {
				os.MkdirAll(destination_dir + "/" + year + "/" + month, 0777)
			}
			switch {
				case (amount != "") && (name != "") && (name != " "): {
				// missing nothing - **********************************								OK
				d := destination_dir+"/"+year+"/"+month+"/"+date+" ($"+amount+") "+name+".pdf"
//				destination_file = R(R(R(R(R(R(d,`(`,`\(`,-1),`)`,`\)`,-1),"  "," ",-1)," .",".",-1),"$", `\$`, -1)," ",`\ `, -1)
				destination_file = R(R(d,"  "," ",-1)," .",".",-1)

				move_the_file(source_file, destination_file)
				}
				case (amount == "") && (name != ""): {
				// missing amount only
				d := destination_dir+"/"+year+"/"+month+"/"+date+" ($"+"NEED AMOUNT"+") "+name+".pdf"
//					destination_file = R(R(R(R(R(R(d,`(`,`\(`,-1),`)`,`\)`,-1),"  "," ",-1)," .",".",-1),"$", `\$`, -1)," ",`\ `, -1)
					destination_file = R(R(d,"  "," ",-1)," .",".",-1)

				move_the_file(source_file, destination_file)
				}
				default : {
				// missing too much to be useful - do nothing
				}
			}
		}
	}
}
////////////////////////////////////////////////////////////////
func move_the_file (source_file, destination_file string) {

	// Check to see if there is a file by that name.  If not, move it.  If so, write error message.

	if _, err := os.Stat(destination_file); err != nil {
		error := os.Rename(source_file, destination_file)
			if error != nil {
			log.Print(error)
			} else {
				fmt.Printf("The file: %s was moved to: %s\n", source_file, destination_file)
			}
	} else {
		log.Print(err)
		fmt.Printf("ERROR - Unable to write the file %s.  Check to see if it exists already.\n", destination_file)
	}
}
////////////////////////////////////////////////////////////////////////////
func main () {
	source_dir, destination_dir, dr_list = get_input()
	rename_files()

}
