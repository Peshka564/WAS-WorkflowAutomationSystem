package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/errors"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/models"
)

type Workflow struct {
	Db *sql.DB
}

// TODO: Wrap these errors
func (repo *Workflow) FindById(id int) (*models.Workflow, error) {
	stmt, err := repo.Db.Prepare("SELECT * FROM workflows WHERE id = ?");
	if err != nil {
		fmt.Printf("Could not form prepared stmt\n");
		return nil, err
	}
	var workflow models.Workflow
	err = stmt.QueryRow(id).Scan(&workflow.Id, &workflow.CreatedAt, &workflow.UpdatedAt, &workflow.Name, &workflow.Active, &workflow.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("No workflow found\n");
			return nil, errs.NotFoundError{EntityName: "Workflow"}
		}
		fmt.Printf("Could not scan row/some other error\n");
		fmt.Println(err)
		return nil, err
	}
	return &workflow, nil
}

func (repo *Workflow) Update(id int, name string) error {
    _, err := repo.Db.Exec("UPDATE workflows SET name = ? WHERE id = ?", name, id)
    return err
}

func (repo *Workflow) Insert(workflow *models.Workflow) error {
	stmt, err := repo.Db.Prepare("INSERT INTO workflows(name, active, user_id) VALUES (?, ?, ?)");
	if err != nil {
		fmt.Printf("Could not form prepared stmt\n");
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(workflow.Name, workflow.Active, workflow.UserId)
	if err != nil {
		fmt.Printf("Could not scan row/some other error\n");
		return err
	}

	newId, err := res.LastInsertId();
	if err != nil {
		fmt.Printf("Coult not get last insert id\n");
		return err
	}

	newWorkflow, err := repo.FindById(int(newId));
	if err != nil {
		fmt.Printf("Could not get smth\n");
		return err
	}
	
	*workflow = *newWorkflow
	return nil;
}

func (repo *Workflow) FindByUserId(userId int64) ([]models.Workflow, error) {
	stmt, err := repo.Db.Prepare("SELECT * FROM workflows WHERE user_id = ? ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.Workflow
	for rows.Next() {
		var w models.Workflow
		if err := rows.Scan(&w.Id, &w.CreatedAt, &w.UpdatedAt, &w.Name, &w.Active, &w.UserId); err != nil {
			return nil, err
		}
		workflows = append(workflows, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workflows, nil
}

func (repo *Workflow) UpdateActiveStatus(id int, active bool) error {
    stmt, err := repo.Db.Prepare("UPDATE workflows SET active = ? WHERE id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    res, err := stmt.Exec(active, id)
    if err != nil {
        return err
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("workflow not found or unauthorized")
    }

    return nil
}