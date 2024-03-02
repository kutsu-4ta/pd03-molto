package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	// Spotify認証のリダイレクトURL
	r.GET("/spotify/authorize", func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("code")

		if state == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": code, "error": "state_mismatch"})
			return
		}

		ok := true
		// c.String(http.StatusOK, "Query Parameter: %s\n", value)
		if ok {
			c.JSON(http.StatusOK, gin.H{"code": code})

			// postリクエストを送信
			doPostRequest(code)
		} else {
			c.JSON(http.StatusOK, gin.H{"query": code, "status": "no value"})
		}
	})

	return r
}

func doPostRequest(code string) {
	// POSTリクエストのボディに含めるデータを定義
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))

	body := strings.NewReader(form.Encode())
	println(body)

	// POSTリクエスト作成
	req, err := http.NewRequest("POST", os.Getenv("SPOTIFY_TOKEN_URL"), body)
	if err != nil {
		log.Fatal(err)
	}

	// ヘッダー
	authToken := basicAuthToken(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	req.Header.Set("Authorization", "Basic "+authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// POSTリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("Failed to send request:", err.Error())
		return
	}

	log.Println(os.Getenv("GIN_PORT"))

	// (実際のアプリケーションでは適切な処理を行う)
	println("Response status:", resp.Status)
	println("Response status:", resp.StatusCode)
}

// Basic認証トークンを生成する関数
func basicAuthToken(clientID, clientSecret string) string {
	auth := clientID + ":" + clientSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {
	// .envファイルから環境変数をロード
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// ログファイル オープン // TODO: 調査が必要 途中からビルドでスタックするようになった
	// file, err := os.Create("gin.log")
	// if err != nil {
	// 	log.Fatal("Failed to create log file:", err)
	// }
	// // Ginのデフォルトのロガーに設定
	// gin.DefaultWriter = file
	// defer file.Close()

	// ルーティング
	r := setupRouter()

	// Listen and Server
	r.Run(":" + os.Getenv("GIN_PORT"))
}
