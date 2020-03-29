package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

type Workspace struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type ProjectMembership struct {
	UserId string `json:"userId"`
}
type Project struct {
	Id string `json:"id"`
	Name string `json:"name"`

	ClientId string `json:"clientId"`
	ClientName string `json:"clientName"`

	Memberships []ProjectMembership `json:"memberships,omitempty"`
}

type TimeEntryTimeInterval struct {
	Start string `json:"start"`
	End string `json:"end"`
	Duration string `json:"duration"`
}
type TimeEntry struct {
	Id string `json:"id"`
	Description string `json:"description"`
	TimeInterval TimeEntryTimeInterval
}

type EnrichedWorkspaceUser struct{
	User

	TotalTime time.Duration `json:"totalTime"`

	TimeEntries []TimeEntry `json:"timeEntries"`
}
type EnrichedWorkspaceProject struct {
	Project

	Users []EnrichedWorkspaceUser `json:"users"`
}
type EnrichedWorkspace struct {
	Workspace
	
	Projects []EnrichedWorkspaceProject `json:"projects"`
}

func loadWorkspaces(httpClient *resty.Client) (workspaces []Workspace, err error) {
	res, err := httpClient.R().Get("workspaces")
	if err != nil {
		return
	}
	if res.StatusCode() != 200 {
		err = errors.New(res.Status())
		return
	}

	err = json.Unmarshal(res.Body(), &workspaces)
	if err != nil {
		return
	}

	return
}

func loadUsers(httpClient *resty.Client, workspaceId string) (users []User, err error) {
	res, err := httpClient.R().Get("workspaces/" + workspaceId + "/users")
	if err != nil {
		return
	}
	if res.StatusCode() != 200 {
		err = errors.New(res.Status())
		return
	}

	err = json.Unmarshal(res.Body(), &users)
	if err != nil {
		return
	}

	return
}

func loadProjects(httpClient *resty.Client, workspaceId string) (projects []Project, err error) {
	res, err := httpClient.R().Get("workspaces/" + workspaceId + "/projects")
	if err != nil {
		return
	}
	if res.StatusCode() != 200 {
		err = errors.New(res.Status())
		return
	}

	err = json.Unmarshal(res.Body(), &projects)
	if err != nil {
		return
	}

	return
}

func loadTimeEntries(httpClient *resty.Client, workspaceId string, userId string, projectId string) (timeEntries []TimeEntry, err error) {
	u := fmt.Sprintf("workspaces/%s/user/%s/time-entries?project=%s&page-size=10000", workspaceId, userId, projectId)
	res, err := httpClient.R().Get(u)
	if err != nil {
		return
	}
	if res.StatusCode() != 200 {
		err = errors.New(res.Status())
		return
	}

	err = json.Unmarshal(res.Body(), &timeEntries)
	if err != nil {
		return
	}

	return
}

// -- Helpers -- //

func durationParse(clockifyDuration string) (duration time.Duration, err error) {
	//	"duration": "PT1M4S", (Example: PT1H30M15S - 1 hour 30 minutes 15 seconds)
	clockifyDuration = clockifyDuration[2:]

	if strings.Contains(clockifyDuration, "H") {
		arr := strings.Split(clockifyDuration, "H")
		if len(arr) == 2 {
			val, err2 := strconv.ParseInt(arr[0], 10, 64)
			if err2 != nil {
				err = err2
				return
			}
			duration += time.Duration(val) * time.Hour
			clockifyDuration = arr[1]
		}
	}

	if strings.Contains(clockifyDuration, "M") {
		arr := strings.Split(clockifyDuration, "M")
		if len(arr) == 2 {
			val, err2 := strconv.ParseInt(arr[0], 10, 64)
			if err2 != nil {
				err = err2
				return
			}
			duration += time.Duration(val) * time.Minute
			clockifyDuration = arr[1]
		}
	}

	if strings.Contains(clockifyDuration, "S") {
		arr := strings.Split(clockifyDuration, "S")
		if len(arr) == 2 {
			val, err2 := strconv.ParseInt(arr[0], 10, 64)
			if err2 != nil {
				err = err2
				return
			}
			duration += time.Duration(val) * time.Second
		}
	}

	return
}

// -- Main -- //

func main() {
	var verboseFlag bool
	flag.BoolVar(&verboseFlag, "v", false, "-v")
	var jsonFlag bool
	flag.BoolVar(&jsonFlag, "json", false, "-json")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	httpClient := resty.New()
	httpClient.HostURL = "https://api.clockify.me/api/v1"
	httpClient.Header.Set("x-api-key", os.Getenv("API_KEY"))

	users := make(map[string]User)
	var enrichedWorkspaces []EnrichedWorkspace

	workspaces, err := loadWorkspaces(httpClient)
	if err != nil {
		log.Fatal(err)
	}

	for _, workspace := range workspaces {
		if !jsonFlag {
			fmt.Println("Workspace: " + workspace.Name)
		}

		enrichedWorkspace := EnrichedWorkspace{
			Workspace: workspace,
		}

		workspaceUsers, err := loadUsers(httpClient, workspace.Id)
		if err != nil {
			log.Fatal(err)
		}
		for _, user := range workspaceUsers {
			users[user.Id] = user
		}

		projects, err := loadProjects(httpClient, workspace.Id)
		if err != nil {
			log.Fatal(err)
		}

		for _, project := range projects {
			if !jsonFlag {
				fmt.Println("\tProject: " + project.Name)
			}

			enrichedWorkspaceProject := EnrichedWorkspaceProject{
				Project: project,
			}

			for _, membership := range project.Memberships {
				if !jsonFlag {
					fmt.Println("\t\tUser: " + users[membership.UserId].Name)
				}

				enrichedWorkspaceUser := EnrichedWorkspaceUser{
					User: users[membership.UserId],
					TimeEntries: []TimeEntry{},
				}

				timeEntries, err := loadTimeEntries(httpClient, workspace.Id, membership.UserId, project.Id)
				if err != nil {
					log.Fatal(err)
				}

				var total time.Duration
				for _, timeEntry := range timeEntries {
					d, err := durationParse(timeEntry.TimeInterval.Duration)
					if err != nil {
						log.Fatal(err)
					}

					total += d
				}

				if !jsonFlag {
					fmt.Println("\t\t\tTotal: " + total.String())
				}

				enrichedWorkspaceUser.TotalTime = total

				if verboseFlag {
					for _, timeEntry := range timeEntries {
						if !jsonFlag {
							fmt.Println("\t\t\tEntry: " + timeEntry.TimeInterval.Duration + " - " + timeEntry.Description)
						}

						enrichedWorkspaceUser.TimeEntries = append(enrichedWorkspaceUser.TimeEntries, timeEntry)
					}
				}

				enrichedWorkspaceProject.Users = append(enrichedWorkspaceProject.Users, enrichedWorkspaceUser)
			}

			enrichedWorkspace.Projects = append(enrichedWorkspace.Projects, enrichedWorkspaceProject)
		}

		enrichedWorkspaces = append(enrichedWorkspaces, enrichedWorkspace)
	}

	// Clean up
	for i := range enrichedWorkspaces {
		for j := range enrichedWorkspaces[i].Projects {
			enrichedWorkspaces[i].Projects[j].Memberships = nil
		}
	}

	enrichedWorkspacesStr, err := json.MarshalIndent(enrichedWorkspaces, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if jsonFlag {
		fmt.Println(string(enrichedWorkspacesStr))
	}
}
