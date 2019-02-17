package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

var backUrl = "http://typhoon-go"

// Tests the healthCheck endpoint
func TestHealthCheck(t *testing.T) {
	if _, err := http.Get(backUrl); err != nil {
		t.Errorf("TestHealthCheck error: %s", err.Error())
	}
	t.Logf("TestHealthCheck ok")
}

// Get a token for user 9999test
func getToken(t *testing.T, client *http.Client) string {
	resp, err := client.Get(backUrl + "/token/9999test")
	if err != nil {
		t.Errorf("TestGetToken error: %s", err.Error())
	}
	rbody := resp.Body
	var respBody bytes.Buffer
	if rbody != nil {
		_, err = io.Copy(&respBody, rbody)
		if err != nil {
			t.Errorf("TestGetToken: %v", err)
		}
		rbody.Close()
	}
	return respBody.String()
}

// do request and log errors if needed
func doRequestAndLog(t *testing.T, client *http.Client, req *http.Request, desc string) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf(desc+" error: %s", err.Error())
	} else if resp.StatusCode > 299 {
		t.Errorf(desc+" status not ok: %d", resp.StatusCode)
	} else {
		log.Printf(desc + " OK\n")
		t.Logf(desc + " OK\n")
	}
	return resp, err
}

// Test the API for the projects
func TestProjectApi(t *testing.T) {
	var projects *[]Project
	var containers *[]Container
	var req *http.Request
	var resp *http.Response
	var err error
	var found bool

	client := &http.Client{}

	// GET /token
	token := getToken(t, client)
	if token == "" {
		t.Errorf("Did not get token: %s", token)
	}

	// // GET /projects
	// req, err = http.NewRequest("GET", backUrl+"/projects", nil)
	// if err != nil {
	// 	t.Errorf("GET /projects error: %s", err.Error())
	// }
	// req.Header.Add("Authorization", "Bearer "+token)
	// resp, err = doRequestAndLog(client, req, "GET /projects", t)

	// projects = new([]Project)
	// json.NewDecoder(resp.Body).Decode(projects)
	// t.Logf("Get response from /projects: %v\n", projects)

	// Prepare test project
	p1 := getTestProject()
	// var testProjectHost = "code"
	// var testProjectPort = "8042"
	p1.UseHttps = false
	// p1.ExternalDomainNames = []string{testProjectHost}

	// POST /projects
	p1Buf := new(bytes.Buffer)
	json.NewEncoder(p1Buf).Encode(p1)
	t.Logf("Made p1Buf: %v\n", p1Buf)
	req, err = http.NewRequest("POST", backUrl+"/projects", p1Buf)
	if err != nil {
		t.Errorf("POST on /projects error: %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err = doRequestAndLog(t, client, req, "POST on /projects")

	// GET /projects
	req, err = http.NewRequest("GET", backUrl+"/projects", nil)
	if err != nil {
		t.Errorf("GET on /projects error: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err = doRequestAndLog(t, client, req, "GET /projects")

	projects = new([]Project)
	json.NewDecoder(resp.Body).Decode(projects)
	t.Logf("Get response from /projects: %v\n", projects)

	// Check if we got our project in the list
	found = false
	p1StrId := ""
	for _, project := range *projects {
		if project.Name == p1.Name {
			if found {
				t.Errorf("Project with name '%s' found multiple times?", p1.Name)
			} else {
				t.Log("Found p1 among list of projects")
				p1StrId = project.Id.Hex()
				found = true
			}
		}
	}
	if !found {
		t.Error("p1 not found in list of projects")
	}

	// Apply docker
	req, err = http.NewRequest("POST", backUrl+"/docker/apply/"+p1StrId, nil)
	if err != nil {
		t.Errorf("POST /docker/apply error: %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	t.Logf("Applying the test project...")
	resp, err = doRequestAndLog(t, client, req, "POST /docker/apply")
	t.Logf("Applying the test project returned")

	// Wait for it to be ready
	log.Println("Waiting for 20 seconds...")
	time.Sleep(20 * time.Second)
	log.Println("End of waiting")

	// // Check if up
	// req, err = http.NewRequest("GET", "http://"+testProjectHost+":"+testProjectPort, nil)
	// if err != nil {
	// 	t.Errorf("GET project homepage error: %s", err.Error())
	// }
	// resp, err = doRequestAndLog(t, client, req, "GET project homepage")

	// // Check if up
	// req, err = http.NewRequest("GET", backUrl+"/docker/status/"+p1StrId, nil)
	// if err != nil {
	// 	t.Errorf("GET /docker/status error: %s", err.Error())
	// }
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Authorization", "Bearer "+token)
	// resp, err = doRequestAndLog(t, client, req, "GET /docker/status")

	// containers = new([]Container)
	// json.NewDecoder(resp.Body).Decode(containers)
	// log.Printf("Get response from /docker/status: %v\n", containers)
	// t.Logf("Get response from /docker/status: %v\n", containers)

	// log.Printf("       Id: %v\n", (*containers)[0].Id)
	// log.Printf("    Image: %v\n", (*containers)[0].Image)
	// log.Printf("   Status: %v\n", (*containers)[0].Status)
	// log.Printf("    State: %v\n", (*containers)[0].State)

	// if len(*containers) == 1 && (*containers)[0].Image == "unittestproject" && (*containers)[0].State == "running" {
	// 	log.Printf("Get response from /docker/status: %v\n", containers)
	// 	t.Logf("Get response from /docker/status: %v\n", containers)
	// } else {
	// 	t.Errorf("Unsatisfying response from /docker/status: %v", containers)
	// }

	// // Debug wait
	// log.Println("Waiting...")
	// time.Sleep(180 * time.Second)
	// log.Println("Waiting stopped")

	// Down docker
	req, err = http.NewRequest("POST", backUrl+"/docker/down/"+p1StrId, nil)
	if err != nil {
		t.Errorf("POST /docker/down error: %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	t.Logf("Downing the test project...")
	resp, err = doRequestAndLog(t, client, req, "POST /docker/down")
	t.Logf("Downing the test project returned")

	// DELETE /projects
	req, err = http.NewRequest("DELETE", backUrl+"/projects/"+p1StrId, nil)
	if err != nil {
		t.Errorf("DELETE /projects/ error: %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err = doRequestAndLog(t, client, req, "DELETE /projects")

	// END
	t.Logf("TestProjectApi finished")
}
