package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"log"
	"regexp"
)

var (
	gitlabCmd = &cobra.Command{
		Use:     "gitlab",
		Aliases: []string{},
		Short:   "Gitlab API commands",
	}
)

var (
	token    string
	fileName string

	treeSearchCmd = &cobra.Command{
		Use:     "tree-search",
		Aliases: []string{},
		Short:   "Check GitLab repos for a file",
		Long:    `Use the GitLab API to check for a given file`,
		Run: func(cmd *cobra.Command, args []string) {
			checkForFile(token, fileName)
		},
	}
)

func init() {
	treeSearchCmd.Flags().StringVar(&token, "token", "t", "GitLab Private Token")
	treeSearchCmd.Flags().StringVar(&fileName, "file-name", "", "The filename to search for")
	err := treeSearchCmd.MarkFlagRequired("token")
	if err != nil {
		log.Fatalf("could not mark flag %v\n as required", token)
	}

	err = treeSearchCmd.MarkFlagRequired("file-name")
	if err != nil {
		log.Fatalf("could not mark flag %v\n as required", fileName)
	}
}

// Check for a fileName in GitLab using the API
func checkForFile(token string, fileName string) []string {
	var foundInProjects []string
	fmt.Printf("Let's search for the file %v\n", fileName)
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.com/api/v4"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	listTreeParams := &gitlab.ListTreeOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 50,
		},
		Recursive: &[]bool{true}[0],
	}

	groups := getGroups(token)
	for groupName, groupId := range groups {
		fmt.Printf("group: %v\t id: %v\n", groupName, groupId)
		projects := getProjects(token, groupName, groupId)
		for projectName, projectId := range projects {
			pid := projectId
			tree, r, err := client.Repositories.ListTree(pid, listTreeParams)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("searching for file %v\t in project %v:%v\n", fileName, projectName, projectId)
			for i, t := range tree {
				fmt.Printf("page: %v\t tree item: %v\t name: %v\t type: %v\n", r.CurrentPage, i, t.Name, t.Type)
				m, err := regexp.Match(fileName, []byte(t.Name))
				if err != nil {
					log.Fatalf("encountered error: %v\n", err)
				}
				if m {
					fmt.Printf("found requested file name: %v in project %v:%v", fileName, projectName, projectId)
					foundInProjects = append(foundInProjects, projectName)
					break
				}
			}

			if r.CurrentPage >= r.TotalPages {
				break
			}

			// Update the page number to get the next page.
			listTreeParams.Page = r.NextPage
		}
	}
	return foundInProjects
}

// Get the groups accessible to the private token source
func getGroups(token string) map[string]int {
	groups := make(map[string]int)
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.com/api/v4"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	g, _, err := client.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	if err != nil {
		log.Fatalf("failed to list groups: %v\n", err)
	}

	fmt.Println("getting group ID's")
	for i, g := range g {
		fmt.Printf("index: %v\t groupName: %v\t groupID: %v\n", i, g.Name, g.ID)
		groups[g.Name] = g.ID
	}
	return groups
}

// Get the projects for the given group
func getProjects(token string, groupName string, groupId int) map[string]int {
	mprojects := make(map[string]int)
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.com/api/v4"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	xp, _, err := client.Groups.ListGroupProjects(groupId, &gitlab.ListGroupProjectsOptions{
		//Archived: &[]bool{false}[0],
	})
	if err != nil {
		log.Fatalf("failed to list projects: %v\n", err)
	}

	fmt.Printf("getting project ID's for group %v\n", groupName)
	for _, p := range xp {
		fmt.Printf("projectName: %v\t projectID: %v\n", p.Name, p.ID)
		mprojects[p.Name] = p.ID
	}

	return mprojects
}
