package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// --- Identifiants Spotify ---
const (
	ClientID     = "a6b4939a7c644ce29650ec3de486c299"
	ClientSecret = "789bc968246240d98a4fa79bdb049c70"
)

// --- Structures de réponse Spotify ---
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type Image struct {
	Url string `json:"url"`
}
type ExternalUrls struct {
	Spotify string `json:"spotify"`
}
type Artist struct {
	Name string `json:"name"`
}
type Album struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ReleaseDate string  `json:"release_date"`
	TotalTracks int     `json:"total_tracks"`
	Images      []Image `json:"images"`
}
type Track struct {
	Name         string       `json:"name"`
	Album        Album        `json:"album"`
	Artists      []Artist     `json:"artists"`
	ExternalUrls ExternalUrls `json:"external_urls"`
}
type AlbumItems struct {
	Items []Album `json:"items"`
}
type TracksItems struct {
	Items []Track `json:"items"`
}
type SearchTracks struct {
	Tracks struct {
		Items []Track `json:"items"`
	} `json:"tracks"`
}

// --- Gestion du token avec cache ---
var cachedToken string
var tokenExpiry time.Time

func getAutoToken() (string, error) {
	if cachedToken != "" && time.Now().Before(tokenExpiry) {
		return cachedToken, nil
	}

	authURL := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(ClientID + ":" + ClientSecret))
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("spotify token error: %s", string(bodyBytes))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return "", err
	}
	cachedToken = tokenResp.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return cachedToken, nil
}

// --- Templates ---
func renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	tmpl := template.Must(template.ParseFiles("template/" + filename))
	_ = tmpl.Execute(w, data)
}

// --- Home ---
func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

// --- Damso: Batterie Faible ---
func DamsoAlbum(w http.ResponseWriter, r *http.Request) {
	token, err := getAutoToken()
	if err != nil {
		http.Error(w, "Impossible d'obtenir le token Spotify", http.StatusInternalServerError)
		return
	}

	// Récupérer les albums de Damso
	albumsURL := "https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie/albums?include_groups=album&limit=50"
	albums, err := getAlbums(albumsURL, token)
	if err != nil {
		http.Error(w, "Erreur Spotify albums", http.StatusInternalServerError)
		return
	}

	// Trouver Batterie Faible
	var bf Album
	for _, a := range albums {
		if strings.EqualFold(a.Name, "Batterie Faible") {
			bf = a
			break
		}
	}
	if bf.ID == "" {
		http.Error(w, "Album 'Batterie Faible' introuvable", http.StatusNotFound)
		return
	}

	// Récupérer les morceaux de l’album
	tracksURL := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks?limit=50", bf.ID)
	tracks, err := getTracks(tracksURL, token)
	if err != nil {
		http.Error(w, "Erreur Spotify tracks", http.StatusInternalServerError)
		return
	}

	// Données envoyées au template
	type DamsoPageData struct {
		AlbumName   string
		AlbumCover  string
		ReleaseDate string
		Tracks      []Track
	}
	page := DamsoPageData{
		AlbumName:   bf.Name,
		AlbumCover:  firstImage(bf.Images),
		ReleaseDate: bf.ReleaseDate,
		Tracks:      tracks,
	}

	renderTemplate(w, "damso.html", page)
}

// --- Laylow: pick one track with details ---
/*func LaylowSong(w http.ResponseWriter, r *http.Request) {
	token, err := getAutoToken()
	if err != nil {
		http.Error(w, "Impossible d'obtenir le token Spotify", http.StatusInternalServerError)
		return
	}

	// Recherche d’un morceau de Laylow (ici Maladresse)
	searchURL := "https://api.spotify.com/v1/search?q=artist:Laylow%20track:Maladresse&type=track&limit=1"
	track, err := searchOneTrack(searchURL, token)
	if err != nil {
		http.Error(w, "Aucun morceau Laylow trouvé", http.StatusNotFound)
		return
	}

	renderTemplate(w, "laylow.html", track)
}*/
func LaylowSong(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LaylowSong handler triggered")

	token, err := getAutoToken()
	if err != nil {
		fmt.Println("Erreur token:", err)
		http.Error(w, "Impossible d'obtenir le token Spotify", http.StatusInternalServerError)
		return
	}

	searchURL := "https://api.spotify.com/v1/search?q=artist:Laylow%20track:Maladresse&type=track&limit=1"
	fmt.Println("Search URL:", searchURL)

	track, err := searchOneTrack(searchURL, token)
	if err != nil {
		fmt.Println("Erreur track:", err)
		http.Error(w, "Aucun morceau Laylow trouvé", http.StatusNotFound)
		return
	}

	fmt.Println("Track trouvé:", track.Name)
	renderTemplate(w, "laylow.html", track)
}

// --- Helpers ---
func getAlbums(apiURL, token string) ([]Album, error) {
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data AlbumItems
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}
	return data.Items, nil
}

func getTracks(apiURL, token string) ([]Track, error) {
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data TracksItems
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}

	return data.Items, nil
}

func searchOneTrack(apiURL, token string) (Track, error) {
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Track{}, err
	}
	defer resp.Body.Close()

	var data SearchTracks
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return Track{}, err
	}
	if len(data.Tracks.Items) == 0 {
		return Track{}, fmt.Errorf("no track found")
	}
	return data.Tracks.Items[0], nil
}

func firstImage(images []Image) string {
	if len(images) == 0 {
		return "/image/placeholder.png"
	}
	return images[0].Url
}
