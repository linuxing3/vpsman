package cmd

import (
	"fmt"
	"log"

	"github.com/linuxing3/vpsman/util"
	"github.com/kjk/notionapi"
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)

var notionCmd = &cobra.Command{
	Use:   "notion",
	Short: "Manage your notion",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("notion called")
		notionMenu()
	},
}

func notionMenu() {
	menu := []string{"编程页面", "默认页面"}
	switch util.LoopInput("请选择: ", menu, true) {
	case 1:
		getCodingPage()
	case 2:
		getDefaultPage()
	default:
		return
	}
}

func getCodingPage() {
	client := &notionapi.Client{}
	client.AuthToken = viper.GetString("main.notion.token")
	pageID := "0eb709506d914ed8b3fe37b856dafcee"
	page, err := client.DownloadPage(pageID)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}
	fmt.Println(page)
}

func getDefaultPage() {
	client := &notionapi.Client{}
	client.AuthToken = viper.GetString("main.notion.token")
	pageID := "0eb709506d914ed8b3fe37b856dafcee"
	page, err := client.DownloadPage(pageID)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}
	res := findSubPageIDs(page)
	fmt.Println(res)
}

func findSubPageIDs(page *notionapi.Page) []string {
	blocks := page.Root().Content
	pageIDs := map[string]struct{}{}
	seen := map[string]struct{}{}
	toVisit := blocks
	for len(toVisit) > 0 {
		block := toVisit[0]
		toVisit = toVisit[1:]
		id := block.ID
		if block.Type == notionapi.BlockPage {
			pageIDs[id] = struct{}{}
			seen[id] = struct{}{}
		}
		for _, b := range block.Content {
			if b == nil {
				continue
			}
			id := block.ID
			if _, ok := seen[id]; ok {
				continue
			}
			toVisit = append(toVisit, b)
		}
	}
	res := []string{}
	for id := range pageIDs {
		res = append(res, id)
	}
	return res
}

func init() {
	rootCmd.AddCommand(notionCmd)
}
