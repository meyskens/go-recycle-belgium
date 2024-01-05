package recyclebelgium

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ZipResponse struct {
	Items []struct {
		City struct {
			Names struct {
				NL string `json:"nl"`
				FR string `json:"fr"`
				DE string `json:"de"`
				EN string `json:"en"`
			} `json:"names"`
			Zipcodes  []string  `json:"zipcodes"`
			CreatedAt time.Time `json:"createdAt"`
			Name      string    `json:"name"`
			UpdatedAt time.Time `json:"updatedAt"`
			ID        string    `json:"id"`
		} `json:"city"`
		Code      string    `json:"code"`
		CreatedAt time.Time `json:"createdAt"`
		Names     []struct {
			NL string `json:"nl"`
			FR string `json:"fr"`
			DE string `json:"de"`
			EN string `json:"en"`
		} `json:"names"`
		UpdatedAt time.Time `json:"updatedAt"`
		ID        string    `json:"id"`
		Available bool      `json:"available"`
	} `json:"items"`
	Total int    `json:"total"`
	Pages int    `json:"pages"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Self  string `json:"self"`
	First string `json:"first"`
	Last  string `json:"last"`
}

func (a *API) GetZipCodes(zipcode string) (ZipResponse, error) {
	auth, err := a.getToken()
	if err != nil {
		return ZipResponse{}, fmt.Errorf("error getting auth token: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.fostplus.be/recycle-public/app/v1/zipcodes?q=%s", url.QueryEscape(zipcode)), nil)
	if err != nil {
		return ZipResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Consumer", "recycleapp.be")
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ZipResponse{}, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return ZipResponse{}, fmt.Errorf("got a %d error: %s", resp.StatusCode, string(b))
	}

	// fetch body
	var body ZipResponse

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return body, fmt.Errorf("error decoding response: %w", err)
	}

	return body, nil
}

type StreetResponse struct {
	Items []struct {
		ID   string `json:"id"`
		City []struct {
			Names struct {
				NL string `json:"nl"`
				FR string `json:"fr"`
				DE string `json:"de"`
				EN string `json:"en"`
			} `json:"names"`
			Zipcodes  []string  `json:"zipcodes"`
			CreatedAt time.Time `json:"createdAt"`
			Name      string    `json:"name"`
			UpdatedAt time.Time `json:"updatedAt"`
			ID        string    `json:"id"`
		} `json:"city"`
		CreatedAt time.Time `json:"createdAt"`
		Deleted   bool      `json:"deleted"`
		Name      string    `json:"name"`
		Names     struct {
			NL string `json:"nl"`
			FR string `json:"fr"`
			DE string `json:"de"`
			EN string `json:"en"`
		} `json:"names"`
		UpdatedAt time.Time `json:"updatedAt"`
		Zipcode   []struct {
			City      string    `json:"city"`
			Code      string    `json:"code"`
			CreatedAt time.Time `json:"createdAt"`
			Names     []struct {
				NL string `json:"nl"`
				FR string `json:"fr"`
				DE string `json:"de"`
				EN string `json:"en"`
			} `json:"names"`
			UpdatedAt time.Time `json:"updatedAt"`
			ID        string    `json:"id"`
		} `json:"zipcode"`
	} `json:"items"`
	Total int    `json:"total"`
	Pages int    `json:"pages"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Self  string `json:"self"`
	First string `json:"first"`
	Last  string `json:"last"`
}

func (a *API) GetStreets(zipcodeID string, street string) (StreetResponse, error) {
	auth, err := a.getToken()
	if err != nil {
		return StreetResponse{}, fmt.Errorf("error getting auth token: %w", err)
	}

	// for reasons unknown this is a POST api with GET parameters
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.fostplus.be/recycle-public/app/v1/streets?q=%s&zipcodes=%s", url.QueryEscape(street), zipcodeID), strings.NewReader("{}"))
	if err != nil {
		return StreetResponse{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Consumer", "recycleapp.be")
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return StreetResponse{}, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return StreetResponse{}, fmt.Errorf("got a %d error: %s", resp.StatusCode, string(b))
	}

	// fetch body
	var body StreetResponse
	err = json.NewDecoder(resp.Body).Decode(&body)

	if err != nil {
		return body, fmt.Errorf("error decoding response: %w", err)
	}

	return body, nil
}
