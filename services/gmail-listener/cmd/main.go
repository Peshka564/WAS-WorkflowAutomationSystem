package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Peshka564/WAS-WorkflowAutomationSystem/services/api/utils"
	pb "github.com/Peshka564/WAS-WorkflowAutomationSystem/shared/proto"
)

const pollInterval int = 20

type GmailListener struct {
	Db           *sql.DB
	UserService  pb.UserServiceClient
	Orchestrator pb.OrchestratorClient
}

type TriggerJob struct {
	NodeId       string
	Credentialid int
	LastCheckAt  time.Time
	// This is used as a hack because there are cases when we can get duplicate gmail emails
	LastMessageId string
}


func (l *GmailListener) Poll() {
	// We get each node from every workflow which is a listener node and we query the gmail API for it
	rows, err := l.Db.Query(`
        SELECT 
            n.id, 
            n.credential_id, 
            COALESCE(t.last_check_at, CAST('1971-01-01 00:00:00' AS DATETIME)), 
            COALESCE(t.last_message_id, '')
        FROM workflow_nodes n
        LEFT JOIN trigger_states t ON n.id = t.node_id
        WHERE n.type = 0 
          AND n.service_name = 'gmail'
        ORDER BY t.last_check_at ASC
        LIMIT 10
	`)
	// AND (t.last_check_at IS NULL OR t.last_check_at < NOW() - INTERVAL 60 SECOND)
	if err != nil {
		log.Printf("DB Error: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var job TriggerJob
		err := rows.Scan(
            &job.NodeId, 
            &job.Credentialid, 
            &job.LastCheckAt, 
            &job.LastMessageId,
        )
		if err != nil {
            log.Printf("Scan error: %v", err)
            continue
        }
		// Note: We don't parallelize this due to rate limits
		l.CheckForNewEmails(job)
	}

	log.Println("Finished polling for now")
}

func (l *GmailListener) CheckForNewEmails(job TriggerJob) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	
	tokenResp, err := l.UserService.GetCredentials(ctx, &pb.GetCredentialsRequest{
		CredentialId: int32(job.Credentialid),
	})
	if err != nil {
		log.Printf("Auth Failed for Node %s: %v", job.NodeId, err)
		return
	}

	token := &oauth2.Token{AccessToken: tokenResp.AccessToken}
	clientConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes: []string{"https://www.googleapis.com/auth/gmail.readonly"},
	}
	client := clientConfig.Client(ctx, token)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Gmail Client Error:", err)
		return
	}

	// Gmail Query: "after:1698300000" (Unix Timestamp)
	query := fmt.Sprintf("after:%d", job.LastCheckAt.Add(-2 * time.Minute).Unix())
	
	listCall := srv.Users.Messages.List("me").Q(query).MaxResults(1)
	res, err := listCall.Do()
	if err != nil {
		log.Printf("Gmail API Error: %v", err)
		return
	}

	if len(res.Messages) == 0 {
		log.Printf("No new emails found for node %s", job.NodeId)
		l.updateCheckpoint(job.NodeId, job.LastMessageId)
		return
	}

	// We have a new email
	messageID := res.Messages[0].Id

	// This is a hack because in some cases we can read an already read email
	if messageID == job.LastMessageId {
		log.Printf("Skipping duplicate message: %s", messageID)
        l.updateCheckpoint(job.NodeId, messageID)
        return
	}

	fullMsg, err := srv.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		log.Println("Error getting messages: ", err)
		return
	}

	// Extract the data from the email
	subjectHeader, _ := utils.Find(fullMsg.Payload.Headers, func(header *gmail.MessagePartHeader) bool { 
		return header.Name == "Subject"
	})
	subject := ""
	if subjectHeader != nil {
		subject = (*subjectHeader).Value
	}

	fromHeader, _ := utils.Find(fullMsg.Payload.Headers, func(header *gmail.MessagePartHeader) bool { 
		return header.Name == "From"
	})
	from := ""
	if fromHeader != nil {
		parsedAddr, err := mail.ParseAddress((*fromHeader).Value)
		if(err == nil) {
			from = parsedAddr.Address
		}
	}

	body := fullMsg.Snippet

	log.Printf("New Email Detected! Subject: %s, From: %s, Body: %s", subject, from, body)

	// Trigger the workflow

	payload := fmt.Sprintf(`{"email_subject": "%s", "email_from": "%s", "email_body": "%s", "id": "%s"}`, 
		subject, from, body, messageID)

	_, err = l.Orchestrator.TriggerWorkflow(ctx, &pb.TriggerRequest{
		TriggerNodeId: job.NodeId,
		InitialPayload: payload,
	})

	if err != nil {
		log.Printf("Failed to trigger workflow: %v", err)
		return
	}

	l.updateCheckpoint(job.NodeId, messageID)
}

func (l *GmailListener) updateCheckpoint(nodeId, msgId string) {
	// Upsert
	query := `
		INSERT INTO trigger_states (node_id, last_check_at, last_message_id)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			last_check_at = VALUES(last_check_at),
            last_message_id = VALUES(last_message_id);
	`
	_, err := l.Db.Exec(query, nodeId, time.Now().UTC(), msgId)
	if err != nil {
        log.Printf("Failed to update checkpoint for %s: %v", nodeId, err)
    }
}


func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/was_api?parseTime=true&loc=UTC")
    if err != nil {
        log.Fatal("Could not connect to db", err);
        return;
    }
	
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Could not load ENV vars", err);
		return;
	}
	
	userConn, _ := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer userConn.Close()
	orchConn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer orchConn.Close()

	listener := &GmailListener{
		Db:           db,
		UserService:  pb.NewUserServiceClient(userConn),
		Orchestrator: pb.NewOrchestratorClient(orchConn),
	}

	log.Printf("Gmail Listener started. Polling every %ds...\n", pollInterval)
	
	// TODO: Webhooks
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	for range ticker.C {
		listener.Poll()
	}
}