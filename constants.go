package main

import (
	"fmt"
)

const (
	//SIMPLE REGULAR EXPRESSIONS
	//datetime_regex
	year_regex  = `(?P<year>\d\d\d\d)`
	month_regex = `(?P<month>\d\d)`
	day_regex   = `(?P<day>\d\d)`

	hour_regex   = `(?P<hour>\d\d)`
	minute_regex = `(?P<minute>\d\d)`
	second_regex = `(?P<second>\d\d)`

	millisecond_regex = `(?P<millisecond>\d+.\d*)`

	//*_command_regex
	ip_regex      = `(?:(?P<ip>\d+[.]\d+[.]\d+[.]\d+)(?::(?P<port>\d+))?)`
	steamid_regex = `(?P<steamid>\d+)`

	by_regex = `(?P<by>\w+)`
	//hotword_regex = `(?P<hotword>/w+)` //Gets any hotword
	admin_hotword_regex  = `(?P<hotword>re)`
	player_hotword_regex = `(?P<hotword>pm)`

	//tick_regex
	gametime_regex     = `(?P<gametime>\d+.\d\d)m`
	fps_regex          = `(?P<fps>\d+.\d\d)`
	heap_regex         = `(?P<heap>\d+.\d)MB`
	maxheap_regex      = `(?P<maxheap>\d+.\d)MB`
	chunks_regex       = `(?P<chunks>\d+)`
	cgo_regex          = `(?P<cgo>\d+)`
	playernum_regex    = `(?P<playernum>\d+)`
	zombienum_regex    = `(?P<zombienum>\d+)`
	entitynum_regex    = `(?P<entitynum>\d+)`
	entitynumpar_regex = `\((?P<entitynumpar>\d+)\)`
	itemnum_regex      = `(?P<itemnum>\d+)`

	//player_regex
	listnum_regex    = `(?P<listnum>\d+)`
	playerid_regex   = `(?P<playerid>\d+)`
	playername_regex = `(?P<playername>\w+)`

	playerx_regex = `(?P<playerx>-?\d+.\d)`
	playery_regex = `(?P<playery>-?\d+.\d)`
	playerz_regex = `(?P<playerz>-?\d+.\d)`
	playeru_regex = `(?P<playeru>-?\d+.\d)`
	playerv_regex = `(?P<playerv>-?\d+.\d)`
	playerw_regex = `(?P<playerw>-?\d+.\d)`

	playerremote_regex  = `(?P<playerremote>(?:True)|(?:False))`
	playerhealth_regex  = `(?P<playerhealth>\d+)`
	playerdeaths_regex  = `(?P<playerdeaths>\d+)`
	playerzkills_regex  = `(?P<playerzkills>\d+)`
	playerpkills_regex  = `(?P<playerpkills>\d+)`
	playerscore_regex   = `(?P<playerscore>\d+)`
	playerlevel_regex   = `(?P<playerlevel>\d+)`
	playersteamid_regex = `(?P<playersteamid>\d+)`
	playerip_regex      = `(?P<playerip>\d+[.]\d+[.]\d+[.]\d+)`
	playerping_regex    = `(?P<playerping>\d+)`

	keystonesnum_regex       = `(?P<keystonesnum>\d+)`
	keystonesprotected_regex = `(?P<keystonesprotected>(False)|(True))`
	keystoneshardness_regex  = `(?P<keystoneshardness>\d+)`
)

var (
	//COMPLEX REGULAR EXPRESSIONS
	date_regex     = fmt.Sprintf(`(?P<date>%s-%s-%s)`, year_regex, month_regex, day_regex)
	time_regex     = fmt.Sprintf(`(?P<time>%s:%s:%s)`, hour_regex, minute_regex, second_regex)
	datetime_regex = fmt.Sprintf(`%sT%s %s`, date_regex, time_regex, millisecond_regex)

	client_regex   = fmt.Sprintf(`(?:client %s)`, steamid_regex)
	ipclient_regex = fmt.Sprintf(`(?:%s|%s)`, ip_regex, client_regex)

	//player_regex
	playerposition_regex = fmt.Sprintf(`(?P<playerposition>\(%s, %s, %s\))`,
		playerx_regex, playery_regex, playerz_regex)
	playerrotation_regex = fmt.Sprintf(`(?P<playerrotation>\(%s, %s, %s\))`,
		playeru_regex, playerv_regex, playerw_regex)

	//TODO Make command a simple regex above?
	//admin_command_regex
	admin_command_regex = fmt.Sprintf(`^%s INF Executing command '%s (?P<command>.+)' (?:by %s )?from %s$`,
		datetime_regex, admin_hotword_regex, by_regex, ipclient_regex)
	//player_command_regex
	player_command_regex = fmt.Sprintf(`^%s INF Denying command '%s (?P<command>.+)' from client %s$`,
		datetime_regex, player_hotword_regex, playername_regex)

	//tick_regex
	tick_regex = fmt.Sprintf(`^%s INF Time: %s FPS: %s Heap: %s Max: %s Chunks: %s CGO: %s Ply: %s Zom: %s Ent: %s %s Items: %s$`,
		datetime_regex, gametime_regex, fps_regex, heap_regex, maxheap_regex, chunks_regex, cgo_regex, playernum_regex, zombienum_regex, entitynum_regex, entitynumpar_regex, itemnum_regex)

	//Each player in a player list (lp) command
	player_regex = fmt.Sprintf(`^%s. id=%s, %s, pos=%s, rot=%s, remote=%s, health=%s, deaths=%s, zombies=%s, players=%s, score=%s, level=%s, steamid=%s, ip=%s, ping=%s$`,
		listnum_regex, playerid_regex, playername_regex, playerposition_regex, playerrotation_regex, playerremote_regex, playerhealth_regex, playerdeaths_regex, playerzkills_regex, playerpkills_regex, playerscore_regex, playerlevel_regex, playersteamid_regex, playerip_regex, playerping_regex)

	//Trigger to prepare for list of land protection blocks
	keystonetrigger_regex = fmt.Sprintf(`^Player "%s \(%s\)" owns %s keystones \(protected: %s, current hardness multiplier: %s\)$`,
		playername_regex, playersteamid_regex, keystonesnum_regex, keystonesprotected_regex, keystoneshardness_regex)
	keystoneendtrigger_regex = fmt.Sprintf(`^Total of %s keystones in the game$`,
		keystonesnum_regex)
)
