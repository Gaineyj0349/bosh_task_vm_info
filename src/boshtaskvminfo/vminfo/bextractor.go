package vminfo

import (
	"strings"
)

func assembleVM(line string) VM{
	vmboshVmName,vmcpiID, vmipAddress, vmagentID, vmstemcellID, vmnetInfo := extractVMdata(line)

	return VM{
		boshVmName: vmboshVmName,
		cpiID: vmcpiID,
		ipAddress: vmipAddress,
		iaasVmName: "",
		agentID: vmagentID,
		stemcellID: vmstemcellID,
		datastoreStemCell: "",
		datastoreVM: "",
		directorID: "",
		directorAndAgentIDFQN: "",
		pingTimeStamp: "",
		pongTimeStamp: "",
		netInfo: vmnetInfo,
		pingSent: false,
		pongRec: false,
	}
}


func extractVMdata(linein string) (string, string, string, string, string, string) {
	//declare variables to send back
	var vmboshVmName, vmagentID, vmstemcellID, vmipAddress, vmNetInfo, vmcpiID string

	lines := strings.Split(linein," ")
	splitMore := false

	for _, element := range lines {
		if(strings.Contains(element, "create_missing_vm")){
			vmboshVmName = element[(strings.Index(element, "(") + 1):len(element)]
		}
		if strings.EqualFold(element,"request:"){
			splitMore = true
		}
		if splitMore{
			split := strings.Split(trimAllUnneeded(element),",")
			nextIsSC := false
			nextIsNI := false
			for _, subelement := range split{
				if(nextIsNI){
					vmNetInfo = subelement[0:strings.Index(subelement,":")]
					nextIsNI = false
				}
				if(nextIsSC){
					vmstemcellID = subelement
					nextIsSC = false
				}
				if(strings.Contains(subelement, "arguments:")){
					vmagentID = subelement[(strings.Index(subelement, ":") + 1):len(subelement)]
					nextIsSC = true
				}
				if(strings.Contains(subelement, "ip:")){
					vmipAddress = subelement[(strings.Index(subelement, ":") + 1):len(subelement)]
				}
				if(strings.Contains(subelement, "ram:")){
					nextIsNI = true
				}
				if(strings.Contains(subelement, "request_id:")){
					vmcpiID = subelement[(strings.Index(subelement, ":") + 1):len(subelement)]
				}
			}
		}
	}
	return vmboshVmName,vmcpiID,vmipAddress,vmagentID,vmstemcellID,vmNetInfo
}


func trimAllUnneeded(linein string) string{

	line := strings.Replace(linein,"\"", "",-1)
	line = strings.Replace(line,"[", "",-1)
	line = strings.Replace(line,"]", "",-1)
	line = strings.Replace(line,"{", "",-1)
	line = strings.Replace(line,"}", "",-1)
	line = strings.Replace(line,"(", "",-1)
	line = strings.Replace(line,")", "",-1)
	return line
}

func cleanAndSliceForDir(linein string) {
	var directorID string
	splitMore := true
	line := trimAllUnneeded(linein)
	lines := strings.Split(line," ")

	for _, element := range lines {
		if splitMore{
			split := strings.Split(element,",")
			for _, subelement := range split {
				if(strings.Contains(subelement,"reply_to")){
					colonIndex := strings.Index(subelement,":")+1
					lastIndexDot := strings.LastIndex(subelement[0:strings.LastIndex(subelement, ".")], ".")
					directorID = subelement[colonIndex:lastIndexDot]
				}
			}
		}
		if(strings.Contains(element,"agent.")){
			splitMore = true
		}
	}
	//update all VMS
	setDirectorAndFQNforVMS(directorID)

}


func isFirstPing(linein string) (bool,string,string){
	var directorFQN string
	var pingTimeStamp string
	splitMore := true
	line := trimAllUnneeded(linein)
	lines := strings.Split(line," ")

	for index, element := range lines {
		if index == 1{
			pingTimeStamp = element
		}
		if splitMore{
			split := strings.Split(element,",")
			for _, subelement := range split {
				if(strings.Contains(subelement,"reply_to")){
					colonIndex := strings.Index(subelement,":")+1
					lastIndexDot := strings.LastIndex(subelement, ".")
					directorFQN = subelement[colonIndex:lastIndexDot]
					if pingSent(directorFQN){
						return false, pingTimeStamp,directorFQN
					}
				}
			}
		}
		if(strings.Contains(element,"agent.")){
			splitMore = true
		}
	}
	return true, pingTimeStamp,directorFQN
}

func pingSent(directorAndAgentIDFQNin string) bool{
	for _, vm := range VMS{
		if strings.EqualFold(directorAndAgentIDFQNin,vm.directorAndAgentIDFQN) {
			return vm.pingSent
		}
	}
	return false
}

func getVMinfoFromPong(linein string) (string, string) {
	var directorFQN string
	var pongTimeStamp string
	splitMore := true
	line := trimAllUnneeded(linein)
	lines := strings.Split(line," ")

	for index, element := range lines {
		if index == 1{
			pongTimeStamp = element
		}
		if splitMore{
			split := strings.Split(element,",")
			for _, subelement := range split {
				if(strings.Contains(subelement,"director.")){
					colonIndex := strings.Index(subelement,":")+1
					lastIndexDot := strings.LastIndex(subelement, ".")
					directorFQN = subelement[colonIndex:lastIndexDot]
					return pongTimeStamp, directorFQN
				}
			}
		}
		if(strings.Contains(element,"agent.")){
			splitMore = true
		}
	}
	return "",""
}

