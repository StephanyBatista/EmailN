package campaign

type Repository interface {
	Save(campaign *Campaign) error
}
