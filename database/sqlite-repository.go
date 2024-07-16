package database

import (
	"fmt"
	"log"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/models"
	"github.com/midedickson/github-service/utils"
	"gorm.io/gorm"
)

type SqliteDBRepository struct {
	DB *gorm.DB
}

func NewSqliteDBRepository(db *gorm.DB) *SqliteDBRepository {
	return &SqliteDBRepository{DB: db}
}

func (s *SqliteDBRepository) CreateUser(createUserPaylod *dto.CreateUserPayloadDTO) (*models.User, error) {
	// Create a user from payload
	existingUser, err := s.GetUser(createUserPaylod.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		// user already exists, update existing record;
		existingUser.FullName = createUserPaylod.FullName
		return existingUser, s.DB.Save(existingUser).Error
	}
	newUser := &models.User{
		Username: createUserPaylod.Username,
		FullName: createUserPaylod.FullName,
	}
	// add users into the pool to get more
	return newUser, s.DB.Create(newUser).Error
}

func (s *SqliteDBRepository) GetUser(username string) (*models.User, error) {
	// Get user by username
	var user models.User
	err := s.DB.Where("username =?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, nil
}

func (s *SqliteDBRepository) StoreRepositoryInfo(remoteRepoInfo *dto.RepositoryInfoResponseDTO, owner *models.User) (*models.Repository, error) {
	//  logic to store repository info in the database

	// check if this remote repository already exists in our database
	existingRepo, err := s.GetRepositoryInfoByRemoteId(remoteRepoInfo.ID)
	if err != nil {
		return nil, err
	}
	if existingRepo != nil {
		// repository already exists, update existing record;
		if existingRepo.RemoteUpdatedAt != remoteRepoInfo.UpdatedAt {
			// but if only there has been an update
			return existingRepo, nil
		}
		existingRepo.Name = remoteRepoInfo.Name
		existingRepo.Description = remoteRepoInfo.Description
		existingRepo.URL = remoteRepoInfo.URL
		existingRepo.Language = remoteRepoInfo.Language
		existingRepo.ForksCount = remoteRepoInfo.ForksCount
		existingRepo.StarsCount = remoteRepoInfo.StarsCount
		existingRepo.OpenIssues = remoteRepoInfo.OpenIssues
		existingRepo.Watchers = remoteRepoInfo.Watchers
		return existingRepo, s.DB.Save(existingRepo).Error
	}
	newRepo := &models.Repository{
		RemoteID:        remoteRepoInfo.ID,
		OwnerID:         owner.ID,
		Name:            remoteRepoInfo.Name,
		Description:     remoteRepoInfo.Description,
		URL:             remoteRepoInfo.HtmlUrl,
		Language:        remoteRepoInfo.Language,
		ForksCount:      remoteRepoInfo.ForksCount,
		StarsCount:      remoteRepoInfo.StarsCount,
		OpenIssues:      remoteRepoInfo.OpenIssues,
		Watchers:        remoteRepoInfo.Watchers,
		RemoteCreatedAt: remoteRepoInfo.CreatedAt,
		RemoteUpdatedAt: remoteRepoInfo.UpdatedAt,
	}
	err = s.DB.Create(newRepo).Error
	if err != nil {
		return nil, err
	}

	return newRepo, nil
}

func (s *SqliteDBRepository) GetRepositoryInfoByRemoteId(remoteID int) (*models.Repository, error) {
	//  logic to retrieve repository info from the database by remote ID
	repo := &models.Repository{}
	err := s.DB.Where("remote_id =?", remoteID).First(repo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return repo, nil
}

func (s *SqliteDBRepository) GetRepository(ownerID uint, repoName string) (*models.Repository, error) {
	//  logic to retrieve repository info from the database by ID
	repo := &models.Repository{}
	err := s.DB.Where("owner_id =?", ownerID).Where("name =?", repoName).Preload("Owner").First(repo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return repo, nil
}

func (s *SqliteDBRepository) SearchRepository(ownerID uint, repoSearchParams *utils.RepositorySearchParams) ([]*models.Repository, error) {
	//  logic to retrieve all repositories from the database
	repos := &[]*models.Repository{}
	dbQueryBuilder := s.DB.Preload("Owner").Where("owner_id =?", ownerID)
	if repoSearchParams.TopStarsCount > 0 {
		dbQueryBuilder = dbQueryBuilder.Order("stars_count DESC").Limit(repoSearchParams.TopStarsCount)
	}
	if repoSearchParams.Name != "" {
		dbQueryBuilder = dbQueryBuilder.Where("name LIKE?", "%"+repoSearchParams.Name+"%")
	}
	if repoSearchParams.Language != "" {
		dbQueryBuilder = dbQueryBuilder.Where("language =?", repoSearchParams.Language)
	}

	err := dbQueryBuilder.Find(&repos).Error
	if err != nil {
		return nil, err
	}
	return *repos, nil
}

func (s *SqliteDBRepository) GetAllRepositories() ([]*models.Repository, error) {
	//  logic to retrieve all repositories from the database
	repos := &[]*models.Repository{}
	err := s.DB.Preload("Owner").Find(&repos).Error
	if err != nil {
		return nil, err
	}
	return *repos, nil
}

func (s *SqliteDBRepository) StoreRepositoryCommits(commitRepoInfos *[]dto.CommitResponseDTO, repoName string, owner *models.User) error {
	//  logic to store commit info in the database
	repo, err := s.GetRepository(owner.ID, repoName)

	if err != nil {
		return err
	}
	if repo == nil {
		return fmt.Errorf("repository not found for owner %v and repo %v", owner.Username, repoName)
	}
	for _, commit := range *commitRepoInfos {
		// check if this commit already exists in our database
		existingCommit, err := s.GetCommitBySHA(commit.SHA)
		if err != nil {
			log.Println("Error in checking existing commits by sha")
			continue
		}
		if existingCommit != nil {
			// commit already exists, skip;
			log.Printf("Commit with SHA: %s already exists; skipping", existingCommit.SHA)
			continue
		}
		newCommit := &models.Commit{
			RepositoryName: repoName,
			SHA:            commit.SHA,
			Message:        commit.Message,
			Author:         commit.Author,
			Date:           commit.Date,
		}
		log.Printf("New commit to be created: %v", newCommit)
		err = s.DB.Create(newCommit).Error
		if err != nil {
			log.Printf("Error in saving commits with SHA: %s", newCommit.SHA)
			return err
		}
	}
	return nil
}

func (s *SqliteDBRepository) GetCommitBySHA(sha string) (*models.Commit, error) {
	commit := &models.Commit{}
	err := s.DB.Where("sha =?", sha).First(commit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return commit, nil
}

func (s *SqliteDBRepository) GetRepositoryCommits(repoName string) ([]*models.Commit, error) {
	//  logic to retrieve commit info from the database by repository name
	commits := &[]*models.Commit{}
	err := s.DB.Where("repository_name =?", repoName).Find(commits).Error
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return *commits, nil
}
