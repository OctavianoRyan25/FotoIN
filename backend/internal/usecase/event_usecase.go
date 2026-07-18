package usecase

import (
	"context"
	"errors"
	"fmt"
	"fotoin/backend/internal/domain"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, photographerID uuid.UUID, title, location, description string, eventDate *time.Time) (*domain.Event, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetEventBySlug(ctx context.Context, slug string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, id, photographerID uuid.UUID, title, location, description string, eventDate *time.Time, isActive bool) error
	DeleteEvent(ctx context.Context, id, photographerID uuid.UUID) error
	ListEventsByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.Event, int64, error)
	ListPublicEvents(ctx context.Context, page, limit int, search string) ([]domain.Event, int64, error)
}

type eventUsecase struct {
	eventRepo domain.EventRepository
}

func NewEventUsecase(eventRepo domain.EventRepository) EventUsecase {
	return &eventUsecase{eventRepo: eventRepo}
}

func (uc *eventUsecase) CreateEvent(ctx context.Context, photographerID uuid.UUID, title, location, description string, eventDate *time.Time) (*domain.Event, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	slug := generateSlug(title)
	existing, _ := uc.eventRepo.FindBySlug(ctx, slug)
	if existing != nil {
		slug = fmt.Sprintf("%s-%s", slug, uuid.New().String()[:8])
	}

	event := &domain.Event{
		PhotographerID: photographerID,
		Title:          title,
		Slug:           slug,
		Location:       location,
		Description:    description,
		EventDate:      eventDate,
		IsActive:       true,
	}

	if err := uc.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *eventUsecase) GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return uc.eventRepo.FindByID(ctx, id)
}

func (uc *eventUsecase) GetEventBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	return uc.eventRepo.FindBySlug(ctx, slug)
}

func (uc *eventUsecase) UpdateEvent(ctx context.Context, id, photographerID uuid.UUID, title, location, description string, eventDate *time.Time, isActive bool) error {
	event, err := uc.eventRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if event.PhotographerID != photographerID {
		return errors.New("unauthorized")
	}

	if title != "" {
		event.Title = title
	}
	event.Location = location
	event.Description = description
	event.EventDate = eventDate
	event.IsActive = isActive

	return uc.eventRepo.Update(ctx, event)
}

func (uc *eventUsecase) DeleteEvent(ctx context.Context, id, photographerID uuid.UUID) error {
	event, err := uc.eventRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if event.PhotographerID != photographerID {
		return errors.New("unauthorized")
	}

	return uc.eventRepo.Delete(ctx, id)
}

func (uc *eventUsecase) ListEventsByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.Event, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.eventRepo.ListByPhotographer(ctx, photographerID, page, limit)
}

func (uc *eventUsecase) ListPublicEvents(ctx context.Context, page, limit int, search string) ([]domain.Event, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.eventRepo.ListPublic(ctx, page, limit, search)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
