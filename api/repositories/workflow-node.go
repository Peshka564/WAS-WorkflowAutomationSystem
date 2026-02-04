package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	errs "github.com/Peshka564/WAS-WorkflowAutomationSystem/errors"
	"github.com/Peshka564/WAS-WorkflowAutomationSystem/models"
)

type WorkflowNode struct {
	Db *sql.DB
}

// TODO: Wrap these errors
func (repo *WorkflowNode) FindById(id int) (*models.WorkflowNode, error) {
	stmt, err := repo.Db.Prepare("SELECT * FROM workflow_nodes WHERE id = ?");
	if err != nil {
		fmt.Printf("Could not form prepared stmt\n");
		return nil, err
	}
	var workflowNode models.WorkflowNode
	err = stmt.QueryRow(id).Scan(&workflowNode.Id, &workflowNode.CreatedAt, &workflowNode.UpdatedAt, &workflowNode.WorkflowId, &workflowNode.TaskName, &workflowNode.WorkflowType, &workflowNode.Position)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("No node found\n");
			return nil, errs.NotFoundError{}
		}
		fmt.Printf("Could not scan row/some other error\n");
		fmt.Println(err)
		return nil, err
	}
	return &workflowNode, nil
}

func (repo *WorkflowNode) FindByWorkflowId(id int) ([]models.WorkflowNode, error) {
	stmt, err := repo.Db.Prepare("SELECT * FROM workflow_nodes WHERE workflow_id = ?");
	if err != nil {
		fmt.Printf("Could not form prepared stmt\n");
		return nil, err
	}
	var workflowNodes []models.WorkflowNode
	rows, err := stmt.Query(id);
	if err != nil {
		fmt.Printf("Could not execute query\n");
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var workflowNode models.WorkflowNode
		err := rows.Scan(&workflowNode.Id, &workflowNode.CreatedAt, &workflowNode.UpdatedAt, &workflowNode.WorkflowId, &workflowNode.TaskName, &workflowNode.WorkflowType, &workflowNode.Position)
		if err != nil {
			fmt.Printf("Could not scan row\n");
			fmt.Println(err)
			return nil, err
		}
		workflowNodes = append(workflowNodes, workflowNode)
	}
	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("No nodes found\n");
			return nil, err
		}
		fmt.Printf("Could not scan row/some other error\n");
		fmt.Println(err)
		return nil, err
	}
	return workflowNodes, nil
}

func (repo *WorkflowNode) InsertMany(workflowNodes []models.WorkflowNode) error {
	sql := "INSERT INTO workflow_nodes(workflow_id, task_name, workflow_type, position) VALUES"
	var inserts []string
    var params []interface{}

    for _, node := range workflowNodes {
        inserts = append(inserts, "(?, ?)")
        params = append(params, node.WorkflowId, node.TaskName, node.WorkflowType, node.Position)
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

func (repo *WorkflowNode) Insert(workflowNode *models.WorkflowNode) error {
	stmt, err := repo.Db.Prepare("INSERT INTO workflow_nodes(workflow_id, task_name, workflow_type, position) VALUES (?, ?, ?, ?)");
	if err != nil {
		fmt.Printf("Could not form prepared stmt\n");
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(workflowNode.WorkflowId, workflowNode.TaskName, workflowNode.WorkflowType, workflowNode.Position)
	if err != nil {
		fmt.Printf("Could not scan row/some other error\n");
		return err
	}

	newId, err := res.LastInsertId();
	if err != nil {
		fmt.Printf("Coult not get last insert id\n");
		return err
	}

	newWorkflowNode, err := repo.FindById(int(newId));
	if err != nil {
		fmt.Printf("Could not get smth\n");
		return err
	}
	
	*workflowNode = *newWorkflowNode
	return nil;
}