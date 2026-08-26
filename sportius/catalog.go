package sportius

// ExtensionID is the stable extension and persistence namespace.
const ExtensionID = "sportius"

type SportID string

const (
	SportBasketball  SportID = "basketball"
	SportChess       SportID = "chess"
	SportFootball    SportID = "football"
	SportGaelic      SportID = "gaelic-games"
	SportHockey      SportID = "hockey"
	SportMultiSport  SportID = "multi-sport"
	SportRunning     SportID = "running"
	SportRugby       SportID = "rugby"
	SportSwimming    SportID = "swimming"
	SportTableTennis SportID = "table-tennis"
	SportTennis      SportID = "tennis"
	SportVolleyball  SportID = "volleyball"
	SportOther       SportID = "other"
)

// SportDefinition is a stable catalogue entry. The ID is persisted; LabelKey
// is translated by each presentation surface.
type SportDefinition struct {
	ID       SportID
	LabelKey string
}

// SportCatalog is intentionally small for the first Telegram selector. The
// open SportID type and "other" entry allow the catalogue to grow without
// storing localised display strings in records.
var SportCatalog = []SportDefinition{
	{ID: SportBasketball, LabelKey: "sportius.sport.basketball"},
	{ID: SportChess, LabelKey: "sportius.sport.chess"},
	{ID: SportFootball, LabelKey: "sportius.sport.football"},
	{ID: SportGaelic, LabelKey: "sportius.sport.gaelic_games"},
	{ID: SportHockey, LabelKey: "sportius.sport.hockey"},
	{ID: SportMultiSport, LabelKey: "sportius.sport.multi_sport"},
	{ID: SportRunning, LabelKey: "sportius.sport.running"},
	{ID: SportRugby, LabelKey: "sportius.sport.rugby"},
	{ID: SportSwimming, LabelKey: "sportius.sport.swimming"},
	{ID: SportTableTennis, LabelKey: "sportius.sport.table_tennis"},
	{ID: SportTennis, LabelKey: "sportius.sport.tennis"},
	{ID: SportVolleyball, LabelKey: "sportius.sport.volleyball"},
	{ID: SportOther, LabelKey: "sportius.sport.other"},
}

type RoleID string

const (
	RolePlayer              RoleID = "player"
	RoleCoach               RoleID = "coach"
	RoleAssistantCoach      RoleID = "assistant-coach"
	RoleTeamManager         RoleID = "team-manager"
	RoleAdministrator       RoleID = "administrator"
	RoleOrganiser           RoleID = "organiser"
	RoleOfficial            RoleID = "official"
	RoleVolunteer           RoleID = "volunteer"
	RoleSupporter           RoleID = "supporter"
	RoleParentGuardian      RoleID = "parent-guardian"
	RoleMedicalWelfare      RoleID = "medical-welfare"
	RoleEquipmentManager    RoleID = "equipment-manager"
	RolePresident           RoleID = "president"
	RoleTreasurer           RoleID = "treasurer"
	RoleAccountant          RoleID = "accountant"
	RoleSecretary           RoleID = "secretary"
	RoleSafeguardingOfficer RoleID = "safeguarding-officer"
	RoleOther               RoleID = "other"
)

type RoleScope string

const (
	RoleScopePersonal RoleScope = "personal"
	RoleScopeTeam     RoleScope = "team"
	RoleScopeClub     RoleScope = "club"
)

type RoleDefinition struct {
	ID              RoleID
	LabelKey        string
	Scopes          []RoleScope
	DefaultPersonal bool
	ImpliesStaff    bool
}

var RoleCatalog = []RoleDefinition{
	{ID: RolePlayer, LabelKey: "sportius.role.player", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam}, DefaultPersonal: true},
	{ID: RoleCoach, LabelKey: "sportius.role.coach", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, DefaultPersonal: true, ImpliesStaff: true},
	{ID: RoleAssistantCoach, LabelKey: "sportius.role.assistant_coach", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleTeamManager, LabelKey: "sportius.role.team_manager", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleAdministrator, LabelKey: "sportius.role.administrator", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleOrganiser, LabelKey: "sportius.role.organiser", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, DefaultPersonal: true, ImpliesStaff: true},
	{ID: RoleOfficial, LabelKey: "sportius.role.official", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam}, DefaultPersonal: true, ImpliesStaff: true},
	{ID: RoleVolunteer, LabelKey: "sportius.role.volunteer", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, DefaultPersonal: true, ImpliesStaff: true},
	{ID: RoleSupporter, LabelKey: "sportius.role.supporter", Scopes: []RoleScope{RoleScopePersonal}, DefaultPersonal: true},
	{ID: RoleParentGuardian, LabelKey: "sportius.role.parent_guardian", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam}, DefaultPersonal: true},
	{ID: RoleMedicalWelfare, LabelKey: "sportius.role.medical_welfare", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleEquipmentManager, LabelKey: "sportius.role.equipment_manager", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}, ImpliesStaff: true},
	{ID: RolePresident, LabelKey: "sportius.role.president", Scopes: []RoleScope{RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleTreasurer, LabelKey: "sportius.role.treasurer", Scopes: []RoleScope{RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleAccountant, LabelKey: "sportius.role.accountant", Scopes: []RoleScope{RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleSecretary, LabelKey: "sportius.role.secretary", Scopes: []RoleScope{RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleSafeguardingOfficer, LabelKey: "sportius.role.safeguarding_officer", Scopes: []RoleScope{RoleScopeClub}, ImpliesStaff: true},
	{ID: RoleOther, LabelKey: "sportius.role.other", Scopes: []RoleScope{RoleScopePersonal, RoleScopeTeam, RoleScopeClub}},
}
