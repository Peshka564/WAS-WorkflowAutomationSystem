package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/errors"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/models"
)

type WorkflowEdge struct {
	Db *sql.DB
}

// TODO: Wrap these errors
func (repo *WorkflowEdge) FindByWorkflowId(id int) ([]models.WorkflowEdge, error) {
	stmt, err := repo.Db.Prepare("SELECT id, created_at, updated_at, node_from, node_to, workflow_id, display_id FROM workflow_edges WHERE workflow_id = ?");
	if err != nil {
		fmt.Printf("Could not form prepared stmt\n");
		return nil, err
	}
	var workflowEdges []models.WorkflowEdge
	rows, err := stmt.Query(id);
	if err != nil {
		fmt.Printf("Could not execute query\n");
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var workflowEdge models.WorkflowEdge
		err := rows.Scan(&workflowEdge.Id, &workflowEdge.CreatedAt, &workflowEdge.UpdatedAt, &workflowEdge.NodeFrom, &workflowEdge.NodeTo, &workflowEdge.WorkflowId, &workflowEdge.DisplayId)
		if err != nil {
			fmt.Printf("Could not scan row\n");
			fmt.Println(err)
			return nil, err
		}
		workflowEdges = append(workflowEdges, workflowEdge)
	}
	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("No edges found\n");
			return nil, errs.NotFoundError{EntityName: "Workflow edge"}
		}
		fmt.Printf("Could not scan row/some other error\n");
		fmt.Println(err)
		return nil, err
	}
	return workflowEdges, nil
}

func (repo *WorkflowEdge) InsertMany(workflowEdges []models.WorkflowEdge) error {
	sql := "INSERT INTO workflow_edges(id, display_id, node_from, node_to, workflow_id) VALUES"
	var inserts []string
    var params []interface{}

    for _, edge := range workflowEdges {
        inserts = append(inserts, "(?, ?, ?, ?, ?)")
        params = append(params, edge.Id, edge.DisplayId, edge.NodeFrom, edge.NodeTo, edge.WorkflowId)
    }

    sql = sql + strings.Join(inserts, ",")

	stmt, err := repo.Db.Prepare(sql);
	if err != nil {
		fmt.Printf("Could not form prepared stmt\n");
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		fmt.Printf("Could not scan row/some other error\n");
		return err
	}
	
	return nil;
}

func (repo *WorkflowEdge) DeleteByWorkflowId(workflowId int) error {
    _, err := repo.Db.Exec("DELETE FROM workflow_edges WHERE workflow_id = ?", workflowId)
    return err
}

func (repo *WorkflowEdge) Delete(id string) error {
    _, err := repo.Db.Exec("DELETE FROM workflow_edges WHERE id = ?", id)
    return err
}
