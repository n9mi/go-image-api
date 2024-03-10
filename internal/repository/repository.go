package repository

import "gorm.io/gorm"

type Repository[T any] struct {
}

func (r *Repository[T]) FindByID(tx *gorm.DB, e *T, ID any) error {
	return tx.First(e, "id = ?", ID).Error
}

func (r *Repository[T]) Create(tx *gorm.DB, e *T) error {
	return tx.Create(e).Error
}

func (r *Repository[T]) Update(tx *gorm.DB, e *T) error {
	return tx.Model(e).Updates(e).Error
}

func (r *Repository[T]) Delete(tx *gorm.DB, e *T) error {
	return tx.Delete(e).Error
}
