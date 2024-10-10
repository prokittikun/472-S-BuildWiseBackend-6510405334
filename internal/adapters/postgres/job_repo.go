// In postgres/job_repository.go
package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type jobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) repositories.JobRepository {
	return &jobRepository{
		db: db,
	}
}

func (jr *jobRepository) CreateJob(ctx context.Context, job *models.Job, materials []models.JobMaterial) error {
	tx, err := jr.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// สร้าง Job
	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO Job (job_id, description, created_at, updated_at)
		VALUES (:job_id, :description, :created_at, :updated_at)
	`, job)

	fmt.Println("err create job", err)
	if err != nil {
		return err
	}

	// สร้าง JobMaterial
	for _, material := range materials {
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO job_material (job_id, material_name, quantity)
			VALUES (:job_id, :material_name, :quantity)
		`, material)

		fmt.Println("err create job job_material", err)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (jr *jobRepository) UpdateJob(ctx context.Context, job *models.Job, materials []models.JobMaterial) error {
	tx, err := jr.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// อัพเดต Job
	_, err = tx.NamedExecContext(ctx, `
		UPDATE Job SET description = :description, updated_at = :updated_at
		WHERE job_id = :job_id
	`, job)
	if err != nil {
		return err
	}

	// ลบ JobMaterial เดิม
	_, err = tx.ExecContext(ctx, "DELETE FROM job_material WHERE job_id = $1", job.JobID)
	if err != nil {
		return err
	}

	// สร้าง JobMaterial ใหม่
	for _, material := range materials {
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO job_material (job_id, material_name, quantity)
			VALUES (:job_id, :material_name, :quantity)
		`, material)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (jr *jobRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	var job models.Job
	err := jr.db.GetContext(ctx, &job, `SELECT * FROM Job WHERE job_id = $1`, id)
	if err != nil {
		return nil, err
	}

	materials, err := jr.GetJobMaterials(ctx, id)
	if err != nil {
		return nil, err
	}
	job.Materials = materials

	return &job, nil
}

func (jr *jobRepository) ListJobs(ctx context.Context) ([]*models.Job, error) {
	var jobs []*models.Job
	err := jr.db.SelectContext(ctx, &jobs, `SELECT * FROM Job`)
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		materials, err := jr.GetJobMaterials(ctx, job.JobID)
		if err != nil {
			return nil, err
		}
		job.Materials = materials
	}

	return jobs, nil
}

func (jr *jobRepository) GetJobMaterials(ctx context.Context, jobID uuid.UUID) ([]models.JobMaterial, error) {
	var materials []models.JobMaterial
	query := `
		SELECT jm.job_id, jm.material_name, jm.quantity, m.type, m.unit_of_measure
		FROM job_material jm
		JOIN Material m ON jm.material_name = m.name
		WHERE jm.job_id = $1
	`
	err := jr.db.SelectContext(ctx, &materials, query, jobID)
	return materials, err
}
func (jr *jobRepository) DeleteJob(ctx context.Context, id uuid.UUID) error {
	_, err := jr.db.ExecContext(ctx, `DELETE FROM Job WHERE job_id = $1`, id)
	return err
}
