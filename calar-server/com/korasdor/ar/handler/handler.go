package handler

import (
    "net/http"
    "fmt"
    "com/korasdor/ar/model"
    "gopkg.in/gomail.v2"
    "github.com/gorilla/mux"
    "com/korasdor/ar/services"
    "strconv"
    "io/ioutil"
    "com/korasdor/ar/utils"
)

var (
    Model *model.Model
)

/**
 * вернуть статичный класс
 */
func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:])
}

/**
 * индексная страница
 */
func IndexHandler(w http.ResponseWriter, r *http.Request) {

    //tmpl, err := template.ParseFiles("templates/index.tmpl", "templates/footer.tmpl", "templates/sidebar.tmpl")
    //if err != nil {
    //    panic(err)
    //}
    //
    //tmpl.ExecuteTemplate(w, "index", Model.CommonData)

    fmt.Fprintf(w, "%s", "index page in construction...");
}

/**
 * состояние базы данных
 */
func DBStateHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "db state is: %s", model.DbSuccess);
}

/**
 * получить файл настроек книг.
 */
func BooksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")

    clientIp := r.Header.Get("X-Forwarded-For")
    country := utils.GetCountry(clientIp)

    booksJsonStr := utils.GetBooksJson(country)

    fmt.Fprint(w, string(booksJsonStr))
}

func UpdateBooksHandler(w http.ResponseWriter, r *http.Request) {
    result := utils.UpdateBooksTemplate()

    if result == true {
        fmt.Fprint(w, "good")
    } else {
        fmt.Fprint(w, "bad")
    }
}

/**
 *
 */
func GetTempHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serialsName := vars["file_name"]

    b, err := ioutil.ReadFile(utils.GetTempDir() + serialsName)
    if err != nil {
        utils.PrintOutput(err.Error())
    }

    fmt.Fprint(w, b)
}

/**
 * создаем серийники и сохраняем в файл
 */
func CreateSerialsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serialsName := vars["serials_name"]
    serialsCount, err := strconv.Atoi(vars["serials_count"])
    if err != nil {
        fmt.Fprintf(w, "%s is incorect value", vars["serials_count"]);
    } else {
        if services.CreateSerialKeys(serialsName, serialsCount) {
            fmt.Fprintf(w, "%s is created", serialsName);
        } else {
            fmt.Fprintf(w, "%s is exist", serialsName);
        }
    }
}

/**
 * получить файл с серийниками
 */
func GetSerialsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serialsName := vars["serials_name"]
    serialFormat := vars["serials_format"]

    if serialFormat == "csv" {
        b, err := ioutil.ReadFile(utils.GetDataDir() + "serials/" + serialsName + "." + serialFormat)
        if err != nil {
            utils.PrintOutput(err.Error())
        }

        w.Header().Set("Content-Type", "text/csv")
        fmt.Fprint(w, string(b))
    } else {
        fmt.Fprintf(w, "%s", "this format if unsupported");
    }
}

/**
 * создаем таблицу
 */
func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serialsName := vars["serials_name"]
    tableName := vars["table_name"]

    if services.CreateTable(tableName, serialsName) {
        fmt.Fprintf(w, "%s table is created", tableName);
    } else {
        fmt.Fprintf(w, "%s", "error");
    }
}

func FillTableHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serialsName := vars["serials_name"]
    tableName := vars["table_name"]
    dealerId := vars["dealer_id"]
    serials, err := services.GetSerials(serialsName)
    if err != nil {
        fmt.Fprintf(w, "%s", "error read file")
    } else {
        if services.FillSerialTable(tableName, serials, dealerId) {
            fmt.Fprintf(w, "%s", "table fill complete")
        } else {
            fmt.Fprintf(w, "%s", "error fill table")
        }
    }

}

/**
 * активация серийника
 */
func ActivateSerialsHandler(w http.ResponseWriter, r *http.Request) {
    result := true

    vars := mux.Vars(r)
    tableName := vars["table_name"]
    serialKey := vars["serial_key"]

    canActivate, tryCount := services.SerialCheck(tableName, serialKey)

    if canActivate {
        result = services.SerialUpdate(tableName, tryCount, serialKey)
    } else {
        result = false
    }

    fmt.Fprintf(w, "{is_activated:%t}", result);
}

/**
 * сбросить серийник
 */
func ResetSerialsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tableName := vars["table_name"]
    serialKey := vars["serial_key"]

    res, rowsAffected := services.ResetSerial(tableName, serialKey)

    if res == true {
        if (rowsAffected > 0) {
            fmt.Fprint(w, "Ключ успешно сброщен");
        } else {
            fmt.Fprint(w, "0 полей обновлено, что то не так");
        }
    } else {
        fmt.Fprint(w, "Произощла ошибка, ключ не сброшен");
    }
}

/**
 * информация о серийнике
 */
func AboutSerialsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tableName := vars["table_name"]
    serialKey := vars["serial_key"]

    result, activatedCount, maxActivation, activatedTime := services.AboutSerial(tableName, serialKey)

    if result == true {
        fmt.Fprintf(w, "Количество активаций: %d\nВремя последней активации: %s\nМаксимальное количество активаций: %d", activatedCount, activatedTime, maxActivation);
    } else {
        fmt.Fprint(w, "в данной таблице, не существует заданный серийный ключ");
    }
}

/**
 * отправляем почту
 */
func SendMailHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    if r.FormValue("name") == "" || r.FormValue("email") == "" || r.FormValue("message") == "" {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "error")

        return
    }

    body := fmt.Sprintf("Имя отправителя: %s,\nПочта отправителя: %s,\nСообщение отправителя: %s", r.FormValue("name"), r.FormValue("email"), r.FormValue("message"))

    m := gomail.NewMessage()
    m.SetHeader("From", "noreply@unimedia.uz")
    m.SetHeader("To", "info@unimedia.uz")
    m.SetHeader("Subject", "Отправлено из формы сайта.")
    m.SetBody("text/plain", body)

    d := gomail.NewDialer("smtp.yandex.ru", 465, "noreply@unimedia.uz", "1q2w3e4r5t")

    // Send the email to Bob, Cora and Dan.
    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }

    fmt.Fprintf(w, "%s", "complete");
}