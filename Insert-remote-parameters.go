package main
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Written by Bob Failla
//
// Description:
//		- Inputs base path from the user
//		- Finds agent.properties files by recursively transcending the base path
//		- Backs up the agent.properties files to a timestamped (RFC3339) backup file
//		- Searches agent.properties line by line
// 		- Does not change existing agent.properties file.  Instead writes to agent.properties.new
// 		- Removes any existing remote properties parameters, except organization parameter
//		- attains computer host name ***
//		- appends remote management parameters to the end of agent.properties.new
//		- increments the port accumulator for the next connector
//		- continues to the next connector until finished updating all connectors
//		- timestamps and logs all actions to a log file dropped into the base directory
//
// Todo:
//		0) Port to Windows - / may need to change to \.
//		1)  Include full safeties so user can run this code without making changes, only logging what would have happened.
//		This allows user to test on each system before they execute
//		2)  Test out the package by Mitchell Hashimoto for accessing the Windows process list
//		This might allow the update of only running connectors.  Currently, all connectors are updated
//		3)  Log the highest port upgraded.
//		4)  Add Windows firewall rule update support
//		5)  Add support to restart ArcSight connectors
//		6)  If windows services cannot be manipulated, replace agent.properties with agent.new and prompt user to restart
//		these services.  If windows services can be manipulated, move the files and trigger the restart
//		7)  Alter non-logging panic conditions to allow continuation to the next connector
//		8)  Better modularize backup and edit function
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
import (
	"os"
	"fmt"
//	"path/filepath"
	"regexp"
	"bufio"
	"io"
	"time"
	"log"
	"strconv"
)
///////////////////////////////////////////////////////////////////////
type fileattr struct {
	fileinfo os.FileInfo
   	path     string
}
///////////////////////////////////////////////////////////////////////
var (
	base_dir string
	port int
)
///////////////////////////////////////////////////////////////////////
func backup_and_edit(path string, f os.FileInfo, _ error) error {
//////////////////////////////////////////////////////////////////////////find directories with agent.properties files
var (
	attr fileattr
)
	attr.fileinfo = f
	attr.path = path
	if attr.fileinfo.IsDir() {
		attr.path += "/user/agent/agent.properties"
		if _, err := os.Stat(attr.path); os.IsNotExist(err) {
			// Function written in the negative - sorry os.IsExist(err) doesn't give same result
		} else {
///////////////////////////////////////////////////////////////////////////////////////backup the agent.properties file
			t := time.Now()
			log.Println(t.Format(time.RFC3339) + " Success: - Found File " + attr.path)
			file, err := os.Open(attr.path)
			if err != nil {
				log.Println(t.Format(time.RFC3339) + " Fatal: Failed to open " + attr.path + " for backup")
				panic(err)
			}
			backup, err := os.Create(attr.path + "." + t.Format(time.RFC3339) + ".bak")
			if err != nil {
				log.Println(t.Format(time.RFC3339) + " Fatal: Failed to create backup destination")
				panic(err)
			}
			io.Copy(backup, file)
			if err != nil {
				log.Println(t.Format(time.RFC3339) + "Fatal:  Failed to make the backup copies")
				panic(err)
			} else {
				backup.Close()
				file.Close()
				log.Println(t.Format(time.RFC3339) + " Success:  successful backup" + attr.path)
			}
/////////////////////////////////////////////////////////////////////////////////////////// Set file and port variables
			file, err = os.Open(attr.path)
			input_file := bufio.NewScanner(file)
			output_file, _ := os.Create(attr.path+ ".new")
////////////////////////// ****************** need error handling for opening the file for writing *****************
			for input_file.Scan() {
				line := input_file.Text()
				remote_management_match, _ := regexp.MatchString("^Remote.Management", line)
				remote_user_match, _ := regexp.MatchString("^Remote.User", line)
				remote_management_organization_match, _ := regexp.MatchString("^Remote.Management.Organization", line)
				remote_match, _ := regexp.MatchString("^Remote", line)
				switch {
					case remote_management_organization_match : {output_file.WriteString(line + "\n") }
					case remote_user_match: {break} //Do not write the line - effective delete
					case remote_management_match : {break}  // Do not write the line - effective delete
					case remote_match: {break}  //Do not write the line - effective delete
					default: output_file.WriteString(line + "\n")  // Copy the line as is
				}
			}
			host_name, _ := os.Hostname()
			output_file.WriteString("Remote.Management.Enabled=True\n")
			host_param := ("Remote.Management.Host=" + host_name)
			output_file.WriteString(host_param + "\n")
			port_param := ("Remote.Management.Port=" + strconv.Itoa(port) + "\n")
			output_file.WriteString((port_param))
			port += 1
			output_file.WriteString("Remote.User=lemon\n")
			log.Println(t.Format(time.RFC3339) + " Success:  Remote parameters added")
		} // End if agent.properties file

	} // End if is directory
	return nil
}  // End backup_and_edit function

func main () {
///////////////////////////////////////////////////////////////////////////////////////////// set up continuous logging
	logfile := "/Users/Bob/Desktop/ArcSight/SmartConnectors/remote_parameter.log"
	logger, err := os.OpenFile(logfile, os.O_RDWR |os.O_CREATE |os.O_APPEND, 0666)
	if err != nil {
		panic(err)
		fmt.Println("unable to open log")
	}
	log.SetOutput(logger)
//////////////////////////////////////////////////////////////////////////////////////////////// get the base directory
	base_dir = "/Users/Bob/Desktop/ArcSight/SmartConnectors"
	fmt.Println("Please enter the base directory for the ArcSight Connectors:")
	fmt.Scanf("%s\n", &base_dir)
	log.Println("Base Directory defined as " + base_dir)
///////////////////////////////////////////////////////// Input the root directory to allow for change during execution
///////////////////////////////////////////////////////////////////////////////////////////////// get the starting port
	fmt.Printf("Please enter the starting port number for this host:")
	fmt.Scanf("%d\n", &port)
	fmt.Println()
	log.Printf("Starting Port number defined by user as: %d\n ", port)
/////////////////////////////////////////////////////////////////////////////// walk the base directory and do the work
	filepath.Walk(base_dir, backup_and_edit)
} // Program end
