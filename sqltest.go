package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {

	type Arrangement struct {
		Dato      string `json:"DATO"`
		Film      string `json:"FILM"`
		Sal       string `json:"SAL"`
		Land      string `json:"LAND"`
		ImdbLink  string `json:"IMDBLINK"`
		ImdbId    string `json:"IMDBID"`
		Stig      string `json:"stig"`
		Paal      string `json:"paal"`
		Terjeb    string `json:"terjeb"`
		Tor       string `json:"tor"`
		Tom       string `json:"tom"`
		Eivind    string `json:"eivind"`
		Olav      string `json:"olav"`
		Finn      string `json:"finn"`
		Molde     string `json:"molde"`
		Pallen    string `json:"pallen"`
		SwingStig string `json:"swing_stig"`
		Bengt     string `json:"bengt"`
	}

	var query string
	var a Arrangement

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	// URL to our queue
	qURL := "https://sqs.eu-west-1.amazonaws.com/947879583077/torsdagskino_sqs"

	query = "select convert(varchar, dato, 23), film, sal, land,"
	query += "	imdblink, case when substring(imdblink,1,5) = 'http:' then substring(imdblink,27,9)"
	query += "	 else substring(imdblink,28,9) end imdbid,"
	query += "      iif(stig is null, '', '\"STIG\": '+ltrim(str(stig))),"
	query += "      iif(paal is null, '', ', \"PAAL\": '+ltrim(str(paal))),"
	query += "      iif(terjeb is null, '', ', \"TERJEB\": '+ltrim(str(terjeb))),"
	query += "      iif(tor is null, '', ', \"TOR\": '+ltrim(str(tor))),"
	query += "      iif(tom is null, '', ', \"TOM\": '+ltrim(str(tom))),"
	query += "      iif(eivind is null, '', ', \"EIVIND\": '+ltrim(str(eivind))),"
	query += "      iif(olav is null, '', ', \"OLAV\": '+ltrim(str(olav))),"
	query += "      iif(finn is null, '', ', \"FINN\": '+ltrim(str(finn))),"
	query += "      iif(molde is null, '', ', \"MOLDE\": '+ltrim(str(molde))),"
	query += "      iif(pallen is null, '', ', \"PALLEN\":'+ltrim(str(pallen))),"
	query += "      iif(swing_stig is null, '', ', \"SWING_STIG\": '+ltrim(str(swing_stig))),"
	query += "      iif(bengt is null, '', ', \"BENGT\": '+ltrim(str(bengt)))"
	query += " from torsdagskino.dbo.filmlogg"

	condb, errdb := sql.Open("mssql", "server=10.21.93.25;user id=reporting;password=reporting;")
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}

	rows, err := condb.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&a.Dato, &a.Film, &a.Sal, &a.Land, &a.ImdbLink, &a.ImdbId,
			&a.Stig, &a.Paal, &a.Terjeb, &a.Tor, &a.Tom, &a.Eivind, &a.Olav, &a.Finn, &a.Molde, &a.Pallen, &a.SwingStig, &a.Bengt)
		if err != nil {
			log.Fatal(err)
		}

		Filmdata := fmt.Sprintf("\"DATO\":\"%s\",\"FILM\":\"%s\",\"SAL\":\"%s\",\"LAND\":\"%s\",\"IMDBLINK\":\"%s\",", a.Dato, a.Film, a.Sal, a.Land, a.ImdbLink)
		Terningkast := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s", a.Stig, a.Paal, a.Terjeb, a.Tor, a.Tom, a.Eivind, a.Olav, a.Finn, a.Molde, a.Pallen, a.SwingStig, a.Bengt)

		ut := fmt.Sprintf("{%s %s}", Filmdata, Terningkast)
		fmt.Println(string(ut))

		result, err := svc.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(10),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"Dato":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(string(a.Dato))},
				"Film":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(string(a.Film))},
				"Sal":      &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(string(a.Sal))},
				"Land":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(string(a.Land))},
				"ImdbLink": &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(string(a.ImdbLink))},
			}, MessageBody: aws.String(string(ut)),
			QueueUrl: &qURL,
		})

		if err != nil {
			fmt.Println("Error", err)
			return
		}
		fmt.Println("Success", *result.MessageId)

	}
	defer condb.Close()

}
