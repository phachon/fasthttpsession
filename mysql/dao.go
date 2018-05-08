package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"strconv"
	"time"
	"fmt"
)

func newSessionDao(dsn string, tableName string) (*sessionDao, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return &sessionDao{}, err
	}
	return &sessionDao{
		mysqlConn: conn,
		tableName: tableName,
	}, nil
}

type sessionDao struct {
	mysqlConn *sql.DB
	tableName string
}

// get session by sessionId
func (dao *sessionDao) getSessionBySessionId(sessionId string) (session map[string][]byte, err error) {

	sqlStr := fmt.Sprintf("SELECT * FROM %s WHERE session_id=?", dao.tableName)
	return dao.getRow(sqlStr, sessionId)
}

// count sessionId
func (dao *sessionDao) countSessions() int {
	sqlStr := fmt.Sprintf("SELECT count(*) as total FROM %s", dao.tableName)
	res, err := dao.getRow(sqlStr)
	if err != nil {
		return 0
	}
	total, _ := strconv.Atoi(string(res["total"]))
	return total
}

// update session by sessionId
func (dao *sessionDao) updateBySessionId(sessionId string, contents string, lastActiveTime int64) (int64, error) {
	sqlStr := fmt.Sprintf("UPDATE %s SET contents=?,last_active=?", dao.tableName)
	return dao.execute(sqlStr, contents, lastActiveTime)
}

// delete session by sessionId
func (dao *sessionDao) deleteBySessionId(sessionId string) (int64, error) {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE session_id=?", dao.tableName)
	return dao.execute(sqlStr, sessionId)
}

// delete session by maxLifeTime
func (dao *sessionDao) deleteSessionByMaxLifeTime(maxLifeTime int64) (int64, error) {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE last_active<=?", dao.tableName)
	lastTime := time.Now().Unix() - maxLifeTime
	return dao.execute(sqlStr, lastTime)
}

// insert new session
func (dao *sessionDao) insert(sessionId string, contents string, lastActiveTime int64) (int64, error) {
	sqlStr := fmt.Sprintf("INSERT INTO %s (session_id, contents, last_active) VALUES (?,?,?)", dao.tableName)
	return dao.execute(sqlStr, sessionId, contents, lastActiveTime)
}

// get rows
// return []map[string][]byte
func (dao *sessionDao) getRows(sql string, args ...interface{}) (results []map[string][]byte, err error) {

	stmt, err := dao.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return
	}
	defer rows.Close()

	cols := []string{}
	cols, err = rows.Columns()
	if err != nil {
		return
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}
	results = []map[string][]byte{}
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return
		}
		row := make(map[string][]byte)
		for k, v := range vals {
			key := cols[k]
			row[key] = v
		}
		results = append(results, row)
	}
	return
}

// get row
func (dao *sessionDao) getRow(sql string, args ...interface{}) (res map[string][]byte, err error) {
	rows, err := dao.getRows(sql, args...)
	if err != nil {
		return
	}
	if len(rows) > 0 {
		return rows[0], nil
	}
	return
}

// execute(insert, update, delete)
func (dao *sessionDao) execute(sql string, args ...interface{}) (int64, error) {
	stmt, err := dao.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	rows, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return rows.RowsAffected()
}