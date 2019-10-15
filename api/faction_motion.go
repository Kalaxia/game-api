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
		tableName struct{} `json:"-" pg:"faction__motions"`

		Id uint32 `json:"id"`
		FactionId uint16 `json:"-"`
		Faction *Faction `json:"faction"`
		AuthorId uint16 `json:"-"`
		Author *Player `json:"author"`
		Type string `json:"type"`
		IsApproved bool `json:"is_approved" pg:",notnull,use_zero"`
		IsProcessed bool `json:"is_processed" pg:",notnull,use_zero"`
		Data map[string]interface{} `json:"data"`
		CreatedAt time.Time `json:"created_at"`
		EndedAt time.Time `json:"ended_at"`
	}
	FactionVote struct {
		tableName struct{} `json:"-" pg:"faction__votes"`
		
		Id uint32 `json:"id"`
		MotionId uint32 `json:"-"`
		Motion *FactionMotion `json:"motion"`
		AuthorId uint16 `json:"-"`
		Author *Player `json:"author"`
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
    defer CatchException(nil)
    factionMotionsJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/motion_types.json")
    if err != nil {
        panic(NewException("Can't open faction motions configuration file", err))
    }
    if err := json.Unmarshal(factionMotionsJSON, &factionMotionsData); err != nil {
        panic(NewException("Can't read faction motions configuration file", err))
	}
	scheduleInProgressMotions()
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

func GetFactionPreviousMotions(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	player := context.Get(r, "player").(*Player)

	faction := getFaction(uint16(factionId))

	if faction.Id != player.Faction.Id {
		panic(NewHttpException(403, "forbidden", nil))
	}
	SendJsonResponse(w, 200, faction.getPreviousMotions())
}

func GetFactionCurrentMotions(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	player := context.Get(r, "player").(*Player)

	faction := getFaction(uint16(factionId))

	if faction.Id != player.Faction.Id {
		panic(NewHttpException(403, "forbidden", nil))
	}
	SendJsonResponse(w, 200, faction.getCurrentMotions())
}

func GetFactionMotion(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	motionId, _ := strconv.ParseUint(mux.Vars(r)["motion_id"], 10, 16)
	player := context.Get(r, "player").(*Player)

	faction := getFaction(uint16(factionId))

	if faction.Id != player.Faction.Id {
		panic(NewHttpException(403, "forbidden", nil))
	}
	SendJsonResponse(w, 200, faction.getMotion(uint32(motionId)))
}

func VoteFactionMotion(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	player := context.Get(r, "player").(*Player)
	factionId, _ := strconv.ParseUint(data["faction_id"], 10, 16)
	motionId, _ := strconv.ParseUint(data["motion_id"], 10, 16)
	option := DecodeJsonRequest(r)["option"].(float64)

	faction := getFaction(uint16(factionId))

	if faction.Id != player.Faction.Id {
		panic(NewHttpException(403, "forbidden", nil))
	}

	motion := faction.getMotion(uint32(motionId))
	vote := motion.vote(player, uint8(option))

	SendJsonResponse(w, 201, vote)
}

func GetFactionVote(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	player := context.Get(r, "player").(*Player)
	factionId, _ := strconv.ParseUint(data["faction_id"], 10, 16)
	motionId, _ := strconv.ParseUint(data["motion_id"], 10, 16)

	faction := getFaction(uint16(factionId))

	if faction.Id != player.Faction.Id {
		panic(NewHttpException(403, "forbidden", nil))
	}

	motion := faction.getMotion(uint32(motionId))

	SendJsonResponse(w, 200, motion.getVote(player))
}

func GetFactionVotes(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	player := context.Get(r, "player").(*Player)
	factionId, _ := strconv.ParseUint(data["faction_id"], 10, 16)
	motionId, _ := strconv.ParseUint(data["motion_id"], 10, 16)

	faction := getFaction(uint16(factionId))

	if faction.Id != player.Faction.Id {
		panic(NewHttpException(403, "forbidden", nil))
	}

	motion := faction.getMotion(uint32(motionId))

	if !motion.IsProcessed {
		panic(NewHttpException(403, "forbidden", nil))
	}

	SendJsonResponse(w, 200, motion.getVotes())
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
	motion.validate()
	if err := Database.Insert(motion); err != nil {
		panic(NewException("Could not create faction motion", err))
	}
	Scheduler.AddTask(uint(time.Until(motion.EndedAt)), func() {
		motion.processResults()
	})
	f.notify(NotificationTypeFaction, "faction.motions.new_motion", map[string]interface{}{
		"motion": motion,
	})
	return motion
}

func (m *FactionMotion) vote(author *Player, option uint8) *FactionVote {
	if time.Now().After(m.EndedAt) {
		panic(NewHttpException(400, "faction.votes.ended_vote", nil))
	}
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

func isMotionTypeValid(mType string) bool {
	m := &FactionMotion{}
	if err := Database.Model(m).Where("type = ?", mType).Where("is_processed = ?", false).Select(); err == nil {
		panic(NewHttpException(403, "faction.motions.currently_voting", nil))
	}

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

func (m *FactionMotion) getVote(author *Player) *FactionVote {
	vote := &FactionVote{}

	if err := Database.Model(vote).Where("motion_id = ?", m.Id).Where("author_id = ?", author.Id).Select(); err != nil {
		panic(NewHttpException(404, "faction.motions.votes.not_found", err))
	}
	return vote
}

func (m *FactionMotion) getVotes() map[uint8]int {
	votes := make([]*FactionVote, 0)

	if err := Database.Model(&votes).Where("motion_id = ?", m.Id).Select(); err != nil {
		panic(NewException("Could not retrieve motion votes", err))
	}

	result := map[uint8]int{
		VoteOptionYes: 0,
		VoteOptionNo: 0,
		VoteOptionNeither: 0,
	}
	for _, v := range votes {
		result[v.Option]++
	}
	return result
}

func (m *FactionMotion) processResults() {
	positiveVotes := m.countVotesByOption(VoteOptionYes)
	nbMembers := m.Faction.countMembers()
	m.IsApproved = positiveVotes >= (nbMembers / 2)
	m.IsProcessed = true
	if m.IsApproved {
		m.apply()
	}
	m.update()
	m.Faction.notify(NotificationTypeFaction, "faction.motions.motion_results", map[string]interface{}{
		"motion": m,
	})
}

func (m *FactionMotion) countVotesByOption(option int) int {
	votes := make([]*FactionVote, 0)
	count, err := Database.Model(&votes).Where("motion_id = ?", m.Id).Where("option = ?", option).Count()
	if err != nil {
		panic(NewException("Could not retrieve motion votes", err))
	}
	return count
}

func (f *Faction) getMotion(id uint32) *FactionMotion {
	motion := &FactionMotion{}
	if err := Database.Model(motion).Relation("Faction").Relation("Author.Faction").Where("faction_motion.id = ?", id).Where("faction_motion.faction_id = ?", f.Id).Select(); err != nil {
		panic(NewHttpException(404, "Motion not found", err))
	}
	return motion
}

func (f *Faction) getCurrentMotions() []*FactionMotion {
	motions := make([]*FactionMotion, 0)
	if err := Database.Model(&motions).Relation("Faction").Relation("Author.Faction").Where("faction_motion.faction_id = ?", f.Id).Where("is_processed = ?", false).Select(); err != nil {
		panic(NewHttpException(404, "Faction motions not found", err))
	}
	return motions
}

func (f *Faction) getPreviousMotions() []*FactionMotion {
	motions := make([]*FactionMotion, 0)
	if err := Database.Model(&motions).Relation("Faction").Relation("Author.Faction").Where("faction_motion.faction_id = ?", f.Id).Where("is_processed = ?", true).Order("ended_at DESC").Select(); err != nil {
		panic(NewHttpException(404, "Faction motions not found", err))
	}
	return motions
}

func scheduleInProgressMotions() {
	motions := make([]*FactionMotion, 0)
	
	if err := Database.Model(&motions).Relation("Author").Relation("Faction").Where("is_processed = ?", false).Select(); err != nil {
		panic(NewException("Faction motions could not be retrieved", err))
	}
	for _, m := range motions {
		if time.Now().After(m.EndedAt) {
			go m.processResults()
			continue
		}
		Scheduler.AddTask(uint(time.Until(m.EndedAt)), func() {
			m.processResults()
		})
	}
}

func (m *FactionMotion) validate() {
	switch (m.Type) {
		case MotionTypePlanetTaxes:
			m.Faction.validatePlanetTaxesMotion(int(m.Data["taxes"].(float64)))
			break
		default:
			panic(NewHttpException(400, "faction.motions.invalid_type", nil))
	}
}

func (m *FactionMotion) apply() {
	switch (m.Type) {
		case MotionTypePlanetTaxes:
			m.Faction.updatePlanetTaxes(int(m.Data["taxes"].(float64)))
			break
		default:
			panic(NewException("Unknown motion type", nil))
	}
}

func (m *FactionMotion) update() {
	if err := Database.Update(m); err != nil {
		panic(NewException("Could not save motion result", err))
	}
}