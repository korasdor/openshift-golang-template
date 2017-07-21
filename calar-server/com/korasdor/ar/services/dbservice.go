package services

import (
    "com/korasdor/ar/model"
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "com/korasdor/ar/utils"
    "time"
    "crypto/sha1"
    "hash"
)

var (
    db *sql.DB
    mHash hash.Hash
)

func InitDb() {
    var err error

    db, err = sql.Open("mysql", utils.GetDBAddress())
    if err != nil {
        fmt.Println(err)

        model.DbSuccess = "error"
    }

    model.DbSuccess = "succes"
    utils.PrintOutput("db success")
}

func CreateSerialsTable(tableName string) bool {
    result := true;

    stmt, err := db.Prepare("CREATE TABLE " + tableName + "(" +
            "serial_id INT NOT NULL AUTO_INCREMENT," +
            "serial_activated SMALLINT NOT NULL," +
            "max_activations SMALLINT NOT NULL," +
            "serial_activated_time VARCHAR(30)," +
            "dealer_id INT," +
            "serial_key VARCHAR(200) NOT NULL," +
            "PRIMARY KEY ( serial_id ))")

    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    defer stmt.Close()

    _, err = stmt.Exec()
    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    return result
}

func FillSerialTable(tableName string, serials []string, dealerId string) bool {
    var values string
    result := true

    for i := 0; i < len(serials); i++ {
        mHash = sha1.New();
        mHash.Write([]byte(serials[i]))

        serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

        if i < len(serials) - 1 {
            values += fmt.Sprintf("('%s',0,3,NULL,%s),", serialHash, dealerId)
        } else {
            values += fmt.Sprintf("('%s',0,3,NULL,%s)", serialHash, dealerId)
        }
    }

    stmt, err := db.Prepare("INSERT INTO " + tableName + "(serial_key,serial_activated,max_activations,serial_activated_time,dealer_id) VALUES " + values)
    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    defer stmt.Close()

    _, err = stmt.Exec()
    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    return result
}

func SerialCheck(tableName string, key string) (bool, int) {
    result := false
    serialActivated := -1

    mHash = sha1.New();
    mHash.Write([]byte(key))
    serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

    rows, err := db.Query("SELECT serial_activated,serial_key,max_activations FROM " + tableName + " WHERE serial_key=?", serialHash)
    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    defer rows.Close()

    if rows.Next() == true {
        result = true
        var serialKey string
        var maxActivations int

        err := rows.Scan(&serialActivated, &serialKey, &maxActivations)
        if err != nil {
            result = false
            utils.PrintOutput(err.Error())
        } else {
            if serialActivated < maxActivations {
                result = true
                serialActivated++;
            } else {
                result = false
            }
        }
    }

    return result, serialActivated
}

func SerialUpdate(tableName string, tryCount int, key string) bool {
    result := true

    stmt, err := db.Prepare("UPDATE " + tableName + " SET serial_activated=?, serial_activated_time=? WHERE serial_key=?")

    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    defer stmt.Close()

    activatedTime := time.Now().Format(time.RFC3339)
    mHash = sha1.New();
    mHash.Write([]byte(key))
    serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

    _, err = stmt.Exec(tryCount, activatedTime, serialHash)
    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    return result
}

func ResetSerial(tableName string, key string) (bool, int64) {
    result := true

    stmt, err := db.Prepare("UPDATE " + tableName + " SET serial_activated=0, serial_activated_time=NULL WHERE serial_key=?")

    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    defer stmt.Close()

    mHash = sha1.New();
    mHash.Write([]byte(key))
    serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

    out, err := stmt.Exec(serialHash)
    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    rowsAffected, _ := out.RowsAffected()

    return result, rowsAffected
}

func AboutSerial(tableName string, key string) (bool, int, int, string) {
    result := false

    var (
        serialActivated int
        maxActivations int
        serialActivatedTime string
    )

    mHash = sha1.New();
    mHash.Write([]byte(key))
    serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

    rows, err := db.Query("SELECT serial_activated, max_activations, COALESCE(serial_activated_time, '') as serial_activated_time FROM " + tableName + " WHERE serial_key=?", serialHash)
    if err != nil {
        result = false
        utils.PrintOutput(err.Error())
    }

    defer rows.Close()

    if rows.Next() == true {
        result = true

        err := rows.Scan(&serialActivated, &maxActivations, &serialActivatedTime)
        if err != nil {
            result = false
            utils.PrintOutput(err.Error())
        } else {
            result = true;
        }
    }

    return result, serialActivated, maxActivations, serialActivatedTime
}

func CloseDb() {
    db.Close()
}
