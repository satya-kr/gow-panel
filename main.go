package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

const port = ":3001"

var services = []string{"mysql", "mariadb", "nodejs", "php-fpm", "nginx"}
var processes = []string{"stop", "start", "status", "reload"}

type UpdateRequest struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type JsonResponse struct {
	Status  bool                   `json:"status"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}

func main() {
	//r := router.Router()
	r := mux.NewRouter()
	r.HandleFunc("/api/update-service", ApiUpdateService).Methods("POST")
	r.HandleFunc("/api/get-services", ApiGetServerStatusTypes).Methods("GET")

	fmt.Printf("Server is lisening on port%s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}
}

func ApiGetServerStatusTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := make(map[string]interface{})
	var strs []string

	for _, val := range services {
		status, err := checkServiceStatus(val)
		if err != nil {
			fmt.Println(err)
		}
		if status {
			strs = append(strs, fmt.Sprintf("%s:active", val))
		} else {
			strs = append(strs, fmt.Sprintf("%s:inactive", val))
		}
	}
	data["services"] = strs

	jsonRes := JsonResponse{
		Status:  true,
		Data:    data,
		Message: "Data loaded successfully!",
	}

	out, err := json.MarshalIndent(jsonRes, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func ApiUpdateService(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	data["services"] = []string{}
	var jsonRes JsonResponse
	var requestData UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	serviceName := strings.TrimSpace(requestData.Service)
	serviceStatus := strings.TrimSpace(requestData.Status)

	check, err := updateService(serviceName, serviceStatus)
	if err != nil {
		fmt.Println(err)
	}
	if check {
		jsonRes = JsonResponse{
			Status:  true,
			Data:    data,
			Message: fmt.Sprintf("The %s service status is updated to %s", requestData.Service, requestData.Status),
		}
	} else {
		jsonRes = JsonResponse{
			Status:  true,
			Data:    data,
			Message: fmt.Sprintf("Faild to update %s service as %s", requestData.Service, requestData.Status),
		}
	}

	out, err := json.MarshalIndent(jsonRes, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func updateService(serviceName, serviceStatus string) (bool, error) {
	// fmt.Sprintf("service:%s and status:%s", serviceName, serviceStatus)

	cmd := exec.Command("sudo", "systemctl", serviceStatus, serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
		return false, err
	}
	return true, nil
}

func checkServiceStatus(serviceName string) (bool, error) {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return strings.TrimSpace(string(output)) == "active", nil
}

//func GetMysqlPID() (string, error) {
//	command := "ps aux | grep mysql"
//	cmd := exec.Command("bash", "-c", command)
//	output, err := cmd.CombinedOutput()
//	if err != nil {
//		fmt.Println("Error:", err)
//		return "", err
//	}
//
//	// Convert the output byte slice to a string
//	outputStr := string(output)
//
//	// Split the output into lines
//	lines := strings.Split(outputStr, "\n")
//
//	// Find the line containing "mysql" and extract the PID
//	var mysqlPID string
//	for _, line := range lines {
//		if strings.Contains(line, "mysql") {
//			fields := strings.Fields(line)
//			if len(fields) >= 2 {
//				mysqlPID = fields[1]
//				break
//			}
//		}
//	}
//
//	if mysqlPID == "" {
//		return "", errors.New("mysql PID not found")
//	} else {
//		return mysqlPID, nil
//	}
//}
