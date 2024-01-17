package album

import "context"

// AlbumService defines the methods for handling CRUD operations on albums
type AlbumService interface {
	CreateAlbum(ctx context.Context, album Album) error
	GetAlbum(ctx context.Context, albumID string) (Album, error)
	UpdateAlbum(ctx context.Context, albumID string, update Album) error
	DeleteAlbum(ctx context.Context, albumID string) error
	GetSongsWithCapitalTitles(ctx context.Context) ([]string, error)
    SearchAlbums(ctx context.Context, searchTerm string) ([]Album, error)

}

// albumService implements AlbumService
type albumService struct {
	repository AlbumRepository
}

// NewAlbumService creates a new instance of albumService
func NewAlbumService(repo AlbumRepository) AlbumService {
	return &albumService{repository: repo}
}

// CreateAlbum adds a new album
func (s *albumService) CreateAlbum(ctx context.Context, album Album) error {
	return s.repository.CreateAlbum(ctx, album)
}

// GetAlbum retrieves an album by ID
func (s *albumService) GetAlbum(ctx context.Context, albumID string) (Album, error) {
	return s.repository.GetAlbum(ctx, albumID)
}

// UpdateAlbum updates an existing album
func (s *albumService) UpdateAlbum(ctx context.Context, albumID string, update Album) error {
	return s.repository.UpdateAlbum(ctx, albumID, update)
}

// DeleteAlbum removes an album by ID
func (s *albumService) DeleteAlbum(ctx context.Context, albumID string) error {
	return s.repository.DeleteAlbum(ctx, albumID)
}
func (s *albumService) GetSongsWithCapitalTitles(ctx context.Context) ([]string, error) {
    return s.repository.GetSongsWithCapitalTitles(ctx)
}
func (s *albumService) SearchAlbums(ctx context.Context, searchTerm string) ([]Album, error) {
    return s.repository.SearchAlbums(ctx, searchTerm)
}