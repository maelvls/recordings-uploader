package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	drivev3 "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	youtubev3 "google.golang.org/api/youtube/v3"
)

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope, drive.DriveFileScope, drive.DriveMetadataReadonlyScope, youtube.YoutubeUploadScope, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	drive, err := drivev3.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}

	l, err := drive.Files.List().
		Q("'1zyxXFzXyKmexQm6Yo9MsWl3XpD3TDrgI' in parents").
		Corpora("allDrives").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Fields("nextPageToken, files(id,name,size)").
		Do()
	if err != nil {
		panic(err)
	}

	for _, d := range l.Files {
		fmt.Printf("%s %s\n", d.Id, d.Name)
	}

	youtube, err := youtubev3.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}

	// 50 is the maximum.
	ytVideos, err := youtube.Search.List([]string{"id,snippet"}).ChannelId("UCNPMkzGrAsQxVUFMPn7n88Q").MaxResults(50).Type("video").Do()
	if err != nil {
		panic(err)
	}

	existingNames := make(map[string]struct{})
	for _, v := range ytVideos.Items {
		existingNames[v.Snippet.Title] = struct{}{}
		fmt.Printf("%s %s\n", v.Id.VideoId, v.Snippet.Title)
	}

	var toBeUploaded []drivev3.File
	for _, f := range l.Files {
		_, ok := existingNames[f.Name]
		if !ok {
			toBeUploaded = append(toBeUploaded, *f)
		}
	}

	fmt.Printf("%d files to be uploaded:\n", len(toBeUploaded))
	for _, f := range toBeUploaded {
		fmt.Printf("%s %s\n", f.Id, f.Name)
	}

	for _, f := range toBeUploaded {
		fmt.Printf("Uploading %s\n", f.Name)

		youtube.Videos.Insert([]string{}, &youtubev3.Video{
			Snippet: &youtubev3.VideoSnippet{
				Title:       f.Name,
				Description: "Recorded with recordings-uploader",
				Tags:        []string{"cert-manager"},
				ChannelId:   "UCNPMkzGrAsQxVUFMPn7n88Q",
			},
			-@))
			@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@																																										@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@																																																																																			@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		})
		_, err := drive.Files.Update(f.Id, &drivev3.File{
			Name: f.Name,
			Parents: []string{
				"1zyxXFzXyKmexQm6Yo9MsWl3XpD3TDrgI",
			},
		}).SupportsAllDrives(true).Do()
		if err != nil {
			panic(err)
		}
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	tokFile := home + "/.config/recordings-uploader.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("\nVisit this URL to get an authorization code: %s\n\n", authURL)
	fmt.Printf("Paste code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
