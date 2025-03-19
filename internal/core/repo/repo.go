//go:generate mockgen -source=repo.go -mock_names=Repo=MockRepository -destination=../../mock/mock_repository.go -package=mock

package repo

import (
	"context"
	"go-fiber-api/internal/core/storage/db"
	"math"

	"gorm.io/gorm"
)

// REF: https://www.ompluscator.com/article/golang/tutorial-generics-gorm/?source=post_page-----7f8891df0934--------------------------------
type GormModel[D any] interface {
	ToDTO() D
	FromDTO(dto D) any
}

type Repo[E GormModel[D], D any] interface {
	FindAll(ctx context.Context) ([]D, error)
	Find(ctx context.Context, specifications ...Specification) ([]D, error)
	FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]D, error)
	FindWithPagination(ctx context.Context, page, limit int, specifications ...Specification) ([]D, PaginationMetadata, error)
	Count(ctx context.Context, specifications ...Specification) (i int64, err error)
	FindByID(ctx context.Context, id any) (D, error)
	Insert(ctx context.Context, dto *D) error
	Update(ctx context.Context, dto *D) error
	Delete(ctx context.Context, dto *D) error
	DeleteById(ctx context.Context, id any) error
}

type repoImpl[E GormModel[D], D any] struct {
	db db.Client
}

func NewRepository[E GormModel[D], D any](db db.Client) Repo[E, D] {
	return &repoImpl[E, D]{
		db: db,
	}
}

func (r *repoImpl[E, D]) getPreWarmDbForSelect(ctx context.Context, specification ...Specification) *gorm.DB {
	dbPrewarm := r.db.WithContext(ctx)
	for _, s := range specification {
		dbPrewarm = dbPrewarm.Where(s.GetQuery(), s.GetValues()...)
	}
	return dbPrewarm
}

func (r *repoImpl[E, D]) FindAll(ctx context.Context) ([]D, error) {
	return r.FindWithLimit(ctx, -1, -1)
}

func (r *repoImpl[E, D]) Find(ctx context.Context, specifications ...Specification) ([]D, error) {
	return r.FindWithLimit(ctx, -1, -1, specifications...)
}

func (r *repoImpl[E, D]) FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]D, error) {
	var entities []E

	dbPrewarm := r.getPreWarmDbForSelect(ctx, specifications...)
	err := dbPrewarm.Limit(limit).Offset(offset).Order("id").Find(&entities).Error

	if err != nil {
		return nil, err
	}

	result := make([]D, 0, len(entities))
	for _, row := range entities {
		result = append(result, row.ToDTO())
	}

	return result, nil
}

func (r *repoImpl[E, D]) FindWithPagination(ctx context.Context, page, limit int, specifications ...Specification) ([]D, PaginationMetadata, error) {
	offset := (page - 1) * limit

	res, err := r.FindWithLimit(ctx, limit, offset, specifications...)
	if err != nil {
		return nil, PaginationMetadata{}, err
	}

	totalCount, err := r.Count(ctx, specifications...)
	if err != nil {
		return nil, PaginationMetadata{}, err
	}

	totalPages := math.Ceil(float64(totalCount) / float64(limit))

	metadata := PaginationMetadata{
		Page:       uint(page),
		PerPage:    uint(limit),
		TotalPages: uint(totalPages),
		TotalItems: uint(totalCount),
	}

	return res, metadata, err
}

func (r *repoImpl[E, D]) Count(ctx context.Context, specifications ...Specification) (i int64, err error) {
	entity := new(E)
	err = r.getPreWarmDbForSelect(ctx, specifications...).Model(entity).Count(&i).Error
	return
}

func (r *repoImpl[E, D]) FindByID(ctx context.Context, id any) (D, error) {
	var entity E
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return entity.ToDTO(), err
	}

	return entity.ToDTO(), nil
}

func (r *repoImpl[E, D]) Insert(ctx context.Context, dto *D) error {
	var entity E
	dao := entity.FromDTO(*dto).(E)

	err := r.db.WithContext(ctx).Create(&dao).Error
	if err != nil {
		return err
	}

	*dto = dao.ToDTO()
	return nil
}

func (r *repoImpl[E, D]) Update(ctx context.Context, dto *D) error {
	var entity E
	model := entity.FromDTO(*dto).(E)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return err
	}

	*dto = model.ToDTO()
	return nil
}

func (r *repoImpl[E, D]) Delete(ctx context.Context, dto *D) error {
	var entity E
	model := entity.FromDTO(*dto).(E)
	err := r.db.WithContext(ctx).Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repoImpl[E, D]) DeleteById(ctx context.Context, id any) error {
	var entity E
	err := r.db.WithContext(ctx).Delete(&entity, &id).Error
	if err != nil {
		return err
	}

	return nil
}
