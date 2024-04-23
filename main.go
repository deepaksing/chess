package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/deepaksing/chess/server"
	"github.com/deepaksing/chess/store/db"
	"github.com/deepaksing/chess/types"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username string
}

type UserResp struct {
	Username string `json:"username"`
}

type MathcResp struct {
	Match_id int `json:"match_id"`
}

func CorsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		corsMiddleware := cors.Default()
		handler := corsMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(c)
		}))
		handler.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func main() {
	//1 .Database
	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	//2. adding tables to the db
	err = dbConn.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	//server
	//1. create api
	//signup api - currently no signup just enter username and we save it
	e := echo.New()
	// corsMiddleware := cors.Default()
	e.Use(CorsMiddleware)

	e.GET("/health", func(c echo.Context) error {
		response := map[string]bool{"status": true}
		return c.JSON(http.StatusOK, response)
	})

	e.POST("/login", func(c echo.Context) error {
		var user User
		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		// Search if user with this username exists
		userMatch, err := dbConn.FindUserByUsername(user.Username)
		if err != nil {
			log.Fatal(err)
		}

		if userMatch != nil {
			fmt.Println("checking password")
			// User exists, check if password is correct
			err := bcrypt.CompareHashAndPassword([]byte(userMatch.HashedPassword), []byte(user.Password))
			fmt.Println(err)
			if err != nil {
				// resp := map[string]string{"error": "This username is already taken or the password is incorrect"}
				// return c.JSON(http.StatusUnauthorized, resp)
				fmt.Println("err")
				return echo.NewHTTPError(http.StatusUnauthorized, "This username is already taken or the password is incorrect")
			}
		} else {
			// User does not exist, create new user
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
			}
			newUser := types.UserTable{
				Username:       user.Username,
				HashedPassword: string(hashedPassword),
			}
			if err := dbConn.CreateUser(&newUser); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
			}
		}

		// Generate JWT token
		accessToken, err := server.CreateJWTToken(user.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
		}

		// Set JWT token as HTTP-only cookie
		cookieExp := time.Now().Add(server.CookieExpDuration)
		server.SetTokenCookie(c, accessToken, cookieExp)

		userResp := UserResponse{
			Username: user.Username,
		}

		return c.JSON(http.StatusOK, userResp)

	})

	e.POST("/join", func(c echo.Context) error {
		var username UserResp
		if err := c.Bind(&username); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		err = dbConn.JoinGame(username.Username)
		if err != nil {
			log.Fatal(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to join the game"})
		}

		//matchmaking logic here
		matchedPlayers, err := dbConn.MatchPlayers(username.Username)
		if err != nil {
			log.Println("Matchmaking failed:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Matchmaking failed"})
		}

		if matchedPlayers == nil {
			log.Println("No other player found")
			return c.JSON(http.StatusNotFound, map[string]string{"message": "No other player found"})
		}
		log.Println("Match found:", username, "vs", matchedPlayers[0])

		return c.JSON(http.StatusOK, map[string]string{"message": "Joined the game successfully"})

	})

	e.POST("/get-status", func(c echo.Context) error {
		fmt.Println("inside status")
		var username UserResp
		if err := c.Bind(&username); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		status, err := dbConn.GetUserStatus(username.Username)
		fmt.Println(status)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "status of given user not found, matching")
		}

		return c.JSON(http.StatusOK, status)

	})

	e.GET("/match", func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("error reading message", err)
				break
			}

			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("Error echoing message:", err)
				break
			}
		}

		return nil
	})

	e.POST("/chess_board", func(c echo.Context) error {
		var matchId MathcResp
		fmt.Println("fetc 1")
		if err := c.Bind(&matchId); err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		fmt.Println("fetcn 2")
		boardState, err := dbConn.GetBoardState(matchId.Match_id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "status of given user not found, matching")
		}

		return c.JSON(http.StatusOK, boardState)
	})

	e.POST("/make_move", func(c echo.Context) error {
		var move types.MoveResp
		if err := c.Bind(&move); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		//save the move
		err := dbConn.SaveBoardMove(move)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "error")
		}
		return nil
	})

	e.Start(":8080")
}
