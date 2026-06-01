package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

type Community struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	CreatedAt string `json:"created_at"`
}

type Entry struct {
	ID          int    `json:"id"`
	CommunityID int    `json:"community_id"`
	UserID      int    `json:"user_id"`
	Nickname    string `json:"nickname"`
	CreatedAt   string `json:"created_at"`
}

type SocialLink struct {
	ID        int    `json:"id"`
	EntryID   int    `json:"entry_id"`
	Platform  string `json:"platform"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ThumbUp struct {
	ID         int    `json:"id"`
	LinkID     int    `json:"link_id"`
	UserID     int    `json:"user_id"`
	Nickname   string `json:"nickname"`
	CreatedAt  string `json:"created_at"`
}

func Init(dbPath string) error {
	var initErr error
	once.Do(func() {
		initErr = initDB(dbPath)
	})
	return initErr
}

func initDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err = createSchema(); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS communities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname TEXT UNIQUE NOT NULL,
		password TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		community_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		nickname TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS social_links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entry_id INTEGER NOT NULL,
		platform TEXT NOT NULL,
		url TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (entry_id) REFERENCES entries(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS thumbs_up (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		link_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		nickname TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(link_id, user_id),
		FOREIGN KEY (link_id) REFERENCES social_links(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS community_admins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		community_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(community_id, user_id)
	);

	CREATE INDEX IF NOT EXISTS idx_entries_community ON entries(community_id);
	CREATE INDEX IF NOT EXISTS idx_social_links_entry ON social_links(entry_id);
	CREATE INDEX IF NOT EXISTS idx_thumbs_up_link ON thumbs_up(link_id);
	CREATE INDEX IF NOT EXISTS idx_community_admins_community ON community_admins(community_id);
	CREATE INDEX IF NOT EXISTS idx_community_admins_user ON community_admins(user_id);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	// Migration: Add password column to users table if it doesn't exist
	migration := `ALTER TABLE users ADD COLUMN password TEXT`
	db.Exec(migration) // Ignore error if column already exists

	return nil
}

func GetDB() *sql.DB {
	return db
}

func SetDB(database *sql.DB) {
	db = database
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func CreateCommunity(name, description string) (int, error) {
	result, err := db.Exec(
		"INSERT INTO communities (name, description) VALUES (?, ?)",
		name, description,
	)
	if err != nil {
		return 0, err
	}
	return getLastInsertID(result)
}

func GetCommunities() ([]Community, error) {
	rows, err := db.Query("SELECT id, name, description, created_at FROM communities ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []Community
	for rows.Next() {
		var c Community
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		communities = append(communities, c)
	}
	return communities, nil
}

func GetCommunityByName(name string) (*Community, error) {
	var c Community
	err := db.QueryRow(
		"SELECT id, name, description, created_at FROM communities WHERE name = ?",
		name,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetOrCreateUser(nickname string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE nickname = ?", nickname).Scan(&userID)
	if err == nil {
		return userID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	result, err := db.Exec("INSERT INTO users (nickname, password) VALUES (?, ?)", nickname, "")
	if err != nil {
		return 0, err
	}
	return getLastInsertID(result)
}

func CreateUser(nickname, password string) (int, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec("INSERT INTO users (nickname, password) VALUES (?, ?)", nickname, hashedPassword)
	if err != nil {
		return 0, err
	}
	return getLastInsertID(result)
}

func VerifyUserCredentials(nickname, password string) (int, bool, error) {
	var userID int
	var hashedPassword string
	err := db.QueryRow("SELECT id, password FROM users WHERE nickname = ?", nickname).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}

	// If no password is set (legacy users), allow access
	if hashedPassword == "" {
		return userID, true, nil
	}

	valid := VerifyPassword(password, hashedPassword)
	return userID, valid, nil
}

func UpdateUserPassword(userID int, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, userID)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateEntry(communityID, userID int, nickname string) (int, error) {
	result, err := db.Exec(
		"INSERT INTO entries (community_id, user_id, nickname) VALUES (?, ?, ?)",
		communityID, userID, nickname,
	)
	if err != nil {
		return 0, err
	}
	return getLastInsertID(result)
}

func GetEntriesByCommunity(communityID int) ([]Entry, error) {
	rows, err := db.Query(
		"SELECT id, community_id, user_id, nickname, created_at FROM entries WHERE community_id = ? ORDER BY created_at DESC",
		communityID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.CommunityID, &e.UserID, &e.Nickname, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func GetEntryByID(entryID int) (*Entry, error) {
	var e Entry
	err := db.QueryRow(
		"SELECT id, community_id, user_id, nickname, created_at FROM entries WHERE id = ?",
		entryID,
	).Scan(&e.ID, &e.CommunityID, &e.UserID, &e.Nickname, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func CreateSocialLink(entryID int, platform, url string) (int, error) {
	result, err := db.Exec(
		"INSERT INTO social_links (entry_id, platform, url) VALUES (?, ?, ?)",
		entryID, platform, url,
	)
	if err != nil {
		return 0, err
	}
	return getLastInsertID(result)
}

func GetSocialLinksByEntry(entryID int) ([]SocialLink, error) {
	rows, err := db.Query(
		"SELECT id, entry_id, platform, url, created_at FROM social_links WHERE entry_id = ? ORDER BY created_at",
		entryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []SocialLink
	for rows.Next() {
		var l SocialLink
		if err := rows.Scan(&l.ID, &l.EntryID, &l.Platform, &l.URL, &l.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}

func AddThumbUp(linkID, userID int, nickname string) error {
	_, err := db.Exec(
		"INSERT OR IGNORE INTO thumbs_up (link_id, user_id, nickname) VALUES (?, ?, ?)",
		linkID, userID, nickname,
	)
	return err
}

func GetThumbsUpByLink(linkID int) ([]ThumbUp, error) {
	rows, err := db.Query(
		"SELECT id, link_id, user_id, nickname, created_at FROM thumbs_up WHERE link_id = ? ORDER BY created_at",
		linkID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var thumbs []ThumbUp
	for rows.Next() {
		var t ThumbUp
		if err := rows.Scan(&t.ID, &t.LinkID, &t.UserID, &t.Nickname, &t.CreatedAt); err != nil {
			return nil, err
		}
		thumbs = append(thumbs, t)
	}
	return thumbs, nil
}

func GetThumbsUpCountByLink(linkID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM thumbs_up WHERE link_id = ?", linkID).Scan(&count)
	return count, err
}

func GetEntriesByUserNickname(nickname string) ([]Entry, error) {
	rows, err := db.Query(
		"SELECT id, community_id, user_id, nickname, created_at FROM entries WHERE nickname = ? ORDER BY created_at DESC",
		nickname,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.CommunityID, &e.UserID, &e.Nickname, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func UpdateEntry(entryID int, nickname string) error {
	_, err := db.Exec("UPDATE entries SET nickname = ? WHERE id = ?", nickname, entryID)
	return err
}

func DeleteEntry(entryID int) error {
	_, err := db.Exec("DELETE FROM entries WHERE id = ?", entryID)
	return err
}

func UpdateSocialLink(linkID int, platform, url string) error {
	_, err := db.Exec("UPDATE social_links SET platform = ?, url = ? WHERE id = ?", platform, url, linkID)
	return err
}

func DeleteSocialLink(linkID int) error {
	_, err := db.Exec("DELETE FROM social_links WHERE id = ?", linkID)
	return err
}

func DeleteUser(userID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete thumbs_up for links belonging to this user
	_, err = tx.Exec(`
		DELETE FROM thumbs_up 
		WHERE link_id IN (
			SELECT id FROM social_links 
			WHERE entry_id IN (SELECT id FROM entries WHERE user_id = ?)
		)`, userID)
	if err != nil {
		return err
	}

	// Delete social_links for this user's entries
	_, err = tx.Exec(`
		DELETE FROM social_links 
		WHERE entry_id IN (SELECT id FROM entries WHERE user_id = ?)`, userID)
	if err != nil {
		return err
	}

	// Delete entries for this user
	_, err = tx.Exec("DELETE FROM entries WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	// Delete the user
	_, err = tx.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetEntryByUserID(userID, communityID int) (*Entry, error) {
	var e Entry
	err := db.QueryRow(
		"SELECT id, community_id, user_id, nickname, created_at FROM entries WHERE user_id = ? AND community_id = ?",
		userID, communityID,
	).Scan(&e.ID, &e.CommunityID, &e.UserID, &e.Nickname, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetUserIDByLinkID(linkID int) (int, error) {
	var userID int
	err := db.QueryRow(`
		SELECT e.user_id 
		FROM social_links sl 
		JOIN entries e ON sl.entry_id = e.id 
		WHERE sl.id = ?
	`, linkID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Community admin functions
func IsCommunityAdmin(communityID, userID int) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM community_admins 
		WHERE community_id = ? AND user_id = ?
	`, communityID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func HasCommunityAdmins(communityID int) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM community_admins 
		WHERE community_id = ?
	`, communityID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func AddCommunityAdmin(communityID, userID int) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO community_admins (community_id, user_id)
		VALUES (?, ?)
	`, communityID, userID)
	return err
}

func RemoveCommunityAdmin(communityID, userID int) error {
	_, err := db.Exec(`
		DELETE FROM community_admins 
		WHERE community_id = ? AND user_id = ?
	`, communityID, userID)
	return err
}

func GetCommunityAdmins(communityID int) ([]User, error) {
	rows, err := db.Query(`
		SELECT u.id, u.nickname, u.created_at
		FROM community_admins ca
		JOIN users u ON ca.user_id = u.id
		WHERE ca.community_id = ?
		ORDER BY ca.added_at ASC
	`, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.CreatedAt); err != nil {
			return nil, err
		}
		admins = append(admins, user)
	}
	return admins, nil
}

func GetUserByNickname(nickname string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, nickname, created_at FROM users WHERE nickname = ?", nickname).Scan(&user.ID, &user.Nickname, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, nickname, created_at FROM users ORDER BY nickname ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Nickname, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func UpdateCommunityName(communityID int, newName string) error {
	_, err := db.Exec("UPDATE communities SET name = ? WHERE id = ?", newName, communityID)
	return err
}

func UpdateCommunityDescription(communityID int, newDescription string) error {
	_, err := db.Exec("UPDATE communities SET description = ? WHERE id = ?", newDescription, communityID)
	return err
}

// VoterCount holds a voter nickname and how many links of an entry they thumbed up.
type VoterCount struct {
	Nickname string `json:"nickname"`
	Count    int    `json:"count"`
}

// GetThumbsUpVotersByEntry returns all voters for an entry aggregated across all its links,
// sorted by count descending.
func GetThumbsUpVotersByEntry(entryID int) ([]VoterCount, error) {
	rows, err := db.Query(`
		SELECT tu.nickname, COUNT(*) as cnt
		FROM thumbs_up tu
		JOIN social_links sl ON tu.link_id = sl.id
		WHERE sl.entry_id = ?
		GROUP BY tu.nickname
		ORDER BY cnt DESC
	`, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var voters []VoterCount
	for rows.Next() {
		var v VoterCount
		if err := rows.Scan(&v.Nickname, &v.Count); err != nil {
			return nil, err
		}
		voters = append(voters, v)
	}
	return voters, nil
}

// UserProfileData holds a user's entry+links for a community plus reciprocity counts.
type UserProfileData struct {
	Entry        *Entry       `json:"entry"`
	Links        []SocialLink `json:"links"`
	ThumbsGiven  int          `json:"thumbs_given"`   // how many thumbs viewer gave to this user
	ThumbsReceived int        `json:"thumbs_received"` // how many thumbs this user gave to viewer
}

// GetUserProfileInCommunity returns a user's links in a community and reciprocity counts
// between that user and the viewer.
func GetUserProfileInCommunity(communityName, nickname, viewerNickname string) (*UserProfileData, error) {
	community, err := GetCommunityByName(communityName)
	if err != nil {
		return nil, err
	}

	var userID int
	if err := db.QueryRow("SELECT id FROM users WHERE nickname = ?", nickname).Scan(&userID); err != nil {
		return nil, err
	}

	entry, err := GetEntryByUserID(userID, community.ID)
	if err != nil {
		return nil, err
	}

	links, err := GetSocialLinksByEntry(entry.ID)
	if err != nil {
		links = []SocialLink{}
	}

	// Thumbs viewer gave to this user (viewer thumbed user's links in this community)
	var thumbsGiven int
	db.QueryRow(`
		SELECT COUNT(*)
		FROM thumbs_up tu
		JOIN social_links sl ON tu.link_id = sl.id
		JOIN entries e ON sl.entry_id = e.id
		WHERE e.community_id = ? AND e.user_id = ? AND tu.nickname = ?
	`, community.ID, userID, viewerNickname).Scan(&thumbsGiven)

	// Thumbs this user gave to viewer (user thumbed viewer's links in this community)
	var thumbsReceived int
	if viewerNickname != "" {
		var viewerUserID int
		if err2 := db.QueryRow("SELECT id FROM users WHERE nickname = ?", viewerNickname).Scan(&viewerUserID); err2 == nil {
			db.QueryRow(`
				SELECT COUNT(*)
				FROM thumbs_up tu
				JOIN social_links sl ON tu.link_id = sl.id
				JOIN entries e ON sl.entry_id = e.id
				WHERE e.community_id = ? AND e.user_id = ? AND tu.nickname = ?
			`, community.ID, viewerUserID, nickname).Scan(&thumbsReceived)
		}
	}

	return &UserProfileData{
		Entry:          entry,
		Links:          links,
		ThumbsGiven:    thumbsGiven,
		ThumbsReceived: thumbsReceived,
	}, nil
}

// UserStatRow holds stats for a single user for the global leaderboard.
type UserStatRow struct {
	Nickname       string `json:"nickname"`
	ThumbsGiven    int    `json:"thumbs_given"`
	ThumbsReceived int    `json:"thumbs_received"`
	Diff           int    `json:"diff"` // given - received
	Color          string `json:"color"` // "green", "orange", "red"
}

// GetGlobalStats returns all users with thumbs given/received sorted red→orange→green.
// Pass communityName="" for global stats, or a community name to filter.
func GetGlobalStats(communityName string) ([]UserStatRow, []UserStatRow, error) {
	var givenQuery, receivedQuery string
	var args []interface{}

	if communityName != "" {
		givenQuery = `
			SELECT tu.nickname, COUNT(*) as cnt
			FROM thumbs_up tu
			JOIN social_links sl ON tu.link_id = sl.id
			JOIN entries e ON sl.entry_id = e.id
			JOIN communities c ON e.community_id = c.id
			WHERE c.name = ?
			GROUP BY tu.nickname`
		receivedQuery = `
			SELECT tu2.nickname as receiver, COUNT(*) as cnt
			FROM thumbs_up tu2
			JOIN social_links sl2 ON tu2.link_id = sl2.id
			JOIN entries e2 ON sl2.entry_id = e2.id
			JOIN communities c2 ON e2.community_id = c2.id
			WHERE c2.name = ?
			GROUP BY receiver`
		args = []interface{}{communityName}
	} else {
		givenQuery = `
			SELECT nickname, COUNT(*) as cnt
			FROM thumbs_up
			GROUP BY nickname`
		receivedQuery = `
			SELECT u.nickname as receiver, COUNT(*) as cnt
			FROM thumbs_up tu
			JOIN social_links sl ON tu.link_id = sl.id
			JOIN entries e ON sl.entry_id = e.id
			JOIN users u ON e.user_id = u.id
			GROUP BY receiver`
	}

	// Build given map
	givenMap := map[string]int{}
	rows, err := db.Query(givenQuery, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var nick string
		var cnt int
		rows.Scan(&nick, &cnt)
		givenMap[nick] = cnt
	}

	// Build received map
	receivedMap := map[string]int{}
	rows2, err := db.Query(receivedQuery, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var nick string
		var cnt int
		rows2.Scan(&nick, &cnt)
		receivedMap[nick] = cnt
	}

	// Collect all unique nicknames from users table (only those with any activity)
	allNicks := map[string]bool{}
	for n := range givenMap {
		allNicks[n] = true
	}
	for n := range receivedMap {
		allNicks[n] = true
	}

	var red, orange, green []UserStatRow
	for nick := range allNicks {
		given := givenMap[nick]
		received := receivedMap[nick]
		diff := given - received
		color := "orange"
		if diff > 5 {
			color = "green"
		} else if diff < -5 {
			color = "red"
		}
		row := UserStatRow{
			Nickname:       nick,
			ThumbsGiven:    given,
			ThumbsReceived: received,
			Diff:           diff,
			Color:          color,
		}
		switch color {
		case "red":
			red = append(red, row)
		case "orange":
			orange = append(orange, row)
		case "green":
			green = append(green, row)
		}
	}

	// Sort each group by diff ascending (worst red first, best green first)
	sortStatRows(red, true)
	sortStatRows(orange, false)
	sortStatRows(green, false)

	// Combined sorted list: red → orange → green
	all := append(append(red, orange...), green...)

	// Top 5 green (highest diff)
	top5 := green
	if len(top5) > 5 {
		top5 = top5[:5]
	}

	return all, top5, nil
}

func sortStatRows(rows []UserStatRow, ascending bool) {
	for i := 1; i < len(rows); i++ {
		for j := i; j > 0; j-- {
			if ascending && rows[j].Diff < rows[j-1].Diff {
				rows[j], rows[j-1] = rows[j-1], rows[j]
			} else if !ascending && rows[j].Diff > rows[j-1].Diff {
				rows[j], rows[j-1] = rows[j-1], rows[j]
			} else {
				break
			}
		}
	}
}

func getLastInsertID(result sql.Result) (int, error) {
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func GetDefaultDBPath() string {
	if dbPath := os.Getenv("FFAF_DB_PATH"); dbPath != "" {
		return dbPath
	}
	return "./ffaf.db"
}