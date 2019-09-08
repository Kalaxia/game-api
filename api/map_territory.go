package api

import(
	"net/http"
	"github.com/gorilla/context"
	"math"
	"time"
)

const(
	TerritoryActionCreation="creation"
	TerritoryActionPlanetGained="planet_gained"
	TerritoryActionPlanetLost="planet_lost"
	TerritoryActionSystemGained="system_gained"
	TerritoryActionSystemLost="system_lost"
	TerritoryActionConquest="conquest"
	TerritoryStatusPledge="pledge"
	TerritoryStatusContest="contest"
)

type(
	Territory struct {
		TableName struct{} `json:"-" sql:"map__territories"`

		Id uint16 `json:"id"`
		MapId uint16 `json:"-"`
		Map *Map `json:"-"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
		Coordinates CoordinatesSlice `json:"coordinates" sql:"type:jsonb" pq:",use_zero"`
		History []*TerritoryHistory `json:"history"`
	}

	TerritoryHistory struct {
		TableName struct{} `json:"-" sql:"map__territory_histories"`

		Id uint16 `json:"id"`
		TerritoryId uint16 `json:"-"`
		Territory *Territory `json:"-"`
		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		Action string `json:"action"`
		Data map[string]interface{} `json:"data" pg:",use_zero"`
		HappenedAt time.Time `json:"happened_at"`
	}
)

func GetMapTerritories(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
	starmap := getServerMap(player.ServerId)
	
	SendJsonResponse(w, 200, starmap.getTerritories())
}

func (m *Map) getTerritories() []*Territory {
	territories := make([]*Territory, 0)
	if err := Database.Model(&territories).Relation("Planet.Player.Faction").Relation("Planet.System").Where("territory.map_id = ?", m.Id).Select(); err != nil {
		panic(NewException("Could not retrieve map territories", err))
	}
	return territories
}

func (p *Planet) createTerritory() *Territory {
	t := &Territory{
		MapId: p.System.MapId,
		Map: p.System.Map,
		PlanetId: p.Id,
		Planet: p,
		Coordinates: make(CoordinatesSlice, 0),
	}
	if err := Database.Insert(t); err != nil {
		panic(NewException("Could not create territory", err))
	}
	t.addHistory(p.Player, TerritoryActionCreation, make(map[string]interface{}, 0))
	p.System.addTerritory(t)
	t.expand()
	return t
}

func (t *Territory) addHistory(p *Player, action string, data map[string]interface{}) *TerritoryHistory {
	th := &TerritoryHistory{
		TerritoryId: t.Id,
		Territory: t,
		PlayerId: p.Id,
		Player: p,
		Action: action,
		Data: data,
		HappenedAt: time.Now(),
	}
	if err := Database.Insert(th); err != nil {
		panic(NewException("Could not create territory history", err))
	}
	t.History = append(t.History, th)
	return th
}

func (t *Territory) getTotalInfluence() uint16 {
	influence := uint16(0)
	for _, st := range t.getSystemTerritories() {
		influence += st.getTotalInfluence()
	}
	return influence
}

func (t *Territory) expand() {
	for {
		t.Coordinates = make(CoordinatesSlice, 0)
		for _, st := range t.getSystemTerritories() {
			t.Coordinates = append(t.Coordinates, st.generateCoordinates()...)
		}
		if t.checkForIncludedSystems() == false {
			break
		}
	}
	t.update()
}

func (t *Territory) checkForIncludedSystems() (hasNewSystem bool) {
	minX, maxX, minY, maxY := t.getCoordLimits()

	systems := make([]*System, 0)
	if err := Database.
		Model(&systems).
		Where("x <= ?", maxX).
		Where("x >= ?", minX).
		Where("y <= ?", maxY).
		Where("y >= ?", minY).
		Where("id != ?", t.Planet.System.Id).
		Where("map_id = ?", t.Planet.System.MapId).
		Select(); err != nil {
		panic(NewException("Could not retrieve included systems", err))
	}
	for _, s := range systems {
		if t.isSystemIn(s) == true {
			hasNewSystem = true
			s.addTerritory(t)
		}
	}
	return
}

func (t *Territory) getCoordLimits() (minX, maxX, minY, maxY float64) {
	var minXInit, maxXInit, minYInit, maxYInit bool
	for _, coords := range t.Coordinates {
		if minXInit == false || coords.X < minX {
			minXInit = true
			minX = coords.X
		}
		if maxXInit == false || coords.X > maxX {
			maxXInit = true
			maxX = coords.X
		}
		if minYInit == false || coords.Y < minY {
			minYInit = true
			minY = coords.Y
		}
		if maxYInit == false || coords.Y > maxY {
			maxYInit = true
			maxY = coords.Y
		}
	}
	return
}

func (t *Territory) isSystemIn(s *System) bool {
    minX := t.Coordinates[0].X
    maxX := t.Coordinates[0].X
    minY := t.Coordinates[0].Y
    maxY := t.Coordinates[0].Y

    for _, p := range t.Coordinates {
        minX = math.Min(p.X, minX)
        maxX = math.Max(p.X, maxX)
        minY = math.Min(p.Y, minY)
        maxY = math.Max(p.Y, maxY)
	}
	systemX := float64(s.X)
	systemY := float64(s.Y)

    if systemX < minX || systemX > maxX || systemY < minY || systemY > maxY {
        return false
    }

    inside := false
    j := len(t.Coordinates) - 1
    for i := 0; i < len(t.Coordinates); i++ {
        if (t.Coordinates[i].Y > systemY) != (t.Coordinates[j].Y > systemY) && systemX < (t.Coordinates[j].X-t.Coordinates[i].X)*(systemY-t.Coordinates[i].Y)/(t.Coordinates[j].Y-t.Coordinates[i].Y)+t.Coordinates[i].X {
            inside = !inside
        }
        j = i
    }
    return inside
}


// func (t *Territory) convexHull() {
// 	sort.Sort(CoordinatesSlice(t.Coordinates));

// 	n := len(t.Coordinates);
// 	hull := CoordinatesSlice{}

// 	for i := 0; i < 2 * n; i++ {
// 		var j int
// 		if i < n {
// 			j = i
// 		} else {
// 			j = 2 * n - 1 - i
// 		}
// 		for len(hull) >= 2 && t.removeMiddle(hull[len(hull) - 2], hull[len(hull) - 1], t.Coordinates[j]) {
// 			hull = hull[1:]
// 		}
// 		hull = append(hull, t.Coordinates[j]);
// 	}
// 	t.Coordinates = hull[1:]
// }

// func (t *Territory) removeMiddle(a, b, c *Coordinates) bool {
// 	cross := (a.X - b.X) * (c.Y - b.Y) - (a.Y - b.Y) * (c.X - b.X);
// 	dot := (a.X - b.X) * (c.X - b.X) + (a.Y - b.Y) * (c.Y - b.Y);
// 	return cross < 0 || cross == 0 && dot <= 0;
// }

func (t *Territory) update() {
	if err := Database.Update(t); err != nil {
		panic(NewException("Could not update territory", err))
	}
}