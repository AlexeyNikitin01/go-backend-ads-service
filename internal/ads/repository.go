package ads

import "context"
//go:generate mockery --output ../tests/mocks --name RepositryAd
type RepositryAd interface {
	ListAds(ctx context.Context) ([]*Ad, error)
	GetAd(ctx context.Context, adID int64) (*Ad, error)
	Add(ctx context.Context, ad *Ad) (int64, error)
	ChangeStatus(ctx context.Context, adID int64, published bool, authorID int64) (*Ad, error)
	Update(ctx context.Context, authorID int64, title string, text string, adID int64) (*Ad, error)
	Search(ctx context.Context, title string) ([]*Ad, error)
	ListAdsAuthor(ctx context.Context, author int64) ([]*Ad, error)
	ListAdsDate(ctx context.Context, day int64) ([]*Ad, error)
	DeleteAd(ctx context.Context, authorID int64, adId int64) (*Ad, error)
}
