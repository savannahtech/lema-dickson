package requester

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/utils"
)

type RepositoryRequester struct {
	http.Client
	rateLimit          int
	rateLimitRemaining int
	rateLimitReset     time.Time
}

func NewRepositoryRequester() *RepositoryRequester {
	return &RepositoryRequester{}
}

// handling rate limit
func (r *RepositoryRequester) checkRateLimit(resp *http.Response) {
	if limit := resp.Header.Get("x-ratelimit-limit"); limit != "" {
		r.rateLimit, _ = strconv.Atoi(limit)
	}
	if remaining := resp.Header.Get("x-ratelimit-remaining"); remaining != "" {
		r.rateLimitRemaining, _ = strconv.Atoi(remaining)
	}
	if reset := resp.Header.Get("x-ratelimit-reset"); reset != "" {
		resetTime, _ := strconv.ParseInt(reset, 10, 64)
		r.rateLimitReset = time.Unix(resetTime, 0)
	}
	log.Printf("Rate limit: %d, Remaining: %d, Reset: %v", r.rateLimit, r.rateLimitRemaining, r.rateLimitReset)

}

func (r *RepositoryRequester) waitForRateLimitReset() {
	if r.rateLimitRemaining == 0 && time.Now().Before(r.rateLimitReset) {
		log.Println("Waiting for rate limit reset")
		time.Sleep(time.Until(r.rateLimitReset))
	}
}

func (r *RepositoryRequester) doRequest(req *http.Request) (*http.Response, error) {
	r.waitForRateLimitReset()

	resp, err := r.Do(req)
	if err != nil {
		log.Printf("Error whilke making request: %v", err)
		return nil, err
	}

	r.checkRateLimit(resp)

	if resp.StatusCode == http.StatusForbidden && resp.Header.Get("x-ratelimit-remaining") == "0" {
		r.waitForRateLimitReset()
		resp, err = r.Do(req)
		if err != nil {
			return nil, err
		}
		r.checkRateLimit(resp)
	}

	return resp, nil
}

func (r *RepositoryRequester) fetchAndDecode(url string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := r.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return utils.ErrRepoNotFound
	}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}

func (r *RepositoryRequester) GetRepositoryInfo(owner, repo string) (*dto.RepositoryInfoResponseDTO, error) {
	// fetch repository info for owner
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	var repository dto.RepositoryInfoResponseDTO
	if err := r.fetchAndDecode(url, &repository); err != nil {
		return nil, err
	}
	return &repository, nil
}
func (r *RepositoryRequester) GetRepositoryCommits(owner, repo string) (*[]dto.CommitResponseDTO, error) {
	// logic to fetch repository commits
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	var commits []dto.CommitResponseDTO

	if err := r.fetchAndDecode(url, &commits); err != nil {
		return nil, err
	}
	return &commits, nil
}

func (r *RepositoryRequester) GetAllUserRepositories(owner string) (*[]dto.RepositoryInfoResponseDTO, error) {
	//  logic to fetch all repositories for a user
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", owner)
	var repositories []dto.RepositoryInfoResponseDTO
	if err := r.fetchAndDecode(url, &repositories); err != nil {
		return nil, err
	}
	return &repositories, nil
}
