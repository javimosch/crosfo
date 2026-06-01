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

func UpdateCommunityName(communityID int, newName string) error {
	_, err := db.Exec("UPDATE communities SET name = ? WHERE id = ?", newName, communityID)
	return err
}

func UpdateCommunityDescription(communityID int, newDescription string) error {
	_, err := db.Exec("UPDATE communities SET description = ? WHERE id = ?", newDescription, communityID)
	return err
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