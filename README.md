# GOLang Twitter and LUIS

This sample shows how to use some libraries in GOLANG to retrieve public twits from an account and validate the information against LUIS (www.luis.ai) from Microsoft.

## Requirements
1. Create your dev twitter account (https://developer.twitter.com/)
2. Create a LUIS account (https://www.luis.ai)
3. (Optional) You can import the LUIS appliction file to understand this sample [TwitterUnderstanding_LuisApplication.json](TwitterUnderstanding_LuisApplication.json)
4. Rename the file appkeysSample.yaml to appkeys.yaml and fillup with Twitter and LUIS keys.


## About the sample

This sample grabs information from @BombeirosPMESP, which reports incidents on Sao Paulo city, and try to export useful information for later processing. This means, in a text informed on twitter can be possible to extract address, time, type of incident and use this in some scenarios, described below:

a)  It can be possible to create a real time incidents maps, showing pins with the location of collision

b) Create analytical data about incidents in city

c) Expand the solution to monitor other accounts and train LUIS to understand each type of Language.

The sample was made for pt-BR

## Technical Highlights

a) The code uses Viper (https://github.com/spf13/viper) to read information from yaml file

b) It runs every 10 seconds to get new information from the account. It will runs continually

c) It uses the package https://github.com/dghubble/go-twitter/twitter to retrieve information on Twitter

d) About LUIS, the idea was to get samples from the account that you want to work and train the *Entities* that you want to detect. On my case, I try to detect the *Intent* and *Entities* **Address** and **Reason**








