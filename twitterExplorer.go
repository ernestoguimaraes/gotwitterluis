package main

import (
	"flag"
	"fmt"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"
)

var (
	viperconsumerkey       = ""
	viperconsumersercret   = ""
	viperaccesstoken       = ""
	viperaccesstokensecret = ""
	viperluisappkey        = ""
	viperluisauthkey       = ""
	viperluisdomain        = ""
	viperpublicaccount     = ""
	channelBuffer          = 5
)

func main() {

	viper.SetConfigName("appkeys")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	s := make(chan string, channelBuffer)

	//Channel to Write Results from LUIS
	go WriteResult(s)

	fmt.Println("Program running...")

	//Executes the code each 1- seconds
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Searching for new tweets --------------- ", t)
			RetrieveTweets(s)
		}
	}()

	//Keeps main open
	runtime.Goexit()
}

//Show on Screen what found o LUIS
func WriteResult(s chan string) {
	for item := range s {
		fmt.Println(item)
	}

}

//Connects to the Twitter account
func RetrieveTweets(s chan string) {
	//Keys App
	viperpublicaccount := viper.GetString("twitterkeys.publicAccount")
	viperconsumerkey := viper.GetString("twitterkeys.consumerKey")
	viperconsumersercret := viper.GetString("twitterkeys.consumersercret")
	viperaccesstoken := viper.GetString("twitterkeys.accesstoken")
	viperaccesstokensecret := viper.GetString("twitterkeys.accesstokensecret")
	viperluisappkey = viper.GetString("luis.appkey")
	viperluisauthkey = viper.GetString("luis.authkey")
	viperluisdomain = viper.GetString("luis.domain")

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", viperconsumerkey, "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", viperconsumersercret, "Twitter Consumer Secret")
	accessToken := flags.String("access-token", viperaccesstoken, "Twitter Access Token")
	accessSecret := flags.String("access-secret", viperaccesstokensecret, "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *accessToken == "" {
		log.Fatal("Application Access Token required")
	}

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// // user timeline
	userTimelineParams := &twitter.UserTimelineParams{ScreenName: viperpublicaccount, Count: channelBuffer}
	tweets, _, err := client.Timelines.UserTimeline(userTimelineParams)

	if err != nil {
		fmt.Println(err)
	}

	var queryText = ""
	for index := 0; index < channelBuffer; index++ {

		queryText = tweets[index].Text
		fmt.Println(queryText)
		//LuisGetResults(queryText, s, viperluisappkey, viperluisauthkey, viperluisdomain)

		//Need this because luis has a throttling mechanism
		time.Sleep(500 * time.Millisecond)
	}
}

//Sends the query to LUIS
func LuisGetResults(c string, s chan string, appkey string, authkey string, domain string) {

	query := c
	endpointUrl := fmt.Sprintf("https://%s/%s?verbose=false&timezoneOffset=-360&subscription-key=%s&q=%s", domain, appkey, authkey, url.QueryEscape(query))
	response, err := http.Get(endpointUrl)

	// 401 - check value of 'subscription-key' - do not use authoring key!
	if err != nil {
		// handle error
		fmt.Println("error from Get")
		log.Fatal(err)
	}

	response2, err2 := ioutil.ReadAll(response.Body)

	if err2 != nil {
		// handle error
		fmt.Println("error from ReadAll")
		log.Fatal(err2)
	}
	s <- string(response2)
}
