package vminfo

import (
	"fmt"
	"strings"
)

var VMS []VM  // an empty list


type VM struct {
	boshVmName string
	cpiID string
	ipAddress string
	iaasVmName string
	agentID string
	stemcellID string
	datastoreStemCell string
	datastoreVM string
	directorID string
	directorAndAgentIDFQN string
	pingTimeStamp string
	pongTimeStamp string
	netInfo string
	pingSent bool
	pongRec bool

	}


func addNewVM(line string) {
//adds a new vm to the array - by calling the extractNewVM method from bextractor.go
	VMS = append(VMS, assembleVM(line))
}

func adjustVMProperties(vm *VM, lineIn string ){
	if(strings.Contains(lineIn, "Storage Options for creating ephemeral disk")){
		line := trimAllUnneeded(lineIn)
		lindex := strings.LastIndex(line, ":")+1
		vm.datastoreVM = strings.TrimSpace(line[lindex:len(line)])
	}
	if(strings.Contains(lineIn, "Creating vm:")){
		line := trimAllUnneeded(lineIn)
		lindex := strings.LastIndex(line, ":")+1
		vm.iaasVmName = strings.TrimSpace(line[lindex:len(line)])
	}
	if(strings.Contains(lineIn, "Searching for stemcell")){
		line := trimAllUnneeded(lineIn)
		lindex := strings.LastIndex(line, " ")+1
		vm.datastoreStemCell = strings.TrimSpace(line[lindex:len(line)])
	}
}

func setDirectorAndFQNforVMS(directorIDin string) {
	//loop through all VMS in list and insert director ID
	for i:=0;i<len(VMS);i++{
		VMS[i].directorID = directorIDin
		VMS[i].directorAndAgentIDFQN = directorIDin+"."+VMS[i].agentID
	}
}

func updateVMPingInfo(pingTimeStamp string,directorAndAgentIDFQNin string) {
	for index, vm := range VMS{
		if strings.EqualFold(directorAndAgentIDFQNin,vm.directorAndAgentIDFQN) {
			VMS[index].pingSent = true
			VMS[index].pingTimeStamp = pingTimeStamp
		}
	}
}

func updateVMPongInfo(pongTimeStamp string,directorAndAgentIDFQNin string) {
	for index, vm := range VMS{
		if strings.EqualFold(directorAndAgentIDFQNin,vm.directorAndAgentIDFQN) {
			VMS[index].pongRec = true
			VMS[index].pongTimeStamp = pongTimeStamp
		}
	}
}

func PrintFileVerbose(){
	for index, _ := range VMS{
		fmt.Print(
			VMS[index].boshVmName,
			"\n\tIP address:\t\t",
			VMS[index].ipAddress,
			"\n\tIaaS ID:\t\t",
			VMS[index].iaasVmName,
			"\n\tVM Datastore:\t\t",
			VMS[index].datastoreVM,
			"\n\tFQ BoshAgent:\t\t",
			VMS[index].directorAndAgentIDFQN,
			"\n\tNetwork/AZ:\t\t",
			VMS[index].netInfo,
			"\n\tStemcell:\t\t",
			VMS[index].stemcellID,
			" from datastore ",
			VMS[index].datastoreStemCell,
			"\n\tFirst Ping:\t\t",
			VMS[index].pingTimeStamp,)

			if VMS[index].pongRec{
				fmt.Print(
				"\n\tPong Received:\t\t",
					VMS[index].pongTimeStamp,
					"\n\n")
				continue
			}

		fmt.Print(
			"\n\tPong Received:\t\t",
			VMS[index].pongRec,
			"\n\n")

	}
}