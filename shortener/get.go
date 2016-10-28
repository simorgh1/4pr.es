package shortener

import (
	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	db "github.com/aws/aws-sdk-go/service/dynamodb"
)

func readShortUrl(surl, tname string) (string, error) {
	var ret string
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("eu-central-1"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
		return ret, err
	}
	svc := db.New(sess)
	out, err := svc.GetItem(&db.GetItemInput{
		ConsistentRead: aws.Bool(false),
		//		AttributesToGet: []*string{aws.String("url"),aws.String("redirect")},
		Key: map[string]*db.AttributeValue{
			"url": &db.AttributeValue{
				S: aws.String(surl),
			},
		},
		TableName: aws.String(tname),
	})
	if err != nil {
		log.Printf("Error getting item on DB for URL %s\n%v", surl, err)
		return ret, err
	}
	ret = out.String()
	return ret, nil
}
