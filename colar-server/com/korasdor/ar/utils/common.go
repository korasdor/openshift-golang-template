package utils

import (
    "fmt"
    "os"
    "github.com/oschwald/geoip2-golang"
    "log"
    "net"
    "io/ioutil"
    "encoding/json"
    "net/http"
)

var (
    TEMPLATES_PATH string = GetDataDir() + "templates/books.json"
)

func GetBind() string {
    var bind string
    if os.Getenv("OPENSHIFT_GO_IP") == "" {
        bind = fmt.Sprintf("%s:%s", "127.0.0.1", "8000")
    } else {
        bind = fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
    }

    return bind
}

func GetDBAddress() string {
    var dbAddress string

    // korasdor:19841986aA@tcp(127.0.0.1:3306)/im
    if os.Getenv("OPENSHIFT_MYSQL_DB_USERNAME") == "" {
        dbAddress = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "korasdor", "19841986aA", "127.0.0.1", "3306", "colar_db")
    } else {
        dbAddress = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("OPENSHIFT_MYSQL_DB_USERNAME"),
            os.Getenv("OPENSHIFT_MYSQL_DB_PASSWORD"),
            os.Getenv("OPENSHIFT_MYSQL_DB_HOST"),
            os.Getenv("OPENSHIFT_MYSQL_DB_PORT"),
            "ar")
    }

    return dbAddress
}

func GetTempDir() string {
    return os.Getenv("OPENSHIFT_TMP_DIR")
}

func GetDataDir() string {
    return os.Getenv("OPENSHIFT_DATA_DIR")
}

func GetCountry(clientIp string) string {
    country := "Us"

    db, err := geoip2.Open("static/data/GeoLite2-Country.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ip := net.ParseIP(clientIp)
    if ip != nil {
        record, err := db.City(ip)
        if err != nil {
            PrintOutput(err.Error())
        } else {
            country = record.Country.IsoCode;
        }
    }

    return country;
}

func UpdateBooksTemplate() bool {
    var result = true;

    response, err := http.Get("http://colarit.com/colar/templates/books.json")
    if err != nil {
        result = false;
        fmt.Println(err)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            result = false;
            fmt.Println(err)
        }

        err = ioutil.WriteFile(TEMPLATES_PATH, contents, 0644)
        if err != nil {
            result = false;
            fmt.Println(err)
        }
    }

    return result;
}

func GetBooksJson(country string) []byte {
    var bookJsonStr []byte

    b, err := ioutil.ReadFile(TEMPLATES_PATH)
    if err != nil {
        PrintOutput(err.Error())
    }

    var dat map[string]interface{}
    if err := json.Unmarshal(b, &dat); err != nil {
        PrintOutput(err.Error())
    } else {

        if country == "UZ" {
            dat["assets_server"] = "http://colarit.com/colar"
            dat["supported_langs"] = []string{"en", "ru", "uz"}
            // dat["assets_server"] = "http://colar.uz"
        } else {
            dat["assets_server"] = "http://colarit.com/colar"
            dat["supported_langs"] = []string{"en", "ru"}
        }

        bookJsonStr, err = json.Marshal(dat)
        if err != nil {
            PrintOutput(err.Error())
        }
    }

    return bookJsonStr

}

func PrintOutput(str string) {
    fmt.Println(str)
}

