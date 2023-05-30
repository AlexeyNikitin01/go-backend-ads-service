package adrepo

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"ads/internal/ads"
)

type keyID int64
type adStructType = *ads.Ad

type AdRepositoryMap struct {
	sync.Mutex
	countID int64
	mapRep map[keyID]adStructType
}

func (r *AdRepositoryMap) Add(ctx context.Context, ad *ads.Ad) (int64, error) {
	r.countID += 1
	ad.ID = r.countID
	r.mapRep[keyID(r.countID)] = ad

	return r.countID, nil
}

func (r *AdRepositoryMap) ChangeStatus(ctx context.Context, adID int64, published bool, authorID int64) (*ads.Ad, error) {
	ad, ok := r.mapRep[keyID(adID)]

	if !ok {
		return nil, fmt.Errorf("is no such ad")
	}
	ad.UpdateDate = time.Now().UTC()
	ad.Published = published

	return ad, nil
}

func (r *AdRepositoryMap) Update(ctx context.Context, authorID int64, title string, text string, adID int64) (*ads.Ad, error) {
	ad, ok := r.mapRep[keyID(adID)]
	
	if !ok {
		return nil, fmt.Errorf("is no such ad")
	}

	ad.UpdateDate = time.Now().UTC()
	ad.Title = title
	ad.Text = text

	return ad, nil
}

func (r *AdRepositoryMap) GetAd(ctx context.Context, adID int64) (*ads.Ad, error) {
	ad, ok := r.mapRep[keyID(adID)]
	if !ok {
		return nil, fmt.Errorf("is no such ad")
	}
	return ad, nil
}

func (r *AdRepositoryMap) ListAds(ctx context.Context) ([]*ads.Ad, error) {
	if r.mapRep == nil {
		return nil, fmt.Errorf("not map repository")
	}
	result := []*ads.Ad{}
	for _, ad := range r.mapRep {
		if ad.Published {
			result = append(result, ad)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("not found ad")
	}

	return result, nil
}

func(r *AdRepositoryMap) Search(ctx context.Context, title string) ([]*ads.Ad, error) {
	var result []*ads.Ad
	for _, i := range r.mapRep {
		if strings.HasPrefix(i.Title, title) {
			result = append(result, i)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("not found ad")
	}
	return result, nil
}

func (r *AdRepositoryMap) ListAdsAuthor(ctx context.Context, author int64) ([]*ads.Ad, error) {
	if r.mapRep == nil {
		return nil, fmt.Errorf("not map repository")
	}
	result := []*ads.Ad{}
	for _, ad := range r.mapRep {
		if ad.AuthorID == author {
			result = append(result, ad)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("not found ad")
	}
	return result, nil
}

func (r *AdRepositoryMap) ListAdsDate(ctx context.Context, day int64) ([]*ads.Ad, error) {
	if r.mapRep == nil {
		return nil, fmt.Errorf("not map repository")
	}
	result := []*ads.Ad{}
	for _, ad := range r.mapRep {
		if int64(ad.CreateDate.Day()) == day {
			result = append(result, ad)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("not found ad")
	}
	return result, nil
}

func (r *AdRepositoryMap) DeleteAd(ctx context.Context, authorID int64, adID int64) (*ads.Ad, error) {
	ad, ok := r.mapRep[keyID(adID)]
	if !ok {
		return nil, fmt.Errorf("not delete")
	}
	fmt.Println(ad.AuthorID, authorID)
	if ad.AuthorID == authorID {
		delete(r.mapRep, keyID(adID))
		return ad, nil
	}

	return nil, fmt.Errorf("didn`t delete")
}

func New() ads.RepositryAd {
	return &AdRepositoryMap{
		countID: -1,
		mapRep: make(map[keyID]adStructType),
		}
}
