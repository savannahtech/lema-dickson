package tasks

import (
	"log"
	"sync"
	"time"
)

func (t *AsyncTask) GetAllRepoForUser(wg *sync.WaitGroup) {
	//  logic to fetch all repositories for the given user
	// Use the GetAllRepoForUserQueue channel to send and recieve the user to and from the worker pool
	defer wg.Done()
	for user := range t.GetAllRepoForUserQueue {

		// if !ok {
		// 	log.Println("GetAllRepoForUserQueue channel closed")
		// 	return
		// }

		// todo: Fetch all repositories for the user
		userRepositories, err := t.requester.GetAllUserRepositories(user.Username)
		if err != nil {
			log.Printf("Error in fetching repositories for user %v: %v", user.Username, err)
			continue
		}
		go func() {
			// using a go routine to optimize the saving of repositories and fetching the repo  commits
			// this will help the worker process tasks from the channel faster for users at scale
			for _, newRepoInfo := range *userRepositories {
				_, err := t.dbRepository.StoreRepositoryInfo(&newRepoInfo, user)
				if err != nil {
					log.Printf("Error in storing repository: %v", err)
					continue
				}
				log.Printf("fetching repository commits for repo: %s...", newRepoInfo.Name)
				remoteCommits, err := t.requester.GetRepositoryCommits(user.Username, newRepoInfo.Name)
				if err != nil {
					log.Printf("Error in fetching commits: %v", err)
					continue
				}
				err = t.dbRepository.StoreRepositoryCommits(remoteCommits, newRepoInfo.Name, user)
				if err != nil {
					log.Printf("Error in saving commits: %v", err)
					continue
				}
			}
			log.Printf("Gotten repositories for user %v", user)
		}()

	}

}

func (t *AsyncTask) FetchNewlyRequestedRepo(wg *sync.WaitGroup) {
	//  logic to fetch a newly requested repo and commits for the given repository
	defer wg.Done()
	log.Println("waiting for newly requested repos...")

	for repoRequest := range t.FetchNewlyRequestedRepoQueue {
		log.Println("checking for newly requested repos...")

		remoteRepoInfo, err := t.requester.GetRepositoryInfo(repoRequest.Username, repoRequest.RepoName)
		if err != nil {
			continue
		}
		user, _ := t.dbRepository.GetUser(repoRequest.Username)
		repo, _ := t.dbRepository.StoreRepositoryInfo(remoteRepoInfo, user)
		go func() {
			log.Printf("fetching repository commits for repo: %s...", repoRequest.RepoName)
			remoteCommits, err := t.requester.GetRepositoryCommits(user.Username, repo.Name)
			if err != nil {
				log.Printf("Error in fetching commits: %v", err)
				return
			}
			err = t.dbRepository.StoreRepositoryCommits(remoteCommits, repo.Name, user)
			if err != nil {
				log.Printf("Error in saving commits: %v", err)
				return
			}
		}()
	}
	log.Println("exiting checking for newly requested repos...")

}

func (t *AsyncTask) CheckForUpdateOnAllRepo(wg *sync.WaitGroup) {
	//  logic to check for updates on all repositories in the database
	defer wg.Done()
	for {
		_, ok := <-t.CheckForUpdateOnAllRepoQueue
		if !ok {
			log.Println("No more signal to check for updates on all repositories")
			return
		}
		allRepos, err := t.dbRepository.GetAllRepositories()
		if err != nil {
			log.Printf("Error in fetching all repositories: %v", err)
			return
		}

		for _, repo := range allRepos {

			log.Printf("Checking for updates on repo: %s...", repo.Name)
			remoteRepoInfo, err := t.requester.GetRepositoryInfo(repo.Owner.Username, repo.Name)
			if err != nil {
				log.Printf("Error in fetching repository info: %v", err)
				continue
			}
			if repo.RemoteUpdatedAt != remoteRepoInfo.UpdatedAt {
				_, err = t.dbRepository.StoreRepositoryInfo(remoteRepoInfo, repo.Owner)
				if err != nil {
					log.Println("Error in updating repository")
				}
			}
			// simulate more processing to reduce wasting ratelimit requests
			time.Sleep(90 * time.Second)

		}
		// trigger the update again after 3days (currently passed as seconds)
		time.Sleep(3 * time.Second)
		go t.AddSignalToCheckForUpdateOnAllRepoQueue()
	}

}
