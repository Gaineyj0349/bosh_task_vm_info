package vminfo

import (
	"strings"
)

func Check(line string){
	if checkForFirstMention(line){
		addNewVM(line)
	}
	adjust,index := checkForVMCreationActivity(line)
	if adjust{
		adjustVMProperties(&VMS[index],line)
	}
	checkForDirAndPing(line)
	checkForPong(line)

}



//this is the first mention of the VM information
func checkForFirstMention(s string) bool{
	if strings.Contains(s,"\"method\":\"create_vm\"") {
		return true
	}
	return false

}

//this is the information about the creation of VM
func checkForVMCreationActivity(s string) (bool,int) {
	//loop through all VMS in list to compare the line in to the cpi of VM	
	for index, VM := range VMS{
		cpiString := "req_id "+VM.cpiID
		if(strings.Contains(s,cpiString)){
			return true, index
		}
	}	
	return false,0
}

//this checks if it is a ping, and if so sets director for all vms. only records first ping per vm
func checkForDirAndPing(s string) {

	if strings.Contains(s,"SENT: agent"){
		// if director isn't set for this vm, this is first ping and we can update all VMs from info here
		if len(VMS[0].directorID) == 0{
			cleanAndSliceForDir(s)
		}
		//is this the first ping for the specific VM in the log?
		first, pingTimeStamp, directorAndAgentIDFQN := isFirstPing(s)
		if first{
			updateVMPingInfo(pingTimeStamp,directorAndAgentIDFQN)
		}
	}
}

//this updates the respected VM upon the pong
func checkForPong(s string) {
	if strings.Contains(s, "pong"){
		pongTimeStamp, directorAndAgentIDFQN := getVMinfoFromPong(s)
		updateVMPongInfo(pongTimeStamp,directorAndAgentIDFQN)
	}
}



