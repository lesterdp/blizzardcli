package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type Character struct {
	LastModified        int     `json:"lastModified"`
	Name                string  `json:"name"`
	Realm               string  `json:"realm"`
	Battlegroup         string  `json:"battlegroup"`
	Class               int     `json:"class"`
	Race                int     `json:"race"`
	Gender              int     `json:"gender"`
	Level               int     `json:"level"`
	AchievementPoints   int     `json:"achievementPoints"`
	Thumbnail           string  `json:"thumbnail"`
	CalcClass           string  `json:"calcClass"`
	Faction             int     `json:"faction"`
	Mounts              *Mounts `json:"mounts"`
	TotalHonorableKills int     `json:"totalHonorableKills"`
}

type Mounts struct {
	NumCollected    int      `json:"numCollected"`
	NumNotCollected int      `json:"numNotCollected"`
	Collected       []*Mount `json:"collected"`
}

type Mount struct {
	Name       string `json:"name"`
	SpellId    int    `json:"spellId"`
	CreatureId int    `json:"creatureId"`
	ItemId     int    `json:"itemId"`
	QualityId  int    `json:"qualityId"`
	Icon       string `json:"icon"`
	IsGround   bool   `json:"isGround"`
	IsFlying   bool   `json:"isFlying"`
	IsAquatic  bool   `json:"isAquatic"`
	IsJumping  bool   `json:"isJumping"`
}

func CreateAccessToken(clientID string, clientSecret string, region string) string {
	requestURL := fmt.Sprintf("https://%s.battle.net/oauth/token", region)
	body := strings.NewReader("grant_type=client_credentials")

	request, err := http.NewRequest(http.MethodPost, requestURL, body)
	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(clientID, clientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	tokenData := Token{}
	err = json.NewDecoder(response.Body).Decode(&tokenData)
	if err != nil {
		log.Fatal(err)
	}
	println(tokenData.AccessToken)
	return tokenData.AccessToken
}

func mountCount() {
	token := CreateAccessToken("{clientID}", "{clientSecret}}", "us")

	character_ptr := flag.String("c", "", "Name of Character")
	realm_ptr := flag.String("r", "", "Name of the Realm")
	//hasmount_ptr := flag.String("has-mount", "", "Name of Mount to check if you have")
	headless_ptr := flag.Bool("has-headless-horseman", false, "Flag to Set for Headless Horseman Check")

	flag.Parse()
	fmt.Println("Character:", *character_ptr)
	fmt.Println("Realm:", *realm_ptr)

	requestURL := fmt.Sprintf("https://us.api.blizzard.com/wow/character/%s/%s?fields=mounts&locale=en_US&access_token=%s", *realm_ptr, *character_ptr, token)

	response, err := http.Get(requestURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	user := Character{}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		log.Fatalln(err)
	}

	collected := user.Mounts.NumCollected
	output := fmt.Sprintf("Mounts Collected: %d", collected)
	println(output)

	if *headless_ptr {
		for index := range user.Mounts.Collected {
			if user.Mounts.Collected[index].Name == "Headless Horseman's Mount" {
				println("HOLY SHIT YOU HAVE THAT MOUNT?!?!?!!!!!!!!11111")
				return
			}
		}
		println("F's in the chat")
	}
	/*
		if hasmount_ptr.Set() {
			for index := range user.Mounts.Collected {
				if user.Mounts.Collected[index].Name == *hasmount_ptr {
					println("HOLY SHIT YOU HAVE THAT MOUNT?!?!?!!!!!!!!11111")
					return
				}
			}
			println("F's in the chat")
		}
	*/

}

func main() {
	mountCount()
}
