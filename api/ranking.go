package api

import(
	"github.com/gorilla/context"
	"net/http"
)

func GetTerritorialRanking(w http.ResponseWriter, r *http.Request) {
	serverId := context.Get(r, "player").(*Player).Faction.ServerId
	
	var res []struct{
		NbSystems int `json:"nb_systems"`
		FactionId uint16 `json:"-"`
		Faction *Faction `json:"faction"`
	}
	if err := Database.Model((*System)(nil)).
		Column("system.faction_id").
		ColumnExpr("count(*) AS nb_systems").
		Join("INNER JOIN map__maps m ON m.id = system.map_id").
		Join("INNER JOIN servers serv ON serv.id = m.server_id").
		Group("system.faction_id").
		Where("system.faction_id IS NOT NULL").
		Where("serv.id = ?", serverId).
		Order("nb_systems DESC").
		Select(&res); err != nil {
			panic(NewException("Could not retrieve territorial ranking", err))
	}
	for i, r := range res {
		res[i].Faction = getFaction(r.FactionId)
	}
	SendJsonResponse(w, 200, res)
}

func GetFinancialRanking(w http.ResponseWriter, r *http.Request) {
	serverId := context.Get(r, "player").(*Player).Faction.ServerId
	ranking := make([]*Faction, 0)
	if err := Database.Model(&ranking).Order("wallet DESC").Where("server_id = ?", serverId).Select(); err != nil {
		panic(NewException("Could not retrieve financial ranking", err))
	}
	SendJsonResponse(w, 200, ranking)
}