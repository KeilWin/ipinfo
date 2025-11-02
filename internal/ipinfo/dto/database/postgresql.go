package database

type PostgreSqlDatabase struct {
	Database
}

func (p *PostgreSqlDatabase) StartUp() {

}

func (p *PostgreSqlDatabase) ShutDown() {

}

func (p *PostgreSqlDatabase) GetIpInfo(ipAddress string) {

}

func NewPostgreSqlDatabase() *PostgreSqlDatabase {
	return &PostgreSqlDatabase{}
}
