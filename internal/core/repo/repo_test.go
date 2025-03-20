package repo_test

import (
	"context"
	repository "go-fiber-api/internal/core/repo"
	database "go-fiber-api/internal/core/storage/db"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Product is a domain entity
type ProductDTO struct {
	ID          uint
	Name        string
	Weight      uint
	IsAvailable bool
}

// Product is DTO used to map Product entity to database
type Product struct {
	ID          uint   `gorm:"primaryKey;column:id"`
	Name        string `gorm:"column:name"`
	Weight      uint   `gorm:"column:weight"`
	IsAvailable bool   `gorm:"column:is_available"`
}

func (g Product) ToDTO() ProductDTO {
	return ProductDTO(g)
}

func (g Product) FromDTO(dto ProductDTO) any {
	return Product(dto)
}

func getDB() (database.Client, error) {
	g, err := gorm.Open(sqlite.Open("file:test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	client := g

	return client, err
}
func TestMain(m *testing.M) {
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(Product{})
	if err != nil {
		log.Fatal(err)
	}
	ret := m.Run()
	os.Exit(ret)
}
func TestGormRepository_Insert(t *testing.T) {
	client, _ := getDB()
	repo := repository.NewRepository[Product, ProductDTO](client)
	ctx := context.Background()

	product := ProductDTO{
		ID:          8,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err := repo.Insert(ctx, &product)
	if err != nil {
		panic(err)
	}

	p, err := repo.FindByID(ctx, product.ID)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, product, p)
}

func TestGormRepository_FindByID(t *testing.T) {
	client, _ := getDB()
	repo := repository.NewRepository[Product, ProductDTO](client)
	ctx := context.Background()

	expected := Product{
		ID:          8,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}

	actual, err := repo.FindByID(ctx, 8)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, expected, actual)
}

func TestGormRepository_Count(t *testing.T) {
	client, _ := getDB()
	repo := repository.NewRepository[Product, ProductDTO](client)
	ctx := context.Background()

	nb, err := repo.Count(ctx)

	if err != nil {
		panic(err)
	}
	if nb != 1 {
		panic("not good count")
	}
}

func TestGormRepository_DeleteByID(t *testing.T) {
	client, _ := getDB()
	repository := repository.NewRepository[Product, ProductDTO](client)
	ctx := context.Background()
	err := repository.DeleteById(ctx, 8)
	if err != nil {
		panic(err)
	}
	_, err = repository.FindByID(ctx, 8)
	if err == nil {
		panic("supposed to be deleted")
	}
}

func TestGormRepository_Find(t *testing.T) {
	client, _ := getDB()
	repo := repository.NewRepository[Product, ProductDTO](client)
	ctx := context.Background()

	product := ProductDTO{
		ID:          1,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	repo.Insert(ctx, &product)
	product2 := ProductDTO{
		ID:          2,
		Name:        "product2",
		Weight:      50,
		IsAvailable: true,
	}
	repo.Insert(ctx, &product2)
	many, err := repo.Find(ctx, repository.GreaterOrEqual("weight", 50))
	if err != nil {
		panic(err)
	}
	if len(many) != 2 {
		panic("should be 2")
	}

	firstActual := []ProductDTO{product, product2}
	assert.Equal(t, many, firstActual)

	product3 := ProductDTO{
		ID:          3,
		Name:        "product3",
		Weight:      250,
		IsAvailable: false,
	}
	repo.Insert(ctx, &product3)

	many, err = repo.Find(ctx, repository.GreaterOrEqual("weight", 90))
	if err != nil {
		panic(err)
	}
	if len(many) != 2 {
		panic("should be 2")
	}
	secondActual := []ProductDTO{product, product3}
	assert.Equal(t, many, secondActual)

	many, err = repo.Find(ctx, repository.And(
		repository.GreaterOrEqual("weight", 90),
		repository.Equal("is_available", true)),
	)
	if err != nil {
		panic(err)
	}
	if len(many) != 1 {
		panic("should be 1")
	}

	thirdActual := []ProductDTO{product}
	assert.Equal(t, thirdActual, many)
}

/*
TODO
Delete (by item)
Update
Find (with sql cond)
*/
