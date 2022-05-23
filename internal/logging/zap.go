package logging

import (
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

type MysqlOutput struct {
	db *sql.DB
}
type Log struct {
	Level   string    `json:"level"`
	Time    time.Time `json:"ts"`
	Caller  string    `json:"caller"`
	Message string    `json:"msg"`
}

func (m *MysqlOutput) Sync() error { return nil }

func (m *MysqlOutput) Write(p []byte) (int, error) {
	var log Log
	err := json.Unmarshal(p, &log)
	timeCreation, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	query := `INSERT INTO logs (log, type, date) values (?,?,?)`
	_, err = m.db.Exec(query, p, log.Level, timeCreation)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func GetLoggerZap(db *sql.DB) *zap.SugaredLogger {
	CreateLogsTable(db)
	writerSyncer := MysqlOutput{db: db}
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, &writerSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	l := logger.Sugar()
	return l
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func CreateLogsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS logs (
	id serial PRIMARY KEY,
	type varchar(50),
	log varchar(500) not null,
	date DATETIME not null			
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
