package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"health_checker/pkg/config"
	"health_checker/pkg/model"
	"health_checker/pkg/utils"
)

var Database *Postgres

type Postgres struct {
	db *sql.DB
}

const (
	Success = "success"
	Failed  = "failed"
)

func (p *Postgres) CreateNewToken(username, token string) error {
	err := p.db.QueryRow(`INSERT INTO tokens(username, token)
	VALUES($1,$2)`, username, token).Err()
	return err
}

func (p *Postgres) GetTokenByUsername(username string) (string, error) {
	var token string
	err := p.db.QueryRow(`select token from tokens where username=$1`, username).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (p *Postgres) FlushUserTokens(username string) error {
	_, err := p.db.Exec(`delete from tokens where username=$1`, username)
	return err
}
func (p *Postgres) CreateNewUser(username, password string) error {
	err := p.db.QueryRow(`INSERT INTO users(username, password)
	VALUES($1,$2)`, username, utils.HashString(password)).Err()
	return err
}

func (p *Postgres) GetUserByID(username string) (*model.User, error) {
	var user model.User
	err := p.db.QueryRow(`select username, password from users where username=$1 `, username).Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (p *Postgres) CreateNewEndpoint(url, username string, threshold int) (int, error) {
	var endpointID int
	err := p.db.QueryRow(`INSERT INTO endpoints(url, threshold, username)
	VALUES($1,$2,$3) RETURNING id`, url, threshold, username).Scan(&endpointID)
	if err != nil {
		return -1, err
	}
	return endpointID, nil
}

func (p *Postgres) GetEndpointsByUsername(username string) ([]*model.EndpointResponse, error) {
	rows, err := p.db.Query(`select id,url from endpoints where username=$1 `, username)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	endpoints := make([]*model.EndpointResponse, 0)
	for rows.Next() {
		var endpoint model.EndpointResponse
		if err2 := rows.Scan(&endpoint.Id, &endpoint.Url); err2 != nil {
			fmt.Println(err2.Error())

			return nil, err2
		}

		endpoints = append(endpoints, &endpoint)
	}

	return endpoints, nil
}

func (p *Postgres) GetAllEndpoints() ([]*model.EndpointResponse, error) {
	rows, err := p.db.Query(`select id,url from endpoints `)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	endpoints := make([]*model.EndpointResponse, 0)
	for rows.Next() {
		var endpoint model.EndpointResponse
		if err2 := rows.Scan(&endpoint.Id, &endpoint.Url); err2 != nil {
			fmt.Println(err2.Error())

			return nil, err2
		}

		endpoints = append(endpoints, &endpoint)
	}

	return endpoints, nil
}

func (p *Postgres) GetEndpointsByThresholdCrossed() ([]*model.EndpointResponse, error) {
	rows, err := p.db.Query(`select id,url from endpoints where failed >= threshold;`)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	endpoints := make([]*model.EndpointResponse, 0)
	for rows.Next() {
		var endpoint model.EndpointResponse
		if err2 := rows.Scan(&endpoint.Id, &endpoint.Url); err2 != nil {
			fmt.Println(err2.Error())

			return nil, err2
		}
		endpoints = append(endpoints, &endpoint)
	}

	return endpoints, nil
}

func (p *Postgres) UpdateEndpointResultByOne(id int, result string) error {
	query := fmt.Sprintf("select %s from endpoints where id=%d", result, id)
	var pastResult int
	err := p.db.QueryRow(query).Scan(&pastResult)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("update endpoints set %s=%d where id=%d", result, pastResult+1, id)
	_, err = p.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) ResetEndpointFailed(id int) error {
	query := fmt.Sprintf("update endpoints set %s=0 where id=%d", "failed", id)
	_, err := p.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetEndpointStatusByID(id int) (int, int, error) {
	var success, failed int
	err := p.db.QueryRow(`select success,failed from endpoints where id=$1 `, id).Scan(&success, &failed)
	if err != nil {
		return -1, -1, err
	}

	return success, failed, nil
}
func (p *Postgres) CreateAlert(endpointID int, alertDesc string) (int, error) {
	var alertID int
	err := p.db.QueryRow(`INSERT INTO alerts(endpoint_id, description)
	VALUES($1,$2) RETURNING id`, endpointID, alertDesc).Scan(&alertID)
	if err != nil {
		return -1, err
	}
	return alertID, nil
}

func (p *Postgres) GetAlertByEndpointID(endpointID int) ([]*model.Alert, error) {
	rows, err := p.db.Query(`select id,description from alerts where endpoint_id=$1 `, endpointID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	alerts := make([]*model.Alert, 0)
	for rows.Next() {
		var alert model.Alert
		if err2 := rows.Scan(&alert.Id, &alert.Description); err2 != nil {
			fmt.Println(err2.Error())

			return nil, err2
		}

		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

func SetupPostgres() {
	postgresConf := config.Conf.Postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgresConf.Host, postgresConf.Port, postgresConf.User, postgresConf.Password, postgresConf.Db)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Postgres successfully connected!")
	Database = &Postgres{db: db}

	// list of tables needed
	_, err = Database.db.Exec(`CREATE TABLE IF NOT EXISTS users (
	username VARCHAR ( 50 ) PRIMARY KEY,
	password VARCHAR ( 50 ) NOT NULL,
    endpoints_num INT default 0
	);`)
	if err != nil {
		panic(errors.Wrap(err, "couldn't create users database"))
	}

	_, err = Database.db.Exec(`CREATE TABLE IF NOT EXISTS endpoints (
	id serial PRIMARY KEY,
	url VARCHAR ( 50 ) NOT NULL,
	threshold int NOT NULL,
    username varchar (50) NOT NULL,
    success int default 0,
    failed int default 0
	);`)
	if err != nil {
		panic(errors.Wrap(err, "couldn't create endpoints database"))
	}

	_, err = Database.db.Exec(`CREATE TABLE IF NOT EXISTS alerts (
	id serial PRIMARY KEY,
	endpoint_id int NOT NULL,
	description VARCHAR ( 50 ) NOT NULL
	);`)
	if err != nil {
		panic(errors.Wrap(err, "couldn't create alerts database"))
	}

	_, err = Database.db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
	username VARCHAR ( 50 ) PRIMARY KEY,
    token varchar (150) NOT NULL
	);`)
	if err != nil {
		panic(errors.Wrap(err, "couldn't create alerts database"))
	}
}
