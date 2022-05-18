package logging

import (
	"database/sql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MysqlOutput struct {
	db *sql.DB
}

func (m *MysqlOutput) Sync() error { return nil }

func (m *MysqlOutput) Write(p []byte) (int, error) {
	query := `INSERT INTO logs (log) values (?)`
	_, err := m.db.Exec(query, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func GetLoggerZap(db *sql.DB) *zap.SugaredLogger {
	query := `CREATE TABLE IF NOT EXISTS logs (
    id serial PRIMARY KEY,
    log varchar(500) not null
	);`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
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
