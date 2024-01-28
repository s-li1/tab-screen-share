package notion

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/joho/godotenv/autoload"
)

func UploadFile() (*http.Response, error) {
	secret := os.Getenv("NOTION_KEY")
	pageId := os.Getenv("NOTION_PAGE_ID")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/pages/%s", os.Getenv("NOTION_URL"), pageId), nil)
	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))
	request.Header.Add("Notion-Version", "2022-06-28")

	if err != nil {
		log.Fatalf("client: error making http request: %s\n", err)
	}

	return http.DefaultClient.Do(request)
}

