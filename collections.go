package recyclebelgium

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CollectionsResponse struct {
	Items []struct {
		Type      string    `json:"type"`
		Timestamp time.Time `json:"timestamp"`
		Fraction  struct {
			National    bool        `json:"national"`
			NationalRef interface{} `json:"nationalRef"`
			DatatankRef interface{} `json:"datatankRef"`
			Name        struct {
				NL string `json:"nl"`
				FR string `json:"fr"`
				EN string `json:"en"`
				DE string `json:"de"`
			} `json:"name"`
			Logo struct {
				Regular struct {
					OneX   string `json:"1x"`
					TwoX   string `json:"2x"`
					ThreeX string `json:"3x"`
				} `json:"regular"`
				Reversed struct {
					OneX   string `json:"1x"`
					TwoX   string `json:"2x"`
					ThreeX string `json:"3x"`
				} `json:"reversed"`
				Name struct {
					NL string `json:"nl"`
					FR string `json:"fr"`
					DE string `json:"de"`
					EN string `json:"en"`
				} `json:"name"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
				V         int       `json:"__v"`
				ID        string    `json:"id"`
			} `json:"logo"`
			Color        string        `json:"color"`
			Variations   []interface{} `json:"variations"`
			Organisation string        `json:"organisation"`
			CreatedAt    time.Time     `json:"createdAt"`
			UpdatedAt    time.Time     `json:"updatedAt"`
			ID           string        `json:"id"`
		} `json:"fraction"`
		ID string `json:"id"`
	} `json:"items"`
	Total int    `json:"total"`
	Pages int    `json:"pages"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Self  string `json:"self"`
	First string `json:"first"`
	Last  string `json:"last"`
}

func (a *API) GetCollections(zipcodeID, streetID, houseNumber string, from, until time.Time, size int) (CollectionsResponse, error) {
	auth, err := a.getToken()
	if err != nil {
		return CollectionsResponse{}, fmt.Errorf("error getting auth token: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.fostplus.be/recycle-public/app/v1/collections?zipcodeId=%s&streetId=%s&houseNumber=%s&fromDate=%s&untilDate=%s&size=%d", zipcodeID, streetID, houseNumber, from.Format("2006-01-02"), until.Format("2006-01-02"), size), nil)
	if err != nil {
		return CollectionsResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Consumer", "recycleapp.be")
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CollectionsResponse{}, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return CollectionsResponse{}, fmt.Errorf("got a %d error: %s", resp.StatusCode, string(b))
	}

	// fetch body
	var body CollectionsResponse

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return body, fmt.Errorf("error decoding response: %w", err)
	}

	return body, nil
}
