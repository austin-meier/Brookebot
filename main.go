package main

import (
	"fmt"
	"os"
	"os/signal"
	"io/ioutil"
	"io"
	"net/http"
	"syscall"
	"strings"
	"math"
	"math/rand"
	"time"
	"regexp"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	filename string
	generalChannel string
	suggestionsChannel string
	rulesChannel string
	testChannel string
	pollsChannel string
	guildID string
	minimumRole string
	roles map[string]float64
	entries []Entry
	reactlist []string
	reactMode bool
	reactemoji *discordgo.Emoji
	david bool
	myID string
)

type Entry struct {
	UserID string `json:"userID"`
	Nick string `json:"nickname"`
	TotalMessages int `json:"totalMessages"`
	AverageLength float64 `json:"avgLength"`
	RankMultiplier float64 `json:"rankMultiplier"`
	Score float64 `json:"score"`
}

func (e Entry) String() string {
	return "Nick: " + e.Nick + " TotalMessages: " + string(e.TotalMessages) + " AverageLength: " + fmt.Sprintf("%.2f", e.AverageLength) + " Score: " + fmt.Sprintf("%.2f", e.Score)
}

func init() {
	Token = "NzAzNDEwNDAxNjc2NjIzODkz.XuTnlQ.m3udiaZ3udOJJWG-K6owmxdrenU" //BOT TOKEN
	generalChannel = "476175990931062787" // General Channel ID
	suggestionsChannel = "753250162067112006" // Suggestions Channel ID
	rulesChannel = "592548422138068993" // Rules Channel ID
	guildID = "476175990931062785" // TheBrookeB Guild ID
	testChannel = "722507153285578872" // Bot Testing Channel ID
	pollsChannel = "826177720639029278" // Polls Channel ID
	filename = "users.json"; //File name for the user data to be stored
	minimumRole = "479802788101226509" //Moderator +
	roles = make(map[string]float64)
	roles["593310655696732163"] = 1 //Golden & up
	roles["593308455247413270"] = 1.15 //Loyalist
	roles["594580477378166784"] = 1.25 //Newbies + Pals
	roles["default"] = 1.25 //No Role 
	david = false // Default state of the David Hyperbleble reaction feature
	reactMode = false //Default state of the specific user message reaction feature
	myID = "171417327961636864" //KingBunz Discord Author ID


	//LOAD SAVED ENTRIES FROM FILE
	jsonFile, err := os.Open(filename)
	if err == nil {
		defer jsonFile.Close()
		d, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal([]byte(d), &entries)
		fmt.Println("Succesfully loaded data")
	} else {
		fmt.Println("File does not exist yet")
	}
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}


	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	//Register the userJoin func as a callback for GuildMemberAdd events.
	dg.AddHandler(userJoin)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

var quotesRegex = regexp.MustCompile(`"(.*?)"`)
var parenthesisRegex = regexp.MustCompile(`\((.*?)\)`)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	//If message is coming from a bot
	if m.Author.Bot == true {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if strings.ToLower(m.Content) == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if strings.ToLower(m.Content) == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	//Add checkmark reactions for server suggestions
	if m.ChannelID == suggestionsChannel {
		if len(m.Mentions)  == 0 {
			_ = s.MessageReactionAdd(m.ChannelID,m.ID, "✅") 
			_ = s.MessageReactionAdd(m.ChannelID,m.ID, "❌")
		}
	}

	//Polling Chanel
	if m.ChannelID == pollsChannel {
		if strings.ToLower(m.Content) == "b!format" {
			s.ChannelMessageSend(m.ChannelID, "POLL FORMAT: **b!poll \"Topic\" \"Option 1\" \"Option 2\"** and so fourth, up to 9 options. \n\nMake sure to include the quotation marks around the topic and options or else they won't appear in the poll.")
		}

		//Allow moderator+ to send help messages
		if strings.Contains(m.Content, ":") {
			if hasPermissions(m.Member.Roles) {
				goto end
			} 
		}
		if strings.Contains(strings.ToLower(m.Content), "b!poll") {
			content := strings.ReplaceAll(m.Content, "“", "\"")
			content = strings.ReplaceAll(content, "”", "\"")

			options := getQuotedStrings(content)
			optionCount := len(options)

			if(optionCount > 9) {
				s.ChannelMessageSend(m.ChannelID, "There can only be a maximum of 9 options for the poll. Please try again.")
			} else {
				var msg strings.Builder
				for i, option := range options {
					if i == 0 {
						msg.WriteString("**" + "Poll: **" + strings.TrimSpace(option) + " \n")
					} else {
						switch i {
							case 1:
								msg.WriteString("\x31\xE2\x83\xA3")
							case 2:
								msg.WriteString("\x32\xE2\x83\xA3")
							case 3:
								msg.WriteString("\x33\xE2\x83\xA3")
							case 4:
								msg.WriteString("\x34\xE2\x83\xA3")
							case 5:
								msg.WriteString("\x35\xE2\x83\xA3")
							case 6:
								msg.WriteString("\x36\xE2\x83\xA3")
							case 7:
								msg.WriteString("\x37\xE2\x83\xA3")
							case 8:
								msg.WriteString("\x38\xE2\x83\xA3")
							case 9:
								msg.WriteString("\x39\xE2\x83\xA3")
						}
						msg.WriteString(" " + strings.TrimSpace(option) + "\n\n")
					}
				}
				msg.WriteString("Poll created by " + m.Author.Mention())
				reactMsg, _ := s.ChannelMessageSend(m.ChannelID, msg.String())
				_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
				for i := 1; i < optionCount; i++ {
					switch i {
						case 1:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x31\xE2\x83\xA3")
						case 2:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x32\xE2\x83\xA3")
						case 3:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x33\xE2\x83\xA3")
						case 4:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x34\xE2\x83\xA3")
						case 5:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x35\xE2\x83\xA3")
						case 6:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x36\xE2\x83\xA3")
						case 7:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x37\xE2\x83\xA3")
						case 8:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x38\xE2\x83\xA3")
						case 9:
							_ = s.MessageReactionAdd(m.ChannelID, reactMsg.ID, "\x39\xE2\x83\xA3")
					}
				}
			}
		} else {
			if hasPermissions(m.Member.Roles) {
				//Let them send a message
			} else {
				_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
			}
		}
		end:
	}

	if david {
		if  m.Author.ID == "330568658374098947" {
			e, err := s.State.Emoji("240281725056319498","650880252310192145")
			if err == nil {
				_ = s.MessageReactionAdd(m.ChannelID, m.ID, e.APIName())
			}
		}
	}

	if strings.ToLower(m.Content) == "b!davidmode" {
		if (m.Author.ID == myID || m.Author.ID == "330568658374098947") {
			if (david) {
				david = false
				s.ChannelMessageSend(m.ChannelID, "DavidMode Disabled")
			} else {
				david = true
				s.ChannelMessageSend(m.ChannelID, "DavidMode Enabled")
			}
		}
	}


	if m.ChannelID == generalChannel {
		length := float64(len(m.Content))
		if length > 1 {
			if !strings.Contains(m.Content, "http") {
				//Get updated name
				nick := m.Author.Username
				if (m.Member.Nick != "") {
					nick = m.Member.Nick
				}

				key := containsUser(entries, m.Author.ID)
				if key < 0 {
					fmt.Println("USER NOT FOUND: ADDING USER")

					e := Entry{
						UserID:	m.Author.ID,
						Nick: nick,
						TotalMessages: 1,
						AverageLength:	length,
						RankMultiplier: getMulitplier(m.Member.Roles),
						Score: 1,
					}
					entries = append(entries, e)

				} else {
					//Update user info
					entries[key].TotalMessages += 1
					entries[key].AverageLength = entries[key].AverageLength * (float64(entries[key].TotalMessages) - 1)/float64(entries[key].TotalMessages) + length/float64(entries[key].TotalMessages)
					entries[key].Nick = nick
					fmt.Println("USER INFO UPDATED")
				}
			}

			//Write to file
			file, _ := json.MarshalIndent(entries, "", " ")
			_ = ioutil.WriteFile(filename, file, 0644)
		}
	}

	if strings.ToLower(m.Content) == "b!debug" {
		if m.Author.ID == myID {
			err := sendUsers(m.ChannelID)
			if err != nil 
			{
				s.ChannelMessageSend(m.ChannelID, "ERROR: " + err)
			}
		}
	}

	if strings.ToLower(m.Content) == "b!loadjson" {
		if m.Author.ID == myID {
			if len(m.Attachments) == 1 {
				err := DownloadFile(filename, m.Attachments[0].URL)
				if err == nil {
					// LOAD ENTRIES FROM FILE
					jsonFile, err := os.Open(filename)
					if err == nil {
						defer jsonFile.Close()
						d, _ := ioutil.ReadAll(jsonFile)
						json.Unmarshal([]byte(d), &entries)
						s.ChannelMessageSend(m.ChannelID, "Succesfully loaded users JSON file")
					} else {
						s.ChannelMessageSend(m.ChannelID, "ERROR: File opening error")
					}	
				} else {
					s.ChannelMessageSend(m.ChannelID, "ERROR: error downloading JSON file from discord")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "ERROR: Wrong format. Only append one users.json atttachment.")
			}
		}
	}

	if strings.ToLower(m.Content) == "b!clear" {
		//DUMP OLD JSON JUST IN CASE
		err := sendUsers(m.ChannelID)
		if err == nil 
		{
			entries = nil
			s.ChannelMessageSend(m.ChannelID, "Entries Cleared")
		} else {
			s.ChannelMessageSend(m.ChannelID, "ERROR: " + err)
		}
	}

	if strings.ToLower(m.Content) == "b!draw" {
		if hasPermissions(m.Member.Roles) {
			if entries != nil {
				//Create weighted total
				total := 0.0
				for key, e := range entries {
					n := math.Pow(e.AverageLength, e.RankMultiplier) * float64(e.TotalMessages)/100
					entries[key].Score = n
					total += n
				}
				rand.Seed(time.Now().UnixNano())
				win := rand.Float64() * total
				fmt.Println(fmt.Sprintf("%.2f",win))
				counter := 0.0
				for _, e := range entries {
					counter += e.Score
					if win <= counter {
						s.ChannelMessageSend(m.ChannelID, "The selected winner is: " + e.Nick)
						break;
					}
				}
			} else { s.ChannelMessageSend(m.ChannelID, "There are no entries in the list") }
		}
	}

	if strings.ToLower(m.Content) == "b!score" {
		if hasPermissions(m.Member.Roles) {
			score()
		}
		s.ChannelMessageSend(m.ChannelID, "Succesfully ran manual scoring update algorithm")
	}

	if strings.ToLower(m.Content) == "b!leaderboard" {
		if hasPermissions(m.Member.Roles) {
			score()
			top := sortScore(entries)

			var str strings.Builder

			for i := 0; i < 5; i++ {
				str.WriteString("**" + fmt.Sprintf("%d",i+1) + ")** " + top[i].Nick + " **|** Score: " + fmt.Sprintf("%.2f", top[i].Score) + "\n" )
			}

			s.ChannelMessageSend(m.ChannelID, str.String())

		}
	}

	//MESSAGE REACTION MODE
	if reactMode {
		if contains(reactlist, m.Author.ID) {
			_ = s.MessageReactionAdd(m.ChannelID, m.ID, reactemoji.APIName())
		}
	}

	if strings.Contains(strings.ToLower(m.Content), "b!reacttarget") {
		if len(m.Mentions) != 0 {
			if contains(reactlist, m.Mentions[0].ID) {
				//Add user to the react list
				reactlist = append(reactlist, m.Mentions[0].ID)
				s.ChannelMessageSend(m.ChannelID, m.Mentions[0].Mention() +" added to reaction target list")
			} else {
				//Remove user from the react list
				reactlist = removeStringFromSlice(reactlist,  m.Mentions[0].ID)
				s.ChannelMessageSend(m.ChannelID, m.Mentions[0].Mention() +" removed from reaction target list")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "You did not mention a user with the command. ex: b!addreacttarget")
		}
	}

	if strings.ToLower(m.Content) == "b!reactmode" {
		//if hasPermissions(m.Member.Roles) {
			reactMode = !reactMode

			if reactMode {
				s.ChannelMessageSend(m.ChannelID, "React Mode Enabled")
				if reactemoji == nil {
					reactemoji, _ = s.State.Emoji("240281725056319498","650880252310192145")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "React Mode Disabled")
			}
		//}
	}

	if strings.Contains(strings.ToLower(m.Content), "b!setreactemoji") {
		//Check to see if Author is me (KingBunz)
		if m.Author.ID == "171417327961636864" {
			params := strings.Fields(m.Content)
			guildID := params[1]
			emojiID := params[2]

			e, err := s.State.Emoji(guildID, emojiID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "ERR: Emoji not found")
			} else {
				reactemoji = e
				s.ChannelMessageSend(m.ChannelID, "Auto-React Emoji succesfully set to " +  e.MessageFormat())
			}
		}
	}

	if strings.Contains(strings.ToLower(m.Content), "b!info") {
		if len(m.Mentions) != 0 {
			key := containsUser(entries, m.Mentions[0].ID)
			if key < 0 {
				s.ChannelMessageSend(m.ChannelID, "User not found on the list")
			} else {
				score()
				str := "Nick:** " + entries[key].Nick + " **TotalMessages:** " + fmt.Sprintf("%d",entries[key].TotalMessages) + " **AverageLength:** " + fmt.Sprintf("%.2f", entries[key].AverageLength) + " **RankMultiplier:** " + fmt.Sprintf("%.2f", entries[key].RankMultiplier) + " **Score:** " + fmt.Sprintf("%.2f", entries[key].Score) + "**"
				s.ChannelMessageSend(m.ChannelID, str)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "User not found")
		}
	}

	if strings.Contains(strings.ToLower(m.Content), "b!reactdraw") {
		params := strings.Fields(m.Content)
		cid := params[1]
		mid := params[2]


		msg, err := s.ChannelMessage(cid, mid)
		if err == nil {
			r := msg.Reactions
			var names []string
			//Get Users
			for _, react := range r {
				users, err := s.MessageReactions(cid, mid, react.Emoji.APIName(), 100, "", "")
				if err == nil {
					for _, user := range users {
						if !contains(names, user.Username) {
							names = append(names, user.Username)
						} 
					}
				} else {
					s.ChannelMessageSend(m.ChannelID,err.Error())
				}
			}

			winner := rand.Intn(len(names))

			s.ChannelMessageSend(m.ChannelID, "Winner: **" + names[winner] + "**")

		} else {
			s.ChannelMessageSend(m.ChannelID, "Message not found. Use format 'b!reactdraw <Channel ID> <Message ID'")
		}
	}

}

func userJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	if m.GuildID == guildID {
		st, err := s.Channel(rulesChannel)
		if err == nil {
			s.ChannelMessageSend(generalChannel, "Welcome to the server " + m.Mention() +"! Please make sure to read the " + st.Mention())
		}
	}
}

func getQuotedStrings(s string) []string {
	ms := quotesRegex.FindAllString(s, -1)
	ss := make([]string, len(ms))
	for i, m := range ms {
		ss[i] = m[1 : len(m)-1] // Note the substring of the match.
	}
	return ss

}

func contains (l []string, u string) bool {
	if len(l) != 0 {
		for _, i := range l {
			if i == u {
				return true
			}
		}
	}
	return false
}

func removeStringFromSlice (s []string, t string) []string {
	for key, user := range s {
		if user == t {
			s[key] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

func containsUser(e []Entry, id string) int {
	for key, i := range e {
		if i.UserID == id {
			return key
		}
	}
	return -1
}

func score() {
	if entries != nil {
		for key, e := range entries {
			n := math.Pow(e.AverageLength, e.RankMultiplier) * float64(e.TotalMessages)/100
			entries[key].Score = n
		}
	} else { fmt.Println("There are no entries in the list") }
}


func getMulitplier(r []string) float64 {
	i := roles["default"]
	for _, role := range r {
		for key, id := range roles {
			if role == key {
				if id < i {
					i = id
				}
			}
		}
	}
	//No Roles Found
	return i;
}

func toJSON(e []Entry) string {
	json, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

func hasPermissions(r []string) bool {
	for _, role := range r {
		if role ==  minimumRole {
			return true
		}
	}
	return false
}

func sortScore(e []Entry) []Entry {
	if len(e) < 2 {
		return e
	}

	left, right := 0, len(e)-1

	pivot := rand.Int() % len(e)

	e[pivot], e[right] = e[right], e[pivot]

	for i, _ := range e {
		if e[i].Score > e[right].Score {
			e[left], e[i] = e[i], e[left]
			left++
		}
	}

	e[left], e[right] = e[right], e[left]

	sortScore(e[:left])
	sortScore(e[left+1:])

	return e
}

func sendUsers(channel string) error {
	score()
	jsonFile, err := os.Open(filename)
	if err == nil {
		defer jsonFile.Close()
		s.ChannelFileSend(channel, "USERS.json", jsonFile)
	}
	return nil
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}