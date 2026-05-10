package services

import (
	"movie-night-planner-backend/internal/models"
	"movie-night-planner-backend/internal/repositories"
	"movie-night-planner-backend/internal/utils"
	"movie-night-planner-backend/pkg/response"

	"github.com/google/uuid"
)

type VoteService struct {
	voteRepo        *repositories.VoteRepository
	eveningRepo     *repositories.EveningRepository
	eveningFilmRepo *repositories.EveningFilmRepository
	userRepo        *repositories.UserRepository
}

type CreateVoteInput struct {
	EveningFilmID uuid.UUID `json:"evening_film_id" validate:"required"`
	Value         int       `json:"value" validate:"required,min=1,max=5"`
}

func NewVoteService(
	voteRepo *repositories.VoteRepository,
	eveningRepo *repositories.EveningRepository,
	eveningFilmRepo *repositories.EveningFilmRepository,
	userRepo *repositories.UserRepository,
) *VoteService {
	return &VoteService{
		voteRepo:        voteRepo,
		eveningRepo:     eveningRepo,
		eveningFilmRepo: eveningFilmRepo,
		userRepo:        userRepo,
	}
}

func (s *VoteService) Create(eveningID uuid.UUID, input CreateVoteInput, userID uuid.UUID) (*models.Vote, error) {
	// Verify evening exists
	_, err := s.eveningRepo.FindByID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	// Verify evening film exists
	_, err = s.eveningFilmRepo.FindByID(input.EveningFilmID)
	if err != nil {
		return nil, utils.NewAppError("FILM_NOT_FOUND", "Film not found", nil, nil)
	}

	// Check if user already voted for this film
	_, err = s.voteRepo.FindByUserIDEveningIDAndFilmID(userID, eveningID, input.EveningFilmID)
	if err == nil {
		return nil, utils.NewAppError("ALREADY_VOTED", "You have already voted for this film", nil, nil)
	}

	// Create vote
	vote := &models.Vote{
		EveningID:     eveningID,
		EveningFilmID: input.EveningFilmID,
		UserID:        userID,
		Value:         input.Value,
	}

	err = s.voteRepo.Create(vote)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to create vote")
	}

	return vote, nil
}

func (s *VoteService) GetVotesForEvening(eveningID uuid.UUID) ([]response.VoteSummary, error) {
	// Verify evening exists
	_, err := s.eveningRepo.FindByID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "EVENING_NOT_FOUND", "Evening not found")
	}

	votes, err := s.voteRepo.FindByEveningID(eveningID)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to get votes")
	}

	// Group votes by film
	filmVotes := make(map[uuid.UUID][]models.Vote)
	for _, vote := range votes {
		filmVotes[vote.EveningFilmID] = append(filmVotes[vote.EveningFilmID], vote)
	}

	// Calculate summaries
	summaries := make([]response.VoteSummary, 0, len(filmVotes))
	for filmID, filmVoteList := range filmVotes {
		film, err := s.eveningFilmRepo.FindByID(filmID)
		if err != nil {
			continue
		}

		totalVotes := len(filmVoteList)
		sum := 0
		distribution := make(map[int]int)

		for _, v := range filmVoteList {
			sum += v.Value
			distribution[v.Value]++
		}

		avg := float64(sum) / float64(totalVotes)

		summaries = append(summaries, response.VoteSummary{
			EveningFilmID: film.ID,
			Title:         film.Title,
			PosterPath:    film.PosterPath,
			TotalVotes:    totalVotes,
			AverageScore:  avg,
			VoteDistribution: response.VoteDistribution{
				One:   distribution[1],
				Two:   distribution[2],
				Three: distribution[3],
				Four:  distribution[4],
				Five:  distribution[5],
			},
		})
	}

	return summaries, nil
}

func (s *VoteService) Update(voteID uuid.UUID, value int, userID uuid.UUID) (*models.Vote, error) {
	vote, err := s.voteRepo.FindByID(voteID)
	if err != nil {
		return nil, utils.WrapError(err, "VOTE_NOT_FOUND", "Vote not found")
	}

	// Verify ownership
	if vote.UserID != userID {
		return nil, utils.NewAppError("FORBIDDEN", "You can only update your own votes", nil, nil)
	}

	vote.Value = value
	err = s.voteRepo.Update(vote)
	if err != nil {
		return nil, utils.WrapError(err, "DATABASE_ERROR", "Failed to update vote")
	}

	return vote, nil
}

func (s *VoteService) Delete(voteID uuid.UUID, userID uuid.UUID) error {
	vote, err := s.voteRepo.FindByID(voteID)
	if err != nil {
		return utils.WrapError(err, "VOTE_NOT_FOUND", "Vote not found")
	}

	// Verify ownership
	if vote.UserID != userID {
		return utils.NewAppError("FORBIDDEN", "You can only delete your own votes", nil, nil)
	}

	err = s.voteRepo.Delete(voteID)
	if err != nil {
		return utils.WrapError(err, "DATABASE_ERROR", "Failed to delete vote")
	}

	return nil
}
