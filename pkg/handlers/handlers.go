package handlers

import (
	"database/sql"
	"encoding/json"
	"ffaf/pkg/database"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	sendJSON(w, map[string]string{"error": message}, statusCode)
}

func HandleCommunities(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		communities, err := database.GetCommunities()
		if err != nil {
			sendError(w, "Failed to fetch communities", http.StatusInternalServerError)
			return
		}
		sendJSON(w, communities, http.StatusOK)
	case "POST":
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			sendError(w, "Name is required", http.StatusBadRequest)
			return
		}

		id, err := database.CreateCommunity(req.Name, req.Description)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint") {
				sendError(w, "Community name already exists", http.StatusConflict)
				return
			}
			sendError(w, "Failed to create community", http.StatusInternalServerError)
			return
		}

		sendJSON(w, map[string]interface{}{"id": id}, http.StatusCreated)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleCommunity(w http.ResponseWriter, r *http.Request) {
	// Extract community name from path: /api/community/{name}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}
	communityName := pathParts[2]

	switch r.Method {
	case "GET":
		community, err := database.GetCommunityByName(communityName)
		if err != nil {
			if err == sql.ErrNoRows {
				sendError(w, "Community not found", http.StatusNotFound)
				return
			}
			sendError(w, "Failed to fetch community", http.StatusInternalServerError)
			return
		}

		entries, err := database.GetEntriesByCommunity(community.ID)
		if err != nil {
			sendError(w, "Failed to fetch entries", http.StatusInternalServerError)
			return
		}

		// Ensure entries is never null
		if entries == nil {
			entries = []database.Entry{}
		}

		type LinkWithThumbs struct {
			database.SocialLink
			ThumbsCount int `json:"thumbs_count"`
		}

		type EntryWithLinks struct {
			database.Entry
			Links []LinkWithThumbs `json:"links"`
		}

		entriesWithLinks := []EntryWithLinks{}
		for _, entry := range entries {
			links, err := database.GetSocialLinksByEntry(entry.ID)
			if err != nil {
				continue
			}

			var linksWithThumbs []LinkWithThumbs
			for _, link := range links {
				count, _ := database.GetThumbsUpCountByLink(link.ID)
				linksWithThumbs = append(linksWithThumbs, LinkWithThumbs{
					SocialLink:  link,
					ThumbsCount: count,
				})
			}

			entriesWithLinks = append(entriesWithLinks, EntryWithLinks{
				Entry: entry,
				Links: linksWithThumbs,
			})
		}

		sendJSON(w, map[string]interface{}{
			"community": community,
			"entries":   entriesWithLinks,
		}, http.StatusOK)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleEntries(w http.ResponseWriter, r *http.Request) {
	// Extract community name from path: /api/entries/{name}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}
	communityName := pathParts[2]

	switch r.Method {
	case "POST":
		var req struct {
			Nickname string `json:"nickname"`
			Password string `json:"password"`
			Links    []struct {
				Platform string `json:"platform"`
				URL      string `json:"url"`
			} `json:"links"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Nickname == "" {
			sendError(w, "Nickname is required", http.StatusBadRequest)
			return
		}

		community, err := database.GetCommunityByName(communityName)
		if err != nil {
			sendError(w, "Community not found", http.StatusNotFound)
			return
		}

		var userID int
		// If password is provided, try to authenticate or create new user with password
		if req.Password != "" {
			// Try to authenticate first
			authUserID, valid, err := database.VerifyUserCredentials(req.Nickname, req.Password)
			if err == nil && valid {
				userID = authUserID
			} else {
				// Create new user with password
				userID, err = database.CreateUser(req.Nickname, req.Password)
				if err != nil {
					sendError(w, "Failed to create user", http.StatusInternalServerError)
					return
				}
			}
		} else {
			// Legacy: no password, use GetOrCreateUser
			userID, err = database.GetOrCreateUser(req.Nickname)
			if err != nil {
				sendError(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		}

		// Check if user already has an entry in this community
		existingEntry, err := database.GetEntryByUserID(userID, community.ID)
		if err == nil && existingEntry != nil {
			// User already has an entry, update it instead
			// Delete existing links
			links, err := database.GetSocialLinksByEntry(existingEntry.ID)
			if err == nil {
				for _, link := range links {
					database.DeleteSocialLink(link.ID)
				}
			}

			// Add new links
			for _, link := range req.Links {
				if link.Platform == "" || link.URL == "" {
					continue
				}
				database.CreateSocialLink(existingEntry.ID, link.Platform, link.URL)
			}

			sendJSON(w, map[string]interface{}{
				"id": existingEntry.ID,
				"updated": true,
				"message": "Your existing entry has been updated",
			}, http.StatusOK)
			return
		}

		// Create new entry
		entryID, err := database.CreateEntry(community.ID, userID, req.Nickname)
		if err != nil {
			sendError(w, "Failed to create entry", http.StatusInternalServerError)
			return
		}

		for _, link := range req.Links {
			if link.Platform == "" || link.URL == "" {
				continue
			}
			database.CreateSocialLink(entryID, link.Platform, link.URL)
		}

		sendJSON(w, map[string]interface{}{
			"id": entryID,
			"updated": false,
			"message": "Entry created successfully",
		}, http.StatusCreated)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleThumbUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		LinkID   int    `json:"link_id"`
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.LinkID == 0 || req.Nickname == "" {
		sendError(w, "link_id and nickname are required", http.StatusBadRequest)
		return
	}

	// Verify credentials if password is provided
	if req.Password != "" {
		_, valid, err := database.VerifyUserCredentials(req.Nickname, req.Password)
		if err != nil || !valid {
			sendError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	userID, err := database.GetOrCreateUser(req.Nickname)
	if err != nil {
		sendError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Check if the link belongs to the user
	linkUserID, err := database.GetUserIDByLinkID(req.LinkID)
	if err != nil {
		sendError(w, "Failed to verify link ownership", http.StatusInternalServerError)
		return
	}

	if linkUserID == userID {
		sendError(w, "You cannot thumbs up your own link", http.StatusBadRequest)
		return
	}

	if err := database.AddThumbUp(req.LinkID, userID, req.Nickname); err != nil {
		sendError(w, "Failed to add thumbs up", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]string{"status": "success"}, http.StatusOK)
}

func HandleThumbsUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	linkIDStr := r.URL.Query().Get("link_id")
	if linkIDStr == "" {
		sendError(w, "link_id is required", http.StatusBadRequest)
		return
	}

	var linkID int
	if _, err := fmt.Sscanf(linkIDStr, "%d", &linkID); err != nil {
		sendError(w, "Invalid link_id", http.StatusBadRequest)
		return
	}

	thumbs, err := database.GetThumbsUpByLink(linkID)
	if err != nil {
		sendError(w, "Failed to fetch thumbs up", http.StatusInternalServerError)
		return
	}

	sendJSON(w, thumbs, http.StatusOK)
}

func HandleUserEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nickname := r.URL.Query().Get("nickname")
	if nickname == "" {
		sendError(w, "nickname is required", http.StatusBadRequest)
		return
	}

	entries, err := database.GetEntriesByUserNickname(nickname)
	if err != nil {
		sendError(w, "Failed to fetch entries", http.StatusInternalServerError)
		return
	}

	sendJSON(w, entries, http.StatusOK)
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Nickname == "" {
		sendError(w, "nickname is required", http.StatusBadRequest)
		return
	}

	// Try to verify existing credentials first
	userID, valid, err := database.VerifyUserCredentials(req.Nickname, req.Password)
	if err != nil {
		sendError(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	if valid {
		// Valid credentials - user exists and password matches
		sendJSON(w, map[string]interface{}{
			"user_id":  userID,
			"nickname": req.Nickname,
			"message": "Login successful",
		}, http.StatusOK)
		return
	}

	// If we get here, either user doesn't exist or password doesn't match
	// Check if user exists
	var existingUserID int
	var existingPassword string
	db := database.GetDB()
	err = db.QueryRow("SELECT id, password FROM users WHERE nickname = ?", req.Nickname).Scan(&existingUserID, &existingPassword)

	if err == sql.ErrNoRows {
		// User doesn't exist - create new user with password
		userID, err = database.CreateUser(req.Nickname, req.Password)
		if err != nil {
			sendError(w, "Failed to create account", http.StatusInternalServerError)
			return
		}
		sendJSON(w, map[string]interface{}{
			"user_id":  userID,
			"nickname": req.Nickname,
			"message": "Account created successfully",
		}, http.StatusOK)
		return
	}

	if err != nil {
		sendError(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// User exists but password doesn't match
	if existingPassword == "" {
		// User has no password set - set it now
		err = database.UpdateUserPassword(existingUserID, req.Password)
		if err != nil {
			sendError(w, "Failed to set password", http.StatusInternalServerError)
			return
		}
		sendJSON(w, map[string]interface{}{
			"user_id":  existingUserID,
			"nickname": req.Nickname,
			"message": "Password set successfully",
		}, http.StatusOK)
		return
	}

	// User has a password but it doesn't match
	sendError(w, "Invalid password", http.StatusUnauthorized)
}

func HandleUpdateEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		EntryID  int             `json:"entry_id"`
		Nickname string          `json:"nickname"`
		Links    []struct {
			ID       int    `json:"id"`
			Platform string `json:"platform"`
			URL      string `json:"url"`
		} `json:"links"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.EntryID == 0 {
		sendError(w, "entry_id is required", http.StatusBadRequest)
		return
	}

	// Update entry nickname
	if err := database.UpdateEntry(req.EntryID, req.Nickname); err != nil {
		sendError(w, "Failed to update entry", http.StatusInternalServerError)
		return
	}

	// Update links
	for _, link := range req.Links {
		if link.ID > 0 {
			if link.Platform == "" || link.URL == "" {
				database.DeleteSocialLink(link.ID)
			} else {
				database.UpdateSocialLink(link.ID, link.Platform, link.URL)
			}
		} else if link.Platform != "" && link.URL != "" {
			database.CreateSocialLink(req.EntryID, link.Platform, link.URL)
		}
	}

	sendJSON(w, map[string]string{"status": "success"}, http.StatusOK)
}

func HandleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		EntryID int `json:"entry_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.EntryID == 0 {
		sendError(w, "entry_id is required", http.StatusBadRequest)
		return
	}

	if err := database.DeleteEntry(req.EntryID); err != nil {
		sendError(w, "Failed to delete entry", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]string{"status": "success"}, http.StatusOK)
}

func HandleCommunityAdmins(w http.ResponseWriter, r *http.Request) {
	// Extract community name from path: /api/community-admins/{name}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}
	communityName := pathParts[2]

	community, err := database.GetCommunityByName(communityName)
	if err != nil {
		sendError(w, "Community not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		admins, err := database.GetCommunityAdmins(community.ID)
		if err != nil {
			sendError(w, "Failed to get admins", http.StatusInternalServerError)
			return
		}
		sendJSON(w, admins, http.StatusOK)

	case "POST":
		var req struct {
			Nickname      string `json:"nickname"`
			Password      string `json:"password"`
			AdminNickname string `json:"admin_nickname"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Verify requester is admin
		requesterUserID, valid, err := database.VerifyUserCredentials(req.Nickname, req.Password)
		if err != nil || !valid {
			sendError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		isAdmin, err := database.IsCommunityAdmin(community.ID, requesterUserID)
		if err != nil || !isAdmin {
			sendError(w, "You must be an admin to add other admins", http.StatusForbidden)
			return
		}

		// Get the user to add as admin
		adminUser, err := database.GetUserByNickname(req.AdminNickname)
		if err != nil {
			sendError(w, "User not found", http.StatusNotFound)
			return
		}

		if err := database.AddCommunityAdmin(community.ID, adminUser.ID); err != nil {
			sendError(w, "Failed to add admin", http.StatusInternalServerError)
			return
		}

		sendJSON(w, map[string]string{"status": "success", "message": "Admin added successfully"}, http.StatusOK)

	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleThumbsBreakdown returns the voter breakdown for a single entry.
// GET /api/thumbs-breakdown/{community}/{entryID}
func HandleThumbsBreakdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// path: /api/thumbs-breakdown/{community}/{entryID}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		sendError(w, "Invalid path — expected /api/thumbs-breakdown/{community}/{entryID}", http.StatusBadRequest)
		return
	}

	var entryID int
	if _, err := fmt.Sscanf(pathParts[3], "%d", &entryID); err != nil {
		sendError(w, "Invalid entry_id", http.StatusBadRequest)
		return
	}

	voters, err := database.GetThumbsUpVotersByEntry(entryID)
	if err != nil {
		sendError(w, "Failed to fetch voter breakdown", http.StatusInternalServerError)
		return
	}
	if voters == nil {
		voters = []database.VoterCount{}
	}
	sendJSON(w, voters, http.StatusOK)
}

// HandleUserProfile returns a user's links in a community plus reciprocity counts.
// GET /api/user-profile?community=c&nickname=n&viewer=v
func HandleUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	community := r.URL.Query().Get("community")
	nickname := r.URL.Query().Get("nickname")
	viewer := r.URL.Query().Get("viewer")

	if community == "" || nickname == "" {
		sendError(w, "community and nickname are required", http.StatusBadRequest)
		return
	}

	profile, err := database.GetUserProfileInCommunity(community, nickname, viewer)
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, "User has no entry in this community", http.StatusNotFound)
			return
		}
		sendError(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}
	sendJSON(w, profile, http.StatusOK)
}

// HandleStats returns the global (or community-filtered) thumbs leaderboard.
// GET /api/stats?community=c
func HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	community := r.URL.Query().Get("community")

	communities, err := database.GetCommunities()
	if err != nil {
		sendError(w, "Failed to fetch communities", http.StatusInternalServerError)
		return
	}

	all, top5, err := database.GetGlobalStats(community)
	if err != nil {
		sendError(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}
	if all == nil {
		all = []database.UserStatRow{}
	}
	if top5 == nil {
		top5 = []database.UserStatRow{}
	}

	sendJSON(w, map[string]interface{}{
		"users":       all,
		"top5_green":  top5,
		"communities": communities,
		"filter":      community,
	}, http.StatusOK)
}

func HandleUpdateCommunity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CommunityName  string `json:"community_name"`
		NewName        string `json:"new_name"`
		NewDescription string `json:"new_description"`
		Nickname       string `json:"nickname"`
		Password       string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	community, err := database.GetCommunityByName(req.CommunityName)
	if err != nil {
		sendError(w, "Community not found", http.StatusNotFound)
		return
	}

	// Verify credentials
	requesterUserID, valid, err := database.VerifyUserCredentials(req.Nickname, req.Password)
	if err != nil || !valid {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if user is admin or if community has no admins yet
	hasAdmins, err := database.HasCommunityAdmins(community.ID)
	if err != nil {
		sendError(w, "Failed to check admin status", http.StatusInternalServerError)
		return
	}

	if hasAdmins {
		isAdmin, err := database.IsCommunityAdmin(community.ID, requesterUserID)
		if err != nil || !isAdmin {
			sendError(w, "You must be an admin to edit this community", http.StatusForbidden)
			return
		}
	} else {
		// No admins yet, make this user the first admin
		if err := database.AddCommunityAdmin(community.ID, requesterUserID); err != nil {
			sendError(w, "Failed to set admin", http.StatusInternalServerError)
			return
		}
	}

	// Update community
	if req.NewName != "" {
		if err := database.UpdateCommunityName(community.ID, req.NewName); err != nil {
			sendError(w, "Failed to update community name", http.StatusInternalServerError)
			return
		}
	}

	if req.NewDescription != "" {
		if err := database.UpdateCommunityDescription(community.ID, req.NewDescription); err != nil {
			sendError(w, "Failed to update community description", http.StatusInternalServerError)
			return
		}
	}

	sendJSON(w, map[string]string{"status": "success", "message": "Community updated successfully"}, http.StatusOK)
}