package main

import (
	"fmt"
	"webauthn/config"
	"webauthn/database"
	"webauthn/model"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

var W *webauthn.WebAuthn

func main() {
	config.InitConfig()
	database.InitDB()
	r := gin.Default()
	wconfig := &webauthn.Config{
		RPDisplayName: "Go Webauthn",                          // Display Name for your site
		RPID:          "webauthn.sylu.site",                   // Generally the FQDN for your site
		RPOrigins:     []string{"https://webauthn.sylu.site"}, // The origin URLs allowed for WebAuthn requests
	}
	var err error
	W, err = webauthn.New(wconfig)
	if err != nil {
		fmt.Println(err)
	}
	r.GET("/get_register", GetRegister)
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/get_login", GetLogin)
	r.Run(":8080")
}
func GetRegister(ctx *gin.Context) {
	name := ctx.Query("name")
	challengeid, _ := protocol.CreateChallenge()

	// if err != nil {
	// 	fmt.Println("get register err", err)
	// }
	DB := database.GetDB()
	var user model.User
	DB.Table("users").Where("name = ?", name).First(&user)
	var userid int64
	DB.Table("users").Count(&userid)
	userid += 1
	newChallenge := model.Challenge{
		Username:  name,
		Challenge: challengeid.String(),
	}
	DB.Table("challenges").Create(&newChallenge)
	if user.ID == 0 {
		ctx.JSON(200, gin.H{
			"code": 200,
			"data": gin.H{
				// "challengeid": challengeid.String(),
				"challengeid": challengeid.String(),
				"userid":      userid,
			},
			"msg": "get register successed",
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": 422,
			"data": "",
			"msg":  "username is in used",
		})
	}
}

func Register(ctx *gin.Context) {
	var webauthnJson protocol.CredentialCreationResponse
	ctx.ShouldBindJSON(&webauthnJson)
	fmt.Println(webauthnJson)
	// var pcc *protocol.ParsedCredentialCreationData
	pcc, _ := webauthnJson.Parse()
	err := pcc.Verify(pcc.Response.CollectedClientData.Challenge, false, W.Config.RPID, []string{pcc.Response.CollectedClientData.Origin})
	if err != nil {
		fmt.Println("err", err)
		ctx.JSON(200, gin.H{
			"code": "422",
			"data": err,
			"msg":  "verify error",
		})
	}
	credential, _ := webauthn.MakeNewCredential(pcc)
	DB := database.GetDB()
	var challenge model.Challenge
	DB.Table("challenges").Where("challenge = ?", pcc.Response.CollectedClientData.Challenge).First(&challenge)

	newUser := model.User{
		Name:         challenge.Username,
		DisplayName:  "displayName",
		Certificates: credential.ID,
	}
	DB.Table("users").Create(&newUser)
	DB.Table("users").Where("name = ?", challenge.Username).First(&newUser)
	newCerd := model.Certificate{
		UserID:     newUser.ID,
		Credention: credential.PublicKey,
	}
	DB.Table("certificates").Create(&newCerd)
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "register success",
		"msg":  "webauthn",
	})
}

func GetLogin(ctx *gin.Context) {
	challengeid, _ := protocol.CreateChallenge()
	DB := database.GetDB()
	name := ctx.Query("name")
	newChallenge := model.Challenge{
		Username:  name,
		Challenge: challengeid.String(),
	}
	DB.Table("challenges").Create(&newChallenge)
	var user model.User
	DB.Table("users").Where("name = ?", name).First(&user)
	if user.ID == 0 {
		ctx.JSON(200, gin.H{
			"code": "422",
			"data": "user is not register",
			"msg":  "error",
		})
	}
	var json_data protocol.PublicKeyCredentialRequestOptions
	json_data.AllowedCredentials = []protocol.CredentialDescriptor{
		{
			Type:            "public-key",
			CredentialID:    user.Certificates,
			Transport:       []protocol.AuthenticatorTransport{"usb", "nfc", "ble"},
			AttestationType: "",
		},
	}
	json_data.Challenge = challengeid
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": json_data,
		"msg":  "msg",
	})
}

func Login(ctx *gin.Context) {

	var webauthnJson protocol.CredentialAssertionResponse
	ctx.ShouldBindJSON(&webauthnJson)
	fmt.Println(webauthnJson)
	pca, _ := webauthnJson.Parse()

	DB := database.GetDB()
	var challenge model.Challenge
	DB.Table("challenges").Where("challenge = ?", pca.Response.CollectedClientData.Challenge).First(&challenge)

	var user model.User
	var cred model.Certificate
	DB.Table("users").Where("name = ?", challenge.Username).First(&user)
	DB.Table("certificates").Where("user_id = ?", user.ID).Scan(&cred)
	fmt.Println(cred.Credention)
	err := pca.Verify(pca.Response.CollectedClientData.Challenge, W.Config.RPID, []string{pca.Response.CollectedClientData.Origin}, "", false, cred.Credention)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": "422",
			"data": err,
			"msg":  "verify error",
		})
	}
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "login success",
		"msg":  "webauthn",
	})
}
