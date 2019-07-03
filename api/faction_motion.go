package api

import(
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
)

type(
	FactionMotion struct {
		TableName struct{} `json:"-" sql:"faction__motions"`

		Id uint32 `json:"id"`
		FactionId uint16 `json:"-"`
		Faction *Faction `json:"faction"`
		AuthorId uint16 `json:"-"`
		Author *Player `json:"player"`
		Type string `json:"type"`
		IsApproved bool `json:"is_approved" sql:",notnull"`
		Data map[string]interface{} `json:"data"`
		CreatedAt time.Time `json:"created_at"`
		EndedAt time.Time `json:"ended_at"`
	}
	FactionVote struct {
		TableName struct{} `json:"-" sql:"faction__votes"`
		
		Id uint32 `json:"id"`
		MotionId uint32 `json:"-"`
		Motion *FactionMotion `json:"motion"`
		AuthorId uint16 `json:"-"`
		Author *Player `json:"player"`
		Option uint8 `json:"option"`
		CreatedAt time.Time `json:"created_at"`
	}
)

var factionMotionsData []string

const(
	MotionTypePlanetTaxes = "planet_taxes"
	MotionTypeRegime = "regime"
	VoteOptionYes = 1
	VoteOptionNo = 2
	VoteOptionNeither = 3
)

func InitFactionMotions() {
    defer CatchException()
    factionMotionsJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/motion_types.json")
    if err != nil {
        panic(NewException("Can't open faction motions configuration file", err))
    }
    if err := json.Unmarshal(factionMotionsJSON, &factionMotionsData); err != nil {
        panic(NewException("Can't read faction motions configuration file", err))
    }
}

func CreateFactionMotion(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	data := DecodeJsonRequest(r)

	player := context.Get(r, "player").(*Player)

	faction := getFaction(uint16(factionId))
	motion := faction.createMotion(
		player,
		data["type"].(string),
		data["data"].(map[string]interface{}),
	)
	SendJsonResponse(w, 201, motion)
}

func GetFactionMotions(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	
	SendJsonResponse(w, 200, getFactionMotions(uint16(factionId)))
}

func VoteFactionMotion(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	player := context.Get(r, "player").(*Player)
	factionId, _ := strconv.ParseUint(data["faction_id"], 10, 16)
	option, _ := strconv.ParseUint(data["option"], 10, 16)

	motion := getFactionMotion(uint16(factionId))
	vote := motion.vote(player, uint8(option))

	SendJsonResponse(w, 201, vote)
}

func (f *Faction) createMotion(author *Player, mType string, data map[string]interface{}) *FactionMotion {
	if !isMotionTypeValid(mType) {
		panic(NewHttpException(400, "factions.motions.invalid_type", nil))
	}
	motion := &FactionMotion{
		Type: mType,
		Data: data,
		FactionId: f.Id,
		Faction: f,
		AuthorId: author.Id,
		Author: author,
		CreatedAt: time.Now(),
		EndedAt: time.Now().Add(time.Hour * time.Duration(f.getSettings(FactionSettingsMotionDuration).Value)),
	}
	if err := Database.Insert(motion); err != nil {
		panic(NewException("Could not create faction motion", err))
	}
	Scheduler.AddTask(uint(time.Until(motion.EndedAt)), func() {
		motion.processResults()
	})
	return motion
}

func (m *FactionMotion) vote(author *Player, option uint8) *FactionVote {
	if m.hasVoted(author) {
		panic(NewHttpException(400, "faction.votes.already_voted", nil))
	}
	vote := &FactionVote{
		MotionId: m.Id,
		Motion: m,
		AuthorId: author.Id,
		Author: author,
		Option: option,
		CreatedAt: time.Now(),
	}
	if err := Database.Insert(vote); err != nil {
		panic(NewException("Could not create faction vote", err))
	}
	return vote
}

func isMotionTypeValid (mType string) bool {
	for _, t := range factionMotionsData {
		if mType == t {
			return true
		}
	}
	return false
}

func (m *FactionMotion) hasVoted(author *Player) bool {
	vote := &FactionVote{}

	if err := Database.Model(vote).Where("motion_id = ?", m.Id).Where("author_id = ?", author.Id).Select(); err != nil {
		return false
	}
	return true
}

func (m *FactionMotion) processResults() {

}

func getFactionMotion(id uint16) *FactionMotion {
	motion := &FactionMotion{}
	if err := Database.Model(motion).Where("id = ?", id).Select(); err != nil {
		panic(NewHttpException(404, "Motion not found", err))
	}
	return motion
}

func getFactionMotions(id uint16) []*FactionMotion {
	motions := make([]*FactionMotion, 0)
	if err := Database.Model(&motions).Where("faction_id = ?", id).Select(); err != nil {
		panic(NewHttpException(404, "Faction motions not found", err))
	}
	return motions
}